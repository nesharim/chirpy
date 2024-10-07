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
	mux.Handle("/", http.FileServer(http.Dir(rootDir)))

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
