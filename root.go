package exchangerates

import (
	"time"
)

// Store is the interface that all exchange rate stores must satisfy
type Store interface {
	GetExchangeRate(from, to string, date string) (float64, error)
	GetMonthExchangeRates(from, to string, year, month int) ([]DateRate, error)
}

// DateRate represents a single exchange rate with its date
type DateRate struct {
	Date time.Time
	Rate float64
}
