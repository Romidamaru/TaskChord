package svc

import (
	"errors"
	gossiper "github.com/pieceowater-dev/lotof.lib.gossiper/v2"
	"gorm.io/gorm"
	"taskchord/internal/pkg/user/ent"
)

type UserService struct {
	db gossiper.Database
}

// NewUserService creates a new UserService instance
func NewUserService(db gossiper.Database) *UserService {
	return &UserService{db: db}
}

// AddUser adds a new user to the database
func (us *UserService) AddUser(userID, username, email, avatarURL string) error {
	var user ent.User

	err := us.db.GetDB().First(&user, "user_id = ?", userID).Error
	if err == nil {
		user.Username = username
		user.Email = email
		user.AvatarURL = avatarURL
		return us.db.GetDB().Save(&user).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		newUser := ent.User{
			UserID:    userID,
			Username:  username,
			Email:     email,
			AvatarURL: avatarURL,
		}
		return us.db.GetDB().Create(&newUser).Error
	}

	return err
}

// GetUserByID retrieves a user by their ID
func (us *UserService) GetUserByID(userID string) (*ent.User, error) {
	var user ent.User
	err := us.db.GetDB().First(&user, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	// Return the user found
	return &user, nil
}
