package persistence

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
)

type eventPersistence struct{}

func NewEventPersistence() repository.EventRepository {
	return &eventPersistence{}
}

func (p eventPersistence) FindAll(context.Context) ([]*model.Event, error) {
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
