package main

import (
	"os"
	"strings"

	"github.com/chrisnappin/flightchecker/pkg/application"
	"github.com/chrisnappin/flightchecker/pkg/domain"
	"github.com/chrisnappin/flightchecker/pkg/framework"
)

func main() {
	mainLogger := framework.NewLogWrapper("airports", true)
	staticDataLoader := framework.NewStaticDataLoader(framework.NewLogWrapper("staticDataLoader", true))
	airportLoader := application.NewAirportLoader(framework.NewLogWrapper("airportLoader", true), staticDataLoader)

	airports, err := airportLoader.LoadMajorAirports()
	if err != nil {
		mainLogger.Fatal(err)
	}

	const countryName = "United Kingdom"
	const regionName = "England"
	filteredAirports := airportLoader.Filter(airports, func(a domain.Airport) bool {
		return a.Country == countryName && a.Region == regionName && !strings.HasPrefix(a.Name, "RAF ")
	})

	mainLogger.Info("Matching Airports\n")
	for _, airport := range filteredAirports {
		mainLogger.Infof("Name: %s, Code: %s, Region: %s", airport.Name, airport.IataCode, airport.Region)
	}

	os.Exit(0)
}
