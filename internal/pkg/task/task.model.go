package task

import (
	gossiper "github.com/pieceowater-dev/lotof.lib.gossiper/v2"
	"log"
	"taskchord/internal/core/config"
	"taskchord/internal/pkg/task/ctrl"
	"taskchord/internal/pkg/task/ent"
	"taskchord/internal/pkg/task/svc"
)

type Module struct {
	Controller *ctrl.TaskController
}

func New() *Module {
	database, err := gossiper.NewDB(
		gossiper.PostgresDB,
		config.Inst().DSN,
		false,
		[]any{ent.Task{}},
	)
	if err != nil {
		log.Fatalf("Failed to create database instance: %v", err)
	}

	return &Module{
		Controller: ctrl.NewTaskController(
			svc.NewTaskService(database),
		),
	}
}
