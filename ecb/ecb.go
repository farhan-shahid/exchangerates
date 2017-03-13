package ecb

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

// Store fetches and stores historical currency exchange data from ecb.europa.eu
type Store struct {
	records          [][]string
	currencyIndexMap map[string]int //maps curreny names to indexes in records
	dateIndexMap     map[string]int //maps dates to indexes in records
}

var once sync.Once //used to fetch from ecb only once

// New returns a new instance of Store
func New() *Store {
	return &Store{}
}

// GetExchangeRate returns exchange rate from the ecb dataset.
// Use from and to for specifying the currencies to convert between
// and date to specify the date of conversion
func (s *Store) GetExchangeRate(from, to string, date string) (float64, error) {
	once.Do(func() { s.fetchData() })
	if s.records == nil {
		return 0, errors.New("Data fetching failed")
	}

	fromValue, err := s.lookup(from, date)
	if err != nil {
		return 0, err
	}

	toValue, err := s.lookup(to, date)
	if err != nil {
		return 0, err
	}

	rate := toValue / fromValue
	rateStr := fmt.Sprintf("%.5f", rate) //reduce precision
	rate, _ = strconv.ParseFloat(rateStr, 5)
	return rate, nil
}

func (s *Store) fetchData() error {
	resp, err := http.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.zip")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	reader, err := zip.NewReader(bytes.NewReader(body), resp.ContentLength)
	if err != nil {
		return err
	}

	file, err := reader.File[0].Open()
	if err != nil {
		return err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	s.records, err = csvReader.ReadAll()
	if err != nil {
		return err
	}

	s.currencyIndexMap = make(map[string]int)

	for i := 1; i < len(s.records[0]); i++ {
		s.currencyIndexMap[s.records[0][i]] = i
	}

	s.dateIndexMap = make(map[string]int)

	for i := 1; i < len(s.records); i++ {
		s.dateIndexMap[s.records[i][0]] = i
	}

	return nil
}

func (s *Store) lookup(curr string, date string) (float64, error) {
	if curr == "EUR" {
		return 1, nil
	}

	dateIndex, ok := s.dateIndexMap[date]
	if !ok {
		return 0, errors.New("date not found")
	}

	currIndex, ok := s.currencyIndexMap[curr]
	if !ok {
		return 0, errors.New("currency " + curr + " not found")
	}

	value, error := strconv.ParseFloat(s.records[dateIndex][currIndex], 64)
	if error != nil {
		return 0, errors.New(curr + " data does not exist for " + date)
	}

	return value, nil
}
