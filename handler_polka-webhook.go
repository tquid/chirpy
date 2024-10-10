package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/tquid/chirpy/internal/auth"
)

type PolkaRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", errors.New("incorrect api key from polka webhook"))
		return
	}
	decoder := json.NewDecoder(r.Body)
	polkaRequest := PolkaRequest{}
	err = decoder.Decode(&polkaRequest)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad request", err)
		return
	}
	if polkaRequest.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	polkaUUID, err := uuid.Parse(polkaRequest.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad request", err)
		return
	}
	err = cfg.userStore.GrantChirpyRed(r.Context(), polkaUUID)
	if err != nil {
		switch err.(type) {
		case *NotFoundError:
			respondWithError(w, http.StatusNotFound, "Not found", err)
			return
		case *DatabaseError:
			respondWithError(w, http.StatusInternalServerError, "Database error", err)
			return
		default:
			respondWithError(w, http.StatusInternalServerError, "Unknown error", err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
