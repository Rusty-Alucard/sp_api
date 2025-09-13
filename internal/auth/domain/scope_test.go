package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScope_String(t *testing.T) {
	tests := []struct {
		name     string
		scope    Scope
		expected string
	}{
		{
			name:     "ReadEvents",
			scope:    ScopeReadEvents,
			expected: "events:read",
		},
		{
			name:     "WriteEvents",
			scope:    ScopeWriteEvents,
			expected: "events:write",
		},
		{
			name:     "ReadTeams",
			scope:    ScopeReadTeams,
			expected: "teams:read",
		},
		{
			name:     "Admin",
			scope:    ScopeAdmin,
			expected: "admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.scope.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlan_String(t *testing.T) {
	tests := []struct {
		name     string
		plan     Plan
		expected string
	}{
		{
			name:     "Free",
			plan:     PlanFree,
			expected: "free",
		},
		{
			name:     "Standard",
			plan:     PlanStandard,
			expected: "standard",
		},
		{
			name:     "Premium",
			plan:     PlanPremium,
			expected: "premium",
		},
		{
			name:     "Enterprise",
			plan:     PlanEnterprise,
			expected: "enterprise",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plan.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlan_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		plan     Plan
		expected bool
	}{
		{
			name:     "Free",
			plan:     PlanFree,
			expected: true,
		},
		{
			name:     "Standard",
			plan:     PlanStandard,
			expected: true,
		},
		{
			name:     "Premium",
			plan:     PlanPremium,
			expected: true,
		},
		{
			name:     "Enterprise",
			plan:     PlanEnterprise,
			expected: true,
		},
		{
			name:     "Invalid",
			plan:     Plan("invalid"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plan.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlan_GetRateLimit(t *testing.T) {
	tests := []struct {
		name     string
		plan     Plan
		expected RateLimit
	}{
		{
			name: "Free",
			plan: PlanFree,
			expected: RateLimit{
				RequestsPerMinute: 10,
				RequestsPerHour:   100,
				RequestsPerDay:    1000,
			},
		},
		{
			name: "Standard",
			plan: PlanStandard,
			expected: RateLimit{
				RequestsPerMinute: 50,
				RequestsPerHour:   1000,
				RequestsPerDay:    10000,
			},
		},
		{
			name: "Premium",
			plan: PlanPremium,
			expected: RateLimit{
				RequestsPerMinute: 200,
				RequestsPerHour:   5000,
				RequestsPerDay:    50000,
			},
		},
		{
			name: "Enterprise",
			plan: PlanEnterprise,
			expected: RateLimit{
				RequestsPerMinute: 1000,
				RequestsPerHour:   50000,
				RequestsPerDay:    500000,
			},
		},
		{
			name: "Invalid",
			plan: Plan("invalid"),
			expected: RateLimit{
				RequestsPerMinute: 1,
				RequestsPerHour:   10,
				RequestsPerDay:    100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plan.GetRateLimit()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlan_GetComplexityLimit(t *testing.T) {
	tests := []struct {
		name     string
		plan     Plan
		expected int
	}{
		{
			name:     "Free",
			plan:     PlanFree,
			expected: 50,
		},
		{
			name:     "Standard",
			plan:     PlanStandard,
			expected: 200,
		},
		{
			name:     "Premium",
			plan:     PlanPremium,
			expected: 1000,
		},
		{
			name:     "Enterprise",
			plan:     PlanEnterprise,
			expected: 5000,
		},
		{
			name:     "Invalid",
			plan:     Plan("invalid"),
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plan.GetComplexityLimit()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRateLimit_GetWindows(t *testing.T) {
	rateLimit := RateLimit{
		RequestsPerMinute: 10,
		RequestsPerHour:   100,
		RequestsPerDay:    1000,
	}

	windows := rateLimit.GetWindows()

	assert.Len(t, windows, 3)
	assert.Equal(t, 10, windows[0].Limit)
	assert.Equal(t, 100, windows[1].Limit)
	assert.Equal(t, 1000, windows[2].Limit)
}

func TestGetScopeForResource(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		expected Scope
	}{
		{
			name:     "Events",
			resource: "events",
			expected: ScopeReadEvents,
		},
		{
			name:     "Teams",
			resource: "teams",
			expected: ScopeReadTeams,
		},
		{
			name:     "Matches",
			resource: "matches",
			expected: ScopeReadMatches,
		},
		{
			name:     "Players",
			resource: "players",
			expected: ScopeReadPlayers,
		},
		{
			name:     "Predictions",
			resource: "predictions",
			expected: ScopeReadPredictions,
		},
		{
			name:     "Unknown",
			resource: "unknown",
			expected: ScopeReadEvents, // Default scope
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetScopeForResource(tt.resource)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetWriteScopeForResource(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		expected Scope
	}{
		{
			name:     "Events",
			resource: "events",
			expected: ScopeWriteEvents,
		},
		{
			name:     "Teams",
			resource: "teams",
			expected: ScopeWriteTeams,
		},
		{
			name:     "Matches",
			resource: "matches",
			expected: ScopeWriteMatches,
		},
		{
			name:     "Players",
			resource: "players",
			expected: ScopeWritePlayers,
		},
		{
			name:     "Predictions",
			resource: "predictions",
			expected: ScopeWritePredictions,
		},
		{
			name:     "Unknown",
			resource: "unknown",
			expected: ScopeWriteEvents, // Default scope
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetWriteScopeForResource(tt.resource)
			assert.Equal(t, tt.expected, result)
		})
	}
}
