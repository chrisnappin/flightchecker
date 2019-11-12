package application

import "github.com/chrisnappin/flightchecker/pkg/domain"

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

//
// Interfaces for application services...
//

// AirportFinder handles being able to load airport datasets
type AirportFinder interface {
	FindAirports(countryName string, regionName string, excludePrefix string) error
	LoadMajorAirports() (map[string]domain.Airport, error)
}
