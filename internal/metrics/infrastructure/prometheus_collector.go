// Package metrics provides Prometheus metrics collection implementation.
package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// SuccessTrue represents successful operation
	SuccessTrue = "true"
	// SuccessFalse represents failed operation
	SuccessFalse = "false"
)

// PrometheusCollector implements MetricsCollector using Prometheus.
type PrometheusCollector struct {
	// GraphQL metrics
	graphqlRequestDuration *prometheus.HistogramVec
	graphqlRequestTotal    *prometheus.CounterVec
	graphqlQueryComplexity *prometheus.HistogramVec
	graphqlComplexityLimit *prometheus.GaugeVec

	// Authentication metrics
	authAttemptsTotal     *prometheus.CounterVec
	tokenValidationsTotal *prometheus.CounterVec

	// Rate limiting metrics
	rateLimitHitsTotal     *prometheus.CounterVec
	rateLimitExceededTotal *prometheus.CounterVec

	// DataLoader metrics
	dataloaderBatchDuration    *prometheus.HistogramVec
	dataloaderBatchSize        *prometheus.HistogramVec
	dataloaderCacheHitsTotal   *prometheus.CounterVec
	dataloaderCacheMissesTotal *prometheus.CounterVec

	// APQ metrics
	apqQueriesTotal     *prometheus.CounterVec
	apqCacheHitsTotal   *prometheus.CounterVec
	apqCacheMissesTotal *prometheus.CounterVec

	// Business metrics
	predictionsCreatedTotal *prometheus.CounterVec
	predictionsResultTotal  *prometheus.CounterVec
	pointsAwardedTotal      *prometheus.CounterVec

	// System metrics
	databaseQueryDuration   *prometheus.HistogramVec
	databaseQueriesTotal    *prometheus.CounterVec
	externalAPICallDuration *prometheus.HistogramVec
	externalAPICallsTotal   *prometheus.CounterVec
}

// NewPrometheusCollector creates a new PrometheusCollector.
func NewPrometheusCollector() *PrometheusCollector {
	return &PrometheusCollector{
		// GraphQL metrics
		graphqlRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "graphql_request_duration_seconds",
				Help:    "GraphQL request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "success", "agent_id"},
		),

		graphqlRequestTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "graphql_requests_total",
				Help: "Total number of GraphQL requests",
			},
			[]string{"operation", "success", "agent_id"},
		),

		graphqlQueryComplexity: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "graphql_query_complexity",
				Help:    "GraphQL query complexity",
				Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
			},
			[]string{"agent_id", "plan"},
		),

		graphqlComplexityLimit: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "graphql_complexity_limit",
				Help: "GraphQL complexity limit for agent",
			},
			[]string{"agent_id", "plan"},
		),

		// Authentication metrics
		authAttemptsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "auth_attempts_total",
				Help: "Total number of authentication attempts",
			},
			[]string{"success", "method"},
		),

		tokenValidationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "token_validations_total",
				Help: "Total number of token validations",
			},
			[]string{"success", "agent_id"},
		),

		// Rate limiting metrics
		rateLimitHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limit_hits_total",
				Help: "Total number of rate limit hits",
			},
			[]string{"agent_id", "limit_type"},
		),

		rateLimitExceededTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limit_exceeded_total",
				Help: "Total number of rate limit exceeded events",
			},
			[]string{"agent_id", "limit_type"},
		),

		// DataLoader metrics
		dataloaderBatchDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "dataloader_batch_duration_seconds",
				Help:    "DataLoader batch duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"loader_type"},
		),

		dataloaderBatchSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "dataloader_batch_size",
				Help:    "DataLoader batch size",
				Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"loader_type"},
		),

		dataloaderCacheHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "dataloader_cache_hits_total",
				Help: "Total number of DataLoader cache hits",
			},
			[]string{"loader_type"},
		),

		dataloaderCacheMissesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "dataloader_cache_misses_total",
				Help: "Total number of DataLoader cache misses",
			},
			[]string{"loader_type"},
		),

		// APQ metrics
		apqQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "apq_queries_total",
				Help: "Total number of APQ queries",
			},
			[]string{"is_persisted", "success"},
		),

		apqCacheHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "apq_cache_hits_total",
				Help: "Total number of APQ cache hits",
			},
			[]string{},
		),

		apqCacheMissesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "apq_cache_misses_total",
				Help: "Total number of APQ cache misses",
			},
			[]string{},
		),

		// Business metrics
		predictionsCreatedTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "predictions_created_total",
				Help: "Total number of predictions created",
			},
			[]string{"agent_id", "prediction_type"},
		),

		predictionsResultTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "predictions_result_total",
				Help: "Total number of prediction results",
			},
			[]string{"agent_id", "correct"},
		),

		pointsAwardedTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "points_awarded_total",
				Help: "Total points awarded",
			},
			[]string{"agent_id"},
		),

		// System metrics
		databaseQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "database_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"query_type", "success"},
		),

		databaseQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "database_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"query_type", "success"},
		),

		externalAPICallDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "external_api_call_duration_seconds",
				Help:    "External API call duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "success"},
		),

		externalAPICallsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "external_api_calls_total",
				Help: "Total number of external API calls",
			},
			[]string{"service", "success"},
		),
	}
}

