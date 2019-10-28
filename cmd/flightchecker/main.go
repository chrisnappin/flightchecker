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

	airports, err := loadAirports()
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

func loadAirports() (map[string]data.Airport, error) {
	loader := data.NewLoader(logwrapper.NewLogger("data", true))
	countries, err := loader.ReadCountries("data/airports/countries.csv")
	if err != nil {
		return nil, err
	}

	regions, err := loader.ReadRegions("data/airports/regions.csv")
	if err != nil {
		return nil, err
	}

	airports, err := loader.ReadAirports("data/airports/airports.csv", countries, regions)
	if err != nil {
		return nil, err
	}

	airportsMap := make(map[string]data.Airport)
	for _, airport := range airports {
		airportsMap[airport.IataCode] = airport
	}
	return airportsMap, nil
}
