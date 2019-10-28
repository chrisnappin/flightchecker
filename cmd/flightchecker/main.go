package main

import (
	"os"

	"github.com/chrisnappin/flightchecker/pkg/arguments"
	"github.com/chrisnappin/flightchecker/pkg/data"
	"github.com/chrisnappin/flightchecker/pkg/logwrapper"
	"github.com/chrisnappin/flightchecker/pkg/skyscanner"
)

func main() {
	logger := logwrapper.NewLogger("flightchecker", true)
	loader := data.NewLoader(logwrapper.NewLogger("data", true))

	airports, err := loader.ReadMajorAirports()
	if err != nil {
		logger.Fatal(err)
	}

	arguments, err := arguments.Load("arguments.json")
	if err != nil {
		logger.Fatal(err)
	}

	originAirport, exists := airports[arguments.Origin]
	if !exists {
		logger.Fatalf("Origin airport code %s unknown", arguments.Origin)
	}

	destinationAirport, exists := airports[arguments.Destination]
	if !exists {
		logger.Fatalf("Destination airport code %s unknown", arguments.Destination)
	}

	logger.Infof("Looking for flights from %s staying for %d nights\n",
		arguments.OutboundDate, arguments.HolidayDuration)
	logger.Infof("from %s (%s) in %s, %s",
		originAirport.Name, originAirport.IataCode, originAirport.Region, originAirport.Country)
	logger.Infof("to %s (%s) in %s, %s",
		destinationAirport.Name, destinationAirport.IataCode, destinationAirport.Region, destinationAirport.Country)
	logger.Infof("for %d adults, %d children, %d infants\n",
		arguments.Adults, arguments.Children, arguments.Infants)

	quoteFinder := skyscanner.NewQuoteFinder(logwrapper.NewLogger("skyscanner", true))
	quoteFinder.FindFlightQuotes(arguments)

	os.Exit(0)
}
