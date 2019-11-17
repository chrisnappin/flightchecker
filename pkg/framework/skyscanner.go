package framework

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/chrisnappin/flightchecker/pkg/domain"
)

// SkyScannerResponse is the top-level result
type SkyScannerResponse struct {
	SessionKey  string
	Query       SkyScannerQuery
	Status      string // "UpdatesPending", "UpdatesComplete"
	Itineraries []SkyScannerItinerary
	Legs        []SkyScannerLeg
	Segments    []SkyScannerSegment
	Carriers    []SkyScannerCarrier
	Agents      []SkyScannerAgent
	Places      []SkyScannerPlace
	Currencies  []SkyScannerCurrency
}

// SkyScannerQuery contains the original input parameters
type SkyScannerQuery struct {
	Country          string
	Currency         string
	Locale           string
	Adults           int
	Children         int
	Infants          int
	OriginPlace      string // actually a Place ID
	DestinationPlace string // actually a Place ID
	OutboundDate     string // YYYY-MM-DD
	InboundDate      string // YYYY-MM-DD
	LocationSchema   string // e,g, "Default"
	CabinClass       string
	GroupPricing     bool
}

// SkyScannerItinerary contains a combination of Legs
type SkyScannerItinerary struct {
	OutboundLegID      string `json:"OutboundLegId"`
	InboundLegID       string `json:"InboundLegId"`
	PricingOptions     []SkyScannerPricingOption
	BookingDetailsLink SkyScannerBookingDetailsLink
}

// SkyScannerPricingOption contains a quote from Agent(s)
type SkyScannerPricingOption struct {
	Agents            []int // Agent ID
	QuoteAgeInMinutes int
	Price             float32 // e.g. 758.42
	DeeplinkURL       string  `json:"DeeplinkUrl"`
}

// SkyScannerBookingDetailsLink contains links to get booking details
type SkyScannerBookingDetailsLink struct {
	URI    string `json:"Uri"` // relative REST query to get the booking
	Body   string // query params to add
	Method string // e.g. "PUT"
}

// SkyScannerLeg contains details of part of an itinery, e.g. the outbound flight
type SkyScannerLeg struct {
	ID                 string `json:"Id"`
	SegmentIds         []int
	OriginStation      int
	DestinationStation int
	Departure          string // e.g. "2019-10-14T12:30:00"
	Arrival            string // e.g. "2019-10-15T12:20:00"
	Duration           int    // minutes
	JourneyMode        string // e.g. "Flight"
	Stops              []int
	Carriers           []int
	OperatingCarriers  []int
	Directionality     string // e.g. "Outbound", "Inbound"
	FlightNumbers      []SkyScannerFlightNumber
}

// SkyScannerFlightNumber can be several for same carrier
type SkyScannerFlightNumber struct {
	FlightNumber string // e.g. "433"
	CarrierID    int    `json:"CarrierId"`
}

// SkyScannerSegment details part of a Leg
type SkyScannerSegment struct {
	ID                 int `json:"Id"`
	OriginStation      int
	DestinationStation int
	DepartureDateTime  string // e.g. "2019-10-14T12:30:00"
	ArrivalDateTime    string // e.g. "2019-10-14T14:06:00"
	Carrier            int
	OperatingCarrier   int
	Duration           int
	FlightNumber       string
	JourneyMode        string // e.g. "Flight"
	Directionality     string // e.g. "Outbound", "Inbound"
}

// SkyScannerCarrier details an Airline
type SkyScannerCarrier struct {
	ID          int    `json:"Id"`
	Code        string `json:"Code,omitempty"`
	Name        string
	ImageURL    string `json:"ImageUrl"`
	DisplayCode string
}

// SkyScannerAgent details who has quoted
type SkyScannerAgent struct {
	ID                 int `json:"Id"`
	Name               string
	ImageURL           string `json:"ImageUrl"`
	Status             string // e.g. "UpdatesPending", "UpdatesComplete"
	OptimisedForMobile bool
	Type               string // e.g. "Airline", "TravelAgent"
}

// SkyScannerPlace details somewhere involved in the quote
type SkyScannerPlace struct {
	ID       int    `json:"Id"`
	ParentID *int   `json:"ParentId,omitempty"`
	Code     string // e.g. Airport IATA code "LHR", ISO country code "GB"
	Type     string // e.g. "Airport", "Country", "City"
	Name     string // e.g. "London Heathrow"
}

