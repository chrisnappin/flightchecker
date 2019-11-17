package framework

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3" // use sqlite3 driver
)

// OpenDatabase deletes the database (if recreate is true), then returns a connection to a new SQLite database, stored
// in the specified database file
func OpenDatabase(filename string, recreate bool) (*sql.DB, error) {
	_, err := os.Stat(filename)
	if err == nil {
		// file exists, so remove it
		err = os.Remove(filename)
		if err != nil {
			return nil, err
		}
	}

	database, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	return database, nil
}
