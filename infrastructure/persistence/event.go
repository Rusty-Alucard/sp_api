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

func (p eventPersistence) FindByIDs(ctx context.Context, ids []string) ([]*model.Event, error) {
	db, err := provider.Connect()
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			_ = closeErr // Ignore close error to avoid masking the original error
		}
	}()

	if len(ids) == 0 {
		return []*model.Event{}, nil
	}

	// Build query with placeholders
	query := "SELECT * FROM events WHERE id IN ("
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		if i > 0 {
			query += ","
		}
		query += "$" + string(rune(i+1))
		args[i] = id
	}
	query += ")"

	rows, err := db.QueryContext(ctx, query, args...)
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
