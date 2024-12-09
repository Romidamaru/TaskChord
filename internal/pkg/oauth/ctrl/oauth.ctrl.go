package ctrl

import (
	"taskchord/internal/pkg/oauth/svc"
)

type OAuthController struct {
	authService *svc.AuthService
}

func NewAuthController(authService *svc.AuthService) *OAuthController {
	return &OAuthController{
		authService: authService,
	}
}

func (c *OAuthController) GetAuthURL() (string, error) {
	// Use AuthService to generate the URL
	return c.authService.GenerateAuthURL()
}

func (c *OAuthController) HandleAuthCallback(code string) (map[string]interface{}, error) {
	// Use AuthService to handle the callback and fetch user info
	return c.authService.HandleAuthCallback(code)
}
