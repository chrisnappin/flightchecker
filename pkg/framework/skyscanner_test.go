package framework

import (
	"testing"

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
	service := skyScannerService{mockLogger}
	actual, err := service.formatSearchPayload(&dummyArguments)

	assert.Equal(t, expected, actual, "Incorrect payload")
	assert.Nil(t, err, "Error not expected")
}

// TestFormatSearchPayload tests formatting the payload of search parameters, with invalid input.
func TestFormatSearchPayload_InvalidDate(t *testing.T) {
	brokenArguments := dummyArguments           // struct of primitives so can copy by value
	brokenArguments.OutboundDate = "01/02/2003" // not YYYY-MM-DD

	mockLogger := &mocks.Logger{}
	service := skyScannerService{mockLogger}
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
	service := skyScannerService{mockLogger}

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
	service := skyScannerService{mockLogger}

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
	service := skyScannerService{mockLogger}

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
	service := skyScannerService{mockLogger}

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
	service := skyScannerService{mockLogger}

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
	service := skyScannerService{mockLogger}

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
	service := skyScannerService{mockLogger}

	mockLogger.On("Debug", mock.Anything)
	mockLogger.On("Debugf", mock.Anything, mock.Anything)

	actual, err := service.PollForQuotes("abc", "test.com", "testKey")
	assert.Error(t, err, "Error expected")
	assert.Nil(t, actual, "No response expected")
	assert.Equal(t, gock.IsDone(), true)
}
