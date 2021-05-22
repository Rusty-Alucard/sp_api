package persistence

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
	"github.com/Rusty-Alucard/sp_api/infrastructure/provider"
)

type teamPersistence struct{}

func NewTeamPersistence() repository.TeamRepository {
	return &teamPersistence{}
}

func (p teamPersistence) FindAll(ctx context.Context) ([]*model.Team, error) {

	db, err := provider.Connect()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM teams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []*model.Team
	for rows.Next() {
		var e model.Team
		rows.Scan(&e.ID, &e.FifaTrigramme, &e.Name, &e.NameJp, &e.Confederation)
		ret = append(ret, &e)
	}

	return ret, nil
}

func (p teamPersistence) Find(ctx context.Context, fifa string) (*model.Team, error) {

	db, err := provider.Connect()
	if err != nil {
		return nil, err
	}

	ret := model.Team{}
	if err := db.QueryRow("SELECT * FROM teams WHERE fifa_trigramma = $1", fifa).Scan(&ret.ID, &ret.FifaTrigramme, &ret.Name, &ret.NameJp, &ret.Confederation); err != nil {
		return nil, err
	}

	return &ret, nil
}
