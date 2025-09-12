// Package persistence provides data access implementations for domain entities.
package persistence

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
	"github.com/Rusty-Alucard/sp_api/infrastructure/provider"
)

type eventPersistence struct{}

// NewEventPersistence creates a new instance of EventRepository.
func NewEventPersistence() repository.EventRepository {
	return &eventPersistence{}
}

func (p eventPersistence) FindAll(ctx context.Context) ([]*model.Event, error) {
	db, err := provider.Connect()
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			_ = closeErr // Ignore close error to avoid masking the original error
		}
	}()

	rows, err := db.QueryContext(ctx, "SELECT * FROM events")
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			_ = closeErr // Ignore close error to avoid masking the original error
		}
	}()

	var events []*model.Event
	for rows.Next() {
		event := &model.Event{}
		if err := rows.Scan(&event.ID, &event.Name, &event.NameJp); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (p eventPersistence) Find(ctx context.Context, id string) (*model.Event, error) {
	db, err := provider.Connect()
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			_ = closeErr // Ignore close error to avoid masking the original error
		}
	}()

	event := &model.Event{}
	err = db.QueryRowContext(ctx, "SELECT * FROM events WHERE id = $1", id).Scan(&event.ID, &event.Name, &event.NameJp)
	if err != nil {
		return nil, err
	}

	return event, nil
}
