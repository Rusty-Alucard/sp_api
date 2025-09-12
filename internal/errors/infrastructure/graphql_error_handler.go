// Package errors provides error handling infrastructure for GraphQL.
package errors

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"

	errors "github.com/Rusty-Alucard/sp_api/internal/errors/domain"
)

const (
	// RetryAfterTimeout represents retry after timeout duration
	RetryAfterTimeout = 30 * time.Second
	// RetryAfterConnectionRefused represents retry after connection refused duration
	RetryAfterConnectionRefused = 60 * time.Second
)

// GraphQLErrorHandler handles GraphQL errors for AI agents.
type GraphQLErrorHandler struct {
	logger Logger
}

// Logger defines the interface for logging.
type Logger interface {
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// NewGraphQLErrorHandler creates a new GraphQLErrorHandler.
func NewGraphQLErrorHandler(logger Logger) *GraphQLErrorHandler {
	return &GraphQLErrorHandler{
		logger: logger,
	}
}

// ErrorPresenter returns a GraphQL error presenter function.
func (h *GraphQLErrorHandler) ErrorPresenter() graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		// Extract request ID from context
		requestID := h.getRequestIDFromContext(ctx)

		// Convert error to structured error
		structuredErr := h.convertToStructuredError(err, requestID)

		// Log the error
		h.logError(structuredErr, ctx)

		// Convert to GraphQL error
		return h.convertToGraphQLError(structuredErr)
	}
}

// convertToStructuredError converts a Go error to a StructuredError.
func (h *GraphQLErrorHandler) convertToStructuredError(err error, requestID string) *errors.StructuredError {
	// Check if it's already a StructuredError
	if structuredErr, ok := err.(*errors.StructuredError); ok {
		if requestID != "" {
			structuredErr.RequestID = requestID
		}
		return structuredErr
	}

	// Check for specific error types and convert them
	switch err.Error() {
	case "context deadline exceeded":
		return errors.NewStructuredError(errors.ErrorCodeServiceUnavailable, "Request timeout").
			WithRequestID(requestID).
			WithRetryAfter(RetryAfterTimeout)
	case "connection refused":
		return errors.NewStructuredError(errors.ErrorCodeServiceUnavailable, "Service temporarily unavailable").
			WithRequestID(requestID).
			WithRetryAfter(RetryAfterConnectionRefused)
	default:
		// Default to internal server error
		return errors.NewStructuredError(errors.ErrorCodeInternalServerError, err.Error()).
			WithRequestID(requestID)
	}
}

// convertToGraphQLError converts a StructuredError to a GraphQL error.
func (h *GraphQLErrorHandler) convertToGraphQLError(err *errors.StructuredError) *gqlerror.Error {
	graphqlErr := &gqlerror.Error{
		Message: err.Message,
		Extensions: map[string]interface{}{
			"code":      string(err.Code),
			"severity":  string(err.Severity),
			"timestamp": err.Timestamp,
			"retryable": err.Retryable,
		},
	}

	// Add optional fields
	if err.RequestID != "" {
		graphqlErr.Extensions["request_id"] = err.RequestID
	}

	if len(err.Path) > 0 {
		path := make(ast.Path, len(err.Path))
		for i, p := range err.Path {
			path[i] = ast.PathName(p)
		}
		graphqlErr.Path = path
	}

	if err.RetryAfter != nil {
		graphqlErr.Extensions["retry_after"] = int(err.RetryAfter.Seconds())
	}

	if len(err.Details) > 0 {
		graphqlErr.Extensions["details"] = err.Details
	}

	return graphqlErr
}

// getRequestIDFromContext extracts the request ID from the context.
func (h *GraphQLErrorHandler) getRequestIDFromContext(ctx context.Context) string {
	if requestID := ctx.Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// logError logs the error with appropriate level based on severity.
func (h *GraphQLErrorHandler) logError(err *errors.StructuredError, ctx context.Context) {
	fields := []interface{}{
		"code", string(err.Code),
		"severity", string(err.Severity),
		"timestamp", err.Timestamp,
		"retryable", err.Retryable,
	}

	if err.RequestID != "" {
		fields = append(fields, "request_id", err.RequestID)
	}

	if len(err.Details) > 0 {
		fields = append(fields, "details", err.Details)
	}

	// Add context information
	if agentID := h.getAgentIDFromContext(ctx); agentID != "" {
		fields = append(fields, "agent_id", agentID)
	}

	// Log with appropriate level
	switch err.Severity {
	case errors.SeverityCritical:
		h.logger.Error(fmt.Sprintf("Critical error: %s", err.Message), fields...)
	case errors.SeverityHigh:
		h.logger.Error(fmt.Sprintf("High severity error: %s", err.Message), fields...)
	case errors.SeverityMedium:
		h.logger.Warn(fmt.Sprintf("Medium severity error: %s", err.Message), fields...)
	case errors.SeverityLow:
		h.logger.Info(fmt.Sprintf("Low severity error: %s", err.Message), fields...)
	default:
		h.logger.Error(fmt.Sprintf("Unknown severity error: %s", err.Message), fields...)
	}
}

// getAgentIDFromContext extracts the agent ID from the context.
func (h *GraphQLErrorHandler) getAgentIDFromContext(ctx context.Context) string {
	if agentClaims := ctx.Value("agent_claims"); agentClaims != nil {
		// Assuming agent claims has an AgentID field
		// This would need to be adapted based on your actual claims structure
		if claims, ok := agentClaims.(map[string]interface{}); ok {
			if agentID, ok := claims["agent_id"].(string); ok {
				return agentID
			}
		}
	}
	return ""
}

// DefaultLogger provides a default logger implementation.
type DefaultLogger struct{}

// Error logs an error message.
func (l *DefaultLogger) Error(msg string, fields ...interface{}) {
	log.Printf("ERROR: %s %v", msg, fields)
}

// Warn logs a warning message.
func (l *DefaultLogger) Warn(msg string, fields ...interface{}) {
	log.Printf("WARN: %s %v", msg, fields)
}

// Info logs an info message.
func (l *DefaultLogger) Info(msg string, fields ...interface{}) {
	log.Printf("INFO: %s %v", msg, fields)
}

// Debug logs a debug message.
func (l *DefaultLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("DEBUG: %s %v", msg, fields)
}
