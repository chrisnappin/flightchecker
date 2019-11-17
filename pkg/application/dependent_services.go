package application

import (
	"github.com/chrisnappin/flightchecker/pkg/domain"
	"github.com/chrisnappin/flightchecker/pkg/framework"
)

//
// Interfaces for framework services...
//

// AirportDataLoader handles being able to load CSV data
type AirportDataLoader interface {
	LoadCountries(filename string) (map[string]string, error)
	LoadRegions(filename string) (map[string]string, error)
	LoadAirports(filename string, countries map[string]string, regions map[string]string) (map[string]domain.Airport, error)
}

// ArgumentsLoader handles being able to load arguments from a JSON file.
type ArgumentsLoader interface {
	Load(filename string) (*domain.Arguments, error)
}

// SkyScannerQuoter handles finding flight quotes from Sky Scanner.
type SkyScannerQuoter interface {
	PollForQuotes(sessionKey string, apiHost string, apiKey string) (*framework.SkyScannerResponse, error) // TODO: convert to domain model
	StartSearch(arguments *domain.Arguments) (string, error)
}

// FlightRepository handles saving and loading flight data
type FlightRepository interface {
	InitialiseSchema() error
	CreateAirports(airports []domain.Airport) error
	ReadAllAirports() ([]domain.Airport, error)
}

//
// Interfaces for application services...
//

// AirportFinder handles being able to load airport datasets
type AirportFinder interface {
	FindAirports(countryName string, regionName string, excludePrefix string) error
	LoadMajorAirports() (map[string]domain.Airport, error)
}
