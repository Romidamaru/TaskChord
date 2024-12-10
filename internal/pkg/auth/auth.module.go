package auth

import (
	gossiper "github.com/pieceowater-dev/lotof.lib.gossiper/v2"
	"log"
	"taskchord/internal/core/config"
	authCtrl "taskchord/internal/pkg/auth/ctrl"
	"taskchord/internal/pkg/auth/svc"
	"taskchord/internal/pkg/user/ctrl"
	userSvc "taskchord/internal/pkg/user/svc"
)

type Module struct {
	AuthController *authCtrl.AuthController
	UserController *ctrl.UserController // Make sure this field is added to the struct
}

func New() *Module {
	database, err := gossiper.NewDB(
		gossiper.PostgresDB,
		config.Inst().DSN,
		false,
		[]any{},
	)
	if err != nil {
		log.Fatalf("Failed to create database instance: %v", err)
	}

	userService := userSvc.NewUserService(database)
	userController := ctrl.NewUserController(userService)

	authService := svc.NewAuthService(userController)
	authCtrl := authCtrl.NewAuthController(authService)

	return &Module{
		AuthController: authCtrl,
		UserController: userController,
	}
}
