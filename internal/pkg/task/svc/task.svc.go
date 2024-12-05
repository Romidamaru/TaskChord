package svc

import (
	"fmt"
	gossiper "github.com/pieceowater-dev/lotof.lib.gossiper/v2"
	"gorm.io/gorm"
	"taskchord/internal/pkg/task/ent"
)

type TaskService struct {
	db gossiper.Database
}

// NewTaskService initializes a new task service
func NewTaskService(db gossiper.Database) *TaskService {
	return &TaskService{db: db}
}

// CreateTask adds a task to the database
func (s *TaskService) CreateTask(guildID, userID, title, description, priority string) (int, error) {
	var maxTaskIdInGuild int

	// Start a transaction to ensure atomicity
	err := s.db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Find the highest TaskIdInGuild for the guild
		err := tx.Model(&ent.Task{}).
			Where("guild_id = ?", guildID).
			Select("COALESCE(MAX(task_id_in_guild), 0)").
			Scan(&maxTaskIdInGuild).Error
		if err != nil {
			return err
		}

		// Increment TaskIdInGuild for the new task
		newTaskIdInGuild := maxTaskIdInGuild + 1

		// Create the new task with the incremented TaskIdInGuild
		task := ent.Task{
			TaskIdInGuild: newTaskIdInGuild,
			GuildID:       guildID,
			UserID:        userID,
			Title:         title,
			Description:   description,
			Priority:      ent.Priority(priority),
		}

		// Save the new task
		if err := tx.Create(&task).Error; err != nil {
			return err
		}

		// Ensure the task is created and its ID is available before we return
		return nil
	})

	if err != nil {
		return 0, err
	}

	// Now retrieve the TaskIdInGuild from the newly created task
	var createdTask ent.Task
	err = s.db.GetDB().Where("guild_id = ? AND user_id = ? AND title = ?", guildID, userID, title).First(&createdTask).Error
	if err != nil {
		return 0, err
	}

	// Return the ID of the newly created task
	return createdTask.TaskIdInGuild, nil
}

// GetTasksByUserID retrieves tasks for a specific user from the database
func (s *TaskService) GetTasksByUserID(guildID string, userID string, id string) ([]ent.Task, error) {
	var tasks []ent.Task
	var err error

	if id != "" { // Check if id is provided (non-empty string)
		// If id is provided, fetch tasks with the specific task ID in the guild
		err = s.db.GetDB().
			Where("user_id = ? AND guild_id = ? AND task_id_in_guild = ?", userID, guildID, id).Find(&tasks).Error
	} else {
		// If no id is provided, fetch all tasks for the user in the guild
		err = s.db.GetDB().
			Where("user_id = ? AND guild_id = ?", userID, guildID).Order("task_id_in_guild DESC").Find(&tasks).Error
	}

	return tasks, err
}

func (s *TaskService) DeleteTask(guildID string, userID string, id string) (string, error) {
	// Find the task by guildID, userID, and task ID (taskIdInGuild)
	var task ent.Task
	err := s.db.GetDB().Where("guild_id = ? AND user_id = ? AND task_id_in_guild = ?", guildID, userID, id).First(&task).Error
	if err != nil {
		// Handle case where task is not found
		return "", fmt.Errorf("task not found: %v", err)
	}

	// Delete the task from the database
	err = s.db.GetDB().Delete(&task).Error
	if err != nil {
		return "", fmt.Errorf("error deleting task: %v", err)
	}

	// Return the ID of the deleted task
	return id, nil
}
