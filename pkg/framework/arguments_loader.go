package framework

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/chrisnappin/flightchecker/pkg/domain"
)

// ArgumentsLoader handles being able to load arguments from a JSON file.
type ArgumentsLoader interface {
	Load(filename string) (*domain.Arguments, error)
}

// argumentsLoaderService handles loading arguments from a JSON file.
type argumentsLoaderService struct {
	logger domain.Logger
}

// NewArgumentsLoader creates a new instance.
func NewArgumentsLoader(logger domain.Logger) ArgumentsLoader {
	return &argumentsLoaderService{logger}
}

// Load reads a JSON file of arguments.
func (l *argumentsLoaderService) Load(filename string) (*domain.Arguments, error) {
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
