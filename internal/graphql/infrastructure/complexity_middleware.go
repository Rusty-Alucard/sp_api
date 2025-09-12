// Package graphql provides GraphQL infrastructure components.
package graphql

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	authapp "github.com/Rusty-Alucard/sp_api/internal/auth/application"
	authinfra "github.com/Rusty-Alucard/sp_api/internal/auth/infrastructure"
	graphqldomain "github.com/Rusty-Alucard/sp_api/internal/graphql/domain"
)

// ComplexityMiddleware provides GraphQL query complexity limiting.
type ComplexityMiddleware struct {
	analyzer *graphqldomain.QueryAnalyzer
}

// NewComplexityMiddleware creates a new ComplexityMiddleware.
func NewComplexityMiddleware() *ComplexityMiddleware {
	analyzer := graphqldomain.NewQueryAnalyzer()

	// Add complexity rules for common fields
	analyzer.AddComplexityRule("Query.events", graphqldomain.ComplexityRule{
		BaseComplexity: 5,
		Multiplier:     "limit",
	})

	analyzer.AddComplexityRule("Query.teams", graphqldomain.ComplexityRule{
		BaseComplexity: 5,
		Multiplier:     "limit",
	})

	analyzer.AddComplexityRule("Event.participants", graphqldomain.ComplexityRule{
		BaseComplexity: 10,
	})

	analyzer.AddComplexityRule("Event.matches", graphqldomain.ComplexityRule{
		BaseComplexity: 10,
		Multiplier:     "limit",
	})

	analyzer.AddComplexityRule("Team.members", graphqldomain.ComplexityRule{
		BaseComplexity: 5,
	})

	analyzer.AddComplexityRule("Match.timelines", graphqldomain.ComplexityRule{
		BaseComplexity: 5,
		Multiplier:     "limit",
	})

	return &ComplexityMiddleware{
		analyzer: analyzer,
	}
}

// ComplexityLimitFunc returns a complexity limit function for gqlgen.
func (m *ComplexityMiddleware) ComplexityLimitFunc() func(ctx context.Context, query string) (int, *gqlerror.Error) {
	return func(ctx context.Context, query string) (int, *gqlerror.Error) {
		// Get agent claims from context
		claims, ok := authinfra.GetAgentClaimsFromContext(ctx)
		if !ok {
			// No authentication, use default limit
			return 50, nil
		}

		// Get complexity limit based on agent's plan
		limit := m.getComplexityLimit(claims)

		// Parse variables from context if available
		variables := m.getVariablesFromContext(ctx)

		// Analyze query complexity
		analysis, err := m.analyzer.AnalyzeQuery(ctx, query, variables)
		if err != nil {
			return 0, &gqlerror.Error{
				Message: "failed to analyze query complexity",
				Extensions: map[string]interface{}{
					"code": "QUERY_ANALYSIS_ERROR",
				},
			}
		}

		// Check if query exceeds limit
		if !analysis.IsWithinLimit(limit) {
			return 0, &gqlerror.Error{
				Message: "query complexity exceeds limit",
				Extensions: map[string]interface{}{
					"code":            "QUERY_COMPLEXITY_EXCEEDED",
					"complexity":      analysis.Complexity,
					"limit":           limit,
					"field_breakdown": analysis.GetComplexityBreakdown(),
				},
			}
		}

		return limit, nil
	}
}

// getComplexityLimit returns the complexity limit for an agent based on their plan.
func (m *ComplexityMiddleware) getComplexityLimit(claims *authapp.AgentClaims) int {
	switch claims.Plan {
	case "free":
		return 50
	case "standard":
		return 200
	case "premium":
		return 1000
	case "enterprise":
		return 5000
	default:
		return 10
	}
}

// getVariablesFromContext extracts variables from the GraphQL context.
func (m *ComplexityMiddleware) getVariablesFromContext(ctx context.Context) map[string]interface{} {
	// Try to get variables from GraphQL context
	if variables := graphql.GetOperationContext(ctx).Variables; variables != nil {
		return variables
	}

	// Try to get variables from HTTP request context
	if req := graphql.GetRequestContext(ctx); req != nil {
		if req.Variables != nil {
			return req.Variables
		}
	}

	return nil
}

// HTTPComplexityMiddleware provides HTTP-level complexity checking.
type HTTPComplexityMiddleware struct {
	analyzer *graphqldomain.QueryAnalyzer
}

// NewHTTPComplexityMiddleware creates a new HTTPComplexityMiddleware.
func NewHTTPComplexityMiddleware() *HTTPComplexityMiddleware {
	analyzer := graphqldomain.NewQueryAnalyzer()

	// Add the same complexity rules as the GraphQL middleware
	analyzer.AddComplexityRule("Query.events", graphqldomain.ComplexityRule{
		BaseComplexity: 5,
		Multiplier:     "limit",
	})

	analyzer.AddComplexityRule("Query.teams", graphqldomain.ComplexityRule{
		BaseComplexity: 5,
		Multiplier:     "limit",
	})

	analyzer.AddComplexityRule("Event.participants", graphqldomain.ComplexityRule{
		BaseComplexity: 10,
	})

	analyzer.AddComplexityRule("Event.matches", graphqldomain.ComplexityRule{
		BaseComplexity: 10,
		Multiplier:     "limit",
	})

	analyzer.AddComplexityRule("Team.members", graphqldomain.ComplexityRule{
		BaseComplexity: 5,
	})

	analyzer.AddComplexityRule("Match.timelines", graphqldomain.ComplexityRule{
		BaseComplexity: 5,
		Multiplier:     "limit",
	})

	return &HTTPComplexityMiddleware{
		analyzer: analyzer,
	}
}

// Middleware returns an HTTP middleware function.
func (m *HTTPComplexityMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only check GraphQL requests
			if r.URL.Path != "/query" && r.URL.Path != "/graphql" {
				next.ServeHTTP(w, r)
				return
			}

			// Parse request body to get query and variables
			var requestBody struct {
				Query     string                 `json:"query"`
				Variables map[string]interface{} `json:"variables"`
			}

			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				http.Error(w, "invalid request body", http.StatusBadRequest)
				return
			}

			// Get agent claims from context
			claims, ok := authinfra.GetAgentClaimsFromContext(r.Context())
			if !ok {
				// No authentication, use default limit
				next.ServeHTTP(w, r)
				return
			}

			// Get complexity limit based on agent's plan
			limit := m.getComplexityLimit(claims)

			// Analyze query complexity
			analysis, err := m.analyzer.AnalyzeQuery(r.Context(), requestBody.Query, requestBody.Variables)
			if err != nil {
				http.Error(w, "failed to analyze query", http.StatusBadRequest)
				return
			}

			// Check if query exceeds limit
			if !analysis.IsWithinLimit(limit) {
				response := map[string]interface{}{
					"errors": []map[string]interface{}{
						{
							"message": "query complexity exceeds limit",
							"extensions": map[string]interface{}{
								"code":            "QUERY_COMPLEXITY_EXCEEDED",
								"complexity":      analysis.Complexity,
								"limit":           limit,
								"field_breakdown": analysis.GetComplexityBreakdown(),
							},
						},
					},
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getComplexityLimit returns the complexity limit for an agent based on their plan.
func (m *HTTPComplexityMiddleware) getComplexityLimit(claims *authapp.AgentClaims) int {
	switch claims.Plan {
	case "free":
		return 50
	case "standard":
		return 200
	case "premium":
		return 1000
	case "enterprise":
		return 5000
	default:
		return 10
	}
}
