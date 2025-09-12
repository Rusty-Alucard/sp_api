// Package repository provides repository interfaces for domain entities.
package repository

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
)

// PlayerRepository defines the interface for player data access operations.
type PlayerRepository interface {
	FindAll(context.Context) ([]*model.Player, error)
	Find(context.Context, string) (*model.Player, error)
	FindByTeamID(context.Context, string) ([]*model.Player, error)
	FindByIDs(context.Context, []string) ([]*model.Player, error)
}
