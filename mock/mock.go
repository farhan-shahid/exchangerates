package mock

import (
	"errors"
)

// Store allows mocking of exchange rate stores for testing
type Store struct {
	OnGetExchangeRate func(from, to string, date string) (float64, error)
}

// New returns a new instance of Store
func New() *Store {
	return &Store{}
}

// GetExchangeRate just calls the OnGetExchangeRate function that is specified by the calling context
func (s *Store) GetExchangeRate(from, to string, date string) (float64, error) {
	if s.OnGetExchangeRate == nil {
		return 0, errors.New("OnGetExchangeRate not set")
	}
	return s.OnGetExchangeRate(from, to, date)
}
