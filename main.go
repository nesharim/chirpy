package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const (
	port    = 8080
	rootDir = "."
)

type ApiConfig struct {
	fileServerHits atomic.Int64
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	apiCfg := ApiConfig{fileServerHits: atomic.Int64{}}
	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(rootDir)))))
	mux.HandleFunc("/healthz", handleReadiness)
	mux.HandleFunc("/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("/reset", apiCfg.handleReset)

	addr := fmt.Sprintf(":%d", port)

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func (cfg *ApiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	text := fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load())
	w.Write([]byte(text))
}

func (cfg *ApiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
