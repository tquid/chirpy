package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/tquid/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid bearer token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed parsing uuid", err)
		return
	}
	chirp, err := cfg.chirpStore.GetChirpById(r.Context(), id)
	if err != nil {
		switch e := err.(type) {
		case *NotFoundError:
			respondWithError(w, http.StatusNotFound, "Not found", e)
			return
		case *DatabaseError:
			respondWithError(w, http.StatusInternalServerError, "Database operation failed", err)
			return
		default:
			respondWithError(w, http.StatusInternalServerError, "Unknown error", err)
			return
		}
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Forbidden", fmt.Errorf("attempt to delete chirp by other than author"))
		return
	}
	err = cfg.chirpStore.DeleteChirp(r.Context(), id)
	if err != nil {
		switch e := err.(type) {
		case *NotFoundError:
			respondWithError(w, http.StatusNotFound, "Not found", e)
			return
		case *DatabaseError:
			respondWithError(w, http.StatusInternalServerError, "Database operation failed", err)
			return
		default:
			respondWithError(w, http.StatusInternalServerError, "Unknown error", err)
			return
		}
	}
	respondWithJSON(w, http.StatusNoContent, *chirp)
}
