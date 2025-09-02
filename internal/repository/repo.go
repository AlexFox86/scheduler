package repository

import (
	"context"

	"github.com/AlexFox86/scheduler/internal/models"
)

// Repo interface for working with storage
type Repo interface {
	AddTask(ctx context.Context, task models.Task) (string, error)
	GetTasks(search string, limit int) ([]*models.Task, error)
	//GetTasksSearch(search string) ([]*models.Task, error)
	//UpdateTask() error
	//DeleteTask() error
	//TaskDone() error
}
