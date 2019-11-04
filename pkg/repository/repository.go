package repository

import (
	"database/sql"

	"github.com/chrisnappin/flightchecker/pkg/framework"
	_ "github.com/mattn/go-sqlite3" // use sqlite3 driver
)

// Repository handles being able to load and save quote data from a local DB
type Repository interface {
	Initialise() error
}

// sqliteRepository handles loading and saving data.
// This is a private struct only instantiable via a factory method.
type sqliteRepository struct {
	logger framework.Logger
}

// NewRepository creates a new instance.
// This is a factory method (aka constructor) returning an interface.
func NewRepository(logger framework.Logger) Repository {
	return &sqliteRepository{logger}
}

// Initialise creates an empty database
func (r *sqliteRepository) Initialise() error {
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
