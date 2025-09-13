// Package auth provides authentication domain models and logic.
package auth

import (
	"time"

	"github.com/Rusty-Alucard/sp_api/internal/shared/domain"
	"github.com/google/uuid"
)

// AgentCreated represents an event when an agent is created.
type AgentCreated struct {
	eventID   string
	AgentID   string    `json:"agent_id"`
	Name      string    `json:"name"`
	Plan      string    `json:"plan"`
	CreatedAt time.Time `json:"created_at"`
}

// NewAgentCreated creates a new AgentCreated event.
func NewAgentCreated(agentID, name, plan string) AgentCreated {
	return AgentCreated{
		eventID:   uuid.New().String(),
		AgentID:   agentID,
		Name:      name,
		Plan:      plan,
		CreatedAt: time.Now(),
	}
}

// EventID returns the event ID for AgentCreated.
func (e AgentCreated) EventID() string {
	return e.eventID
}

// EventType returns the event type for AgentCreated.
func (e AgentCreated) EventType() string {
	return "agent.created"
}

// OccurredAt returns when the event occurred.
func (e AgentCreated) OccurredAt() time.Time {
	return e.CreatedAt
}

// EventData returns the event data as a map.
func (e AgentCreated) EventData() map[string]interface{} {
	return map[string]interface{}{
		"agent_id":   e.AgentID,
		"name":       e.Name,
		"plan":       e.Plan,
		"created_at": e.CreatedAt,
	}
}

// ScopeAdded represents an event when a scope is added to an agent.
type ScopeAdded struct {
	eventID string
	AgentID string    `json:"agent_id"`
	Scope   string    `json:"scope"`
	AddedAt time.Time `json:"added_at"`
}

// NewScopeAdded creates a new ScopeAdded event.
func NewScopeAdded(agentID, scope string) ScopeAdded {
	return ScopeAdded{
		eventID: uuid.New().String(),
		AgentID: agentID,
		Scope:   scope,
		AddedAt: time.Now(),
	}
}

// EventID returns the event ID for ScopeAdded.
func (e ScopeAdded) EventID() string {
	return e.eventID
}

// EventType returns the event type for ScopeAdded.
func (e ScopeAdded) EventType() string {
	return "agent.scope_added"
}

// OccurredAt returns when the event occurred.
func (e ScopeAdded) OccurredAt() time.Time {
	return e.AddedAt
}

// EventData returns the event data as a map.
func (e ScopeAdded) EventData() map[string]interface{} {
	return map[string]interface{}{
		"agent_id": e.AgentID,
		"scope":    e.Scope,
		"added_at": e.AddedAt,
	}
}

// ScopeRemoved represents an event when a scope is removed from an agent.
type ScopeRemoved struct {
	eventID   string
	AgentID   string    `json:"agent_id"`
	Scope     string    `json:"scope"`
	RemovedAt time.Time `json:"removed_at"`
}

// NewScopeRemoved creates a new ScopeRemoved event.
func NewScopeRemoved(agentID, scope string) ScopeRemoved {
	return ScopeRemoved{
		eventID:   uuid.New().String(),
		AgentID:   agentID,
		Scope:     scope,
		RemovedAt: time.Now(),
	}
}

// EventID returns the event ID for ScopeRemoved.
func (e ScopeRemoved) EventID() string {
	return e.eventID
}

// EventType returns the event type for ScopeRemoved.
func (e ScopeRemoved) EventType() string {
	return "agent.scope_removed"
}

// OccurredAt returns when the event occurred.
func (e ScopeRemoved) OccurredAt() time.Time {
	return e.RemovedAt
}

// EventData returns the event data as a map.
func (e ScopeRemoved) EventData() map[string]interface{} {
	return map[string]interface{}{
		"agent_id":   e.AgentID,
		"scope":      e.Scope,
		"removed_at": e.RemovedAt,
	}
}

// AgentDeactivated represents an event when an agent is deactivated.
type AgentDeactivated struct {
	eventID       string
	AgentID       string    `json:"agent_id"`
	DeactivatedAt time.Time `json:"deactivated_at"`
}

// NewAgentDeactivated creates a new AgentDeactivated event.
func NewAgentDeactivated(agentID string) AgentDeactivated {
	return AgentDeactivated{
		eventID:       uuid.New().String(),
		AgentID:       agentID,
		DeactivatedAt: time.Now(),
	}
}

// EventID returns the event ID for AgentDeactivated.
func (e AgentDeactivated) EventID() string {
	return e.eventID
}

// EventType returns the event type for AgentDeactivated.
func (e AgentDeactivated) EventType() string {
	return "agent.deactivated"
}

// OccurredAt returns when the event occurred.
func (e AgentDeactivated) OccurredAt() time.Time {
	return e.DeactivatedAt
}

// EventData returns the event data as a map.
func (e AgentDeactivated) EventData() map[string]interface{} {
	return map[string]interface{}{
		"agent_id":       e.AgentID,
		"deactivated_at": e.DeactivatedAt,
	}
}

// AgentActivated represents an event when an agent is activated.
type AgentActivated struct {
	eventID     string
	AgentID     string    `json:"agent_id"`
	ActivatedAt time.Time `json:"activated_at"`
}

// NewAgentActivated creates a new AgentActivated event.
func NewAgentActivated(agentID string) AgentActivated {
	return AgentActivated{
		eventID:     uuid.New().String(),
		AgentID:     agentID,
		ActivatedAt: time.Now(),
	}
}

// EventID returns the event ID for AgentActivated.
func (e AgentActivated) EventID() string {
	return e.eventID
}

// EventType returns the event type for AgentActivated.
func (e AgentActivated) EventType() string {
	return "agent.activated"
}

// OccurredAt returns when the event occurred.
func (e AgentActivated) OccurredAt() time.Time {
	return e.ActivatedAt
}

// EventData returns the event data as a map.
func (e AgentActivated) EventData() map[string]interface{} {
	return map[string]interface{}{
		"agent_id":     e.AgentID,
		"activated_at": e.ActivatedAt,
	}
}

// Ensure all events implement the domain.Event interface
var _ domain.Event = (*AgentCreated)(nil)
var _ domain.Event = (*ScopeAdded)(nil)
var _ domain.Event = (*ScopeRemoved)(nil)
var _ domain.Event = (*AgentDeactivated)(nil)
var _ domain.Event = (*AgentActivated)(nil)
