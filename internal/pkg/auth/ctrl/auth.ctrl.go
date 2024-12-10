package ctrl

import (
	"net/http"
	"taskchord/internal/pkg/auth/svc"
)

type AuthController struct {
	authService *svc.AuthService
}

func NewAuthController(authService *svc.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) GetAuthURL() (string, error) {
	// Use AuthService to generate the URL
	return c.authService.GenerateAuthURL()
}

func (c *AuthController) HandleAuthCallback(w http.ResponseWriter, r *http.Request, code string) (map[string]interface{}, error) {
	// Use AuthService to handle the callback and fetch user info
	return c.authService.HandleAuthCallback(w, r, code)
}
