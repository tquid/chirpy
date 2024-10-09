package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	err := uuid.Validate(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid uuid", err)
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed parsing uuid", err)
		return
	}
	chirp, err := cfg.chirpStore.GetChirpById(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "not found", err)
		return
	}
	respondWithJSON(w, http.StatusOK, *chirp)
}
