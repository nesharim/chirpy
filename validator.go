package main

import (
	"strings"
)

func isValidLength(msg string) bool {
	const maxChirpLength = 140

	return len(msg) < maxChirpLength
}

// TODO: should this be here. Would be best to move it to a utils file.
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
