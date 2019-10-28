package main

import (
	"os"
	"strings"

	"github.com/chrisnappin/flightchecker/pkg/data"
	"github.com/chrisnappin/flightchecker/pkg/logwrapper"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logwrapper.NewLogger("airports", true)
	loader := data.NewLoader(logwrapper.NewLogger("data", true))

	loadAndFilterAirports(logger, loader)

	os.Exit(0)
}

func loadAndFilterAirports(logger *logrus.Entry, loader *data.Loader) {
	airports, err := loader.ReadMajorAirports()
	if err != nil {
		logger.Fatal(err)
	}

	const countryName = "United Kingdom"
	const regionName = "England"
	filteredAirports := loader.Filter(airports, func(a data.Airport) bool {
		return a.Country == countryName && a.Region == regionName && !strings.HasPrefix(a.Name, "RAF ")
	})

	logger.Info("Matching Airports\n")
	for _, airport := range filteredAirports {
		logger.Infof("Name: %s, Code: %s, Region: %s", airport.Name, airport.IataCode, airport.Region)
	}
}
