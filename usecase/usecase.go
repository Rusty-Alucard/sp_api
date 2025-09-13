// Package usecase provides use case implementations for business logic.
package usecase

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
)

// EventUseCase defines the interface for event business logic operations.
type EventUseCase interface {
	FindAll(context.Context) ([]*model.Event, error)
	Find(context.Context, string) (*model.Event, error)
}

// TeamUseCase defines the interface for team business logic operations.
type TeamUseCase interface {
	FindAll(context.Context) ([]*model.Team, error)
	Find(context.Context, string) (*model.Team, error)
}

type eventUseCase struct {
	repo repository.EventRepository
}

type teamUseCase struct {
	repo repository.TeamRepository
}

// NewEventUseCase creates a new instance of EventUseCase.
func NewEventUseCase(r repository.EventRepository) EventUseCase {
	return &eventUseCase{
		repo: r,
	}
}

// NewTeamUseCase creates a new instance of TeamUseCase.
func NewTeamUseCase(r repository.TeamRepository) TeamUseCase {
	return &teamUseCase{
		repo: r,
	}
}

func (u eventUseCase) FindAll(ctx context.Context) (events []*model.Event, err error) {
	events, err = u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (u eventUseCase) Find(ctx context.Context, id string) (event *model.Event, err error) {
	event, err = u.repo.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (u teamUseCase) FindAll(ctx context.Context) (teams []*model.Team, err error) {
	teams, err = u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (u teamUseCase) Find(ctx context.Context, fifa string) (team *model.Team, err error) {
	team, err = u.repo.Find(ctx, fifa)
	if err != nil {
		return nil, err
	}
	return team, nil
}
