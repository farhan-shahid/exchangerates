package exchangerates

type Store interface {
	GetExchangeRate(to string, from string, date string) (string, error)
}
