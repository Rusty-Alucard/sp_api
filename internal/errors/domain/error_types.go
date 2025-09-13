// Package errors provides structured error handling for AI agents.
package errors

import (
	"fmt"
	"time"
)

// ErrorCode represents a standardized error code.
type ErrorCode string

const (
	// ErrorCodeUnauthorized indicates authentication is required
	ErrorCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	// ErrorCodeInvalidToken indicates the provided token is invalid
	ErrorCodeInvalidToken       ErrorCode = "INVALID_TOKEN"
	// ErrorCodeTokenExpired indicates the token has expired
	ErrorCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	// ErrorCodeInsufficientScopes indicates insufficient permissions
	ErrorCodeInsufficientScopes ErrorCode = "INSUFFICIENT_SCOPES"

	// ErrorCodeRateLimitExceeded indicates rate limit has been exceeded
	ErrorCodeRateLimitExceeded       ErrorCode = "RATE_LIMIT_EXCEEDED"
	// ErrorCodeQueryComplexityExceeded indicates query complexity limit exceeded
	ErrorCodeQueryComplexityExceeded ErrorCode = "QUERY_COMPLEXITY_EXCEEDED"

	// ErrorCodeValidationFailed indicates input validation failed
	ErrorCodeValidationFailed     ErrorCode = "VALIDATION_FAILED"
	ErrorCodeInvalidInput         ErrorCode = "INVALID_INPUT"
	ErrorCodeMissingRequiredField ErrorCode = "MISSING_REQUIRED_FIELD"

	// ErrorCodeBusinessRuleViolation indicates a business rule was violated
	ErrorCodeBusinessRuleViolation ErrorCode = "BUSINESS_RULE_VIOLATION"
	ErrorCodeResourceNotFound      ErrorCode = "RESOURCE_NOT_FOUND"
	ErrorCodeResourceConflict      ErrorCode = "RESOURCE_CONFLICT"
	ErrorCodeInsufficientPoints    ErrorCode = "INSUFFICIENT_POINTS"

	// ErrorCodeInternalServerError indicates an internal server error occurred
	ErrorCodeInternalServerError  ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeServiceUnavailable   ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeDatabaseError        ErrorCode = "DATABASE_ERROR"
	ErrorCodeExternalServiceError ErrorCode = "EXTERNAL_SERVICE_ERROR"
)

// ErrorSeverity represents the severity level of an error.
type ErrorSeverity string

const (
	// SeverityLow indicates low severity error
	SeverityLow      ErrorSeverity = "low"
	// SeverityMedium indicates medium severity error
	SeverityMedium   ErrorSeverity = "medium"
	// SeverityHigh indicates high severity error
	SeverityHigh     ErrorSeverity = "high"
	// SeverityCritical indicates critical severity error
	SeverityCritical ErrorSeverity = "critical"
)

// StructuredError represents a structured error for AI agents.
type StructuredError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Severity   ErrorSeverity          `json:"severity"`
	Timestamp  time.Time              `json:"timestamp"`
	RequestID  string                 `json:"request_id,omitempty"`
	Path       []string               `json:"path,omitempty"`
	Retryable  bool                   `json:"retryable"`
	RetryAfter *time.Duration         `json:"retry_after,omitempty"`
}

// Error implements the error interface.
func (e *StructuredError) Error() string {
	return e.Message
}

// NewStructuredError creates a new StructuredError.
func NewStructuredError(code ErrorCode, message string) *StructuredError {
	return &StructuredError{
		Code:      code,
		Message:   message,
		Details:   make(map[string]interface{}),
		Severity:  getSeverityForCode(code),
		Timestamp: time.Now(),
		Retryable: isRetryableCode(code),
	}
}

// WithDetails adds details to the error.
func (e *StructuredError) WithDetails(details map[string]interface{}) *StructuredError {
	e.Details = details
	return e
}

// WithRequestID adds a request ID to the error.
func (e *StructuredError) WithRequestID(requestID string) *StructuredError {
	e.RequestID = requestID
	return e
}

// WithPath adds a path to the error.
func (e *StructuredError) WithPath(path []string) *StructuredError {
	e.Path = path
	return e
}

// WithRetryAfter adds a retry after duration to the error.
func (e *StructuredError) WithRetryAfter(duration time.Duration) *StructuredError {
	e.RetryAfter = &duration
	return e
}

