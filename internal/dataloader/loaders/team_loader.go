// Package loaders provides specific DataLoader implementations.
package loaders

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
	dataloader "github.com/Rusty-Alucard/sp_api/internal/dataloader/domain"
	infrastructure "github.com/Rusty-Alucard/sp_api/internal/dataloader/infrastructure"
)

// TeamLoader provides DataLoader functionality for teams.
type TeamLoader struct {
	loader dataloader.Loader[string, *model.Team]
}

// NewTeamLoader creates a new TeamLoader.
func NewTeamLoader(teamRepo repository.TeamRepository) *TeamLoader {
	batchFunc := func(ctx context.Context, fifaCodes []string) ([]*model.Team, error) {
		// Convert string keys to the format expected by repository
		teams, err := teamRepo.FindByFifaCodes(ctx, fifaCodes)
		if err != nil {
			return nil, err
		}

		// Create a map for easy lookup
		teamMap := make(map[string]*model.Team)
		for _, team := range teams {
			teamMap[team.FifaTrigramme] = team
		}

		// Return teams in the same order as requested keys
		result := make([]*model.Team, len(fifaCodes))
		for i, fifaCode := range fifaCodes {
			if team, exists := teamMap[fifaCode]; exists {
				result[i] = team
			}
		}

		return result, nil
	}

	config := dataloader.DefaultConfig()
	config.MaxBatch = 50 // Teams are relatively small, can batch more

	loader := infrastructure.NewLoader(batchFunc, config)

	return &TeamLoader{
		loader: loader,
	}
}

// Load loads a team by FIFA code.
func (tl *TeamLoader) Load(ctx context.Context, fifaCode string) (*model.Team, error) {
	return tl.loader.Load(ctx, fifaCode)
}

// LoadMany loads multiple teams by FIFA codes.
func (tl *TeamLoader) LoadMany(ctx context.Context, fifaCodes []string) ([]*model.Team, error) {
	return tl.loader.LoadMany(ctx, fifaCodes)
}

// Clear clears a team from the cache.
func (tl *TeamLoader) Clear(fifaCode string) {
	tl.loader.Clear(fifaCode)
}

// ClearAll clears all teams from the cache.
func (tl *TeamLoader) ClearAll() {
	tl.loader.ClearAll()
}

// Prime primes the cache with a team.
func (tl *TeamLoader) Prime(fifaCode string, team *model.Team) {
	tl.loader.Prime(fifaCode, team)
}

// EventLoader provides DataLoader functionality for events.
type EventLoader struct {
	loader dataloader.Loader[string, *model.Event]
}

// createBatchFunc creates a generic batch function for loading entities by IDs.
func createBatchFunc[T any](
	repo interface{ FindByIDs(ctx context.Context, ids []string) ([]T, error) },
	getID func(T) string,
) func(ctx context.Context, ids []string) ([]T, error) {
	return func(ctx context.Context, entityIDs []string) ([]T, error) {
		entities, err := repo.FindByIDs(ctx, entityIDs)
		if err != nil {
			return nil, err
		}

		// Create a map for easy lookup
		entityMap := make(map[string]T)
		for _, entity := range entities {
			entityMap[getID(entity)] = entity
		}

		// Return entities in the same order as requested keys
		result := make([]T, len(entityIDs))
		for i, entityID := range entityIDs {
			if entity, exists := entityMap[entityID]; exists {
				result[i] = entity
			}
		}

		return result, nil
	}
}

// NewEventLoader creates a new EventLoader.
func NewEventLoader(eventRepo repository.EventRepository) *EventLoader {
	batchFunc := createBatchFunc(
		eventRepo,
		func(event *model.Event) string { return event.ID },
	)

	config := dataloader.DefaultConfig()
	config.MaxBatch = 20 // Events might be larger, batch fewer

	loader := infrastructure.NewLoader(batchFunc, config)

	return &EventLoader{
		loader: loader,
	}
}

// Load loads an event by ID.
func (el *EventLoader) Load(ctx context.Context, eventID string) (*model.Event, error) {
	return el.loader.Load(ctx, eventID)
}

// LoadMany loads multiple events by IDs.
func (el *EventLoader) LoadMany(ctx context.Context, eventIDs []string) ([]*model.Event, error) {
	return el.loader.LoadMany(ctx, eventIDs)
}

// Clear clears an event from the cache.
func (el *EventLoader) Clear(eventID string) {
	el.loader.Clear(eventID)
}

// ClearAll clears all events from the cache.
func (el *EventLoader) ClearAll() {
	el.loader.ClearAll()
}

// Prime primes the cache with an event.
func (el *EventLoader) Prime(eventID string, event *model.Event) {
	el.loader.Prime(eventID, event)
}

// PlayerLoader provides DataLoader functionality for players.
type PlayerLoader struct {
	loader dataloader.Loader[string, *model.Player]
}

// NewPlayerLoader creates a new PlayerLoader.
func NewPlayerLoader(playerRepo repository.PlayerRepository) *PlayerLoader {
	batchFunc := createBatchFunc(
		playerRepo,
		func(player *model.Player) string { return player.ID },
	)

	config := dataloader.DefaultConfig()
	config.MaxBatch = 100 // Players are small, can batch many

	loader := infrastructure.NewLoader(batchFunc, config)

	return &PlayerLoader{
		loader: loader,
	}
}

// Load loads a player by ID.
func (pl *PlayerLoader) Load(ctx context.Context, playerID string) (*model.Player, error) {
	return pl.loader.Load(ctx, playerID)
}

// LoadMany loads multiple players by IDs.
func (pl *PlayerLoader) LoadMany(ctx context.Context, playerIDs []string) ([]*model.Player, error) {
	return pl.loader.LoadMany(ctx, playerIDs)
}

// Clear clears a player from the cache.
func (pl *PlayerLoader) Clear(playerID string) {
	pl.loader.Clear(playerID)
}

// ClearAll clears all players from the cache.
func (pl *PlayerLoader) ClearAll() {
	pl.loader.ClearAll()
}

// Prime primes the cache with a player.
func (pl *PlayerLoader) Prime(playerID string, player *model.Player) {
	pl.loader.Prime(playerID, player)
}
