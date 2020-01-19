package application

import (
	"fmt"
	"time"

	"github.com/chrisnappin/flightchecker/pkg/domain"
)

// QuoteForFlightsService handles finding quotes for flights.
type QuoteForFlightsService struct {
	logger           domain.Logger
	loader           ArgumentsLoader
	finder           AirportFinder
	skyScannerQuoter SkyScannerQuoter
	flightRepository FlightRepository
}

// NewQuoteForFlightsService creates a new instance.
func NewQuoteForFlightsService(logger domain.Logger, loader ArgumentsLoader, finder AirportFinder,
	skyScannerQuoter SkyScannerQuoter, flightRepository FlightRepository) *QuoteForFlightsService {
	return &QuoteForFlightsService{logger, loader, finder, skyScannerQuoter, flightRepository}
}

// QuoteForFlights finds some quotes for flights defined in the arguments.
func (service *QuoteForFlightsService) QuoteForFlights(argumentsFilename string) {

	airports, err := service.finder.LoadMajorAirports()
	if err != nil {
		service.logger.Fatal(err)
	}

	arguments, originAirport, destinationAirport, err := service.loadArguments(argumentsFilename, airports)
	if err != nil {
		service.logger.Fatal(err)
	}

	service.logger.Infof("Looking for flights from %s staying for %d nights",
		arguments.OutboundDate, arguments.HolidayDuration)
	service.logger.Infof("from %s (%s) in %s, %s",
		originAirport.Name, originAirport.IataCode, originAirport.Region, originAirport.Country)
	service.logger.Infof("to %s (%s) in %s, %s",
		destinationAirport.Name, destinationAirport.IataCode, destinationAirport.Region, destinationAirport.Country)
	service.logger.Infof("for %d adults, %d children, %d infants",
		arguments.Adults, arguments.Children, arguments.Infants)

	// TODO: do a cached-data only option where the schema is preserved
	err = service.flightRepository.InitialiseSchema()
	if err != nil {
		service.logger.Fatal(err)
	}

	err = service.flightRepository.CreateAirports(domain.AirportMapValues(airports))
	if err != nil {
		service.logger.Fatal(err)
	}

	/*
	 * The way the skyscanner API works is that we first make our search,
	 * then poll for results.
	 */
	sessionKey, err := service.skyScannerQuoter.StartSearch(arguments)
	if err != nil {
		service.logger.Fatal(err)
	}

	/*
	 * In practice, initial polls return partial results and have status of "UpdatesPending"
	 * Then after typically 20-30 seconds we get a fully populated result with status of "UpdatesComplete".
	 */
	var response *domain.Quote
	for index := 0; index < 6; index++ {

		service.logger.Debugf("Poll %d...", index)
		response, err = service.skyScannerQuoter.PollForQuotes(sessionKey, arguments.APIHost, arguments.APIKey, airports)
		if err != nil {
			service.logger.Fatal(err)
		}

		service.logger.Debugf("Polled for quotes, status is %t, found %d itineries",
			response.Complete, len(response.Itineraries))

		if response.Complete {
			service.logger.Debugf("Quotes are complete...")
			break
		}

		time.Sleep(10 * time.Second)
	}

	if response.Complete {
		service.outputQuotes(response)
	} else {
		service.logger.Fatalf("Quotes not completed in time")
	}
}

// loadArguments attempts to load the details to quote for, and returns the arguments, origin airport, dest airport,
// or an error.
func (service *QuoteForFlightsService) loadArguments(argumentsFilename string, airports map[string]domain.Airport) (
	*domain.Arguments, *domain.Airport, *domain.Airport, error) {
	arguments, err := service.loader.Load(argumentsFilename)
	if err != nil {
		return nil, nil, nil, err
	}

	originAirport, exists := airports[arguments.Origin]
	if !exists {
		return nil, nil, nil, fmt.Errorf("Origin airport code %s unknown", arguments.Origin)
	}

	destinationAirport, exists := airports[arguments.Destination]
	if !exists {
		return nil, nil, nil, fmt.Errorf("Destination airport code %s unknown", arguments.Destination)
	}

	return arguments, &originAirport, &destinationAirport, nil
}

func (service *QuoteForFlightsService) outputQuotes(response *domain.Quote) {
	const dayTimeFormat = "2006-01-02 15:04"
	service.logger.Infof("Quote completed, found %d flights", len(response.Itineraries))
	for _, itinerary := range response.Itineraries {
		service.logger.Infof("Flight with %s (%s) is %s",
			itinerary.SupplierName, itinerary.SupplierType, formatPrice(itinerary.Amount))

		outboundJourney := itinerary.OutboundJourney
		service.logger.Infof("Outbound Journey takes %s", formatFlightDuration(outboundJourney.Duration))
		for index, flight := range outboundJourney.Flights {
			service.logger.Infof("Outbound flight %d is flight %s%s (%s) from %s (%s) to %s (%s)",
				index+1, flight.FlightNumber.CarrierCode, flight.FlightNumber.FlightNumber,
				flight.FlightNumber.CarrierName,
				flight.StartAirport.Name, flight.StartAirport.IataCode,
				flight.DestinationAirport.Name, flight.DestinationAirport.IataCode)
			service.logger.Infof("%s to %s",
				flight.StartTime.Format(dayTimeFormat),
				flight.DestinationTime.Format(dayTimeFormat))
		}

		inboundJourney := itinerary.InboundJourney
		service.logger.Infof("Inbound Journey takes %s", formatFlightDuration(inboundJourney.Duration))
		for index, flight := range inboundJourney.Flights {
			service.logger.Infof("Inbound flight %d is flight %s%s (%s) from %s (%s) to %s (%s)",
				index+1, flight.FlightNumber.CarrierCode, flight.FlightNumber.FlightNumber,
				flight.FlightNumber.CarrierName,
				flight.StartAirport.Name, flight.StartAirport.IataCode,
				flight.DestinationAirport.Name, flight.DestinationAirport.IataCode)
			service.logger.Infof("%s to %s",
				flight.StartTime.Format(dayTimeFormat),
				flight.DestinationTime.Format(dayTimeFormat))
		}
	}
}

func formatPrice(amount int) string {
	return fmt.Sprintf("%.2f", float64(amount)/100.0)
}

func formatFlightDuration(duration time.Duration) string {
	minutes := duration.Minutes()
	return fmt.Sprintf("%.f hrs, %d mins", minutes/60.0, int(minutes)%60)
}
