package framework

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/chrisnappin/flightchecker/pkg/domain"
)

// StaticDataLoader handles being able to load static data
type StaticDataLoader interface {
	LoadCountries(filename string) (map[string]string, error)
	LoadRegions(filename string) (map[string]string, error)
	LoadAirports(filename string, countries map[string]string, regions map[string]string) (map[string]domain.Airport, error)
}

// staticDataLoader handles loading static data.
type staticDataLoaderService struct {
	logger domain.Logger
}

// NewStaticDataLoader creates a new instance.
func NewStaticDataLoader(logger domain.Logger) StaticDataLoader {
	return &staticDataLoaderService{logger}
}

// LoadCountries returns a map of countries, keyed by country code.
// The data is read from the specified CSV file.
func (l *staticDataLoaderService) LoadCountries(filename string) (map[string]string, error) {
	const indexCode = 1
	const indexName = 2

	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	countries := make(map[string]string)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		countries[line[indexCode]] = line[indexName]
	}
	l.logger.Debugf("Read %d countries", len(countries))
	return countries, nil
}

// LoadRegions returns a map of regions, keyed by region code.
// The data is read from the specified CSV file.
func (l *staticDataLoaderService) LoadRegions(filename string) (map[string]string, error) {
	const indexCode = 1
	const indexName = 3

	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	regions := make(map[string]string)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		regions[line[indexCode]] = line[indexName]
	}
	l.logger.Debugf("Read %d regions", len(regions))
	return regions, nil
}

// LoadAirports returns a slice of Airports, with country and region names populated.
// The data is read from the specified CSV file.
func (l *staticDataLoaderService) LoadAirports(filename string, countries map[string]string, regions map[string]string) (map[string]domain.Airport, error) {
	const indexName = 3
	const indexCountryCode = 8
	const indexRegionCode = 9
	const indexIataCode = 13

	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	airports := make(map[string]domain.Airport)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		iataCode := line[indexIataCode]

		// only include major airports (assigned IATA codes), miss out header row
		if len(iataCode) > 0 && iataCode != "iata_code" {
			countryName, exists := countries[line[indexCountryCode]]
			if !exists {
				l.logger.Fatalf("Countries missing name for code %s", line[indexCountryCode])
			}

			regionName, exists := regions[line[indexRegionCode]]
			if !exists {
				l.logger.Fatalf("Regions missing name for code %s", line[indexRegionCode])
			}

			airports[iataCode] = domain.Airport{
				Name:     line[indexName],
				IataCode: iataCode,
				Country:  countryName,
				Region:   regionName,
			}
		}
	}
	l.logger.Debugf("Read %d airports", len(airports))
	return airports, nil
}
