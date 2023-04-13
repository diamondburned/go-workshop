package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"libdb.so/acm-go-present/cmd/bookstore/bookstore/api"
	"libdb.so/acm-go-present/cmd/bookstore/bookstore/db"
)

var addr = ":8080"

func main() {
	flag.StringVar(&addr, "addr", addr, "address to listen on")
	flag.Parse()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join(os.TempDir(), "bookstore.db")
	}

	db, err := db.NewSQLite(dbPath)
	if err != nil {
		log.Fatalln("creating database:", err)
	}

	// Create a new Bookstore API handler. This implicitly implements
	// http.Handler, so we can pass it to http.ListenAndServe.
	//
	// API, NOT Api!!! Go style is to use all caps for acronyms.
	// See Effective Go (https://go.dev/doc/effective_go#mixed-caps).
	bookAPI := api.NewBookstoreHandler(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger) // log all requests
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Version our API! What if we wanna change it?!
	r.Mount("/api/v0", bookAPI)

	// Start the server.
	log.Println("listening on", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalln(err)
	}
}
