package framework

import (
	"database/sql"

	"github.com/chrisnappin/flightchecker/pkg/domain"
)

// FlightRepository handles CRUD operations on flight data.
type FlightRepository struct {
	logger domain.Logger
	db     *sql.DB
}

// NewFlightRepository creates a new instance.
func NewFlightRepository(logger domain.Logger, db *sql.DB) *FlightRepository {
	return &FlightRepository{logger, db}
}

// InitialiseSchema populates a blank repository with the schema - ie empty tables.
func (repo *FlightRepository) InitialiseSchema() error {
	_, err := withTransaction(repo.db, func(tx *sql.Tx) (interface{}, error) {
		statement, err := repo.db.Prepare(
			"CREATE TABLE IF NOT EXISTS airport (code TEST PRIMARY KEY, name TEXT, region TEXT, country TEXT)")
		if err != nil {
			return nil, err
		}
		_, err = statement.Exec()
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

// CreateAirports inserts all specified airports into the repository.
func (repo *FlightRepository) CreateAirports(airports []domain.Airport) error {
	_, err := withTransaction(repo.db, func(tx *sql.Tx) (interface{}, error) {
		// inserts all values in the array, in one transaction
		statement, err := tx.Prepare("INSERT INTO airport (code, name, region, country) VALUES (?, ?, ?, ?)")
		if err != nil {
			return nil, err
		}

		for _, airport := range airports {
			_, err = statement.Exec(airport.IataCode, airport.Name, airport.Region, airport.Country)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	return err
}

// ReadAllAirports reads all airports from the repository.
func (repo *FlightRepository) ReadAllAirports() ([]domain.Airport, error) {
	airports, err := withTransaction(repo.db, func(tx *sql.Tx) (interface{}, error) {
		rows, err := tx.Query("SELECT code, name, region, country FROM airport")
		if err != nil {
			return nil, err
		}
		airports := make([]domain.Airport, 0)
		var code string
		var name string
		var region string
		var country string
		for rows.Next() {
			err = rows.Scan(&code, &name, &region, &country)
			if err != nil {
				return nil, err
			}
			airports = append(airports, domain.Airport{Name: name, IataCode: code, Country: country, Region: region})
			repo.logger.Infof("%s %s %s %s\n", code, name, region, country)
		}
		return airports, nil
	})
	if err != nil {
		return nil, err
	}
	return airports.([]domain.Airport), nil
}

// withTransaction starts a transaction, passes it to a callback, then commits or rolls it back based on if an error is
// returned from the callback function.
func withTransaction(db *sql.DB, callback func(transaction *sql.Tx) (interface{}, error)) (interface{}, error) {
	transaction, err := db.Begin()
	if err != nil {
		return nil, err
	}

	result, err := callback(transaction)
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}
	return result, nil
}
