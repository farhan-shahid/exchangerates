package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/farhan-shahid/exchangerates"
	"github.com/farhan-shahid/exchangerates/chart"
	"github.com/farhan-shahid/exchangerates/ecb"
	"github.com/farhan-shahid/exchangerates/ecbsql"
	"github.com/farhan-shahid/exchangerates/googlefinance"
	"github.com/farhan-shahid/exchangerates/mock"
)

// Make sure Store interface is being satisfied
var (
	_ exchangerates.Store = (*ecb.Store)(nil)
	_ exchangerates.Store = (*mock.Store)(nil)
	_ exchangerates.Store = (*googlefinance.Store)(nil)
	_ exchangerates.Store = (*ecbsql.Store)(nil)
)

func main() {
	var (
		storename = flag.String("store", "ecbsql", "the store to be used")
		from      = flag.String("from", "EUR", "the currency to convert from")
		to        = flag.String("to", "USD", "the currency to convert to")
		date      = flag.String("date", "2017-03-02", "the date for which to get exchange rate")
		getchart  = flag.Bool("getchart", false, "set to true to get exchange rate chart")
		month     = flag.Int("month", 1, "the month for which to get exchange rate chart")
		year      = flag.Int("year", 2017, "the year for which to get exchange rate chart")
	)
	flag.Parse()

	var s exchangerates.Store
	if *storename == "ecbsql" {
		s = ecbsql.New()
	} else if *storename == "ecb" {
		s = ecb.New()
	} else if *storename == "googlefinance" {
		s = googlefinance.New()
	} else {
		log.Fatal("Invalid store")
	}

	if !*getchart {
		val, err := s.GetExchangeRate(*from, *to, *date)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(val)
	} else {
		rates, err := s.GetMonthExchangeRates(*from, *to, *year, *month)
		if err != nil {
			log.Fatal(err)
		}
		filename := "chart.png"
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		err = chart.MakeRateChart(*from, *to, rates, file)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(filename + " has been saved to current directory")
	}
}
