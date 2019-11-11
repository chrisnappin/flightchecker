package application

import (
	"errors"
	"testing"

	"github.com/chrisnappin/flightchecker/pkg/domain"
	"github.com/chrisnappin/flightchecker/pkg/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAirportDataLoader struct {
	mock.Mock
}

func (mock *MockAirportDataLoader) LoadCountries(filename string) (map[string]string, error) {
	args := mock.Called(filename)
	if args.Get(0) != nil {
		return args.Get(0).(map[string]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *MockAirportDataLoader) LoadRegions(filename string) (map[string]string, error) {
	args := mock.Called(filename)
	if args.Get(0) != nil {
		return args.Get(0).(map[string]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *MockAirportDataLoader) LoadAirports(filename string, countries map[string]string, regions map[string]string) (map[string]domain.Airport, error) {
	args := mock.Called(filename)
	if args.Get(0) != nil {
		return args.Get(0).(map[string]domain.Airport), args.Error(1)
	}
	return nil, args.Error(1)
}

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
	Name:     "AAAA",
	Region:   "BBBB",
	Country:  "CCCC",
	IataCode: "AA",
}

var airport2 = domain.Airport{
	Name:     "XXXX",
	Region:   "YYYY",
	Country:  "ZZZZ",
	IataCode: "XX",
}

var dummyAirports = map[string]domain.Airport{
	"AA": airport1,
	"XX": airport2,
}

var logger = framework.NewLogWrapper("test", true)

// TestLoadMajorAirports_HappyPath tests LoadMajorAirports when all is good.
func TestLoadMajorAirports_HappyPath(t *testing.T) {
	mockLoader := &MockAirportDataLoader{}
	service := NewFindAirportsService(logger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(dummyCountries, nil)
	mockLoader.On("LoadRegions", mock.Anything).Return(dummyRegions, nil)
	mockLoader.On("LoadAirports", mock.Anything, mock.Anything, mock.Anything).Return(dummyAirports, nil)

	result, err := service.LoadMajorAirports()
	assert.Equal(t, dummyAirports, result, "Wrong results")
	assert.Nil(t, err, "Expected no error")
}

// TestLoadMajorAirports_HappyPath tests LoadMajorAirports when countries error.
func TestLoadMajorAirports_CountriesFail(t *testing.T) {
	mockLoader := &MockAirportDataLoader{}
	service := NewFindAirportsService(logger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(nil, errors.New("Oops"))

	result, err := service.LoadMajorAirports()
	assert.Nil(t, result, "Expected no result")
	assert.Error(t, err, "Expected an error")
}

// TestLoadMajorAirports_HappyPath tests LoadMajorAirports when regions error.
func TestLoadMajorAirports_RegionsFail(t *testing.T) {
	mockLoader := &MockAirportDataLoader{}
	service := NewFindAirportsService(logger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(dummyCountries, nil)
	mockLoader.On("LoadRegions", mock.Anything).Return(nil, errors.New("Oops"))

	result, err := service.LoadMajorAirports()
	assert.Nil(t, result, "Expected no result")
	assert.Error(t, err, "Expected an error")
}

// TestLoadMajorAirports_HappyPath tests LoadMajorAirports when airports error.
func TestLoadMajorAirports_AirportFail(t *testing.T) {
	mockLoader := &MockAirportDataLoader{}
	service := NewFindAirportsService(logger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(dummyCountries, nil)
	mockLoader.On("LoadRegions", mock.Anything).Return(dummyRegions, nil)
	mockLoader.On("LoadAirports", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Oops"))

	result, err := service.LoadMajorAirports()
	assert.Nil(t, result, "Expected no result")
	assert.Error(t, err, "Expected an error")
}

// TestFilter_HappyPath tests filter when all is well.
func TestFilter_HappyPath(t *testing.T) {
	mockLoader := &MockAirportDataLoader{}
	service := NewFindAirportsService(logger, mockLoader)

	expected := []domain.Airport{airport2}
	result := service.filter(dummyAirports, func(airport domain.Airport) bool {
		return airport.IataCode == "XX"
	})
	assert.Equal(t, result, expected, "Wrong result")
}

// TestFindAirports_HappyPath tests FindAirports with a prefix.
// TODO: test the logger output
func TestFindAirports_WithPrefix(t *testing.T) {
	mockLoader := &MockAirportDataLoader{}
	service := NewFindAirportsService(logger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(dummyCountries, nil)
	mockLoader.On("LoadRegions", mock.Anything).Return(dummyRegions, nil)
	mockLoader.On("LoadAirports", mock.Anything, mock.Anything, mock.Anything).Return(dummyAirports, nil)

	err := service.FindAirports("AA", "BB", "CC")
	assert.Nil(t, err, "Expected no error")
}

// TestFindAirports_HappyPath tests FindAirports without a prefix.
// TODO: test the logger output
func TestFindAirports_WithoutPrefix(t *testing.T) {
	mockLoader := &MockAirportDataLoader{}
	service := NewFindAirportsService(logger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(dummyCountries, nil)
	mockLoader.On("LoadRegions", mock.Anything).Return(dummyRegions, nil)
	mockLoader.On("LoadAirports", mock.Anything, mock.Anything, mock.Anything).Return(dummyAirports, nil)

	err := service.FindAirports("AA", "BB", "")
	assert.Nil(t, err, "Expected no error")
}

// TestFindAirports_HappyPath tests FindAirports when it fails.
func TestFindAirports_Fails(t *testing.T) {
	mockLoader := &MockAirportDataLoader{}
	service := NewFindAirportsService(logger, mockLoader)

	mockLoader.On("LoadCountries", mock.Anything).Return(nil, errors.New("Oops"))

	err := service.FindAirports("AA", "BB", "")
	assert.Error(t, err, "Expected an error")
}
