package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID uuid.UUID `json:"user_id"`
		Body   string    `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)
	params.Body = cleaned

	chirp, err := cfg.chirpStore.CreateChirp(r.Context(), params.Body, params.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create chirp", err)
	}
	respondWithJSON(w, 201, chirp)
}
