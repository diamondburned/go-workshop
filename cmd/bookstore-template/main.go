package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var addr = ":8080"

func main() {
	flag.StringVar(&addr, "addr", addr, "address to listen on")
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.Logger) // log all requests
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// TODO: implement API handlers
	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", http.StatusNotImplemented)
	})

	// Start the server.
	log.Println("listening on", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalln(err)
	}
}
