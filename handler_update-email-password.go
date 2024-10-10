package main

import (
	"encoding/json"
	"net/http"

	"github.com/tquid/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpdateEmailPassword(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := LoginParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	user, err := cfg.userStore.UpdateLoginCredentials(r.Context(), userID, params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}
