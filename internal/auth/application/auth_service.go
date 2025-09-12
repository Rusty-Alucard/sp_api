// Package auth provides authentication application services.
package auth

import (
	"context"
	"errors"
	"time"

	auth "github.com/Rusty-Alucard/sp_api/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles authentication operations.
type AuthService struct {
	agentRepo   AgentRepository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

// NewAuthService creates a new AuthService.
func NewAuthService(agentRepo AgentRepository, jwtSecret []byte, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		agentRepo:   agentRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

// AgentClaims represents JWT claims for an agent.
type AgentClaims struct {
	AgentID string   `json:"agent_id"`
	Name    string   `json:"name"`
	Scopes  []string `json:"scopes"`
	Plan    string   `json:"plan"`
	jwt.RegisteredClaims
}

// AuthenticateAgent authenticates an agent using client credentials.
func (s *AuthService) AuthenticateAgent(ctx context.Context, clientID, clientSecret string) (*auth.Agent, error) {
	agent, err := s.agentRepo.FindByClientID(ctx, clientID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if agent.ClientSecret() != clientSecret {
		return nil, errors.New("invalid credentials")
	}

	if !agent.IsActive() {
		return nil, errors.New("agent is inactive")
	}

	// Update last used timestamp
	agent.UpdateLastUsed()
	if err := s.agentRepo.Save(ctx, agent); err != nil {
		// Log error but don't fail authentication
		// TODO: Add proper logging
	}

	return agent, nil
}

// GenerateToken generates a JWT token for an agent.
func (s *AuthService) GenerateToken(agent *auth.Agent) (string, error) {
	now := time.Now()
	claims := &AgentClaims{
		AgentID: agent.ClientID(),
		Name:    agent.Name(),
		Scopes:  s.scopesToStrings(agent.Scopes()),
		Plan:    string(agent.Plan()),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "sp-api",
			Subject:   agent.ClientID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims.
func (s *AuthService) ValidateToken(tokenString string) (*AgentClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AgentClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AgentClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new token for an existing valid token.
func (s *AuthService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	agent, err := s.agentRepo.FindByClientID(ctx, claims.AgentID)
	if err != nil {
		return "", errors.New("agent not found")
	}

	if !agent.IsActive() {
		return "", errors.New("agent is inactive")
	}

	return s.GenerateToken(agent)
}

// HasScope checks if the agent has a specific scope.
func (s *AuthService) HasScope(claims *AgentClaims, scope auth.Scope) bool {
	for _, claimScope := range claims.Scopes {
		if claimScope == string(scope) {
			return true
		}
	}
	return false
}

// CanAccess checks if the agent can access a specific resource.
func (s *AuthService) CanAccess(claims *AgentClaims, resource string) bool {
	requiredScope := auth.GetScopeForResource(resource)
	return s.HasScope(claims, requiredScope)
}

// CanWrite checks if the agent can write to a specific resource.
func (s *AuthService) CanWrite(claims *AgentClaims, resource string) bool {
	requiredScope := auth.GetWriteScopeForResource(resource)
	return s.HasScope(claims, requiredScope)
}

// GetComplexityLimit returns the complexity limit for the agent's plan.
func (s *AuthService) GetComplexityLimit(claims *AgentClaims) int {
	plan := auth.Plan(claims.Plan)
	return plan.GetComplexityLimit()
}

// GetRateLimit returns the rate limit for the agent's plan.
func (s *AuthService) GetRateLimit(claims *AgentClaims) auth.RateLimit {
	plan := auth.Plan(claims.Plan)
	return plan.GetRateLimit()
}

// scopesToStrings converts domain scopes to strings.
func (s *AuthService) scopesToStrings(scopes []auth.Scope) []string {
	result := make([]string, len(scopes))
	for i, scope := range scopes {
		result[i] = string(scope)
	}
	return result
}

// AgentRepository defines the interface for agent data access.
type AgentRepository interface {
	FindByClientID(ctx context.Context, clientID string) (*auth.Agent, error)
	Save(ctx context.Context, agent *auth.Agent) error
	FindAll(ctx context.Context) ([]*auth.Agent, error)
	Delete(ctx context.Context, clientID string) error
}
