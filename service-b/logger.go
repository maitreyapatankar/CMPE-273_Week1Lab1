package main

import (
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func withLogging(service string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK} // default 200
		h(sw, r)
		log.Printf("service=%s method=%s path=%s status=%d latency_ms=%d",
			service, r.Method, r.URL.Path, sw.status, time.Since(start).Milliseconds())
	}
}
