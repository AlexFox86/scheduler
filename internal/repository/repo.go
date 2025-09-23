package repository

import (
	"github.com/AlexFox86/scheduler/internal/models"
)

// Repo interface for working with storage
type Repo interface {
	Add(task models.Task) (string, error)
	Get(id string) (models.Task, error)
	GetBySearch(search string, limit int) ([]*models.Task, error)
	UpdateRecord(task *models.Task) error
	Delete(id string) error
	UpdateDate(task *models.Task) error
}
