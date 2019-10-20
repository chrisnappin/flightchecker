package data

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Airport includes details of each airport.
type Airport struct {
	Name     string
	IataCode string
	Country  string
	Region   string
}

// Loader handles loading data.
type Loader struct {
	logger *logrus.Logger
}

// NewLoader creates a new instance.
func NewLoader(logger *logrus.Logger) *Loader {
	return &Loader{logger}
}

// ReadCountries returns a map of countries, keyed by country code.
// The data is read from the specified CSV file.
func (l *Loader) ReadCountries(filename string) (map[string]string, error) {
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

// ReadRegions returns a map of regions, keyed by region code.
// The data is read from the specified CSV file.
func (l *Loader) ReadRegions(filename string) (map[string]string, error) {
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

// ReadAirports returns a slice of Airports, with country and region names populated.
// The data is read from the specified CSV file.
func (l *Loader) ReadAirports(filename string, countries map[string]string, regions map[string]string) ([]Airport, error) {
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
	var airports []Airport
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

			airports = append(airports, Airport{
				Name:     line[indexName],
				IataCode: iataCode,
				Country:  countryName,
				Region:   regionName,
			})
		}
	}
	l.logger.Debugf("Read %d airports", len(airports))
	return airports, nil
}
