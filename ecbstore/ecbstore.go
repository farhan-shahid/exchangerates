package ecbstore

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ECBStore struct {
	records          [][]string
	currencyIndexMap map[string]int //maps curreny names to indexes in records
	dateIndexMap     map[string]int //maps dates to indexes in records
}

func (s *ECBStore) fetchData() error {
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

func (s *ECBStore) lookup(curr string, date string) (float64, error) {
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

func (s *ECBStore) GetExchangeRate(from, to string, date string) (string, error) {
	if s.records == nil {
		err := s.fetchData()
		if err != nil {
			return "", err
		}
	}

	fromValue, err := s.lookup(from, date)
	if err != nil {
		return "", err
	}

	toValue, err := s.lookup(to, date)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.5f", toValue/fromValue), nil
}
