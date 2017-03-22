package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/farhan-shahid/exchangerates"
	"github.com/gorilla/mux"
)

func getRateHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("serving", req.URL)

	var store exchangerates.Store

	vars := mux.Vars(req)
	if vars["store"] == "ecb" {
		store = ec
	} else if vars["store"] == "googlefinance" {
		store = goog
	} else if vars["store"] == "mock" {
		store = moc
	} else {
		http.Error(w, vars["store"]+" is not a valid store", http.StatusBadRequest)
		return
	}

	from, to, date, err := getRateFormValues(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rate, err := store.GetExchangeRate(from, to, date)
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

func getRateFormValues(w http.ResponseWriter, req *http.Request) (from, to, date string, err error) {
	from = req.FormValue("from")
	if from == "" {
		err = errors.New(`missing "from" URL parameter`)
		return
	}

	to = req.FormValue("to")
	if to == "" {
		err = errors.New(`missing "to" URL parameter`)
		return
	}

	date = time.Now().AddDate(0, 0, -1).UTC().Format("2006-01-02") // default date is yesterday's date
	if req.FormValue("date") != "" {
		_, err = time.Parse("2006-01-02", req.FormValue("date"))
		if err == nil {
			date = req.FormValue("date")
		} else {
			err = errors.New(`incorrect date format, should be similar to 2016-03-28`)
			return
		}
	}
	return
}
