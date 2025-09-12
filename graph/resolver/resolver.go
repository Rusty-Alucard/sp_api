// Package resolver provides GraphQL resolvers for the API.
package resolver

import (
	"github.com/Rusty-Alucard/sp_api/usecase"
)

// Resolver contains the use cases for GraphQL operations.
type Resolver struct {
	EventUseCase usecase.EventUseCase
	TeamUseCase  usecase.TeamUseCase
}
