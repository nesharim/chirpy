package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nesharim/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	secret         string
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load env file: %v", err)
		return
	}

	secret := os.Getenv("SECRET")
	if len(secret) == 0 {
		log.Println("SECRET value needs to be set in env")
		return
	}

	dbURL := os.Getenv("DB_URL")
	if len(dbURL) == 0 {
		log.Println("DB_URL value needs to be set in env")
		return
	}

	platform := os.Getenv("PLATFORM")
	if len(platform) == 0 {
		log.Println("PLATFORM value need to be set in env")
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Failed to open sql connection: %v", err)
		return
	}

	dbQueries := database.New(db)

	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
		platform:       platform,
		secret:         secret,
	}

	mux := http.NewServeMux()
	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.handlerReset)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /api/chirps", apiConfig.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.handlerGetChirpByID)
	mux.HandleFunc("POST /api/chirps", apiConfig.handlerCreateChirps)

	mux.HandleFunc("POST /api/users", apiConfig.handlerCreateUser)

	mux.HandleFunc("POST /api/login", apiConfig.handlerLogin)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}
}
