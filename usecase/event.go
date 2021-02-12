package usecase

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
)

type EventUseCase interface {
	FindAll(context.Context) ([]*model.Event, error)
	Find(context.Context, string) (*model.Event, error)
}

type eventUseCase struct {
	repo repository.EventRepository
}

func NewEventUseCase(r repository.EventRepository) EventUseCase {
	return &eventUseCase{
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
