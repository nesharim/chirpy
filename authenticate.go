package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nesharim/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type paramaters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	}

	type response struct {
		User
		Token         string `json:"token"`
		RefereshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(req.Body)
	params := paramaters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode login creds", err)
		return
	}

	dbUser, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	isValidPassword, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil || !isValidPassword {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expirationTime := time.Hour

	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate jwt token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
		Token: accessToken,
	})
}