// getSeverityForCode returns the severity level for an error code.
func getSeverityForCode(code ErrorCode) ErrorSeverity {
	switch code {
	case ErrorCodeUnauthorized, ErrorCodeInvalidToken, ErrorCodeTokenExpired:
		return SeverityMedium
	case ErrorCodeInsufficientScopes:
		return SeverityMedium
	case ErrorCodeRateLimitExceeded, ErrorCodeQueryComplexityExceeded:
		return SeverityMedium
	case ErrorCodeValidationFailed, ErrorCodeInvalidInput, ErrorCodeMissingRequiredField:
		return SeverityLow
	case ErrorCodeBusinessRuleViolation, ErrorCodeResourceNotFound, ErrorCodeResourceConflict:
		return SeverityMedium
	case ErrorCodeInsufficientPoints:
		return SeverityLow
	case ErrorCodeInternalServerError, ErrorCodeServiceUnavailable, ErrorCodeDatabaseError:
		return SeverityHigh
	case ErrorCodeExternalServiceError:
		return SeverityMedium
	default:
		return SeverityMedium
	}
}

// isRetryableCode determines if an error code is retryable.
func isRetryableCode(code ErrorCode) bool {
	switch code {
	case ErrorCodeUnauthorized, ErrorCodeInvalidToken, ErrorCodeTokenExpired, ErrorCodeInsufficientScopes:
		return false
	case ErrorCodeRateLimitExceeded, ErrorCodeQueryComplexityExceeded:
		return true
	case ErrorCodeValidationFailed, ErrorCodeInvalidInput, ErrorCodeMissingRequiredField:
		return false
	case ErrorCodeBusinessRuleViolation, ErrorCodeResourceNotFound, ErrorCodeResourceConflict:
		return false
	case ErrorCodeInsufficientPoints:
		return false
	case ErrorCodeInternalServerError, ErrorCodeServiceUnavailable, ErrorCodeDatabaseError:
		return true
	case ErrorCodeExternalServiceError:
		return true
	default:
		return false
	}
}

// Predefined error constructors for common scenarios

// NewUnauthorizedError creates an unauthorized error.
func NewUnauthorizedError(message string) *StructuredError {
	if message == "" {
		message = "Authentication required"
	}
	return NewStructuredError(ErrorCodeUnauthorized, message)
}

// NewInvalidTokenError creates an invalid token error.
func NewInvalidTokenError() *StructuredError {
	return NewStructuredError(ErrorCodeInvalidToken, "Invalid or malformed token")
}

// NewTokenExpiredError creates a token expired error.
func NewTokenExpiredError() *StructuredError {
	return NewStructuredError(ErrorCodeTokenExpired, "Token has expired")
}

// NewInsufficientScopesError creates an insufficient scopes error.
func NewInsufficientScopesError(requiredScopes []string) *StructuredError {
	return NewStructuredError(ErrorCodeInsufficientScopes, "Insufficient permissions").
		WithDetails(map[string]interface{}{
			"required_scopes": requiredScopes,
		})
}

// NewRateLimitExceededError creates a rate limit exceeded error.
func NewRateLimitExceededError(retryAfter time.Duration, limit, remaining int) *StructuredError {
	return NewStructuredError(ErrorCodeRateLimitExceeded, "Rate limit exceeded").
		WithDetails(map[string]interface{}{
			"limit":     limit,
			"remaining": remaining,
		}).
		WithRetryAfter(retryAfter)
}

// NewQueryComplexityExceededError creates a query complexity exceeded error.
func NewQueryComplexityExceededError(complexity, limit int, breakdown map[string]int) *StructuredError {
	return NewStructuredError(ErrorCodeQueryComplexityExceeded, "Query complexity exceeds limit").
		WithDetails(map[string]interface{}{
			"complexity":      complexity,
			"limit":           limit,
			"field_breakdown": breakdown,
		})
}

// NewValidationError creates a validation error.
func NewValidationError(message string, fieldErrors map[string]string) *StructuredError {
	return NewStructuredError(ErrorCodeValidationFailed, message).
		WithDetails(map[string]interface{}{
			"field_errors": fieldErrors,
		})
}

// NewResourceNotFoundError creates a resource not found error.
func NewResourceNotFoundError(resourceType, resourceID string) *StructuredError {
	return NewStructuredError(ErrorCodeResourceNotFound, fmt.Sprintf("%s not found", resourceType)).
		WithDetails(map[string]interface{}{
			"resource_type": resourceType,
			"resource_id":   resourceID,
		})
}

// NewInsufficientPointsError creates an insufficient points error.
func NewInsufficientPointsError(required, current int) *StructuredError {
	return NewStructuredError(ErrorCodeInsufficientPoints, "Insufficient points").
		WithDetails(map[string]interface{}{
			"required_points": required,
			"current_points":  current,
		})
}

// NewInternalServerError creates an internal server error.
func NewInternalServerError(message string) *StructuredError {
	if message == "" {
		message = "An internal server error occurred"
	}
	return NewStructuredError(ErrorCodeInternalServerError, message)
}

