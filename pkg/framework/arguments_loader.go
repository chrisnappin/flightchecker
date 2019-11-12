package framework

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/chrisnappin/flightchecker/pkg/domain"
)

// argumentsLoaderService handles loading arguments from a JSON file.
type argumentsLoaderService struct {
	logger domain.Logger
}

// NewArgumentsLoader creates a new instance.
func NewArgumentsLoader(logger domain.Logger) *argumentsLoaderService {
	return &argumentsLoaderService{logger}
}

// Load reads a JSON file of arguments.
func (service *argumentsLoaderService) Load(filename string) (*domain.Arguments, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var arguments domain.Arguments

	err = json.Unmarshal(bytes, &arguments)
	if err != nil {
		return nil, err
	}

	return &arguments, nil
}
