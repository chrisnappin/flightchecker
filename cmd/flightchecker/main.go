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

	airports, err := finder.LoadMajorAirports()
	if err != nil {
		mainLogger.Fatal(err)
	}

	argumentsLoader := framework.NewArgumentsLoader(framework.NewLogWrapper("argumentsLoader", true))
	arguments, err := argumentsLoader.Load("arguments.json")
	if err != nil {
		mainLogger.Fatal(err)
	}

	originAirport, exists := airports[arguments.Origin]
	if !exists {
		mainLogger.Fatalf("Origin airport code %s unknown", arguments.Origin)
	}

	destinationAirport, exists := airports[arguments.Destination]
	if !exists {
		mainLogger.Fatalf("Destination airport code %s unknown", arguments.Destination)
	}

	mainLogger.Infof("Looking for flights from %s staying for %d nights\n",
		arguments.OutboundDate, arguments.HolidayDuration)
	mainLogger.Infof("from %s (%s) in %s, %s",
		originAirport.Name, originAirport.IataCode, originAirport.Region, originAirport.Country)
	mainLogger.Infof("to %s (%s) in %s, %s",
		destinationAirport.Name, destinationAirport.IataCode, destinationAirport.Region, destinationAirport.Country)
	mainLogger.Infof("for %d adults, %d children, %d infants\n",
		arguments.Adults, arguments.Children, arguments.Infants)

	sqliteRepository := framework.NewSQLiteRepository(framework.NewLogWrapper("sqliteRepository", true))
	quoteRepository := application.NewQuoteRepository(framework.NewLogWrapper("QuoteRepository", true), sqliteRepository)
	err = quoteRepository.Initialise()
	if err != nil {
		mainLogger.Fatal(err)
	}

	//quoteFinder := skyscanner.NewQuoteFinder(logwrapper.NewLogger("skyscanner", true))
	//quoteFinder.FindFlightQuotes(arguments)

	os.Exit(0)
}
