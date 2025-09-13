// Package auth provides authentication domain models and logic.
package auth

import (
	"errors"
	"time"

	"github.com/Rusty-Alucard/sp_api/internal/shared/domain"
)

// Agent represents an AI agent that can access the API.
type Agent struct {
	*domain.BaseAggregate
	clientID     string
	clientSecret string
	name         string
	scopes       []Scope
	plan         Plan
	isActive     bool
	lastUsedAt   *time.Time
}

// NewAgent creates a new Agent aggregate.
func NewAgent(clientID, clientSecret, name string, plan Plan) (*Agent, error) {
	if clientID == "" {
		return nil, errors.New("client ID cannot be empty")
	}
	if clientSecret == "" {
		return nil, errors.New("client secret cannot be empty")
	}
	if name == "" {
		return nil, errors.New("agent name cannot be empty")
	}

	agent := &Agent{
		BaseAggregate: &domain.BaseAggregate{
			ID:        clientID,
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		clientID:     clientID,
		clientSecret: clientSecret,
		name:         name,
		scopes:       []Scope{},
		plan:         plan,
		isActive:     true,
	}

	agent.AddDomainEvent(NewAgentCreated(clientID, name, string(plan)))

	return agent, nil
}

// ClientID returns the agent's client ID.
func (a *Agent) ClientID() string {
	return a.clientID
}

// ClientSecret returns the agent's client secret.
func (a *Agent) ClientSecret() string {
	return a.clientSecret
}

// Name returns the agent's name.
func (a *Agent) Name() string {
	return a.name
}

// Scopes returns the agent's scopes.
func (a *Agent) Scopes() []Scope {
	return a.scopes
}

// Plan returns the agent's plan.
func (a *Agent) Plan() Plan {
	return a.plan
}

// IsActive returns whether the agent is active.
func (a *Agent) IsActive() bool {
	return a.isActive
}

// LastUsedAt returns when the agent was last used.
func (a *Agent) LastUsedAt() *time.Time {
	return a.lastUsedAt
}

// AddScope adds a scope to the agent.
func (a *Agent) AddScope(scope Scope) error {
	if !a.isActive {
		return errors.New("cannot add scope to inactive agent")
	}

	for _, existingScope := range a.scopes {
		if existingScope == scope {
			return nil // Scope already exists
		}
	}

	a.scopes = append(a.scopes, scope)
	a.UpdatedAt = time.Now()

	a.AddDomainEvent(NewScopeAdded(a.clientID, string(scope)))

	return nil
}

// RemoveScope removes a scope from the agent.
func (a *Agent) RemoveScope(scope Scope) error {
	if !a.isActive {
		return errors.New("cannot remove scope from inactive agent")
	}

	for i, existingScope := range a.scopes {
		if existingScope == scope {
			a.scopes = append(a.scopes[:i], a.scopes[i+1:]...)
			a.UpdatedAt = time.Now()

			a.AddDomainEvent(NewScopeRemoved(a.clientID, string(scope)))

			return nil
		}
	}

	return errors.New("scope not found")
}

// UpdateLastUsed updates the last used timestamp.
func (a *Agent) UpdateLastUsed() {
	now := time.Now()
	a.lastUsedAt = &now
	a.UpdatedAt = now
}

// Deactivate deactivates the agent.
func (a *Agent) Deactivate() error {
	if !a.isActive {
		return errors.New("agent is already inactive")
	}

	a.isActive = false
	a.UpdatedAt = time.Now()

	a.AddDomainEvent(NewAgentDeactivated(a.clientID))

	return nil
}

// Activate activates the agent.
func (a *Agent) Activate() error {
	if a.isActive {
		return errors.New("agent is already active")
	}

	a.isActive = true
	a.UpdatedAt = time.Now()

	a.AddDomainEvent(NewAgentActivated(a.clientID))

	return nil
}

// HasScope checks if the agent has a specific scope.
func (a *Agent) HasScope(scope Scope) bool {
	for _, existingScope := range a.scopes {
		if existingScope == scope {
			return true
		}
	}
	return false
}

// CanAccess checks if the agent can access a specific resource.
func (a *Agent) CanAccess(resource string) bool {
	if !a.isActive {
		return false
	}

	// Check if agent has the required scope for the resource
	requiredScope := GetScopeForResource(resource)
	return a.HasScope(requiredScope)
}
