package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 1 * time.Second}

func health(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func callEcho(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")

	url := fmt.Sprintf("http://127.0.0.1:8080/echo?msg=%s", msg)
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("service=B endpoint=/call-echo status=error error=%q latency_ms=%d", err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service_b": "ok",
			"service_a": "unavailable",
			"error":     err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("service=B outbound_call=serviceA path=/echo bad_status=%d", resp.StatusCode)

		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service_b": "ok",
			"service_a": "bad_status",
			"status":    resp.StatusCode,
		})
		return
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("service=B outbound_call=serviceA path=/echo decode_error=%q", err.Error())

		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service_b": "ok",
			"service_a": "invalid_response",
			"error":     err.Error(),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"service_b": "ok",
		"service_a": data,
	})
}

func main() {
	http.HandleFunc("/health", withLogging("B", health))
	http.HandleFunc("/call-echo", withLogging("B", callEcho))

	log.Println("service=B listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
