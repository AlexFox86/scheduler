package sqlite

import (
	"database/sql"
	"os"
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
	db *sql.DB
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

	r.db, err = sql.Open("sqlite", dbFile)
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
