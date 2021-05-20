package persistence

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
)

type teamPersistence struct{}

func NewTeamPersistence(cfg Config) repository.TeamRepository {
	return &teamPersistence{}
}

func (p teamPersistence) FindAll(ctx context.Context) ([]*model.Team, error) {
	// 日本
	team1 := model.Team{}
	team1.ID = "1"
	team1.FifaTrigramme = "JPN"
	team1.Name = "Japan"
	team1.NameJp = "日本"
	team1.Confederation = "AFC"
	team1.JoinedEvents = []*model.Event{}
	team1.Coach = nil
	team1.Members = []*model.Player{}

	// ブラジル
	team2 := model.Team{}
	team2.ID = "2"
	team2.FifaTrigramme = "BRA"
	team2.Name = "Brazil"
	team2.NameJp = "ブラジル"
	team2.Confederation = "CONMEBOL"
	team2.JoinedEvents = []*model.Event{}
	team2.Coach = nil
	team2.Members = []*model.Player{}

	return []*model.Team{&team1, &team2}, nil
}

func (p teamPersistence) Find(ctx context.Context, id string) (*model.Team, error) {
	// 日本
	team1 := model.Team{}
	team1.ID = "1"
	team1.FifaTrigramme = "JPN"
	team1.Name = "Japan"
	team1.NameJp = "日本"
	team1.Confederation = "AFC"
	team1.JoinedEvents = []*model.Event{}
	team1.Coach = nil
	team1.Members = []*model.Player{}

	return &team1, nil
}
