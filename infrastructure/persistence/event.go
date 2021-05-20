package persistence

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/config"
	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
)

type eventPersistence struct{}

func NewEventPersistence(cfg Config) repository.EventRepository {
	if cfg != nil {
		// TODO
	}
	return &eventPersistence{}
}

func (p eventPersistence) FindAll(ctx context.Context) ([]*model.Event, error) {
	// ロシア・ワールドカップ
	event1 := model.Event{}
	event1.ID = "1"
	event1.Name = "2018 FIFA World Cup"
	event1.NameJp = "2018 FIFA ワールドカップ"
	event1.Participants = []*model.Team{}
	event1.Matches = []*model.Match{}

	// カタール・ワールドカップ
	event2 := model.Event{}
	event2.ID = "2"
	event2.Name = "2022 FIFA World Cup"
	event2.NameJp = "2022 FIFA ワールドカップ"
	event2.Participants = []*model.Team{}
	event2.Matches = []*model.Match{}

	return []*model.Event{&event1, &event2}, nil
}

func (p eventPersistence) Find(ctx context.Context, id string) (*model.Event, error) {
	event := model.Event{}
	event.ID = "1"
	event.Name = "2018 FIFA World Cup"
	event.NameJp = "2018 FIFA ワールドカップ"
	event.Participants = []*model.Team{}
	event.Matches = []*model.Match{}

	return &event, nil
}
