package ctrl

import (
	"fmt"
	"log"
	"taskchord/internal/pkg/task/ent"
	"taskchord/internal/pkg/task/svc"
)

type TaskController struct {
	taskService *svc.TaskService
}

// NewTaskController creates a new task controller
func NewTaskController(taskService *svc.TaskService) *TaskController {
	return &TaskController{taskService: taskService}
}

// CreateTask delegates the task creation to the service layer
func (c *TaskController) CreateTask(guildID, userID, title, description, priority string) (int, error) {
	// Call the service layer to create the task and return the taskIdInGuild
	taskIdInGuild, err := c.taskService.CreateTask(guildID, userID, title, description, priority)
	if err != nil {
		log.Println("Controller error:", err)
		return 0, err
	}

	// Return the taskIdInGuild along with nil error
	return taskIdInGuild, nil
}

func (c *TaskController) UpdateTask(guildID, userID, title, description, priority, id string) (int, error) {
	// Validate the task ID
	if id == "" {
		log.Println("Controller error: Task ID is required")
		return 0, fmt.Errorf("task ID is required")
	}

	// Ensure at least one field is provided for updating
	if title == "" && description == "" && priority == "" {
		log.Println("Controller error: At least one field (title, description, or priority) must be provided for update")
		return 0, fmt.Errorf("at least one field (title, description, or priority) must be provided for update")
	}

	// Optional: Validate priority if provided
	if priority != "" {
		validPriorities := map[string]bool{"High": true, "Medium": true, "Low": true}
		if !validPriorities[priority] {
			log.Println("Controller error: Invalid priority value")
			return 0, fmt.Errorf("invalid priority value")
		}
	}

	// Call the service layer to update the task
	taskIdInGuild, err := c.taskService.UpdateTask(guildID, userID, title, description, priority, id)
	if err != nil {
		log.Println("Controller error:", err)
		return 0, err
	}

	// Return the updated task ID along with nil error
	return taskIdInGuild, nil
}

// GetTasksByUserID retrieves tasks for a specific user
func (c *TaskController) GetTasksByUserID(guildID string, userID string, id string) ([]ent.Task, error) {
	return c.taskService.GetTasksByUserID(guildID, userID, id)
}

func (c *TaskController) DeleteTask(guildID string, userID string, id string) (string, error) {
	taskIdInGuild, err := c.taskService.DeleteTask(guildID, userID, id)
	if err != nil {
		log.Println("Controller error:", err)
		return "", err
	}

	return taskIdInGuild, nil
}