// RecordGraphQLRequest records GraphQL request metrics.
func (c *PrometheusCollector) RecordGraphQLRequest(ctx context.Context, operation string, complexity int, duration time.Duration, success bool) {
	agentID := c.getAgentIDFromContext(ctx)
	successStr := SuccessTrue
	if !success {
		successStr = SuccessFalse
	}

	c.graphqlRequestDuration.WithLabelValues(operation, successStr, agentID).Observe(duration.Seconds())
	c.graphqlRequestTotal.WithLabelValues(operation, successStr, agentID).Inc()
}

// RecordGraphQLQueryComplexity records GraphQL query complexity metrics.
func (c *PrometheusCollector) RecordGraphQLQueryComplexity(ctx context.Context, complexity int, limit int) {
	agentID := c.getAgentIDFromContext(ctx)
	plan := c.getAgentPlanFromContext(ctx)

	c.graphqlQueryComplexity.WithLabelValues(agentID, plan).Observe(float64(complexity))
	c.graphqlComplexityLimit.WithLabelValues(agentID, plan).Set(float64(limit))
}

// RecordAuthenticationAttempt records authentication attempt metrics.
func (c *PrometheusCollector) RecordAuthenticationAttempt(ctx context.Context, success bool, agentID string) {
	successStr := SuccessTrue
	if !success {
		successStr = SuccessFalse
	}

	c.authAttemptsTotal.WithLabelValues(successStr, "client_credentials").Inc()
}

// RecordTokenValidation records token validation metrics.
func (c *PrometheusCollector) RecordTokenValidation(ctx context.Context, success bool, agentID string) {
	successStr := SuccessTrue
	if !success {
		successStr = SuccessFalse
	}

	c.tokenValidationsTotal.WithLabelValues(successStr, agentID).Inc()
}

// RecordRateLimitHit records rate limit hit metrics.
func (c *PrometheusCollector) RecordRateLimitHit(ctx context.Context, agentID string, limitType string) {
	c.rateLimitHitsTotal.WithLabelValues(agentID, limitType).Inc()
}

// RecordRateLimitExceeded records rate limit exceeded metrics.
func (c *PrometheusCollector) RecordRateLimitExceeded(ctx context.Context, agentID string, limitType string) {
	c.rateLimitExceededTotal.WithLabelValues(agentID, limitType).Inc()
}

// RecordDataLoaderBatch records DataLoader batch metrics.
func (c *PrometheusCollector) RecordDataLoaderBatch(ctx context.Context, loaderType string, batchSize int, duration time.Duration) {
	c.dataloaderBatchDuration.WithLabelValues(loaderType).Observe(duration.Seconds())
	c.dataloaderBatchSize.WithLabelValues(loaderType).Observe(float64(batchSize))
}

