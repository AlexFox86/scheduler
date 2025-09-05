package repository

import (
	"github.com/AlexFox86/scheduler/internal/models"
)

// Repo interface for working with storage
type Repo interface {
	AddTask(task models.Task) (string, error)
	GetTasks(search string, limit int) ([]*models.Task, error)
	GetTask(id string) (models.Task, error)
	UpdateTask(task *models.Task) error
	//DeleteTask() error
	//TaskDone() error
}
