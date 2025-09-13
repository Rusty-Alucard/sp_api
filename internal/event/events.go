// Package event provides domain events for the application.
package event

import "time"

// Created represents an event creation domain event.
type Created struct {
	ID        string
	EventName string
	CreatedAt time.Time
}

// EventID returns the event ID.
func (e Created) EventID() string {
	return e.ID
}

// EventType returns the event type.
func (e Created) EventType() string {
	return "EventCreated"
}

// OccurredAt returns when the event occurred.
func (e Created) OccurredAt() time.Time {
	return e.CreatedAt
}

// EventData returns the event data.
func (e Created) EventData() map[string]interface{} {
	return map[string]interface{}{
		"event_id":   e.ID,
		"event_name": e.EventName,
		"created_at": e.CreatedAt,
	}
}

// Started represents an event start domain event.
type Started struct {
	ID        string
	EventName string
	StartedAt time.Time
}

// EventID returns the event ID.
func (e Started) EventID() string {
	return e.ID
}

// EventType returns the event type.
func (e Started) EventType() string {
	return "EventStarted"
}

// OccurredAt returns when the event occurred.
func (e Started) OccurredAt() time.Time {
	return e.StartedAt
}

// EventData returns the event data.
func (e Started) EventData() map[string]interface{} {
	return map[string]interface{}{
		"event_id":   e.ID,
		"event_name": e.EventName,
		"started_at": e.StartedAt,
	}
}

// Completed represents an event completion domain event.
type Completed struct {
	ID          string
	EventName   string
	CompletedAt time.Time
}

// EventID returns the event ID.
func (e Completed) EventID() string {
	return e.ID
}

// EventType returns the event type.
func (e Completed) EventType() string {
	return "EventCompleted"
}

// OccurredAt returns when the event occurred.
func (e Completed) OccurredAt() time.Time {
	return e.CompletedAt
}

// EventData returns the event data.
func (e Completed) EventData() map[string]interface{} {
	return map[string]interface{}{
		"event_id":     e.ID,
		"event_name":   e.EventName,
		"completed_at": e.CompletedAt,
	}
}

// Cancelled represents an event cancellation domain event.
type Cancelled struct {
	ID           string
	EventName    string
	CancelledAt  time.Time
	Cancellation string
}

// EventID returns the event ID.
func (e Cancelled) EventID() string {
	return e.ID
}

// EventType returns the event type.
func (e Cancelled) EventType() string {
	return "EventCancelled"
}

// OccurredAt returns when the event occurred.
func (e Cancelled) OccurredAt() time.Time {
	return e.CancelledAt
}

// EventData returns the event data.
func (e Cancelled) EventData() map[string]interface{} {
	return map[string]interface{}{
		"event_id":     e.ID,
		"event_name":   e.EventName,
		"cancelled_at": e.CancelledAt,
		"reason":       e.Cancellation,
	}
}


