package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

// PrintGET does an HTTP PrintGET request and prints the response.
func PrintGET(url string) {
	Do("GET", url, "")
}

// Do does an HTTP request and prints the response.
func Do(method, url string, body string) {
	if !strings.Contains(url, "://") {
		url = "http://" + url
	}

	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url, r)
	if err != nil {
		log.Panicln("Error creating request:", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicln("Error sending request:", err)
	}
	defer resp.Body.Close()

	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Panicln("Error dumping response:", err)
	}

	fmt.Println(string(b))
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/", handleIndex)
	go func() { log.Fatalln(http.ListenAndServe(":12345", r)) }()

	PrintGET("localhost:12345")
}

type indexResponse struct {
	Message string `json:"message"` // must be capitalized, more on this later
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j := json.NewEncoder(w)
	j.Encode(indexResponse{Message: "Hello, 世界!"})
}
