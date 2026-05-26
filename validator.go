package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func isValidLength(msg string) bool {
	const maxChirpLength = 140

	return len(msg) > maxChirpLength
}

func cleanMessage(msg string) string {
	words := strings.Split(msg, " ")
	for i, word := range words {
		lWord := strings.ToLower(word)
		if lWord == "kerfuffle" || lWord == "sharbert" || lWord == "fornax" {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

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

	if isValidLength(params.Body) {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedMsg := cleanMessage(params.Body)

	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respondWithJSON(w, http.StatusOK, successResponse{
		CleanedBody: cleanedMsg,
	})
}
