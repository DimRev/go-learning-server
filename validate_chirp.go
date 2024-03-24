package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140

	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	forbiddenWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	w.Header().Add("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	chirp := parameters{}

	err := decoder.Decode(&chirp)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	if len(chirp.Body) > maxChirpLength {
		RespondWithError(w, 400, "Chirp is too long")
		return
	}

	chirpWords := strings.Split(chirp.Body, " ")
	for i, word := range chirpWords {
		if forbiddenWords[strings.ToLower(word)] {
			chirpWords[i] = "****"
		}
	}
	cleanedBody := strings.Join(chirpWords, " ")

	RespondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanedBody,
	})
}
