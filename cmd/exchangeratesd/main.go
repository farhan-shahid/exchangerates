package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/farhan-shahid/exchangerates"
	"github.com/farhan-shahid/exchangerates/ecb"
	"github.com/farhan-shahid/exchangerates/googlefinance"
)

type rateResp struct {
	Rate float64
}

var (
	ec   exchangerates.Store = ecb.New()
	goog exchangerates.Store = googlefinance.New()
)

func main() {
	addr := flag.String("addr", "localhost:7777", "the address of the server")
	flag.Parse()

	http.HandleFunc("/googlefinance", handleGoogleFinance)
	http.HandleFunc("/ecb", handleEcb)

	log.Fatal(http.ListenAndServe(*addr, nil))
	fmt.Println("serving on " + *addr)
}

func handleEcb(w http.ResponseWriter, req *http.Request) {
	log.Println("serving", req.URL)

	from, to, date, success := getFormValues(w, req)
	if !success {
		return
	}

	rate, err := ec.GetExchangeRate(from, to, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(&rateResp{Rate: rate})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGoogleFinance(w http.ResponseWriter, req *http.Request) {
	log.Println("serving", req.URL)

	from, to, date, success := getFormValues(w, req)
	if !success {
		return
	}

	rate, err := goog.GetExchangeRate(from, to, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(&rateResp{Rate: rate})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getFormValues(w http.ResponseWriter, req *http.Request) (from, to, date string, success bool) {
	from = req.FormValue("from")
	if from == "" {
		http.Error(w, `missing "from" URL parameter`, http.StatusBadRequest)
		return "", "", "", false
	}

	to = req.FormValue("to")
	if to == "" {
		http.Error(w, `missing "to" URL parameter`, http.StatusBadRequest)
		return "", "", "", false
	}

	date = time.Now().AddDate(0, 0, -1).UTC().Format("2006-01-02") // default date is yesterday's date
	if req.FormValue("date") != "" {
		if match, _ := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", req.FormValue("date")); match {
			date = req.FormValue("date")
		} else {
			http.Error(w, `incorrect date format, should be similar to 2016-03-28`, http.StatusBadRequest)
			return "", "", "", false
		}
	}

	success = true
	return
}
