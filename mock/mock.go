package mock

import (
	"errors"

	"github.com/farhan-shahid/exchangerates"
)

// Store allows mocking of exchange rate stores for testing
type Store struct {
	OnGetExchangeRate       func(from, to string, date string) (float64, error)
	OnGetMonthExchangeRates func(from, to string, year, month int) ([]exchangerates.DateRate, error)
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

// GetMonthExchangeRates just calls the OnGetMonthExchangeRates function that is specified by the calling context
func (s *Store) GetMonthExchangeRates(from, to string, year, month int) ([]exchangerates.DateRate, error) {
	if s.OnGetMonthExchangeRates == nil {
		return nil, errors.New("OnGetMonthExchangeRates not set")
	}
	return s.OnGetMonthExchangeRates(from, to, year, month)
}
