package main

import (
	"os"

	"github.com/chrisnappin/flightchecker/pkg/arguments"
	"github.com/chrisnappin/flightchecker/pkg/logwrapper"
	"github.com/chrisnappin/flightchecker/pkg/skyscanner"
)

func main() {
	logger := logwrapper.NewLogger("flightchecker", true)

	arguments, err := arguments.Load("arguments.json")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Read arguments: %+v", arguments)

	quoteFinder := skyscanner.NewQuoteFinder(logwrapper.NewLogger("skyscanner", true))
	quoteFinder.FindFlightQuotes(arguments)

	os.Exit(0)
}
