package event

import (
	"time"

	"github.com/Rusty-Alucard/sp_api/internal/shared/domain"
)

// EventCreated represents an event creation domain event
type EventCreated struct {
	domain.BaseDomainEvent
}

// NewEventCreated creates a new event created domain event
func NewEventCreated(eventID domain.AggregateID, name EventName, nameJP EventNameJP) EventCreated {
	return EventCreated{
		BaseDomainEvent: domain.NewBaseDomainEvent(
			generateEventID(),
			eventID,
			"EventCreated",
			map[string]interface{}{
				"name":   name.Value(),
				"nameJP": nameJP.Value(),
			},
		),
	}
}

// EventStarted represents an event start domain event
type EventStarted struct {
	domain.BaseDomainEvent
}

// NewEventStarted creates a new event started domain event
func NewEventStarted(eventID domain.AggregateID) EventStarted {
	return EventStarted{
		BaseDomainEvent: domain.NewBaseDomainEvent(
			generateEventID(),
			eventID,
			"EventStarted",
			map[string]interface{}{
				"startedAt": time.Now(),
			},
		),
	}
}

// EventCompleted represents an event completion domain event
type EventCompleted struct {
	domain.BaseDomainEvent
}

// NewEventCompleted creates a new event completed domain event
func NewEventCompleted(eventID domain.AggregateID) EventCompleted {
	return EventCompleted{
		BaseDomainEvent: domain.NewBaseDomainEvent(
			generateEventID(),
			eventID,
			"EventCompleted",
			map[string]interface{}{
				"completedAt": time.Now(),
			},
		),
	}
}

// EventCancelled represents an event cancellation domain event
type EventCancelled struct {
	domain.BaseDomainEvent
	reason string
}

// NewEventCancelled creates a new event cancelled domain event
func NewEventCancelled(eventID domain.AggregateID, reason string) EventCancelled {
	return EventCancelled{
		BaseDomainEvent: domain.NewBaseDomainEvent(
			generateEventID(),
			eventID,
			"EventCancelled",
			map[string]interface{}{
				"reason":      reason,
				"cancelledAt": time.Now(),
			},
		),
		reason: reason,
	}
}

// Reason returns the cancellation reason
func (e EventCancelled) Reason() string {
	return e.reason
}

// TeamAddedToEvent represents a team being added to an event
type TeamAddedToEvent struct {
	domain.BaseDomainEvent
	teamID TeamID
}

// NewTeamAddedToEvent creates a new team added to event domain event
func NewTeamAddedToEvent(eventID domain.AggregateID, teamID TeamID) TeamAddedToEvent {
	return TeamAddedToEvent{
		BaseDomainEvent: domain.NewBaseDomainEvent(
			generateEventID(),
			eventID,
			"TeamAddedToEvent",
			map[string]interface{}{
				"teamID":  teamID.Value(),
				"addedAt": time.Now(),
			},
		),
		teamID: teamID,
	}
}

// TeamID returns the team ID
func (e TeamAddedToEvent) TeamID() TeamID {
	return e.teamID
}

// TeamRemovedFromEvent represents a team being removed from an event
type TeamRemovedFromEvent struct {
	domain.BaseDomainEvent
	teamID TeamID
}

// NewTeamRemovedFromEvent creates a new team removed from event domain event
func NewTeamRemovedFromEvent(eventID domain.AggregateID, teamID TeamID) TeamRemovedFromEvent {
	return TeamRemovedFromEvent{
		BaseDomainEvent: domain.NewBaseDomainEvent(
			generateEventID(),
			eventID,
			"TeamRemovedFromEvent",
			map[string]interface{}{
				"teamID":    teamID.Value(),
				"removedAt": time.Now(),
			},
		),
		teamID: teamID,
	}
}

// TeamID returns the team ID
func (e TeamRemovedFromEvent) TeamID() TeamID {
	return e.teamID
}

// MatchAddedToEvent represents a match being added to an event
type MatchAddedToEvent struct {
	domain.BaseDomainEvent
	matchID MatchID
}

// NewMatchAddedToEvent creates a new match added to event domain event
func NewMatchAddedToEvent(eventID domain.AggregateID, matchID MatchID) MatchAddedToEvent {
	return MatchAddedToEvent{
		BaseDomainEvent: domain.NewBaseDomainEvent(
			generateEventID(),
			eventID,
			"MatchAddedToEvent",
			map[string]interface{}{
				"matchID": matchID.Value(),
				"addedAt": time.Now(),
			},
		),
		matchID: matchID,
	}
}

// MatchID returns the match ID
func (e MatchAddedToEvent) MatchID() MatchID {
	return e.matchID
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

