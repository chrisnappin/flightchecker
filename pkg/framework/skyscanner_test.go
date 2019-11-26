package framework

import (
	"testing"
	"time"

	"github.com/chrisnappin/flightchecker/mocks"
	"github.com/chrisnappin/flightchecker/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	gock "gopkg.in/h2non/gock.v1"
)

var dummyArguments = domain.Arguments{
	Origin:          "LHR",
	Destination:     "LAX",
	Adults:          2,
	Children:        2,
	Infants:         0,
	OutboundDate:    "2019-11-01",
	HolidayDuration: 9,
	APIHost:         "test.com",
	APIKey:          "testKey",
}

// TestFormatSearchPayload tests formatting the payload of search parameters, with valid input.
func TestFormatSearchPayload_AllValid(t *testing.T) {
	expected := "inboundDate=2019-11-10&cabinClass=economy&children=2&infants=0&country=GB&currency=GBP&locale=en-GB" +
		"&originPlace=LHR-sky&destinationPlace=LAX-sky&outboundDate=2019-11-01&adults=2&groupPricing=true"

	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}
	actual, err := service.formatSearchPayload(&dummyArguments)

	assert.Equal(t, expected, actual, "Incorrect payload")
	assert.Nil(t, err, "Error not expected")
}

// TestFormatSearchPayload tests formatting the payload of search parameters, with invalid input.
func TestFormatSearchPayload_InvalidDate(t *testing.T) {
	brokenArguments := dummyArguments           // struct of primitives so can copy by value
	brokenArguments.OutboundDate = "01/02/2003" // not YYYY-MM-DD

	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}
	actual, err := service.formatSearchPayload(&brokenArguments)

	assert.Equal(t, "", actual, "No payload expected")
	assert.Error(t, err, "Error expected")
}

// TestStartSearch_HappyPath tests starting a search, when the response is success.
func TestStartSearch_HappyPath(t *testing.T) {
	defer gock.Off()

	gock.New("https://test.com").
		Post("/apiservices/pricing/v1.0").
		HeaderPresent("x-rapidapi-host").
		HeaderPresent("x-rapidapi-key").
		Reply(201).
		AddHeader("Location", "https://test.com/aaa/bbb/ccc/abc")

	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}

	mockLogger.On("Debug", mock.Anything)
	mockLogger.On("Infof", mock.Anything, mock.Anything)

	sessionKey, err := service.StartSearch(&dummyArguments)
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, "abc", sessionKey, "Invalid session key")
	assert.Equal(t, gock.IsDone(), true)
}

// TestStartSearch_NoLocation tests starting a search, when the response doesn't include a Location header.
func TestStartSearch_NoLocation(t *testing.T) {
	defer gock.Off()

	gock.New("https://test.com").
		Post("/apiservices/pricing/v1.0").
		HeaderPresent("x-rapidapi-host").
		HeaderPresent("x-rapidapi-key").
		Reply(201)

	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}

	mockLogger.On("Debug", mock.Anything)
	mockLogger.On("Infof", mock.Anything, mock.Anything)

	sessionKey, err := service.StartSearch(&dummyArguments)
	assert.Error(t, err, "Error expected")
	assert.Equal(t, "", sessionKey, "No session key expected")
	assert.Equal(t, gock.IsDone(), true)
}

// TestStartSearch_NoSessionKey tests starting a search, when the response has no session key.
func TestStartSearch_NoSessionKey(t *testing.T) {
	defer gock.Off()

	gock.New("https://test.com").
		Post("/apiservices/pricing/v1.0").
		HeaderPresent("x-rapidapi-host").
		HeaderPresent("x-rapidapi-key").
		Reply(201).
		AddHeader("Location", "wibble") // no / character...

	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}

	mockLogger.On("Debug", mock.Anything)
	mockLogger.On("Infof", mock.Anything, mock.Anything)

	sessionKey, err := service.StartSearch(&dummyArguments)
	assert.Error(t, err, "Error expected")
	assert.Equal(t, "", sessionKey, "No session key expected")
	assert.Equal(t, gock.IsDone(), true)
}

// TestStartSearch_Rejected tests starting a search, when the response is unauthorized.
func TestStartSearch_Rejected(t *testing.T) {
	defer gock.Off()

	gock.New("https://test.com").
		Post("/apiservices/pricing/v1.0").
		Reply(401)

	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}

	mockLogger.On("Debug", mock.Anything)
	mockLogger.On("Infof", mock.Anything, mock.Anything)
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything)

	sessionKey, err := service.StartSearch(&dummyArguments)
	assert.Error(t, err, "Error expected")
	assert.Equal(t, "", sessionKey, "No session key expected")
	assert.Equal(t, gock.IsDone(), true)
}

