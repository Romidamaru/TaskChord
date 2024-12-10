package ctrl

import (
	"log"
	"taskchord/internal/pkg/user/ent"
	"taskchord/internal/pkg/user/svc"
)

type UserController struct {
	UserService *svc.UserService
}

// NewUserController creates a new UserController instance
func NewUserController(userService *svc.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// AddUser handles adding a new user
func (uc *UserController) AddUser(userID, username, email, avatarURL string) error {
	// Add user through the service
	err := uc.UserService.AddUser(userID, username, email, avatarURL)
	if err != nil {
		log.Printf("Error adding user: %v", err)
		return err
	}
	return nil
}

// GetUserByID retrieves a user by their ID
func (uc *UserController) GetUserByID(userID string) (*ent.User, error) {
	// Call the service to get user by ID
	user, err := uc.UserService.GetUserByID(userID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return nil, err
	}
	return user, nil
}
