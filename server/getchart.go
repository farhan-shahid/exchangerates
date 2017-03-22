package server

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/farhan-shahid/exchangerates"
)

func getChartHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("serving", req.URL)

	from, to, month, year, err := getChartFormValues(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rates, err := ec.GetMonthExchangeRates(from, to, year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exchangerates.MakeRateChart(from, to, rates, w)
}

func getChartFormValues(w http.ResponseWriter, req *http.Request) (from, to string, month, year int, err error) {
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

	if req.FormValue("month") != "" {
		month, err = strconv.Atoi(req.FormValue("month"))
		if err != nil {
			return
		}
	} else {
		err = errors.New(`missing "month" URL parameter`)
		return
	}

	if req.FormValue("year") != "" {
		year, err = strconv.Atoi(req.FormValue("year"))
		if err != nil {
			return
		}
	} else {
		err = errors.New(`missing "year" URL parameter`)
		return
	}
	return
}
