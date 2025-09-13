// Package domain provides shared domain concepts and base implementations.
package domain

import "time"

// Event represents a domain event interface.
type Event interface {
	EventID() string
	EventType() string
	OccurredAt() time.Time
	EventData() map[string]interface{}
}
