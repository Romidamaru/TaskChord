package ent

import (
	"gorm.io/gorm"
	"time"
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
	TaskIdInGuild int      `gorm:"not null" json:"task_id_in_guild"` // Task ID within a guild
	UserID        string   `gorm:"not null" json:"user_id"`
	ExecutorID    string   `gorm:"not null" json:"executor_id"`
	GuildID       string   `gorm:"not null;index" json:"guild_id"`                    // Indexed for grouping tasks by guild
	Title         string   `gorm:"not null" json:"title"`                             // Title of the task
	Priority      Priority `gorm:"type:varchar(20);default:'Medium'" json:"priority"` // Priority of the task (High, Medium, Low)
	Description   string   `gorm:"type:text" json:"description"`                      // Task description
}

type TaskMeta struct {
	gorm.Model
	Deadline  time.Time `json:"deadline"`
	Estimated uint      `json:"estimated"` //TODO: make notifier when deadline - estimated time
	Reminder  time.Time `json:"reminder"`  //TODO: make presets for reminder: 10 min, 15, 30, 60, day, custom
}
