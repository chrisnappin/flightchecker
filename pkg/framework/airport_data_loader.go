package framework

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/chrisnappin/flightchecker/pkg/domain"
)

// AirportDataLoaderService handles loading airport data from static CSV files.
type AirportDataLoaderService struct {
	logger domain.Logger
}

// NewAirportDataLoader creates a new instance.
func NewAirportDataLoader(logger domain.Logger) *AirportDataLoaderService {
	return &AirportDataLoaderService{logger}
}

// LoadCountries returns a map of countries, keyed by country code.
// The data is read from the specified CSV file.
func (service *AirportDataLoaderService) LoadCountries(filename string) (map[string]string, error) {
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
	service.logger.Debugf("Read %d countries", len(countries))
	return countries, nil
}

// LoadRegions returns a map of regions, keyed by region code.
// The data is read from the specified CSV file.
func (service *AirportDataLoaderService) LoadRegions(filename string) (map[string]string, error) {
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
	service.logger.Debugf("Read %d regions", len(regions))
	return regions, nil
}

// LoadAirports returns a slice of Airports, with country and region names populated.
// The data is read from the specified CSV file.
func (service *AirportDataLoaderService) LoadAirports(filename string, countries map[string]string,
	regions map[string]string) (map[string]domain.Airport, error) {
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
				service.logger.Fatalf("Countries missing name for code %s", line[indexCountryCode])
			}

			regionName, exists := regions[line[indexRegionCode]]
			if !exists {
				service.logger.Fatalf("Regions missing name for code %s", line[indexRegionCode])
			}

			airports[iataCode] = domain.Airport{
				Name:     line[indexName],
				IataCode: iataCode,
				Country:  countryName,
				Region:   regionName,
			}
		}
	}
	service.logger.Debugf("Read %d airports", len(airports))
	return airports, nil
}
