package application

import (
	"errors"
	"testing"

	"github.com/chrisnappin/flightchecker/mocks"
	"github.com/chrisnappin/flightchecker/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var dummyCountries = map[string]string{
	"AA": "Country 1",
	"BB": "Country 2",
}

var dummyRegions = map[string]string{
	"AA-AA": "Region 1",
	"AA-BB": "Region 2",
	"BB-AA": "Region 3",
	"BB-BB": "Region 4",
}

var airport1 = domain.Airport{
	Name:     "Airport1",
	Region:   "Region1",
	Country:  "Country1",
	IataCode: "Code1",
}

var airport2 = domain.Airport{
	Name:     "Airport2",
	Region:   "Region2",
	Country:  "Country2",
	IataCode: "Code2",
}

var dummyAirports = map[string]domain.Airport{
	airport1.IataCode: airport1,
	airport2.IataCode: airport2,
}

// TestLoadMajorAirports_HappyPath tests LoadMajorAirports when all is good.
func TestLoadMajorAirports_HappyPath(t *testing.T) {
	mockLogger := &mocks.Logger{}
	mockLoader := &mocks.AirportDataLoader{}
	service := NewFindAirportsService(mockLogger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(dummyCountries, nil)
	mockLoader.On("LoadRegions", mock.Anything).Return(dummyRegions, nil)
	mockLoader.On("LoadAirports", mock.Anything, mock.Anything, mock.Anything).Return(dummyAirports, nil)

	result, err := service.LoadMajorAirports()
	assert.Equal(t, dummyAirports, result, "Wrong results")
	assert.Nil(t, err, "Expected no error")
}

// TestLoadMajorAirports_HappyPath tests LoadMajorAirports when countries error.
func TestLoadMajorAirports_CountriesFail(t *testing.T) {
	mockLogger := &mocks.Logger{}
	mockLoader := &mocks.AirportDataLoader{}
	service := NewFindAirportsService(mockLogger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(nil, errors.New("Oops"))

	result, err := service.LoadMajorAirports()
	assert.Nil(t, result, "Expected no result")
	assert.Error(t, err, "Expected an error")
}

// TestLoadMajorAirports_HappyPath tests LoadMajorAirports when regions error.
func TestLoadMajorAirports_RegionsFail(t *testing.T) {
	mockLogger := &mocks.Logger{}
	mockLoader := &mocks.AirportDataLoader{}
	service := NewFindAirportsService(mockLogger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(dummyCountries, nil)
	mockLoader.On("LoadRegions", mock.Anything).Return(nil, errors.New("Oops"))

	result, err := service.LoadMajorAirports()
	assert.Nil(t, result, "Expected no result")
	assert.Error(t, err, "Expected an error")
}

// TestLoadMajorAirports_HappyPath tests LoadMajorAirports when airports error.
func TestLoadMajorAirports_AirportFail(t *testing.T) {
	mockLogger := &mocks.Logger{}
	mockLoader := &mocks.AirportDataLoader{}
	service := NewFindAirportsService(mockLogger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(dummyCountries, nil)
	mockLoader.On("LoadRegions", mock.Anything).Return(dummyRegions, nil)
	mockLoader.On("LoadAirports", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Oops"))

	result, err := service.LoadMajorAirports()
	assert.Nil(t, result, "Expected no result")
	assert.Error(t, err, "Expected an error")
}

// TestFindAirports_HappyPath tests FindAirports with a prefix.
func TestFindAirports_WithPrefix(t *testing.T) {
	mockLogger := &mocks.Logger{}
	mockLoader := &mocks.AirportDataLoader{}
	service := NewFindAirportsService(mockLogger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(dummyCountries, nil)
	mockLoader.On("LoadRegions", mock.Anything).Return(dummyRegions, nil)
	mockLoader.On("LoadAirports", mock.Anything, mock.Anything, mock.Anything).Return(dummyAirports, nil)

	mockLogger.On("Info", "Matching Airports")
	// no matching result logged

	err := service.FindAirports("Country1", "Region1", "A") // filters out the result
	assert.Nil(t, err, "Expected no error")
	mockLoader.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// TestFindAirports_HappyPath tests FindAirports without a prefix.
func TestFindAirports_WithoutPrefix(t *testing.T) {
	mockLogger := &mocks.Logger{}
	mockLoader := &mocks.AirportDataLoader{}
	service := NewFindAirportsService(mockLogger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(dummyCountries, nil)
	mockLoader.On("LoadRegions", mock.Anything).Return(dummyRegions, nil)
	mockLoader.On("LoadAirports", mock.Anything, mock.Anything, mock.Anything).Return(dummyAirports, nil)

	mockLogger.On("Info", "Matching Airports")
	mockLogger.On("Infof", "Name: %s, Code: %s, Region: %s", airport1.Name, airport1.IataCode, airport1.Region)

	err := service.FindAirports("Country1", "Region1", "") // doesn't filter out the result
	assert.Nil(t, err, "Expected no error")
	mockLoader.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// TestFindAirports_HappyPath tests FindAirports when it fails.
func TestFindAirports_Fails(t *testing.T) {
	mockLogger := &mocks.Logger{}
	mockLoader := &mocks.AirportDataLoader{}
	service := NewFindAirportsService(mockLogger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(nil, errors.New("Oops"))

	err := service.FindAirports("AA", "BB", "")
	assert.Error(t, err, "Expected an error")
}
