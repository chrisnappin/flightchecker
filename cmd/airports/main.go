package main

import (
	"os"
	"strings"

	"github.com/chrisnappin/flightchecker/pkg/data"
	"github.com/chrisnappin/flightchecker/pkg/logwrapper"
)

func main() {
    logger := logwrapper.NewLogger("airports", true)

	loader := data.NewLoader(logwrapper.NewLogger("data", true))

	countries, err := loader.ReadCountries("data/airports/countries.csv")
	if err != nil {
		logger.Fatal(err)
	}

	regions, err := loader.ReadRegions("data/airports/regions.csv")
	if err != nil {
		logger.Fatal(err)
	}

	airports, err := loader.ReadAirports("data/airports/airports.csv", countries, regions)
	if err != nil {
		logger.Fatal(err)
	}

    filteredAirports := findEnglishNonMilitaryAirports(airports)

    logger.Info("Matching Airports\n",)
    for _, airport := range filteredAirports {
		logger.Infof("Name: %s, Code: %s, Region: %s", airport.Name, airport.IataCode, airport.Region)
	}

	os.Exit(0)
}

func filter(airports []data.Airport, f func(data.Airport) bool) []data.Airport {
	filteredValues := make([]data.Airport, 0)
	for _, value := range airports {
		if f(value) {
            filteredValues = append(filteredValues, value)
		}
	}
	return filteredValues
}

func findEnglishNonMilitaryAirports(airports []data.Airport) []data.Airport {
	countryName := "United Kingdom"
	regionName := "England"
    return filter(airports, func(a data.Airport) bool {
		return a.Country == countryName && a.Region == regionName && !strings.HasPrefix(a.Name, "RAF ")
	})
}
