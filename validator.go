package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	type bodyParams struct {
		Body string `json:"body"`
	}
	params := bodyParams{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decoding params failed", err)
		return
	}

	const maxChirpLength = 140

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	type successResponse struct {
		Valid bool `json:"valid"`
	}

	respondWithJSON(w, http.StatusOK, successResponse{
		Valid: true,
	})
}
