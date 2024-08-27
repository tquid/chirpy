package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func cleaner(profane string) string {
	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(profane, " ")
	clean := []string{}
	for _, word := range words {
		nasty := false
		if _, ok := profanities[strings.ToLower(word)]; ok {
			nasty = true
		}
		if nasty {
			clean = append(clean, "****")
		} else {
			clean = append(clean, word)
		}
	}
	return strings.Join(clean, " ")
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
	}

	type ValidateResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, 400, "bad request")
		return
	}
	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	vr := ValidateResponse{
		CleanedBody: cleaner(chirp.Body),
	}
	log.Printf("Responding with %v", vr)
	respondWithJSON(w, 200, vr)
}
