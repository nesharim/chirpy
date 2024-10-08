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
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)

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
	text := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits.Load())
	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

func (cfg *ApiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileServerHits.Store(0)
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