// TestPollForQuote_HappyPath tests polling for quotes, when the response is success.
func TestPollForQuote_HappyPath(t *testing.T) {
	defer gock.Off()

	expected := SkyScannerResponse{
		SessionKey: "abc",
		Query: SkyScannerQuery{
			Country:  "GB",
			Currency: "GBP",
			Locale:   "en-GB",
		},
	}

	gock.New("https://test.com/apiservices/pricing/uk2/v1.0/%7Babc%7D").
		MatchParam("pageIndex", "0").
		MatchParam("pageSize", "10").
		HeaderPresent("x-rapidapi-host").
		HeaderPresent("x-rapidapi-key").
		Reply(200).
		JSON(expected)

	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}

	mockLogger.On("Debug", mock.Anything)
	mockLogger.On("Debugf", mock.Anything, mock.Anything)

	actual, err := service.PollForQuotes("abc", "test.com", "testKey")
	assert.Nil(t, err, "No error expected")
	assert.EqualValues(t, &expected, actual, "Invalid response")
	assert.Equal(t, gock.IsDone(), true)
}

// TestPollForQuote_ServerError tests polling for quotes, when the response is a server error.
func TestPollForQuote_ServerError(t *testing.T) {
	defer gock.Off()

	gock.New("https://test.com/apiservices/pricing/uk2/v1.0/%7Babc%7D").
		MatchParam("pageIndex", "0").
		MatchParam("pageSize", "10").
		Reply(500).
		BodyString("Oops")

	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}

	mockLogger.On("Debug", mock.Anything)
	mockLogger.On("Debugf", mock.Anything, mock.Anything)
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything)

	actual, err := service.PollForQuotes("abc", "test.com", "testKey")
	assert.Error(t, err, "Error expected")
	assert.Nil(t, actual, "No response expected")
	assert.Equal(t, gock.IsDone(), true)
}

// TestPollForQuote_InvalidResponse tests polling for quotes, when the response is invalid JSON.
func TestPollForQuote_InvalidResponse(t *testing.T) {
	defer gock.Off()

	gock.New("https://test.com/apiservices/pricing/uk2/v1.0/%7Babc%7D").
		MatchParam("pageIndex", "0").
		MatchParam("pageSize", "10").
		Reply(200).
		BodyString("{\"wibble\":1234,") // un-terminated JSON

	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}

	mockLogger.On("Debug", mock.Anything)
	mockLogger.On("Debugf", mock.Anything, mock.Anything)

	actual, err := service.PollForQuotes("abc", "test.com", "testKey")
	assert.Error(t, err, "Error expected")
	assert.Nil(t, actual, "No response expected")
	assert.Equal(t, gock.IsDone(), true)
}

// TestConvertToDomain_Empty tests converting to domain values, when the response is empty.
func TestConvertToDomain_Empty(t *testing.T) {
	input := SkyScannerResponse{}
	expected := []*domain.Itinerary{}

	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}

	actual, err := service.convertToDomain(&input)
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, expected, actual, "Wrong output")
}

// TestConvertToDomain_Populated tests converting to domain values, when the response is populated and valid.
func TestConvertToDomain_PopulatedValid(t *testing.T) {
	expected := []*domain.Itinerary{
		&domain.Itinerary{
			SupplierName: "Agent1",
			SupplierType: "Airline",
			Amount:       10099,
			OutboundJourney: &domain.Journey{
				ID:        "leg1",
				Direction: domain.Outbound,
				Flights: []*domain.Flight{
					&domain.Flight{
						ID:                 "10",
						FlightNumber:       nil,
						StartAirport:       nil,
						StartTime:          time.Date(2019, time.October, 14, 8, 35, 0, 0, time.UTC),
						DestinationAirport: nil,
						DestinationTime:    time.Date(2019, time.October, 14, 9, 30, 0, 0, time.UTC),
						Duration:           55 * time.Minute,
					},
				},
				Duration:  65 * time.Minute,
				StartTime: time.Date(2019, time.October, 14, 8, 30, 0, 0, time.UTC),
				EndTime:   time.Date(2019, time.October, 14, 9, 35, 0, 0, time.UTC),
			},
			InboundJourney: &domain.Journey{
				ID:        "leg2",
				Direction: domain.Inbound,
				Flights: []*domain.Flight{
					&domain.Flight{
						ID:                 "20",
						FlightNumber:       nil,
						StartAirport:       nil,
						StartTime:          time.Date(2019, time.October, 16, 10, 20, 0, 0, time.UTC),
						DestinationAirport: nil,
						DestinationTime:    time.Date(2019, time.October, 16, 11, 30, 0, 0, time.UTC),
						Duration:           70 * time.Minute,
					},
				},
				Duration:  80 * time.Minute,
				StartTime: time.Date(2019, time.October, 16, 10, 15, 0, 0, time.UTC),
				EndTime:   time.Date(2019, time.October, 16, 11, 35, 0, 0, time.UTC),
			},
		},
	}

	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}

	actual, err := service.convertToDomain(getExampleResponse(valid))
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, expected, actual, "Wrong output")
}

