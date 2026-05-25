package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getHitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, "Hits: %v", cfg.fileserverHits.Load()); err != nil {
		log.Fatalf("Failed to get total hits: %v", err)
	}
}

func (cfg *apiConfig) resetAPIConfigHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
		log.Printf("Failed to write to body: %v", err)
	}
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apicfg := &apiConfig{}
	mux := http.NewServeMux()

	mux.Handle("/app/", apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /healthz", healthCheckHandler)
	mux.HandleFunc("GET /metrics", apicfg.getHitsHandler)
	mux.HandleFunc("POST /reset", apicfg.resetAPIConfigHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}
}
