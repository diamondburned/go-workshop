package main

import (
	"encoding/json"
	"log"
	"net/http"

	"libdb.so/go-workshop/lib/httprint"
)

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/", handleIndex)
	go func() { log.Fatalln(http.ListenAndServe(":12345", r)) }()

	httprint.GET("localhost:12345")
}

type indexResponse struct {
	Message string `json:"message"` // must be capitalized, more on this later
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j := json.NewEncoder(w)
	j.Encode(indexResponse{Message: "Hello, 世界!"})
}
