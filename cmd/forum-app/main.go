package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()

	server := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on %s", server.Addr)

	err := server.ListenAndServe()
	log.Fatal(err)
}
