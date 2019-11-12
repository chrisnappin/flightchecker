package main

import (
	"os"

	"github.com/chrisnappin/flightchecker/pkg/application"
	"github.com/chrisnappin/flightchecker/pkg/framework"
)

func main() {
	mainLogger := framework.NewLogWrapper("flightchecker", true)
	loader := framework.NewAirportDataLoader(framework.NewLogWrapper("airportDataLoader", true))
	finder := application.NewFindAirportsService(framework.NewLogWrapper("airportLoader", true), loader)
	argumentsLoader := framework.NewArgumentsLoader(framework.NewLogWrapper("argumentsLoader", true))

	sqliteRepository := framework.NewSQLiteRepository(framework.NewLogWrapper("sqliteRepository", true))
	quoteRepository := application.NewQuoteRepository(framework.NewLogWrapper("QuoteRepository", true), sqliteRepository)
	err := quoteRepository.Initialise()
	if err != nil {
		mainLogger.Fatal(err)
	}

	skyscanner := framework.NewSkyScannerQuoter(framework.NewLogWrapper("skyscannerQuoter", true))
	flightQuoter := application.NewQuoteForFlightsService(framework.NewLogWrapper("quoteForFlights", true), argumentsLoader, finder, skyscanner)
	flightQuoter.QuoteForFlights("arguments.json")

	os.Exit(0)
}
