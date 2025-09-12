// Package auth provides authentication domain models and logic.
package auth

import "time"

// Rate limit constants
const (
	FreeRequestsPerMinute    = 10
	FreeRequestsPerHour      = 100
	FreeRequestsPerDay       = 1000
	StandardRequestsPerMinute = 50
	StandardRequestsPerHour   = 1000
	StandardRequestsPerDay    = 10000
	PremiumRequestsPerMinute  = 200
	PremiumRequestsPerHour    = 5000
	PremiumRequestsPerDay     = 50000
	EnterpriseRequestsPerMinute = 1000
	EnterpriseRequestsPerHour   = 50000
	EnterpriseRequestsPerDay    = 500000
	ReadOnlyRequestsPerHour     = 10
	ReadOnlyRequestsPerDay      = 100
	HoursPerDay                 = 24
)

// Scope represents a permission scope for API access.
type Scope string

const (
	// ScopeReadEvents allows reading event data
	ScopeReadEvents Scope = "events:read"
	// ScopeWriteEvents allows writing event data
	ScopeWriteEvents Scope = "events:write"
	// ScopeReadTeams allows reading team data
	ScopeReadTeams Scope = "teams:read"
	// ScopeWriteTeams allows writing team data
	ScopeWriteTeams Scope = "teams:write"
	// ScopeReadMatches allows reading match data
	ScopeReadMatches Scope = "matches:read"
	// ScopeWriteMatches allows writing match data
	ScopeWriteMatches Scope = "matches:write"
	// ScopeReadPlayers allows reading player data
	ScopeReadPlayers Scope = "players:read"
	// ScopeWritePlayers allows writing player data
	ScopeWritePlayers Scope = "players:write"
	// ScopeReadPredictions allows reading prediction data
	ScopeReadPredictions Scope = "predictions:read"
	// ScopeWritePredictions allows writing prediction data
	ScopeWritePredictions Scope = "predictions:write"
	// ScopeAdmin allows administrative access
	ScopeAdmin Scope = "admin"
)

// Plan represents an agent's subscription plan.
type Plan string

const (
	// PlanFree provides basic access with limited rate limits
	PlanFree Plan = "free"
	// PlanStandard provides standard access with moderate rate limits
	PlanStandard Plan = "standard"
	// PlanPremium provides premium access with high rate limits
	PlanPremium Plan = "premium"
	// PlanEnterprise provides enterprise access with unlimited rate limits
	PlanEnterprise Plan = "enterprise"
)

// String returns the string representation of the scope.
func (s Scope) String() string {
	return string(s)
}

// IsValid checks if the scope is valid.
func (s Scope) IsValid() bool {
	validScopes := []Scope{
		ScopeReadEvents, ScopeWriteEvents,
		ScopeReadTeams, ScopeWriteTeams,
		ScopeReadMatches, ScopeWriteMatches,
		ScopeReadPlayers, ScopeWritePlayers,
		ScopeReadPredictions, ScopeWritePredictions,
		ScopeAdmin,
	}

	for _, validScope := range validScopes {
		if s == validScope {
			return true
		}
	}
	return false
}

// String returns the string representation of the plan.
func (p Plan) String() string {
	return string(p)
}

// IsValid checks if the plan is valid.
func (p Plan) IsValid() bool {
	validPlans := []Plan{
		PlanFree, PlanStandard, PlanPremium, PlanEnterprise,
	}

	for _, validPlan := range validPlans {
		if p == validPlan {
			return true
		}
	}
	return false
}

// GetRateLimit returns the rate limit for a plan.
func (p Plan) GetRateLimit() RateLimit {
	switch p {
	case PlanFree:
		return RateLimit{
			RequestsPerMinute: FreeRequestsPerMinute,
			RequestsPerHour:   FreeRequestsPerHour,
			RequestsPerDay:    FreeRequestsPerDay,
		}
	case PlanStandard:
		return RateLimit{
			RequestsPerMinute: StandardRequestsPerMinute,
			RequestsPerHour:   StandardRequestsPerHour,
			RequestsPerDay:    StandardRequestsPerDay,
		}
	case PlanPremium:
		return RateLimit{
			RequestsPerMinute: PremiumRequestsPerMinute,
			RequestsPerHour:   PremiumRequestsPerHour,
			RequestsPerDay:    PremiumRequestsPerDay,
		}
	case PlanEnterprise:
		return RateLimit{
			RequestsPerMinute: EnterpriseRequestsPerMinute,
			RequestsPerHour:   EnterpriseRequestsPerHour,
			RequestsPerDay:    EnterpriseRequestsPerDay,
		}
	default:
		return RateLimit{
			RequestsPerMinute: 1,
			RequestsPerHour:   10,
			RequestsPerDay:    100,
		}
	}
}

// GetComplexityLimit returns the GraphQL query complexity limit for a plan.
func (p Plan) GetComplexityLimit() int {
	switch p {
	case PlanFree:
		return 50
	case PlanStandard:
		return 200
	case PlanPremium:
		return 1000
	case PlanEnterprise:
		return 5000
	default:
		return 10
	}
}

// RateLimit represents rate limiting configuration.
type RateLimit struct {
	RequestsPerMinute int
	RequestsPerHour   int
	RequestsPerDay    int
}

// Window represents a time window for rate limiting.
type Window struct {
	Duration time.Duration
	Limit    int
}

// GetWindows returns the rate limit windows.
func (r RateLimit) GetWindows() []Window {
	return []Window{
		{Duration: time.Minute, Limit: r.RequestsPerMinute},
		{Duration: time.Hour, Limit: r.RequestsPerHour},
		{Duration: 24 * time.Hour, Limit: r.RequestsPerDay},
	}
}

// GetScopeForResource returns the required scope for a resource.
func GetScopeForResource(resource string) Scope {
	switch resource {
	case "events":
		return ScopeReadEvents
	case "teams":
		return ScopeReadTeams
	case "matches":
		return ScopeReadMatches
	case "players":
		return ScopeReadPlayers
	case "predictions":
		return ScopeReadPredictions
	default:
		return ScopeReadEvents // Default scope
	}
}

// GetWriteScopeForResource returns the required write scope for a resource.
func GetWriteScopeForResource(resource string) Scope {
	switch resource {
	case "events":
		return ScopeWriteEvents
	case "teams":
		return ScopeWriteTeams
	case "matches":
		return ScopeWriteMatches
	case "players":
		return ScopeWritePlayers
	case "predictions":
		return ScopeWritePredictions
	default:
		return ScopeWriteEvents // Default scope
	}
}
