package tasks

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AlexFox86/scheduler/internal/models"
	"github.com/AlexFox86/scheduler/internal/repository"
)

const dateFmt = "20060102"

// Service provides methods for the task scheduler
type Service struct {
	repo repository.Repo
}

// New creates a new scheduler service
func New(repo repository.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) repeatDays(now time.Time, parsedDate time.Time, parts []string) (time.Time, error) {
	if len(parts) == 1 {
		return time.Time{}, fmt.Errorf("missing day count in 'repeat' parameter")
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil || days <= 0 || days > 400 {
		return time.Time{}, fmt.Errorf("invalid day count in 'repeat' parameter")
	}

	for {
		parsedDate = parsedDate.AddDate(0, 0, days)
		if parsedDate.After(now) {
			return parsedDate, nil
		}
	}
}

func (s *Service) repeatYears(now time.Time, parsedDate time.Time) (time.Time, error) {
	for {
		parsedDate = parsedDate.AddDate(1, 0, 0)
		if parsedDate.After(now) {
			return parsedDate, nil
		}
	}
}

// NextDate calculates the next date for the task according to the specified rule
func (s *Service) NextDate(now time.Time, startDate string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("'repeat' parameter not found")
	}

	parsedDate, err := time.Parse(dateFmt, startDate)
	if err != nil {
		return "", fmt.Errorf("invalid 'startDate' format")
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid 'repeat' format")
	}

	switch parts[0] {
	case "d":
		parsedDate, err = s.repeatDays(now, parsedDate, parts)
		if err != nil {
			return "", err
		}

	case "y":
		parsedDate, err = s.repeatYears(now, parsedDate)
		if err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf("unsupported 'repeat' type")
	}

	return parsedDate.Format(dateFmt), nil
}

func (s *Service) checkDate(task *models.Task) error {
	now := time.Now()
	today := now.Format(dateFmt)

	if task.Date == "" {
		task.Date = now.Format(dateFmt)
		return nil
	}

	parsedDate, err := time.Parse(dateFmt, task.Date)
	if err != nil {
		return fmt.Errorf("invalid 'Date' param")
	}

	if task.Date == today {
		return nil
	}

	if parsedDate.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFmt)
		} else {
			nextDate, err := s.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = nextDate
		}
	}

	return nil
}

// AddTask adds a task to the database
func (s *Service) AddTask(task models.Task) (string, error) {
	if task.Title == "" {
		return "", fmt.Errorf("empty 'Title' param")
	}

	err := s.checkDate(&task)
	if err != nil {
		return "", err
	}

	id, err := s.repo.AddTask(task)
	if err != nil {
		return "", err
	}
	return id, nil
}

// Tasks returns records from the database
// The number of records is limited by the 'limit' parameter
func (s *Service) Tasks(search string, limit int) ([]*models.Task, error) {
	tasks, err := s.repo.GetTasks(search, limit)
	if err != nil {
		return tasks, err
	}
	return tasks, nil
}

// Task returns record from the database by id
func (s *Service) Task(id string) (models.Task, error) {
	tasks, err := s.repo.GetTask(id)
	if err != nil {
		return tasks, err
	}
	return tasks, nil
}

// Update updates a record in the database
func (s *Service) Update(task *models.Task) error {
	if task.Title == "" {
		return fmt.Errorf("empty 'Title' param")
	}

	err := s.checkDate(task)
	if err != nil {
		return err
	}

	err = s.repo.UpdateTask(task)
	if err != nil {
		return err
	}

	return nil
}
