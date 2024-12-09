package svc

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"log"
	"taskchord/internal/core/config"
)

var DiscordEndpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/api/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}

type AuthService struct {
	oauthConfig *oauth2.Config
}

func NewAuthService() *AuthService {
	return &AuthService{
		oauthConfig: &oauth2.Config{
			ClientID:     config.Inst().DiscordClientID,
			ClientSecret: config.Inst().DiscordSecret,
			Endpoint:     DiscordEndpoint, // Use manually defined endpoint
			RedirectURL:  config.Inst().DiscordRedirect,
			Scopes:       []string{"identify", "email"},
		},
	}
}

// GenerateAuthURL creates a Discord OAuth URL for redirection
func (s *AuthService) GenerateAuthURL() (string, error) {
	if s.oauthConfig == nil {
		return "", errors.New("OAuth configuration is missing")
	}
	// Generate the URL where users will be redirected to Discord for login
	return s.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline), nil
}

// HandleAuthCallback processes the Discord OAuth callback
func (s *AuthService) HandleAuthCallback(code string) (map[string]interface{}, error) {
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
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("Failed to decode user info: %v", err)
		return nil, errors.New("failed to decode user information")
	}

	// Optionally, save the user info and token in the database
	// if err := s.SaveUser(token, userInfo); err != nil {
	//	log.Printf("Failed to save user: %v", err)
	// }

	return userInfo, nil
}

// SaveUser saves the user details and token (optional)
// func (s *AuthService) SaveUser(token *oauth2.Token, userInfo map[string]interface{}) error {
// 	log.Printf("Saving user with token: %v and userInfo: %v", token, userInfo)
// 	// Add logic to save user details and tokens in your database if needed.
// 	return nil
// }
