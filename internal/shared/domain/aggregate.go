// Package domain provides shared domain concepts and base implementations.
package domain

import "time"

// AggregateRoot represents the base interface for domain aggregates.
type AggregateRoot interface {
	GetID() string
	GetVersion() int
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDomainEvents() []Event
	ClearDomainEvents()
}

// BaseAggregate provides a base implementation for domain aggregates.
type BaseAggregate struct {
	ID           string
	Version      int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DomainEvents []Event
}

// GetID returns the aggregate ID.
func (a *BaseAggregate) GetID() string {
	return a.ID
}

// GetVersion returns the aggregate version.
func (a *BaseAggregate) GetVersion() int {
	return a.Version
}

// GetCreatedAt returns the creation timestamp.
func (a *BaseAggregate) GetCreatedAt() time.Time {
	return a.CreatedAt
}

// GetUpdatedAt returns the last update timestamp.
func (a *BaseAggregate) GetUpdatedAt() time.Time {
	return a.UpdatedAt
}

// GetDomainEvents returns all domain events.
func (a *BaseAggregate) GetDomainEvents() []Event {
	return a.DomainEvents
}

// ClearDomainEvents clears all domain events.
func (a *BaseAggregate) ClearDomainEvents() {
	a.DomainEvents = nil
}

// AddDomainEvent adds a domain event to the aggregate.
func (a *BaseAggregate) AddDomainEvent(event Event) {
	a.DomainEvents = append(a.DomainEvents, event)
}
