package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/farhan-shahid/exchangerates"
	"github.com/farhan-shahid/exchangerates/ecb"
	"github.com/farhan-shahid/exchangerates/mock"
)

// Make sure Store interface is being satisfied
var (
	_ exchangerates.Store = (*ecb.Store)(nil)
	_ exchangerates.Store = (*mock.Store)(nil)
	//_ exchangerates.Store = (*xe.Store)(nil)
)

func main() {
	var (
		storename = flag.String("store", "xe", "the store to be used")
		from      = flag.String("from", "EUR", "the currency to convert from")
		to        = flag.String("to", "USD", "the currency to convert to")
		date      = flag.String("date", "2017-03-02", "the date for which to get exchange rate")
	)
	flag.Parse()

	var s exchangerates.Store
	if *storename == "ecb" {
		s = ecb.New()
		/*} else if *storename == "xe" {
		s = xe.New()*/
	} else {
		log.Fatal("Invalid store")
	}

	val, err := s.GetExchangeRate(*from, *to, *date)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
}
