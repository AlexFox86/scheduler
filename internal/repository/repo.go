package repository

import (
	"github.com/AlexFox86/scheduler/internal/models"
)

// Repo interface for working with storage
type Repo interface {
	//NextDate(nd models.NextDate) (string, error)
	CreateTask(task models.Task) (int64, error)
	GetTaskById(id string) (models.Task, error)
	//GetTasks(search string) (models.ListTasks, error)
	UpdateTask(task models.Task) error
	DeleteTask(id string) error
	TaskDone(id string) error
}
