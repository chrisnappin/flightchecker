package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"gitlab.com/chrisnappin/flightchecker/pkg/arguments"
	"gitlab.com/chrisnappin/flightchecker/pkg/skyscanner"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)
	logger.Info("Hello world...")

	// loader := data.NewLoader(logger)

	// countries, err := loader.ReadCountries("data/airports/countries.csv")
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// regions, err := loader.ReadRegions("data/airports/regions.csv")
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// airports, err := loader.ReadAirports("data/airports/airports.csv", countries, regions)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	arguments, err := arguments.Load("arguments.json")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Read arguments: %+v", arguments)

	quoteFinder := skyscanner.NewQuoteFinder(logger)
	quoteFinder.FindFlightQuotes(arguments)

	os.Exit(0)
}
