// Package repository provides repository interfaces for domain entities.
package repository

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
)

// EventRepository defines the interface for event data access operations.
type EventRepository interface {
	FindAll(context.Context) ([]*model.Event, error)
	Find(context.Context, string) (*model.Event, error)
}
