package main

import (
	"flag"
	"fmt"

	"github.com/farhan-shahid/exchangerates"
	"github.com/farhan-shahid/exchangerates/ecbstore"
	"github.com/farhan-shahid/exchangerates/mockstore"
)

func main() {
	storename := flag.String("store", "ecbstore", "the store to be used")
	from := flag.String("from", "EUR", "the currency to convert from")
	to := flag.String("to", "USD", "the currency to convert to")
	date := flag.String("date", "2017-03-02", "the date for which to get exchange rate")
	flag.Parse()

	var s exchangerates.Store
	if *storename == "ecbstore" {
		s = &ecbstore.ECBStore{}
	} else if *storename == "mockstore" {
		s = &mockstore.MockStore{}
	} else {
		fmt.Println("Invalid store")
		return
	}

	val, err := s.GetExchangeRate(*from, *to, *date)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(val)
}
