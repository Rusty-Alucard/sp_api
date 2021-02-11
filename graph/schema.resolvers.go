package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	//	"fmt"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/graph/generated"
)

func (r *queryResolver) Events(ctx context.Context) ([]*model.Event, error) {
	return r.EventUseCase.FindAll(ctx)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
