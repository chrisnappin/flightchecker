package main

import (
	"os"

	"github.com/chrisnappin/flightchecker/pkg/application"
	"github.com/chrisnappin/flightchecker/pkg/framework"
)

func main() {
	logger := framework.NewLogWrapper("airports", true)
	loader := framework.NewAirportDataLoader(framework.NewLogWrapper("airportDataLoader", true))
	service := application.NewFindAirportsService(framework.NewLogWrapper("findAirportsService", true), loader)

	const countryName = "United Kingdom"
	const regionName = "England"
	const exclude = "RAF "
	err := service.FindAirports(countryName, regionName, exclude)
	if err != nil {
		logger.Fatal(err)
	}

	os.Exit(0)
}
