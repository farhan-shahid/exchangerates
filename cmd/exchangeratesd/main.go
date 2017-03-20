package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/farhan-shahid/exchangerates/server"
)

func main() {
	addr := flag.String("addr", "localhost:7777", "the address of the server")
	flag.Parse()

	r := server.GetRouter()

	//http.Handle("/", r)
	//log.Fatal(http.ListenAndServe(*addr, nil))
	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	fmt.Println("serving on " + *addr)
	log.Fatal(srv.ListenAndServe())
}
