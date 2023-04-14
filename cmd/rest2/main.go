package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

// PrintPOST does an HTTP POST request and prints the response.
func PrintPOST(url string, body string) {
	Do("POST", url, body)
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
	r.HandleFunc("/echo", handleEcho)
	go func() { log.Fatalln(http.ListenAndServe(":12345", r)) }()

	// POST to the endpoint with a null body.
	PrintPOST("localhost:12345/echo", "null")
}

type echoRequest struct {
	Message string `json:"message"`
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	var req echoRequest
	// Decode the request but tee it to stdout as well.
	body := io.TeeReader(r.Body, os.Stderr)
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Go says %q\n", req.Message)
}
