package server

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/farhan-shahid/exchangerates"
	"github.com/farhan-shahid/exchangerates/ecb"
	"github.com/farhan-shahid/exchangerates/googlefinance"
	"github.com/farhan-shahid/exchangerates/mock"
	"github.com/gorilla/mux"
)

type rateResp struct {
	Rate float64
}

var (
	ec   exchangerates.Store = ecb.New()
	goog exchangerates.Store = googlefinance.New()
	moc                      = mock.New()
)

type loggingHandler struct {
	w io.Writer
	h http.Handler
}

func (l *loggingHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	l.w.Write([]byte(r.URL.String()))
	l.h.ServeHTTP(rw, r)
}

// LoggingHandler returns a http.Handler that logs requests before calling given handler's ServeHTTP method
func LoggingHandler(w io.Writer, h http.Handler) http.Handler {
	return &loggingHandler{w: w, h: h}
}

// Server type manages routes for accessing exchange rates over http
type Server struct {
	h http.Handler
}

// New returns a *Server with the necessary routing handler(s) attached
func New() *Server {
	r := mux.NewRouter()
	r.HandleFunc("/{store}", handleStoreReq)
	return &Server{h: r}
}

// AddLogging adds request logging using LoggingHandler
func (s *Server) AddLogging(w io.Writer) {
	s.h = LoggingHandler(w, s.h)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.h.ServeHTTP(w, r)
}

func handleStoreReq(w http.ResponseWriter, req *http.Request) {
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

	from, to, date, chart, month, year, err := getFormValues(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !chart {
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
	} else {
		rates, err := store.GetMonthExchangeRates(from, to, year, month)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		exchangerates.MakeRateChart(from, to, rates, w)
	}
}

func getFormValues(w http.ResponseWriter, req *http.Request) (from, to, date string, chart bool, month, year int, err error) {
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
		if match, _ := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", req.FormValue("date")); match {
			date = req.FormValue("date")
		} else {
			err = errors.New(`incorrect date format, should be similar to 2016-03-28`)
			return
		}
	}

	if req.FormValue("chart") != "" {
		chart = true
	}

	if req.FormValue("month") != "" {
		month, err = strconv.Atoi(req.FormValue("month"))
		if err != nil {
			return
		}
	}

	if req.FormValue("year") != "" {
		year, err = strconv.Atoi(req.FormValue("year"))
		if err != nil {
			return
		}
	}

	if chart && (month == 0 || year == 0) {
		err = errors.New(`month and year are required for chart`)
		return
	}
	return
}
