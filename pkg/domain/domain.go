package domain

// Airport includes details of each airport.
type Airport struct {
	Name     string
	IataCode string
	Country  string
	Region   string
}

// AirportMapFilter filters a map of airports, returning an array of values that pass the filter function.
func AirportMapFilter(airports map[string]Airport, f func(Airport) bool) []Airport {
	filteredValues := make([]Airport, 0)
	for _, value := range airports {
		if f(value) {
			filteredValues = append(filteredValues, value)
		}
	}
	return filteredValues
}

// AirportMapValues returns all values of a map of airports.
func AirportMapValues(airports map[string]Airport) []Airport {
	values := make([]Airport, 0)
	for _, value := range airports {
		values = append(values, value)
	}
	return values
}

// Arguments encapsulates all quote criteria and supporting info needed.
type Arguments struct {
	Origin          string // IATA airport code
	Destination     string // IATA airport code
	Adults          int    // adults are over 16
	Children        int    // children are 1-16
	Infants         int    // infants are 0-12 months
	OutboundDate    string // must be YYYY-MM-DD
	HolidayDuration int    // in nights
	APIHost         string // from your rapidapi account
	APIKey          string // from your rapidapi account
}
