package user

import (
	gossiper "github.com/pieceowater-dev/lotof.lib.gossiper/v2"
	"log"
	"taskchord/internal/core/config"
	"taskchord/internal/pkg/user/ctrl"
	"taskchord/internal/pkg/user/ent"
	"taskchord/internal/pkg/user/svc"
)

type Module struct {
	Controller *ctrl.UserController
}

func New() *Module {
	database, err := gossiper.NewDB(
		gossiper.PostgresDB,
		config.Inst().DSN,
		false,
		[]any{ent.User{}},
	)
	if err != nil {
		log.Fatalf("Failed to create database instance: %v", err)
	}

	return &Module{
		Controller: ctrl.NewUserController(
			svc.NewUserService(database),
		),
	}
}
