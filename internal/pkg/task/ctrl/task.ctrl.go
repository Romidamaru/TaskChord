package ctrl

import (
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

// GetTasksByUserID retrieves tasks for a specific user
func (c *TaskController) GetTasksByUserID(guildID string, userID string) ([]ent.Task, error) {
	return c.taskService.GetTasksByUserID(guildID, userID)
}
