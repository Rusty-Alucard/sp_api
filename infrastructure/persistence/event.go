package persistence

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
	"github.com/Rusty-Alucard/sp_api/infrastructure/provider"
)

type eventPersistence struct{}

func NewEventPersistence() repository.EventRepository {

	return &eventPersistence{}
}

func (p eventPersistence) FindAll(ctx context.Context) ([]*model.Event, error) {
	db, err := provider.Connect()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []*model.Event
	for rows.Next() {
		var e model.Event
		rows.Scan(&e.ID, &e.Name, &e.NameJp)
		ret = append(ret, &e)
	}

	return ret, nil
}

func (p eventPersistence) Find(ctx context.Context, id string) (*model.Event, error) {

	db, err := provider.Connect()
	if err != nil {
		return nil, err
	}

	ret := model.Event{}
	if err := db.QueryRow("SELECT * FROM events WHERE id = $1", id).Scan(&ret.ID, &ret.Name, &ret.NameJp); err != nil {
		return nil, err
	}

	return &ret, nil
}
