package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
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
	sh("sleep 0.1 && curl -s -D - localhost:12345")
}

type indexResponse struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j := json.NewEncoder(w)
	j.SetIndent("", "  ") // make it pretty
	j.Encode(indexResponse{
		Time:    time.Now(),
		Message: "Hello, 世界!",
	})
}
