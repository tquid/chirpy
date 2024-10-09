package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.chirpStore.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to get chirps", err)
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
