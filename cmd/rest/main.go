package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func sh(cmd string) {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/", handleIndex)
	go func() { log.Fatalln(http.ListenAndServe(":12345", r)) }()

	sh("httpie -p hb localhost:12345")
}

type indexResponse struct {
	Message string `json:"message"` // must be capitalized, more on this later
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j := json.NewEncoder(w)
	j.Encode(indexResponse{Message: "Hello, 世界!"})
}
