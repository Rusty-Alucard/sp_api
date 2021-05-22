package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Rusty-Alucard/sp_api/domain/model"
	"github.com/Rusty-Alucard/sp_api/graph/generated"
)

func (r *queryResolver) Events(ctx context.Context) ([]*model.Event, error) {
	return r.EventUseCase.FindAll(ctx)
}

func (r *queryResolver) Event(ctx context.Context, id string) (*model.Event, error) {
	return r.EventUseCase.Find(ctx, id)
}

func (r *queryResolver) Teams(ctx context.Context) ([]*model.Team, error) {
	return r.TeamUseCase.FindAll(ctx)
}

func (r *queryResolver) Team(ctx context.Context, fifa string) (*model.Team, error) {
	return r.TeamUseCase.Find(ctx, fifa)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
