package oauth

import (
	"taskchord/internal/pkg/oauth/ctrl"
	"taskchord/internal/pkg/oauth/svc"
)

type Module struct {
	Controller *ctrl.OAuthController
}

func New() *Module {
	// Initialize the AuthService directly
	authService := svc.NewAuthService()

	// Initialize the controller with the auth service
	controller := ctrl.NewAuthController(authService)

	// Initialize and return the module
	return &Module{
		Controller: controller,
	}
}
