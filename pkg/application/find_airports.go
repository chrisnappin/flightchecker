package application

import (
	"strings"

	"github.com/chrisnappin/flightchecker/pkg/domain"
)

// FindAirportsService handles finding a range of airports.
type FindAirportsService struct {
	logger domain.Logger
	loader AirportDataLoader
}

// NewFindAirportsService creates a new instance.
func NewFindAirportsService(logger domain.Logger, loader AirportDataLoader) *FindAirportsService {
	return &FindAirportsService{logger, loader}
}

// FindAirports logs all airports within the specified country and region, excluding any matching the prefix (if set).
func (service *FindAirportsService) FindAirports(countryName string, regionName string, excludePrefix string) error {
	airports, err := service.LoadMajorAirports()
	if err != nil {
		return err
	}

	filteredAirports := domain.AirportMapFilter(airports, func(a domain.Airport) bool {
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
func (service *FindAirportsService) LoadMajorAirports() (map[string]domain.Airport, error) {
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
