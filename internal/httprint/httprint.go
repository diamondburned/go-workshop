package httprint

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

// GET does an HTTP GET request and prints the response.
func GET(url string) {
	Do("GET", url, "")
}

// POST does an HTTP POST request and prints the response.
func POST(url string, body string) {
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
