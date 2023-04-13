package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func sh(cmd string) {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Run()
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/echo", handleEcho)
	go func() { log.Fatalln(http.ListenAndServe(":12345", r)) }()

	// POST to the endpoint with a null body.
	sh(`httpie -p hb POST http://localhost:12345/echo <<< null`)
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
