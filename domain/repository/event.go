package repository

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
)

type EventRepository interface {
	FindAll(context.Context) ([]*model.Event, error)
}
