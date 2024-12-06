package svc

import (
	"errors"
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
func (s *TaskService) CreateTask(guildID, userID, title, description, priority string, executorID string) (int, error) {
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
			ExecutorID:    executorID,
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

func (s *TaskService) UpdateTask(guildID, userID, title, description, priority, executorID, id string) (int, error) {
	// Start a transaction to ensure atomicity
	err := s.db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Fetch the existing task by guild ID, user ID, and task ID
		var task ent.Task
		err := tx.Where("guild_id = ? AND user_id = ? AND task_id_in_guild = ?", guildID, userID, id).First(&task).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("task with ID %s does not exist", id)
			}
			return err
		}

		// Update only non-empty fields
		if title != "" {
			task.Title = title
		}
		if description != "" {
			task.Description = description
		}
		if priority != "" {
			task.Priority = ent.Priority(priority)
		}
		if executorID != "" { // Update executor if provided
			task.ExecutorID = executorID
		}

		// Save the changes
		if err := tx.Save(&task).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	// Retrieve the updated task to confirm the changes and get the ID
	var updatedTask ent.Task
	err = s.db.GetDB().Where("guild_id = ? AND user_id = ? AND task_id_in_guild = ?", guildID, userID, id).First(&updatedTask).Error
	if err != nil {
		return 0, err
	}

	return updatedTask.TaskIdInGuild, nil
}

// GetTasksByUserID retrieves tasks for a specific user from the database
func (s *TaskService) GetTasksByUserID(guildID string, userID string, id string) ([]ent.Task, error) {
	var tasks []ent.Task
	var err error

	if id != "" { // If a specific task ID is provided
		err = s.db.GetDB().
			Where("(user_id = ? OR executor_id = ?) AND guild_id = ? AND task_id_in_guild = ?", userID, userID, guildID, id).
			Find(&tasks).Error
	} else { // Fetch all tasks for the user (as author or executor) in the guild
		err = s.db.GetDB().
			Where("(user_id = ? OR executor_id = ?) AND guild_id = ?", userID, userID, guildID).
			Order("task_id_in_guild ASC").
			Find(&tasks).Error
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
