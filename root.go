package exchangerates

// Store is the interface that all exchange rate stores must satisfy
type Store interface {
	GetExchangeRate(to string, from string, date string) (float64, error)
}
