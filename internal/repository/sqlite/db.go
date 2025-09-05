package sqlite

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/AlexFox86/scheduler/internal/models"
	"github.com/jmoiron/sqlx"
)

const schema = `CREATE TABLE scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8) NOT NULL DEFAULT "",
				title VARCHAR(128) NOT NULL DEFAULT "",
				comment TEXT NOT NULL DEFAULT "",
				repeat VARCHAR(128) NOT NULL DEFAULT ""
				);
				CREATE INDEX index_date ON scheduler(date);`

// Repo the structure for working with SQLite database
type Repo struct {
	db *sqlx.DB
}

// NewSQLiteRepo creates a new object of 'SQLiteRepo' type
// and returns a pointer to it.
func NewSQLiteRepo() *Repo {
	return &Repo{}
}

// Init initializes the database
func (r *Repo) Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	r.db, err = sqlx.Connect("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err := r.db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close closes the database
func (r *Repo) Close() {
	r.db.Close()
}

// AddTask adds a task to the database
func (r *Repo) AddTask(task models.Task) (string, error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (:date, :title, :comment, :repeat)`

	res, err := r.db.NamedExec(query, task)
	if err != nil {
		return "", fmt.Errorf("failed to add task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get id: %w", err)
	}

	return strconv.FormatInt(id, 10), nil
}

// GetTasks returns records from the database
// The number of records is limited by the 'limit' parameter
func (r *Repo) GetTasks(search string, limit int) ([]*models.Task, error) {
	var query string
	params := []any{}

	if search == "" {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		params = append(params, limit)
	} else if t, err := time.Parse("02.01.2006", search); err == nil {
		time := t.Format(`20060102`)
		query = `SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
		params = append(params, time, limit)
	} else {
		search = fmt.Sprintf("%%%s%%", search)
		query = `SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
		params = append(params, search, search, limit)
	}

	return r.getTasksQuery(query, params)
}

func (r *Repo) getTasksQuery(query string, params []any) ([]*models.Task, error) {
	tasks := []*models.Task{}
	err := r.db.Select(&tasks, query, params...)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask returns record from the database by id
func (r *Repo) GetTask(id string) (models.Task, error) {
	task := models.Task{}
	query := `SELECT * FROM scheduler WHERE id = ? ORDER BY date`

	err := r.db.Get(&task, query, id)

	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

// UpdateTask updates a record in the database
func (r *Repo) UpdateTask(task *models.Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`

	res, err := r.db.NamedExec(query, task)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("task with id %s not found", task.ID)
	}

	return nil
}
