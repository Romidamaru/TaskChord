package ent

import (
	"gorm.io/gorm"
)

// Priority represents the allowed values for the Priority field.
type Priority string

const (
	High   Priority = "High"
	Medium Priority = "Medium"
	Low    Priority = "Low"
)

// String method for the Priority type, to provide string representation of each priority.
func (p Priority) String() string {
	switch p {
	case High:
		return "High"
	case Medium:
		return "Medium"
	case Low:
		return "Low"
	default:
		return "Unknown"
	}
}

// Task represents a task model for GORM
type Task struct {
	gorm.Model
	UserID      string   `gorm:"not null" json:"user_id"`                           // Foreign key from User table
	Title       string   `gorm:"not null" json:"title"`                             // Title of the task
	Priority    Priority `gorm:"type:varchar(20);default:'Medium'" json:"priority"` // Priority of the task (High, Medium, Low)
	Description string   `gorm:"type:text" json:"description"`                      // Task description (long text)
}
