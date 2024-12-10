package svc

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"taskchord/internal/core/config"
	"taskchord/internal/pkg/user/ctrl"
)

var DiscordEndpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/api/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}

type AuthService struct {
	oauthConfig    *oauth2.Config
	userController *ctrl.UserController
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userController *ctrl.UserController) *AuthService {
	return &AuthService{
		oauthConfig: &oauth2.Config{
			ClientID:     config.Inst().DiscordClientID,
			ClientSecret: config.Inst().DiscordSecret,
			Endpoint:     DiscordEndpoint,
			RedirectURL:  config.Inst().DiscordRedirect,
			Scopes:       []string{"identify", "email"},
		},
		userController: userController,
	}
}

// GenerateAuthURL creates a Discord OAuth URL for redirection
func (s *AuthService) GenerateAuthURL() (string, error) {
	if s.oauthConfig == nil {
		return "", errors.New("OAuth configuration is missing")
	}
	return s.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline), nil
}

// HandleAuthCallback processes the Discord OAuth callback
func (s *AuthService) HandleAuthCallback(w http.ResponseWriter, r *http.Request, code string) (map[string]interface{}, error) {
	if code == "" {
		return nil, errors.New("authorization code is missing")
	}

	// Exchange the authorization code for a token
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		return nil, errors.New("failed to exchange token")
	}

	// Fetch user information
	client := s.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		log.Printf("Failed to fetch user info: %v", err)
		return nil, errors.New("failed to fetch user information")
	}
	defer resp.Body.Close()

	// Decode user information from the response
	var userInfo struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("Failed to decode user info: %v", err)
		return nil, errors.New("failed to decode user information")
	}

	// Save user information in the database
	err = s.userController.AddUser(userInfo.ID, userInfo.Username, userInfo.Email, userInfo.Avatar)
	if err != nil {
		log.Printf("Failed to save user: %v", err)
		return nil, errors.New("failed to save user information")
	}

	// Return user information as a map (for internal use if needed)
	userData := map[string]interface{}{
		"id":       userInfo.ID,
		"username": userInfo.Username,
		"email":    userInfo.Email,
		"avatar":   userInfo.Avatar,
	}

	// Redirect to the homepage (or any other page) after successful authentication
	http.Redirect(w, r, "http://localhost:3000", http.StatusSeeOther)

	return userData, nil
}
