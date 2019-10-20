package skyscanner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/chrisnappin/flightchecker/pkg/arguments"
)

// Response is the top-level result
type Response struct {
	SessionKey  string
	Query       Query
	Status      string // "UpdatesPending", "UpdatesComplete"
	Itineraries []Itinerary
	Legs        []Leg
	Segments    []Segment
	Carriers    []Carrier
	Agents      []Agent
	Places      []Place
	Currencies  []Currency
}

// Query contains the original input parameters
type Query struct {
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

// Itinerary contains a combination of Legs
type Itinerary struct {
	OutboundLegID      string `json:"OutboundLegId"`
	InboundLegID       string `json:"InboundLegId"`
	PricingOptions     []PricingOption
	BookingDetailsLink BookingDetailsLink
}

// PricingOption contains a quote from Agent(s)
type PricingOption struct {
	Agents            []int // Agent ID
	QuoteAgeInMinutes int
	Price             float32 // e.g. 758.42
	DeeplinkURL       string  `json:"DeeplinkUrl"`
}

// BookingDetailsLink contains links to get booking details
type BookingDetailsLink struct {
	URI    string `json:"Uri"` // relative REST query to get the booking
	Body   string // query params to add
	Method string // e.g. "PUT"
}

// Leg contains details of part of an itinery, e.g. the outbound flight
type Leg struct {
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
	FlightNumbers      []FlightNumber
}

// FlightNumber can be several for same carrier
type FlightNumber struct {
	FlightNumber string // e.g. "433"
	CarrierID    int    `json:"CarrierId"`
}

// Segment details part of a Leg
type Segment struct {
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

// Carrier details an Airline
type Carrier struct {
	ID          int    `json:"Id"`
	Code        string `json:"Code,omitempty"`
	Name        string
	ImageURL    string `json:"ImageUrl"`
	DisplayCode string
}

// Agent details who has quoted
type Agent struct {
	ID                 int `json:"Id"`
	Name               string
	ImageURL           string `json:"ImageUrl"`
	Status             string // e.g. "UpdatesPending", "UpdatesComplete"
	OptimisedForMobile bool
	Type               string // e.g. "Airline", "TravelAgent"
}

// Place details somewhere involved in the quote
type Place struct {
	ID       int    `json:"Id"`
	ParentID *int   `json:"ParentId,omitempty"`
	Code     string // e.g. Airport IATA code "LHR", ISO country code "GB"
	Type     string // e.g. "Airport", "Country", "City"
	Name     string // e.g. "London Heathrow"
}

// Currency details how to format monetary values
type Currency struct {
	Code                        string // e.g. "GBP",
	Symbol                      string // e.g. "£"
	ThousandsSeparator          string // e.g. ","
	DecimalSeparator            string // e.g. "."
	SymbolOnLeft                bool
	SpaceBetweenAmountAndSymbol bool
	RoundingCoefficient         int // e.g. 0
	DecimalDigits               int // e.g. 2
}

// QuoteFinder handles finding quotes.
type QuoteFinder struct {
	logger *logrus.Logger
}

// NewQuoteFinder creates a new instance.
func NewQuoteFinder(logger *logrus.Logger) *QuoteFinder {
	return &QuoteFinder{logger}
}

const (
	quotesCompleteStatus = "UpdatesComplete"
)

// FindFlightQuotes calls the skyscanner API to find some flight quotes.
func (q *QuoteFinder) FindFlightQuotes(arguments *arguments.Arguments) {

	/*
	 * The way the skyscanner API works is that we first make our search,
	 * then poll for results.
	 */
	sessionKey, err := q.startSearch(arguments)
	if err != nil {
		q.logger.Fatal(err)
	}

	/*
	 * In practice, initial polls return partial results and have status of "UpdatesPending"
	 * Then after typically 20-30 seconds we get a fully populated result with status of "UpdatesComplete".
	 */
	var response *Response
	for index := 0; index < 6; index++ {

		q.logger.Debugf("Poll %d...\n", index)
		response, err = q.pollForQuotes(sessionKey, arguments.APIHost, arguments.APIKey)
		if err != nil {
			q.logger.Fatal(err)
		}

		q.logger.Debugf("Polled for quotes, status is %s, found %d itineries\n", response.Status, len(response.Itineraries))

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

func (q *QuoteFinder) outputQuotes(response *Response) {
	q.logger.Infof("Quote completed, found %d flights", len(response.Itineraries))
	for _, itinerary := range response.Itineraries {
		for _, pricingOption := range itinerary.PricingOptions {
			q.logger.Infof("Flight costs %.2f\n", pricingOption.Price)
		}
	}
}

// pollForQuotes calls the skyscanner "Poll session results" operation, to look for quotes
func (q *QuoteFinder) pollForQuotes(sessionKey string, apiHost string, apiKey string) (*Response, error) {
	pageIndex := 0
	pageSize := 10

	q.logger.Debugf("GET first page of %d quotes...\n", pageSize)
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

	q.logger.Debugln("Response received...")
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// q.logger.Debugf("Response code: %d, content length: %d\n", res.StatusCode, res.ContentLength)
	// fmt.Println("Headers are:")
	// for key, values := range res.Header {
	// 	for index, value := range values {
	// 		fmt.Printf("    %s [%d] => %s\n", key, index, value)
	// 	}
	// }
	// fmt.Printf("Response body is: [%s]\n", string(body))

	if res.StatusCode != 200 {
		q.logInvalidResponse(res)
		return nil, errors.New("Request rejected")
	}

	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("\nSuccessfully unmarshalled to JSON\n\n\nValue is %v", r)
	return &r, nil
}

// startSearch calls the skyscanner "Create session" operation, which returns a session key.
func (q *QuoteFinder) startSearch(arguments *arguments.Arguments) (string, error) {
	url := "https://" + arguments.APIHost + "/apiservices/pricing/v1.0"

	country := "GB"
	currency := "GBP"
	locale := "en-GB"
	cabinClass := "economy" // economy, premiumeconomy, business, first
	//groupPricing := false // group = price for all, false = price for 1 adult

	payload := strings.NewReader(fmt.Sprintf("inboundDate=%s&cabinClass=%s&children=%d&infants=%d&country=%s&"+
		"currency=%s&locale=%s&originPlace=%s&destinationPlace=%s&outboundDate=%s&adults=%d",
		arguments.InboundDate, cabinClass, arguments.Children, arguments.Infants, country, currency, locale,
		arguments.Origin, arguments.Destination, arguments.OutboundDate, arguments.Adults))

	q.logger.Debugln("POST flight search to create session...")
	req, err := http.NewRequest("POST", url, payload)
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
	//fmt.Println("Response received...")
	// defer res.Body.Close()
	// body, err := ioutil.ReadAll(res.Body)
	// _, err = ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return "", nil
	// }

	if res.StatusCode != 201 {
		q.logInvalidResponse(res)
		return "", errors.New("Request rejected")
	}
	//fmt.Printf("Response code: %d, content length: %d\n", res.StatusCode, res.ContentLength)
	//fmt.Println("Headers are:")
	// for key, values := range res.Header {
	// 	for index, value := range values {
	// 		fmt.Printf("    %s [%d] => %s\n", key, index, value)
	// 	}
	// }
	// fmt.Printf("\nResponse body is: [%s]", string(body))

	// in practice this returns 201, with body of "{}"
	// Location [0] http header in response is of the form:
	// http://partners.api.skyscanner.net/apiservices/pricing/uk2/v1.0/c1b3deed-0419-4296-a4e3-3b5afa6b8ea9
	// can't actually use this (gives a 403) but need to extract the session key from the end of the url
	// then perform a GET using this to get the actual quotes...
	location := res.Header.Get("Location")
	key := location[strings.LastIndex(location, "/")+1 : len(location)]
	q.logger.Infof("The session key is %s\n", key)

	return key, nil
}

func (q *QuoteFinder) logInvalidResponse(res *http.Response) {
	q.logger.Errorf("Request rejected with %s", res.Status)

	if res.ContentLength != 0 {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			q.logger.Error(err)
		} else {
			q.logger.Errorf("Response was: %s", body)
		}
	}
}
