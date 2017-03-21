package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/farhan-shahid/exchangerates/server"
)

func main() {
	addr := flag.String("addr", "localhost:7777", "the address of the server")
	flag.Parse()

	s := server.New()
	srv := &http.Server{
		Addr:    *addr,
		Handler: s,
	}
	log.Printf("Serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
