package tasks

import (
	"fmt"
	"math"
	"slices"
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

func (s *Service) repeatDay(now time.Time, parsedDate time.Time, repeat []string) (time.Time, error) {
	if len(repeat) < 2 {
		return time.Time{}, fmt.Errorf("missing day count in 'repeat' parameter")
	}

	days, err := strconv.Atoi(repeat[1])
	if err != nil || days <= 0 || days > 400 {
		return time.Time{}, fmt.Errorf("incorrect day value in 'repeat' parameter")
	}

	for {
		parsedDate = parsedDate.AddDate(0, 0, days)
		if parsedDate.After(now) {
			return parsedDate, nil
		}
	}
}

func (s *Service) repeatYear(now time.Time, parsedDate time.Time) (time.Time, error) {
	for {
		parsedDate = parsedDate.AddDate(1, 0, 0)
		if parsedDate.After(now) {
			return parsedDate, nil
		}
	}
}

func getNumDay(day time.Weekday) int {
	if day == 0 {
		day = 7
	}

	return int(day)
}

func (s *Service) repeatWeek(now time.Time, repeat []string) (time.Time, error) {
	if len(repeat) < 2 {
		return time.Time{}, fmt.Errorf("missing day count in 'repeat' parameter")
	}

	nowDay := getNumDay(now.Weekday())
	dayNums := strings.Split(repeat[1], ",")
	minDay := math.MaxInt32

	for _, dayStr := range dayNums {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return time.Time{}, err
		}

		if day < 1 || day > 7 {
			return time.Time{}, fmt.Errorf("incorrect day in 'repeat' parameter")
		}

		diff := (day - nowDay + 7) % 7
		if diff == 0 {
			diff = 7
		}

		if diff < minDay {
			minDay = diff
		}
	}

	return now.AddDate(0, 0, minDay), nil
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func (s *Service) getDaysFromRepeat(repeat []string) ([]int, error) {
	var days []int
	for s := range strings.SplitSeq(repeat[1], ",") {
		day, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		if day < -2 || day == 0 || day > 31 {
			return nil, fmt.Errorf("invalid day [-2,-1, 1...31]")
		}
		days = append(days, day)
	}

	return days, nil
}

func (s *Service) getMonthsFromRepeat(repeat []string) ([]int, error) {
	var months []int
	if len(repeat) == 3 {
		for s := range strings.SplitSeq(repeat[2], ",") {
			month, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			if month < 1 || month > 12 {
				return nil, fmt.Errorf("invalid month [1:12]")
			}
			months = append(months, month)
		}
	} else {
		for i := 1; i <= 12; i++ {
			months = append(months, i)
		}
	}

	return months, nil
}

func calcDateForMonth(now time.Time, parsedDate time.Time, days []int, months []int) (time.Time, error) {
	var current time.Time
	if parsedDate.After(now) {
		current = parsedDate.AddDate(0, 0, 1)
	} else {
		current = now.AddDate(0, 0, 1)
	}

	nearest := parsedDate

	for range 2 {
		for month := 1; month <= 12; month++ {
			if !slices.Contains(months, month) {
				continue
			}

			for _, rule := range days {
				var day int
				switch rule {
				case -1:
					day = daysInMonth(current.Year(), time.Month(month))
				case -2:
					day = daysInMonth(current.Year(), time.Month(month)) - 1
				default:
					day = rule
				}

				if day < 1 || day > daysInMonth(current.Year(), time.Month(month)) {
					continue
				}

				date := time.Date(current.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)

				if date.After(now) && date.After(parsedDate) {
					if nearest.Equal(parsedDate) || nearest.After(date) {
						nearest = date
					}
				}
			}
		}

		current = time.Date(current.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	return nearest, nil
}

func (s *Service) repeatMonth(now time.Time, parsedDate time.Time, repeat []string) (time.Time, error) {
	if len(repeat) < 2 {
		return time.Time{}, fmt.Errorf("missing day count in 'repeat' parameter")
	}

	days, err := s.getDaysFromRepeat(repeat)
	if err != nil {
		return time.Time{}, err
	}

	months, err := s.getMonthsFromRepeat(repeat)
	if err != nil {
		return time.Time{}, err
	}

	return calcDateForMonth(now, parsedDate, days, months)
}

// NextDate calculates the next date for the task according to the specified rule
func (s *Service) NextDate(nowStr string, startDate string, repeat string) (string, error) {
	var err error
	var now time.Time

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFmt, nowStr)
		if err != nil {
			return "", fmt.Errorf("invalid 'now' parameter")
		}
	}

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
		parsedDate, err = s.repeatDay(now, parsedDate, parts)
		if err != nil {
			return "", err
		}

	case "y":
		parsedDate, err = s.repeatYear(now, parsedDate)
		if err != nil {
			return "", err
		}

	case "w":
		parsedDate, err = s.repeatWeek(now, parts)
		if err != nil {
			return "", err
		}

	case "m":
		parsedDate, err = s.repeatMonth(now, parsedDate, parts)
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
			nextDate, err := s.NextDate(now.Format(dateFmt), task.Date, task.Repeat)
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

	id, err := s.repo.Add(task)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetTasks returns records from the database
// The number of records is limited by the 'limit' parameter
func (s *Service) GetTasks(search string, limit int) ([]*models.Task, error) {
	if limit < 1 || limit > 50 {
		return nil, fmt.Errorf("invalid 'limit' param")
	}

	tasks, err := s.repo.GetBySearch(search, limit)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTask returns record from the database by id
func (s *Service) GetTask(id string) (models.Task, error) {
	if id == "" {
		return models.Task{}, fmt.Errorf("empty 'id' param")
	}

	tasks, err := s.repo.Get(id)
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

	err = s.repo.UpdateRecord(task)
	if err != nil {
		return err
	}

	return nil
}

// DoneTask confirms the completion of the task
func (s *Service) DoneTask(id string) error {
	if id == "" {
		return fmt.Errorf("empty 'id' param")
	}

	task, err := s.repo.Get(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		err = s.repo.Delete(id)
		if err != nil {
			return err
		}
		return nil
	}

	now := time.Now().Format(dateFmt)
	newDate, err := s.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	task.Date = newDate

	err = s.repo.UpdateDate(&task)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTask deletes record from the database by id
func (s *Service) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("empty 'id' param")
	}

	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
