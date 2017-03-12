package mockstore

type MockStore struct {
}

func (s *MockStore) GetExchangeRate(from, to string, date string) (string, error) {

	return "1", nil
}
