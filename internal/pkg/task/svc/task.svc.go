package svc

import (
	gossiper "github.com/pieceowater-dev/lotof.lib.gossiper/v2"
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
func (s *TaskService) CreateTask(userID, title, description, priority string) error {
	task := ent.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		Priority:    ent.Priority(priority),
	}

	return s.db.GetDB().Create(&task).Error
}

// GetTasksByUserID retrieves tasks for a specific user from the database
func (s *TaskService) GetTasksByUserID(userID string) ([]ent.Task, error) {
	var tasks []ent.Task
	err := s.db.GetDB().Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}
