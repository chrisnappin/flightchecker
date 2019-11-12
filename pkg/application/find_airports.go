package application

import (
	"strings"

	"github.com/chrisnappin/flightchecker/pkg/domain"
)

// AirportDataLoader handles being able to load airport data
type AirportDataLoader interface {
	LoadCountries(filename string) (map[string]string, error)
	LoadRegions(filename string) (map[string]string, error)
	LoadAirports(filename string, countries map[string]string, regions map[string]string) (map[string]domain.Airport, error)
}

// findAirportsService handles finding a range of airports.
type findAirportsService struct {
	logger domain.Logger
	loader AirportDataLoader
}

// NewFindAirportsService creates a new instance.
func NewFindAirportsService(logger domain.Logger, staticDataLoader AirportDataLoader) *findAirportsService {
	return &findAirportsService{logger, staticDataLoader}
}

// FindAirports logs all airports within the specified country and region, excluding any matching the prefix (if set).
func (service *findAirportsService) FindAirports(countryName string, regionName string, excludePrefix string) error {
	airports, err := service.LoadMajorAirports()
	if err != nil {
		return err
	}

	filteredAirports := service.filter(airports, func(a domain.Airport) bool {
		if excludePrefix != "" {
			return a.Country == countryName && a.Region == regionName && !strings.HasPrefix(a.Name, excludePrefix)
		}
		return a.Country == countryName && a.Region == regionName
	})

	service.logger.Info("Matching Airports")
	for _, airport := range filteredAirports {
		service.logger.Infof("Name: %s, Code: %s, Region: %s", airport.Name, airport.IataCode, airport.Region)
	}

	return nil
}

// LoadMajorAirports returns a map of all major airports, keyed by IATA code
func (service *findAirportsService) LoadMajorAirports() (map[string]domain.Airport, error) {
	countries, err := service.loader.LoadCountries("data/airports/countries.csv")
	if err != nil {
		return nil, err
	}

	regions, err := service.loader.LoadRegions("data/airports/regions.csv")
	if err != nil {
		return nil, err
	}

	airports, err := service.loader.LoadAirports("data/airports/airports.csv", countries, regions)
	if err != nil {
		return nil, err
	}
	return airports, nil
}

// filter returns an array of values that pass a filter function
func (service *findAirportsService) filter(airports map[string]domain.Airport, f func(domain.Airport) bool) []domain.Airport {
	filteredValues := make([]domain.Airport, 0)
	for _, value := range airports {
		if f(value) {
			filteredValues = append(filteredValues, value)
		}
	}
	return filteredValues
}
