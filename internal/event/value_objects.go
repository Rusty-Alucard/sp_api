// Package event provides domain events for the application.
package event

import (
	"fmt"
	"strings"
)

// Name represents an event name value object.
type Name struct {
	value string
}

// NewName creates a new Name value object.
func NewName(name string) (Name, error) {
	const minLength = 3
	const maxLength = 100

	if len(strings.TrimSpace(name)) < minLength {
		return Name{}, fmt.Errorf("event name must be at least %d characters long", minLength)
	}

	if len(name) > maxLength {
		return Name{}, fmt.Errorf("event name must not exceed %d characters", maxLength)
	}

	return Name{value: strings.TrimSpace(name)}, nil
}

// Value returns the string value of the Name.
func (n Name) Value() string {
	return n.value
}

// NameJP represents a Japanese event name value object.
type NameJP struct {
	value string
}

// NewNameJP creates a new NameJP value object.
func NewNameJP(name string) (NameJP, error) {
	const minLength = 3
	const maxLength = 100

	if len(strings.TrimSpace(name)) < minLength {
		return NameJP{}, fmt.Errorf("event name (JP) must be at least %d characters long", minLength)
	}

	if len(name) > maxLength {
		return NameJP{}, fmt.Errorf("event name (JP) must not exceed %d characters", maxLength)
	}

	return NameJP{value: strings.TrimSpace(name)}, nil
}

// Value returns the string value of the NameJP.
func (n NameJP) Value() string {
	return n.value
}

// Status represents the status of an event.
type Status int

const (
	// StatusScheduled indicates the event is scheduled.
	StatusScheduled Status = iota
	// StatusInProgress indicates the event is in progress.
	StatusInProgress
	// StatusCompleted indicates the event is completed.
	StatusCompleted
	// StatusCancelled indicates the event is cancelled.
	StatusCancelled
)

// String returns the string representation of the Status.
func (s Status) String() string {
	switch s {
	case StatusScheduled:
		return "scheduled"
	case StatusInProgress:
		return "in_progress"
	case StatusCompleted:
		return "completed"
	case StatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}
