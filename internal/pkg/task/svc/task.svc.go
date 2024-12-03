package svc

import (
	"TaskChord/internal/pkg/task/ent"
	gossiper "github.com/pieceowater-dev/lotof.lib.gossiper/v2"
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
