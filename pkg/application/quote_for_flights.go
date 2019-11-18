package application

import (
	"fmt"
	"time"

	"github.com/chrisnappin/flightchecker/pkg/domain"
	"github.com/chrisnappin/flightchecker/pkg/framework"
)

type quoteForFlightsService struct {
	logger           domain.Logger
	loader           ArgumentsLoader
	finder           AirportFinder
	skyScannerQuoter SkyScannerQuoter
	flightRepository FlightRepository
}

// NewQuoteForFlightsService creates a new instance.
func NewQuoteForFlightsService(logger domain.Logger, loader ArgumentsLoader, finder AirportFinder,
	skyScannerQuoter SkyScannerQuoter, flightRepository FlightRepository) *quoteForFlightsService {
	return &quoteForFlightsService{logger, loader, finder, skyScannerQuoter, flightRepository}
}

const (
	quotesCompleteStatus = "UpdatesComplete"
)

// QuoteForFlights finds some quotes for flights defined in the arguments.
func (service *quoteForFlightsService) QuoteForFlights(argumentsFilename string) {

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
	var response *framework.SkyScannerResponse
	for index := 0; index < 6; index++ {

		service.logger.Debugf("Poll %d...", index)
		response, err = service.skyScannerQuoter.PollForQuotes(sessionKey, arguments.APIHost, arguments.APIKey) // TODO: convert to domain model
		if err != nil {
			service.logger.Fatal(err)
		}

		service.logger.Debugf("Polled for quotes, status is %s, found %d itineries",
			response.Status, len(response.Itineraries))

		if response.Status == quotesCompleteStatus {
			service.logger.Debugf("Quotes are complete...")
			break
		}

		time.Sleep(10 * time.Second)
	}

	if response.Status == quotesCompleteStatus {
		service.outputQuotes(response)
	} else {
		service.logger.Fatalf("Quotes not completed in time, status is still %s", response.Status)
	}
}

// loadArguments attempts to load the details to quote for, and returns the arguments, origin airport, dest airport,
// or an error.
func (service *quoteForFlightsService) loadArguments(argumentsFilename string, airports map[string]domain.Airport) (
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

func (service *quoteForFlightsService) outputQuotes(response *framework.SkyScannerResponse) {
	// maps agent id to agent name
	agents := make(map[int]string)
	for _, agent := range response.Agents {
		agents[agent.ID] = agent.Name
	}

	// maps leg id to Leg
	legs := make(map[string]framework.SkyScannerLeg)
	for _, leg := range response.Legs {
		legs[leg.ID] = leg
	}

	// maps segment id to Segment
	segments := make(map[int]framework.SkyScannerSegment)
	for _, segment := range response.Segments {
		segments[segment.ID] = segment
	}

	// maps carrier id to Carrier
	carriers := make(map[int]framework.SkyScannerCarrier)
	for _, carrier := range response.Carriers {
		carriers[carrier.ID] = carrier
	}

	service.logger.Infof("Quote completed, found %d flights", len(response.Itineraries))
	for _, itinerary := range response.Itineraries {
		for _, pricingOption := range itinerary.PricingOptions {
			for _, agentID := range pricingOption.Agents {
				service.logger.Infof("Flight with %s is %.2f", agents[agentID], pricingOption.Price)

				outboundLeg := legs[itinerary.OutboundLegID]
				service.logger.Infof("Outbound Leg: from %s to %s (%d minutes)",
					outboundLeg.Departure, outboundLeg.Arrival, outboundLeg.Duration)
				for index, segmentID := range outboundLeg.SegmentIds {
					segment := segments[segmentID]
					service.logger.Infof("Outbound segment %d is flight %s from %s to %s (%d minutes) with %s (%s)",
						index, segment.FlightNumber, segment.DepartureDateTime, segment.ArrivalDateTime,
						segment.Duration, carriers[segment.Carrier].Name, carriers[segment.Carrier].Code)
				}

				inboundLeg := legs[itinerary.InboundLegID]
				service.logger.Infof("Inbound Leg: from %s to %s (%d minutes)",
					inboundLeg.Departure, inboundLeg.Arrival, inboundLeg.Duration)
				for index, segmentID := range inboundLeg.SegmentIds {
					segment := segments[segmentID]
					service.logger.Infof("Inbound segment %d is flight %s from %s to %s (%d minutes) with %s (%s)",
						index, segment.FlightNumber, segment.DepartureDateTime, segment.ArrivalDateTime,
						segment.Duration, carriers[segment.Carrier].Name, carriers[segment.Carrier].Code)
				}
			}
		}
	}
}
