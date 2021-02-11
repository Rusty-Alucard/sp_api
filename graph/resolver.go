package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/Rusty-Alucard/sp_api/usecase"
)

type Resolver struct {
	EventUseCase usecase.EventUseCase
}
