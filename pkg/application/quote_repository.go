package application

import (
	"github.com/chrisnappin/flightchecker/pkg/domain"
	"github.com/chrisnappin/flightchecker/pkg/framework"
)

// QuoteRepository handles being able to load and save quote data
type QuoteRepository interface {
	Initialise() error
}

// quoteRepositoryService handles loading and saving data.
type quoteRepositoryService struct {
	logger           domain.Logger
	sqliteRepository framework.SQLiteRepository
}

// NewQuoteRepository creates a new instance.
func NewQuoteRepository(logger domain.Logger, sqliteRepository framework.SQLiteRepository) QuoteRepository {
	return &quoteRepositoryService{logger, sqliteRepository}
}

// Initialise creates an empty database
func (r *quoteRepositoryService) Initialise() error {
	return r.sqliteRepository.Initialise()
}
