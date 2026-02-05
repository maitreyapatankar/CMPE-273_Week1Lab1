package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	msg := r.URL.Query().Get("msg")
	_ = json.NewEncoder(w).Encode(map[string]string{"echo": msg})
}

func main() {
	http.HandleFunc("/health", withLogging("A", health))
	http.HandleFunc("/echo", withLogging("A", echo))

	log.Println("service=A listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
