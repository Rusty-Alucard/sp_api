// Package domain provides domain-driven design base implementations.
package domain

import (
	"fmt"
	"time"
)

// AggregateID represents a unique identifier for an aggregate
type AggregateID string

// String returns the string representation of the ID
func (id AggregateID) String() string {
	return string(id)
}

// DomainEvent represents a domain event
type DomainEvent interface {
	EventID() string
	AggregateID() AggregateID
	EventType() string
	OccurredAt() time.Time
	EventData() map[string]interface{}
}

// BaseDomainEvent provides common functionality for domain events
type BaseDomainEvent struct {
	eventID     string
	aggregateID AggregateID
	eventType   string
	occurredAt  time.Time
	eventData   map[string]interface{}
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(eventID string, aggregateID AggregateID, eventType string, eventData map[string]interface{}) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:     eventID,
		aggregateID: aggregateID,
		occurredAt:  time.Now(),
		eventType:   eventType,
		eventData:   eventData,
	}
}

// EventID returns the event ID
func (e BaseDomainEvent) EventID() string {
	return e.eventID
}

// AggregateID returns the aggregate ID
func (e BaseDomainEvent) AggregateID() AggregateID {
	return e.aggregateID
}

// EventType returns the event type
func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

// OccurredAt returns when the event occurred
func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EventData returns the event data
func (e BaseDomainEvent) EventData() map[string]interface{} {
	return e.eventData
}

// Aggregate represents a domain aggregate root
type Aggregate interface {
	ID() AggregateID
	Version() int
	Events() []DomainEvent
	ClearEvents()
	AddEvent(event DomainEvent)
}

// BaseAggregate provides common functionality for aggregates
type BaseAggregate struct {
	id      AggregateID
	version int
	events  []DomainEvent
}

// NewBaseAggregate creates a new base aggregate
func NewBaseAggregate(id AggregateID) BaseAggregate {
	return BaseAggregate{
		id:      id,
		version: 0,
		events:  make([]DomainEvent, 0),
	}
}

// ID returns the aggregate ID
func (a BaseAggregate) ID() AggregateID {
	return a.id
}

// Version returns the aggregate version
func (a BaseAggregate) Version() int {
	return a.version
}

// Events returns all domain events
func (a BaseAggregate) Events() []DomainEvent {
	return a.events
}

// ClearEvents clears all domain events
func (a BaseAggregate) ClearEvents() {
	a.events = make([]DomainEvent, 0)
}

// AddEvent adds a domain event
func (a BaseAggregate) AddEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

// IncrementVersion increments the aggregate version
func (a *BaseAggregate) IncrementVersion() {
	a.version++
}

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

// Error implements the error interface
func (e DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string, details map[string]interface{}) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
