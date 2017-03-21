package server

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
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

// Server type manages routes for accessing exchange rates over http
type Server struct {
	router *mux.Router
}

// New returns a *Server with the necessary routing handler(s) attached
func New() *Server {
	r := mux.NewRouter()
	r.HandleFunc("/{store}", handleStoreReq)
	return &Server{router: r}
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
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
	}

	from, to, date, success := getFormValues(w, req)
	if !success {
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
