package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	port    = 8080
	rootDir = "."
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(rootDir))))
	mux.HandleFunc("/healthz", handleReadiness)

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

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
