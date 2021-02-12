package repository

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
)

type TeamRepository interface {
	FindAll(context.Context) ([]*model.Team, error)
	Find(context.Context, string) (*model.Team, error)
}
