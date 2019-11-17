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
	skyscanner := framework.NewSkyScannerService(framework.NewLogWrapper("skyscannerQuoter", true))

	recreateDatabase := true // TODO - set this via command line flag
	db, err := framework.OpenDatabase("./data/flightchecker.db", recreateDatabase)
	if err != nil {
		mainLogger.Fatal(err)
	}
	defer db.Close()

	flightRepository := framework.NewFlightRepository(framework.NewLogWrapper("sqliteRepository", true), db)
	flightQuoter := application.NewQuoteForFlightsService(framework.NewLogWrapper("quoteForFlights", true),
		argumentsLoader, finder, skyscanner, flightRepository)

	flightQuoter.QuoteForFlights("arguments.json")

	os.Exit(0)
}
