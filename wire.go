// +build wireinject

package main

import (
	"github.com/Rusty-Alucard/sp_api/infrastructure/persistence"
	"github.com/Rusty-Alucard/sp_api/usecase"
	"github.com/google/wire"
)

func InitializeEventUseCase(cfg Config) usecase.EventUseCase {
	wire.Build(
		usecase.NewEventUseCase,
		persistence.NewEventPersistence,
	)
	return nil
}

func InitializeTeamUseCase(cfg Config) usecase.TeamUseCase {
	wire.Build(
		usecase.NewTeamUseCase,
		persistence.NewTeamPersistence,
	)
	return nil
}
