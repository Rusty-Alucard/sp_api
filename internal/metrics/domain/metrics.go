// Package metrics provides metrics collection for AI agent monitoring.
package metrics

import (
	"context"
	"time"
)

// MetricsCollector defines the interface for metrics collection.
type MetricsCollector interface {
	// GraphQL metrics
	RecordGraphQLRequest(ctx context.Context, operation string, complexity int, duration time.Duration, success bool)
	RecordGraphQLQueryComplexity(ctx context.Context, complexity int, limit int)

	// Authentication metrics
	RecordAuthenticationAttempt(ctx context.Context, success bool, agentID string)
	RecordTokenValidation(ctx context.Context, success bool, agentID string)

	// Rate limiting metrics
	RecordRateLimitHit(ctx context.Context, agentID string, limitType string)
	RecordRateLimitExceeded(ctx context.Context, agentID string, limitType string)

	// DataLoader metrics
	RecordDataLoaderBatch(ctx context.Context, loaderType string, batchSize int, duration time.Duration)
	RecordDataLoaderCacheHit(ctx context.Context, loaderType string)
	RecordDataLoaderCacheMiss(ctx context.Context, loaderType string)

	// APQ metrics
	RecordAPQQuery(ctx context.Context, isPersisted bool, success bool)
	RecordAPQCacheHit(ctx context.Context)
	RecordAPQCacheMiss(ctx context.Context)

	// Business metrics
	RecordPredictionCreated(ctx context.Context, agentID string, predictionType string)
	RecordPredictionResult(ctx context.Context, agentID string, correct bool, points int)

	// System metrics
	RecordDatabaseQuery(ctx context.Context, queryType string, duration time.Duration, success bool)
	RecordExternalAPICall(ctx context.Context, service string, duration time.Duration, success bool)
}

// MetricLabels represents labels for metrics.
type MetricLabels map[string]string

// GraphQLRequestMetrics represents GraphQL request metrics.
type GraphQLRequestMetrics struct {
	Operation  string
	Complexity int
	Duration   time.Duration
	Success    bool
	AgentID    string
	RequestID  string
}

// AuthenticationMetrics represents authentication metrics.
type AuthenticationMetrics struct {
	Success bool
	AgentID string
	Method  string
}

// RateLimitMetrics represents rate limit metrics.
type RateLimitMetrics struct {
	AgentID   string
	LimitType string
	Exceeded  bool
}

// DataLoaderMetrics represents DataLoader metrics.
type DataLoaderMetrics struct {
	LoaderType string
	BatchSize  int
	Duration   time.Duration
	CacheHit   bool
}

// APQMetrics represents APQ metrics.
type APQMetrics struct {
	IsPersisted bool
	Success     bool
	CacheHit    bool
}

// BusinessMetrics represents business metrics.
type BusinessMetrics struct {
	AgentID        string
	PredictionType string
	Correct        bool
	Points         int
}

// SystemMetrics represents system metrics.
type SystemMetrics struct {
	QueryType string
	Duration  time.Duration
	Success   bool
	Service   string
}