// RecordDataLoaderCacheHit records DataLoader cache hit metrics.
func (c *PrometheusCollector) RecordDataLoaderCacheHit(ctx context.Context, loaderType string) {
	c.dataloaderCacheHitsTotal.WithLabelValues(loaderType).Inc()
}

// RecordDataLoaderCacheMiss records DataLoader cache miss metrics.
func (c *PrometheusCollector) RecordDataLoaderCacheMiss(ctx context.Context, loaderType string) {
	c.dataloaderCacheMissesTotal.WithLabelValues(loaderType).Inc()
}

// RecordAPQQuery records APQ query metrics.
func (c *PrometheusCollector) RecordAPQQuery(ctx context.Context, isPersisted bool, success bool) {
	persistedStr := "false"
	if isPersisted {
		persistedStr = "true"
	}

	successStr := SuccessTrue
	if !success {
		successStr = SuccessFalse
	}

	c.apqQueriesTotal.WithLabelValues(persistedStr, successStr).Inc()
}

// RecordAPQCacheHit records APQ cache hit metrics.
func (c *PrometheusCollector) RecordAPQCacheHit(ctx context.Context) {
	c.apqCacheHitsTotal.WithLabelValues().Inc()
}

// RecordAPQCacheMiss records APQ cache miss metrics.
func (c *PrometheusCollector) RecordAPQCacheMiss(ctx context.Context) {
	c.apqCacheMissesTotal.WithLabelValues().Inc()
}

// RecordPredictionCreated records prediction creation metrics.
func (c *PrometheusCollector) RecordPredictionCreated(ctx context.Context, agentID string, predictionType string) {
	c.predictionsCreatedTotal.WithLabelValues(agentID, predictionType).Inc()
}

// RecordPredictionResult records prediction result metrics.
func (c *PrometheusCollector) RecordPredictionResult(ctx context.Context, agentID string, correct bool, points int) {
	correctStr := "true"
	if !correct {
		correctStr = "false"
	}

	c.predictionsResultTotal.WithLabelValues(agentID, correctStr).Inc()
	c.pointsAwardedTotal.WithLabelValues(agentID).Add(float64(points))
}

// RecordDatabaseQuery records database query metrics.
func (c *PrometheusCollector) RecordDatabaseQuery(ctx context.Context, queryType string, duration time.Duration, success bool) {
	successStr := SuccessTrue
	if !success {
		successStr = SuccessFalse
	}

	c.databaseQueryDuration.WithLabelValues(queryType, successStr).Observe(duration.Seconds())
	c.databaseQueriesTotal.WithLabelValues(queryType, successStr).Inc()
}

// RecordExternalAPICall records external API call metrics.
func (c *PrometheusCollector) RecordExternalAPICall(ctx context.Context, service string, duration time.Duration, success bool) {
	successStr := SuccessTrue
	if !success {
		successStr = SuccessFalse
	}

	c.externalAPICallDuration.WithLabelValues(service, successStr).Observe(duration.Seconds())
	c.externalAPICallsTotal.WithLabelValues(service, successStr).Inc()
}

// getAgentIDFromContext extracts the agent ID from the context.
func (c *PrometheusCollector) getAgentIDFromContext(ctx context.Context) string {
	if agentClaims := ctx.Value("agent_claims"); agentClaims != nil {
		if claims, ok := agentClaims.(map[string]interface{}); ok {
			if agentID, ok := claims["agent_id"].(string); ok {
				return agentID
			}
		}
	}
	return "unknown"
}

// getAgentPlanFromContext extracts the agent plan from the context.
func (c *PrometheusCollector) getAgentPlanFromContext(ctx context.Context) string {
	if agentClaims := ctx.Value("agent_claims"); agentClaims != nil {
		if claims, ok := agentClaims.(map[string]interface{}); ok {
			if plan, ok := claims["plan"].(string); ok {
				return plan
			}
		}
	}
	return "unknown"
}

