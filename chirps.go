package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nesharim/chirpy/internal/database"
)

type chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, req *http.Request) {
	type userPayload struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	payload := userPayload{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&payload); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode chirp", err)
		return
	}

	if !isValidLength(payload.Body) {
		respondWithError(w, http.StatusBadRequest, "Chirp exceeds max length", nil)
		return
	}

	cleanedChirp := cleanMessage(payload.Body)

	newChirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: payload.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to insert chirp in table", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.dbQueries.ListChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps", err)
		return
	}

	var chirpsList []chirp

	for _, item := range chirps {
		chirpsList = append(chirpsList, chirp{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Body:      item.Body,
			UserID:    item.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsList)
}
