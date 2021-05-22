package usecase

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/domain/repository"
)

type TeamUseCase interface {
	FindAll(context.Context) ([]*model.Team, error)
	Find(context.Context, string) (*model.Team, error)
}

type teamUseCase struct {
	repo repository.TeamRepository
}

func NewTeamUseCase(r repository.TeamRepository) TeamUseCase {
	return &teamUseCase{
		repo: r,
	}
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
