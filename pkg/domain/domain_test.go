package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var airport1 = Airport{
	Name:     "Airport1",
	Region:   "Region1",
	Country:  "Country1",
	IataCode: "Code1",
}

var airport2 = Airport{
	Name:     "Airport2",
	Region:   "Region2",
	Country:  "Country2",
	IataCode: "Code2",
}

var dummyAirports = map[string]Airport{
	airport1.IataCode: airport1,
	airport2.IataCode: airport2,
}

// TestAirportMapFilter_Matching tests filter when an entry is found.
func TestAirportMapFilter_Matching(t *testing.T) {
	expected := []Airport{airport2}
	result := AirportMapFilter(dummyAirports, func(airport Airport) bool {
		return airport.IataCode == airport2.IataCode
	})
	assert.Equal(t, result, expected, "Wrong result")
}

// TestAirportMapFilter_NotMatching tests filter when an entry is not found.
func TestAirportMapFilter_NotMatching(t *testing.T) {
	expected := []Airport{}
	result := AirportMapFilter(dummyAirports, func(airport Airport) bool {
		return airport.IataCode == "XYZ"
	})
	assert.Equal(t, result, expected, "Wrong result")
}

// TestAirportMapValues_Populated tests when entries are populated.
func TestAirportMapValues_Populated(t *testing.T) {
	expected := []Airport{airport1, airport2}
	result := AirportMapValues(dummyAirports)
	assert.Equal(t, result, expected, "Wrong result")
}

// TestAirportMapValues_Empty tests when entries are not populated.
func TestAirportMapValues_Empty(t *testing.T) {
	input := make(map[string]Airport)
	expected := []Airport{}
	result := AirportMapValues(input)
	assert.Equal(t, result, expected, "Wrong result")
}
