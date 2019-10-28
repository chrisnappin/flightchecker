package arguments

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

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

// Load reads a JSON file of arguments.
func Load(filename string) (*Arguments, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var arguments Arguments

	err = json.Unmarshal(bytes, &arguments)
	if err != nil {
		return nil, err
	}

	return &arguments, nil
}
