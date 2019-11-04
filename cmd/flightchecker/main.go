package main

import (
	"os"

	"github.com/chrisnappin/flightchecker/pkg/application"
	"github.com/chrisnappin/flightchecker/pkg/framework"
	"github.com/chrisnappin/flightchecker/pkg/repository"
	//"github.com/chrisnappin/flightchecker/pkg/skyscanner"
)

func main() {
	mainLogger := framework.NewLogger("flightchecker", true)
	staticDataLoader := framework.NewStaticDataLoader(framework.NewLogger("staticDataLoader", true))
	airportLoader := application.NewAirportLoader(framework.NewLogger("airportLoader", true), staticDataLoader)

	airports, err := airportLoader.LoadMajorAirports()
	if err != nil {
		mainLogger.Fatal(err)
	}

	argumentsLoader := framework.NewArgumentsLoader(framework.NewLogger("argumentsLoader", true))
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

	flightRepository := repository.NewRepository(framework.NewLogger("repository", true))
	err = flightRepository.Initialise()
	if err != nil {
		mainLogger.Fatal(err)
	}

	//quoteFinder := skyscanner.NewQuoteFinder(logwrapper.NewLogger("skyscanner", true))
	//quoteFinder.FindFlightQuotes(arguments)

	os.Exit(0)
}
