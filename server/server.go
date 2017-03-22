package server

import (
	"io"
	"net/http"

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
	r.HandleFunc("/chart", getChartHandler)
	r.HandleFunc("/{store}", getRateHandler)
	return &Server{h: r}
}

// AddLogging adds request logging using LoggingHandler
func (s *Server) AddLogging(w io.Writer) {
	s.h = LoggingHandler(w, s.h)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.h.ServeHTTP(w, r)
}
