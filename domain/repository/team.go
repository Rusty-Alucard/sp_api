// Package repository provides repository interfaces for domain entities.
package repository

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
)

// TeamRepository defines the interface for team data access operations.
type TeamRepository interface {
	FindAll(context.Context) ([]*model.Team, error)
	Find(context.Context, string) (*model.Team, error)
}
