package ent

import (
	"gorm.io/gorm"
	"taskchord/internal/pkg/task/ent"
)

// User represents a user model for GORM
type User struct {
	gorm.Model
	UserID    string     `gorm:"unique;not null" json:"user_id"`      // Unique identifier from Discord OAuth
	Username  string     `gorm:"not null" json:"username"`            // Discord username
	Email     string     `gorm:"not null" json:"email"`               // User's email (if available)
	AvatarURL string     `json:"avatar_url"`                          // URL to user's avatar (optional)
	Tasks     []ent.Task `gorm:"foreignKey:UserID;references:UserID"` // Tasks linked to the user
}
