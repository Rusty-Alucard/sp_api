package usecase

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
)

type EventUseCase interface {
	FindAll(context.Context) ([]*model.Event, error)
}

type eventUseCase struct {
	eventRepository repository.EventRepository
}

func NewEventUseCase(r repository.EventRepository) EventUseCase {
	return &eventUseCase{
		eventRepository: r,
	}
}

func (u eventUseCase) FindAll(ctx context.Context) (events []*model.Event, err error) {
	events, err = u.eventRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return events, nil
}
