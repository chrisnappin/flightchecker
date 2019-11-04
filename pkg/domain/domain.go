package domain

// Airport includes details of each airport.
type Airport struct {
	Name     string
	IataCode string
	Country  string
	Region   string
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
