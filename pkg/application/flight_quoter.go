package application

import (
	"time"

	"github.com/chrisnappin/flightchecker/pkg/domain"
	"github.com/chrisnappin/flightchecker/pkg/framework"
)

// FlightQuoter handles quoting for flights.
type FlightQuoter interface {
	Quote(*domain.Arguments)
}

type flightQuoterService struct {
	logger           framework.Logger
	skyScannerQuoter framework.SkyScannerQuoter
}

// NewFlightQuoter creates a new instance.
func NewFlightQuoter(logger framework.Logger, skyScannerQuoter framework.SkyScannerQuoter) FlightQuoter {
	return &flightQuoterService{logger, skyScannerQuoter}
}

const (
	quotesCompleteStatus = "UpdatesComplete"
)

// Quote find some quotes for flights defined in the arguments.
func (q *flightQuoterService) Quote(arguments *domain.Arguments) {

	/*
	 * The way the skyscanner API works is that we first make our search,
	 * then poll for results.
	 */
	sessionKey, err := q.skyScannerQuoter.StartSearch(arguments)
	if err != nil {
		q.logger.Fatal(err)
	}

	/*
	 * In practice, initial polls return partial results and have status of "UpdatesPending"
	 * Then after typically 20-30 seconds we get a fully populated result with status of "UpdatesComplete".
	 */
	var response *framework.SkyScannerResponse
	for index := 0; index < 6; index++ {

		q.logger.Debugf("Poll %d...", index)
		response, err = q.skyScannerQuoter.PollForQuotes(sessionKey, arguments.APIHost, arguments.APIKey)
		if err != nil {
			q.logger.Fatal(err)
		}

		q.logger.Debugf("Polled for quotes, status is %s, found %d itineries", response.Status, len(response.Itineraries))

		if response.Status == quotesCompleteStatus {
			q.logger.Debugf("Quotes are complete...")
			break
		}

		time.Sleep(10 * time.Second)
	}

	if response.Status == quotesCompleteStatus {
		q.outputQuotes(response)
	} else {
		q.logger.Fatalf("Quotes not completed in time, status is still %s", response.Status)
	}
}

func (q *flightQuoterService) outputQuotes(response *framework.SkyScannerResponse) {
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

	q.logger.Infof("Quote completed, found %d flights", len(response.Itineraries))
	for _, itinerary := range response.Itineraries {
		for _, pricingOption := range itinerary.PricingOptions {
			for _, agentID := range pricingOption.Agents {
				q.logger.Infof("Flight with %s is %.2f", agents[agentID], pricingOption.Price)

				outboundLeg := legs[itinerary.OutboundLegID]
				q.logger.Infof("Outbound Leg: from %s to %s (%d minutes)",
					outboundLeg.Departure, outboundLeg.Arrival, outboundLeg.Duration)
				for index, segmentID := range outboundLeg.SegmentIds {
					segment := segments[segmentID]
					q.logger.Infof("Outbound segment %d is flight %s from %s to %s (%d minutes) with %s (%s)",
						index, segment.FlightNumber, segment.DepartureDateTime, segment.ArrivalDateTime,
						segment.Duration, carriers[segment.Carrier].Name, carriers[segment.Carrier].Code)
				}

				inboundLeg := legs[itinerary.InboundLegID]
				q.logger.Infof("Inbound Leg: from %s to %s (%d minutes)",
					inboundLeg.Departure, inboundLeg.Arrival, inboundLeg.Duration)
				for index, segmentID := range inboundLeg.SegmentIds {
					segment := segments[segmentID]
					q.logger.Infof("Inbound segment %d is flight %s from %s to %s (%d minutes) with %s (%s)",
						index, segment.FlightNumber, segment.DepartureDateTime, segment.ArrivalDateTime,
						segment.Duration, carriers[segment.Carrier].Name, carriers[segment.Carrier].Code)
				}
			}
		}
	}
}
