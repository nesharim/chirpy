package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type userEmail struct {
		Email string `json:"email"`
	}

	usrEml := userEmail{}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&usrEml); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode req body for user email", err)
		return
	}

	newUser, err := cfg.dbQueries.CreateUser(req.Context(), usrEml.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create new user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	})
}
