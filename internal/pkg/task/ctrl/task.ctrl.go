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
func (c *TaskController) CreateTask(userID, title, description, priority string) error {
	err := c.taskService.CreateTask(userID, title, description, priority)
	if err != nil {
		log.Println("Controller error!", err)
	}

	return err
}

// GetTasksByUserID retrieves tasks for a specific user
func (c *TaskController) GetTasksByUserID(userID string) ([]ent.Task, error) {
	return c.taskService.GetTasksByUserID(userID)
}
