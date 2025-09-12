// Package persistence provides data access implementations for domain entities.
package persistence

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
	"github.com/Rusty-Alucard/sp_api/infrastructure/provider"
)

type teamPersistence struct{}

// NewTeamPersistence creates a new instance of TeamRepository.
func NewTeamPersistence() repository.TeamRepository {
	return &teamPersistence{}
}

func (p teamPersistence) FindAll(ctx context.Context) ([]*model.Team, error) {
	db, err := provider.Connect()
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			_ = closeErr // Ignore close error to avoid masking the original error
		}
	}()

	rows, err := db.QueryContext(ctx, "SELECT * FROM teams")
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			_ = closeErr // Ignore close error to avoid masking the original error
		}
	}()

	var teams []*model.Team
	for rows.Next() {
		team := &model.Team{}
		if err := rows.Scan(&team.ID, &team.FifaTrigramme, &team.Name, &team.NameJp, &team.Confederation); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func (p teamPersistence) Find(ctx context.Context, fifa string) (*model.Team, error) {
	db, err := provider.Connect()
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			_ = closeErr // Ignore close error to avoid masking the original error
		}
	}()

	query := "SELECT * FROM teams WHERE fifa_trigramme = $1"
	team := &model.Team{}
	err = db.QueryRowContext(ctx, query, fifa).Scan(
		&team.ID, &team.FifaTrigramme, &team.Name, &team.NameJp, &team.Confederation)
	if err != nil {
		return nil, err
	}

	return team, nil
}
