package event

import (
	"fmt"
	"strings"

	"github.com/Rusty-Alucard/sp_api/internal/shared/domain"
)

// EventName represents the name of an event
type EventName struct {
	domain.StringValueObject
}

// NewEventName creates a new event name value object
func NewEventName(name string) (EventName, error) {
	if len(strings.TrimSpace(name)) < 3 {
		return EventName{}, domain.NewDomainError(
			"INVALID_EVENT_NAME",
			"Event name must be at least 3 characters long",
			map[string]interface{}{"name": name},
		)
	}

	if len(name) > 100 {
		return EventName{}, domain.NewDomainError(
			"INVALID_EVENT_NAME",
			"Event name must be at most 100 characters long",
			map[string]interface{}{"name": name},
		)
	}

	return EventName{
		StringValueObject: domain.NewStringValueObject(strings.TrimSpace(name)),
	}, nil
}

// EventNameJP represents the Japanese name of an event
type EventNameJP struct {
	domain.StringValueObject
}

// NewEventNameJP creates a new Japanese event name value object
func NewEventNameJP(name string) (EventNameJP, error) {
	if len(strings.TrimSpace(name)) < 3 {
		return EventNameJP{}, domain.NewDomainError(
			"INVALID_EVENT_NAME_JP",
			"Japanese event name must be at least 3 characters long",
			map[string]interface{}{"name": name},
		)
	}

	if len(name) > 100 {
		return EventNameJP{}, domain.NewDomainError(
			"INVALID_EVENT_NAME_JP",
			"Japanese event name must be at most 100 characters long",
			map[string]interface{}{"name": name},
		)
	}

	return EventNameJP{
		StringValueObject: domain.NewStringValueObject(strings.TrimSpace(name)),
	}, nil
}

// EventStatus represents the status of an event
type EventStatus int

const (
	EventStatusScheduled EventStatus = iota
	EventStatusActive
	EventStatusCompleted
	EventStatusCancelled
)

// String returns the string representation of the event status
func (s EventStatus) String() string {
	switch s {
	case EventStatusScheduled:
		return "SCHEDULED"
	case EventStatusActive:
		return "ACTIVE"
	case EventStatusCompleted:
		return "COMPLETED"
	case EventStatusCancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}

// IsValid checks if the event status is valid
func (s EventStatus) IsValid() bool {
	return s >= EventStatusScheduled && s <= EventStatusCancelled
}

// CanTransitionTo checks if the event can transition to the given status
func (s EventStatus) CanTransitionTo(newStatus EventStatus) bool {
	switch s {
	case EventStatusScheduled:
		return newStatus == EventStatusActive || newStatus == EventStatusCancelled
	case EventStatusActive:
		return newStatus == EventStatusCompleted || newStatus == EventStatusCancelled
	case EventStatusCompleted, EventStatusCancelled:
		return false
	default:
		return false
	}
}

// TeamID represents a team identifier
type TeamID struct {
	domain.StringValueObject
}

// NewTeamID creates a new team ID value object
func NewTeamID(id string) (TeamID, error) {
	if len(strings.TrimSpace(id)) == 0 {
		return TeamID{}, domain.NewDomainError(
			"INVALID_TEAM_ID",
			"Team ID cannot be empty",
			map[string]interface{}{"id": id},
		)
	}

	return TeamID{
		StringValueObject: domain.NewStringValueObject(strings.TrimSpace(id)),
	}, nil
}

// MatchID represents a match identifier
type MatchID struct {
	domain.StringValueObject
}

// NewMatchID creates a new match ID value object
func NewMatchID(id string) (MatchID, error) {
	if len(strings.TrimSpace(id)) == 0 {
		return MatchID{}, domain.NewDomainError(
			"INVALID_MATCH_ID",
			"Match ID cannot be empty",
			map[string]interface{}{"id": id},
		)
	}

	return MatchID{
		StringValueObject: domain.NewStringValueObject(strings.TrimSpace(id)),
	}, nil
}

// Score represents a match score
type Score struct {
	home int
	away int
}

// NewScore creates a new score value object
func NewScore(home, away int) (Score, error) {
	if home < 0 || away < 0 {
		return Score{}, domain.NewDomainError(
			"INVALID_SCORE",
			"Score values cannot be negative",
			map[string]interface{}{"home": home, "away": away},
		)
	}

	return Score{home: home, away: away}, nil
}

// Home returns the home team score
func (s Score) Home() int {
	return s.home
}

// Away returns the away team score
func (s Score) Away() int {
	return s.away
}

// IsDraw checks if the match is a draw
func (s Score) IsDraw() bool {
	return s.home == s.away
}

// Winner returns the winning team (nil if draw)
func (s Score) Winner() *TeamID {
	if s.home > s.away {
		// Return home team ID (would need to be passed in context)
		return nil
	} else if s.away > s.home {
		// Return away team ID (would need to be passed in context)
		return nil
	}
	return nil
}

// String returns the string representation of the score
func (s Score) String() string {
	return fmt.Sprintf("%d-%d", s.home, s.away)
}

// Equals compares two scores for equality
func (s Score) Equals(other domain.ValueObject) bool {
	if other == nil {
		return false
	}
	if otherScore, ok := other.(Score); ok {
		return s.home == otherScore.home && s.away == otherScore.away
	}
	return false
}

// MatchTime represents match time
type MatchTime struct {
	minutes      int
	extraMinutes int
}

// NewMatchTime creates a new match time value object
func NewMatchTime(minutes, extraMinutes int) (MatchTime, error) {
	if minutes < 0 || minutes > 120 {
		return MatchTime{}, domain.NewDomainError(
			"INVALID_MATCH_TIME",
			"Match minutes must be between 0 and 120",
			map[string]interface{}{"minutes": minutes},
		)
	}

	if extraMinutes < 0 || extraMinutes > 30 {
		return MatchTime{}, domain.NewDomainError(
			"INVALID_MATCH_TIME",
			"Extra minutes must be between 0 and 30",
			map[string]interface{}{"extraMinutes": extraMinutes},
		)
	}

	return MatchTime{
		minutes:      minutes,
		extraMinutes: extraMinutes,
	}, nil
}

// Minutes returns the regular minutes
func (t MatchTime) Minutes() int {
	return t.minutes
}

// ExtraMinutes returns the extra minutes
func (t MatchTime) ExtraMinutes() int {
	return t.extraMinutes
}

// TotalMinutes returns the total minutes
func (t MatchTime) TotalMinutes() int {
	return t.minutes + t.extraMinutes
}

// String returns the string representation
func (t MatchTime) String() string {
	if t.extraMinutes > 0 {
		return fmt.Sprintf("%d+%d'", t.minutes, t.extraMinutes)
	}
	return fmt.Sprintf("%d'", t.minutes)
}

// Equals compares two match times for equality
func (t MatchTime) Equals(other domain.ValueObject) bool {
	if other == nil {
		return false
	}
	if otherTime, ok := other.(MatchTime); ok {
		return t.minutes == otherTime.minutes && t.extraMinutes == otherTime.extraMinutes
	}
	return false
}