// SkyScannerCurrency details how to format monetary values
type SkyScannerCurrency struct {
	Code                        string // e.g. "GBP",
	Symbol                      string // e.g. "£"
	ThousandsSeparator          string // e.g. ","
	DecimalSeparator            string // e.g. "."
	SymbolOnLeft                bool
	SpaceBetweenAmountAndSymbol bool
	RoundingCoefficient         int // e.g. 0
	DecimalDigits               int // e.g. 2
}

// SkyScannerService handles calling the sky scanner API.
type SkyScannerService struct {
	logger domain.Logger
}

// NewSkyScannerService creates a new instance.
func NewSkyScannerService(logger domain.Logger) *SkyScannerService {
	return &SkyScannerService{logger}
}

// PollForQuotes calls the skyscanner "Poll session results" operation, to look for quotes
func (service *SkyScannerService) PollForQuotes(sessionKey string, apiHost string, apiKey string) (*SkyScannerResponse, error) {
	const pageIndex = 0
	const pageSize = 10

	service.logger.Debugf("GET first page of %d quotes...", pageSize)
	url := fmt.Sprintf("https://%s/apiservices/pricing/uk2/v1.0/%%7B%s%%7D?pageIndex=%d&pageSize=%d",
		apiHost, sessionKey, pageIndex, pageSize)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-host", apiHost)
	req.Header.Add("x-rapidapi-key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	service.logger.Debug("Response received...")
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		service.logInvalidResponse(res)
		return nil, errors.New("Request rejected")
	}

	var r SkyScannerResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// StartSearch calls the skyscanner "Create session" operation, which returns a session key.
func (service *SkyScannerService) StartSearch(arguments *domain.Arguments) (string, error) {
	url := "https://" + arguments.APIHost + "/apiservices/pricing/v1.0"

	service.logger.Debug("POST flight search to create session...")
	payload, err := service.formatSearchPayload(arguments)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return "", err
	}

	req.Header.Add("x-rapidapi-host", arguments.APIHost)
	req.Header.Add("x-rapidapi-key", arguments.APIKey)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil
	}

	if res.StatusCode != 201 {
		service.logInvalidResponse(res)
		return "", errors.New("Request rejected")
	}

	// in practice this returns 201, with body of "{}"
	// Location [0] http header in response is of the form:
	// http://partners.api.skyscanner.net/apiservices/pricing/uk2/v1.0/c1b3deed-0419-4296-a4e3-3b5afa6b8ea9
	// can't actually use this (gives a 403) but need to extract the session key from the end of the url
	// then perform a GET using this to get the actual quotes...
	location := res.Header.Get("Location")
	if location != "" {
		lastIndex := strings.LastIndex(location, "/")
		if lastIndex != -1 {
			key := location[lastIndex+1 : len(location)]
			service.logger.Infof("The session key is %s", key)
			return key, nil
		}
		return "", errors.New("No session key found in URL: " + location)
	}
	return "", errors.New("No Location returned in response")
}

func (service *SkyScannerService) formatSearchPayload(arguments *domain.Arguments) (string, error) {
	const country = "GB"
	const currency = "GBP"
	const locale = "en-GB"
	const cabinClass = "economy"    // economy, premiumeconomy, business, first
	const groupPricing = true       // true = price for all, false = price for 1 adult
	const dateFormat = "2006-01-02" // i.e. YYYY-MM-DD

	holidayStartDate, err := time.Parse(dateFormat, arguments.OutboundDate)
	if err != nil {
		return "", err
	}
	holidayEndDate := holidayStartDate.AddDate(0, 0, arguments.HolidayDuration)

	return fmt.Sprintf("inboundDate=%s&cabinClass=%s&children=%d&infants=%d&country=%s&"+
		"currency=%s&locale=%s&originPlace=%s-sky&destinationPlace=%s-sky&outboundDate=%s&adults=%d&groupPricing=%t",
		holidayEndDate.Format(dateFormat), cabinClass, arguments.Children, arguments.Infants, country, currency, locale,
		arguments.Origin, arguments.Destination, arguments.OutboundDate, arguments.Adults, groupPricing), nil

}

func (service *SkyScannerService) logInvalidResponse(res *http.Response) {
	service.logger.Errorf("Request rejected with %s", res.Status)

	if res.ContentLength != 0 {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			service.logger.Error(err)
		} else {
			service.logger.Errorf("Response was: %s", body)
		}
	}
}
