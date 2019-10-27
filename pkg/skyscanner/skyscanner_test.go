package skyscanner

import (
	"testing"

	"github.com/chrisnappin/flightchecker/pkg/arguments"
	"github.com/chrisnappin/flightchecker/pkg/logwrapper"
	"github.com/stretchr/testify/assert"
)

func TestFormatSearchPayload(t *testing.T) {
	arguments := arguments.Arguments{
		Origin:       "LHR",
		Destination:  "LAX",
		Adults:       2,
		Children:     2,
		Infants:      0,
		OutboundDate: "2019-11-01",
		InboundDate:  "2019-11-10",
		APIHost:      "test.com",
		APIKey:       "testKey",
	}
	quoteFinder := NewQuoteFinder(logwrapper.NewLogger("skyscanner", true))
	actual := quoteFinder.formatSearchPayload(&arguments)

	expected := "inboundDate=2019-11-10&cabinClass=economy&children=2&infants=0&country=GB&currency=GBP&locale=en-GB" +
		"&originPlace=LHR-sky&destinationPlace=LAX-sky&outboundDate=2019-11-01&adults=2&groupPricing=true"

	assert.Equal(t, expected, actual, "Incorrect payload")
}
