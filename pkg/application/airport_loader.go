package application

import (
	"github.com/chrisnappin/flightchecker/pkg/domain"
	"github.com/chrisnappin/flightchecker/pkg/framework"
)

// AirportLoader handles being able to load airport data
type AirportLoader interface {
	LoadMajorAirports() (map[string]domain.Airport, error)
	Filter(airports map[string]domain.Airport, f func(domain.Airport) bool) []domain.Airport
}

// airportLoaderService handles loading data.
type airportLoaderService struct {
	logger           framework.Logger
	staticDataLoader framework.StaticDataLoader
}

// NewAirportLoader creates a new instance.
func NewAirportLoader(logger framework.Logger, staticDataLoader framework.StaticDataLoader) AirportLoader {
	return &airportLoaderService{logger, staticDataLoader}
}

// LoadMajorAirports returns a map of all major airports, keyed by IATA code
func (l *airportLoaderService) LoadMajorAirports() (map[string]domain.Airport, error) {
	countries, err := l.staticDataLoader.LoadCountries("data/airports/countries.csv")
	if err != nil {
		return nil, err
	}

	regions, err := l.staticDataLoader.LoadRegions("data/airports/regions.csv")
	if err != nil {
		return nil, err
	}

	airports, err := l.staticDataLoader.LoadAirports("data/airports/airports.csv", countries, regions)
	if err != nil {
		return nil, err
	}
	return airports, nil
}

// Filter returns an array of values that pass a filter function
func (l *airportLoaderService) Filter(airports map[string]domain.Airport, f func(domain.Airport) bool) []domain.Airport {
	filteredValues := make([]domain.Airport, 0)
	for _, value := range airports {
		if f(value) {
			filteredValues = append(filteredValues, value)
		}
	}
	return filteredValues
}