// TestConvertToDomain_Errors tests converting to domain values, when the response contains various types of errors.
func TestConvertToDomain_Errors(t *testing.T) {
	mockLogger := &mocks.Logger{}
	service := SkyScannerService{mockLogger}

	testCases := []struct {
		option responseOption
	}{
		{missingAgent},               // response references a non-existant agent
		{missingOutboundLeg},         // esponse references a non-existant outbound leg
		{invalidOutboundLegStart},    // response outbound leg has an invalid format start time
		{invalidOutboundLegEnd},      // response outbound leg has an invalid format arrival time
		{missingInboundLeg},          // response references a non-existant inbound leg
		{invalidInboundLegStart},     // response inbound leg has an invalid format start time
		{invalidInboundLegEnd},       // response inbound leg has an invalid format arrival time
		{missingOutboundSegment},     // response outbound leg is missing
		{invalidInboundSegmentStart}, // segment within inbound leg has an invalid format start date time value
		{invalidOutboundSegmentEnd},  // segment within outbound leg has an invalid format end date time value
	}

	for _, testrCase := range testCases {
		actual, err := service.convertToDomain(getExampleResponse(testrCase.option))
		assert.Error(t, err, "Error expected")
		assert.Nil(t, actual, "No output expected")
	}
}

type responseOption int

const (
	valid responseOption = iota
	missingAgent
	missingOutboundLeg
	invalidOutboundLegStart
	invalidOutboundLegEnd
	missingInboundLeg
	invalidInboundLegStart
	invalidInboundLegEnd
	missingOutboundSegment
	invalidInboundSegmentStart
	invalidOutboundSegmentEnd
)

// getExampleResponse returns an example skyscanner format response
func getExampleResponse(option responseOption) *SkyScannerResponse {
	response := SkyScannerResponse{
		Itineraries: []SkyScannerItinerary{
			SkyScannerItinerary{
				OutboundLegID: "leg1",
				InboundLegID:  "leg2",
				PricingOptions: []SkyScannerPricingOption{
					SkyScannerPricingOption{
						Agents: []int{1111},
						Price:  100.99,
					},
				},
			},
		},
		Legs: []SkyScannerLeg{
			SkyScannerLeg{
				ID:             "leg1",
				SegmentIds:     []int{10},
				Departure:      "2019-10-14T08:30:00",
				Arrival:        "2019-10-14T09:35:00",
				Duration:       65,
				Directionality: "Outbound",
			},
			SkyScannerLeg{
				ID:             "leg2",
				SegmentIds:     []int{20},
				Departure:      "2019-10-16T10:15:00",
				Arrival:        "2019-10-16T11:35:00",
				Duration:       80,
				Directionality: "Inbound",
			},
		},
		Segments: []SkyScannerSegment{
			SkyScannerSegment{
				ID:                10,
				DepartureDateTime: "2019-10-14T08:35:00",
				ArrivalDateTime:   "2019-10-14T09:30:00",
				Carrier:           2222,
				Duration:          55,
				FlightNumber:      "123",
			},
			SkyScannerSegment{
				ID:                20,
				DepartureDateTime: "2019-10-16T10:20:00",
				ArrivalDateTime:   "2019-10-16T11:30:00",
				Carrier:           2222,
				Duration:          70,
				FlightNumber:      "456",
			},
		},
		Agents: []SkyScannerAgent{
			SkyScannerAgent{
				ID:   1111,
				Name: "Agent1",
				Type: "Airline",
			},
		},
	}

	switch option {
	case missingAgent:
		response.Agents = []SkyScannerAgent{}

	case missingOutboundLeg:
		response.Legs = response.Legs[1:1]

	case invalidOutboundLegStart:
		response.Legs[0].Departure = "wibble"

	case invalidOutboundLegEnd:
		response.Legs[0].Arrival = "wibble"

	case missingInboundLeg:
		response.Legs = response.Legs[0:0]

	case invalidInboundLegStart:
		response.Legs[1].Departure = "wibble"

	case invalidInboundLegEnd:
		response.Legs[1].Arrival = "wibble"

	case missingOutboundSegment:
		response.Segments = response.Segments[1:1]

	case invalidInboundSegmentStart:
		response.Segments[1].DepartureDateTime = "wibble"

	case invalidOutboundSegmentEnd:
		response.Segments[0].ArrivalDateTime = "wibble"
	}

	return &response
}
