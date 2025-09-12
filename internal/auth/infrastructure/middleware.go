// Package auth provides authentication infrastructure components.
package auth

import (
	"context"
	"net/http"
	"strings"

	authapp "github.com/Rusty-Alucard/sp_api/internal/auth/application"
	authdomain "github.com/Rusty-Alucard/sp_api/internal/auth/domain"
)

// ContextKey represents a context key type.
type ContextKey string

const (
	// AgentClaimsKey is the context key for agent claims.
	AgentClaimsKey ContextKey = "agent_claims"
	// AgentKey is the context key for agent.
	AgentKey ContextKey = "agent"
)

// AuthMiddleware provides authentication middleware for HTTP handlers.
type AuthMiddleware struct {
	authService *authapp.AuthService
}

// NewAuthMiddleware creates a new AuthMiddleware.
func NewAuthMiddleware(authService *authapp.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware that requires authentication.
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractBearerToken(r)
		if token == "" {
			m.writeUnauthorizedResponse(w, "missing authorization token")
			return
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			m.writeUnauthorizedResponse(w, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), AgentClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireScope middleware that requires a specific scope.
func (m *AuthMiddleware) RequireScope(scope authdomain.Scope) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(AgentClaimsKey).(*authapp.AgentClaims)
			if !ok {
				m.writeUnauthorizedResponse(w, "authentication required")
				return
			}

			if !m.authService.HasScope(claims, scope) {
				m.writeForbiddenResponse(w, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireResourceAccess middleware that requires access to a specific resource.
func (m *AuthMiddleware) RequireResourceAccess(resource string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(AgentClaimsKey).(*authapp.AgentClaims)
			if !ok {
				m.writeUnauthorizedResponse(w, "authentication required")
				return
			}

			if !m.authService.CanAccess(claims, resource) {
				m.writeForbiddenResponse(w, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireWriteAccess middleware that requires write access to a specific resource.
func (m *AuthMiddleware) RequireWriteAccess(resource string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(AgentClaimsKey).(*authapp.AgentClaims)
			if !ok {
				m.writeUnauthorizedResponse(w, "authentication required")
				return
			}

			if !m.authService.CanWrite(claims, resource) {
				m.writeForbiddenResponse(w, "insufficient write permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuth middleware that optionally authenticates the request.
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractBearerToken(r)
		if token != "" {
			claims, err := m.authService.ValidateToken(token)
			if err == nil {
				ctx := context.WithValue(r.Context(), AgentClaimsKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// extractBearerToken extracts the bearer token from the Authorization header.
func (m *AuthMiddleware) extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	const bearerParts = 2 // Split into 2 parts: "Bearer" and token
	parts := strings.SplitN(authHeader, " ", bearerParts)
	if len(parts) != bearerParts || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// writeUnauthorizedResponse writes an unauthorized response.
func (m *AuthMiddleware) writeUnauthorizedResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error": "unauthorized", "message": "` + message + `"}`))
}

// writeForbiddenResponse writes a forbidden response.
func (m *AuthMiddleware) writeForbiddenResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error": "forbidden", "message": "` + message + `"}`))
}

// GetAgentClaimsFromContext extracts agent claims from the request context.
func GetAgentClaimsFromContext(ctx context.Context) (*authapp.AgentClaims, bool) {
	claims, ok := ctx.Value(AgentClaimsKey).(*authapp.AgentClaims)
	return claims, ok
}

// GetAgentFromContext extracts agent from the request context.
func GetAgentFromContext(ctx context.Context) (*authdomain.Agent, bool) {
	agent, ok := ctx.Value(AgentKey).(*authdomain.Agent)
	return agent, ok
}
