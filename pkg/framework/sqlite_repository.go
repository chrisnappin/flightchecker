package framework

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // use sqlite3 driver
)

// SQLiteRepository handles being able to load and save quote data from a local DB
type SQLiteRepository interface {
	Initialise() error
}

// sqliteRepositoryService handles loading and saving data.
type sqliteRepositoryService struct {
	logger Logger
}

// NewSQLiteRepository creates a new instance.
func NewSQLiteRepository(logger Logger) SQLiteRepository {
	return &sqliteRepositoryService{logger}
}

// Initialise creates an empty database
func (r *sqliteRepositoryService) Initialise() error {
	database, err := sql.Open("sqlite3", "./data/flightchecker.db")
	if err != nil {
		return err
	}

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS airport (code TEST PRIMARY KEY, name TEXT, region TEXT, country TEXT)")
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	statement, err = database.Prepare("INSERT INTO airport (code, name, region, country) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec("EMA", "East Mids Airport", "England", "UK")
	if err != nil {
		return err
	}

	rows, err := database.Query("SELECT code, name, region, country FROM airport")
	if err != nil {
		return err
	}
	var code string
	var name string
	var region string
	var country string
	for rows.Next() {
		err = rows.Scan(&code, &name, &region, &country)
		if err != nil {
			return err
		}
		r.logger.Infof("%s %s %s %s\n", code, name, region, country)
	}
	return nil
}
