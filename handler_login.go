package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := LoginParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "json decoder: %v", err)
		return
	}
	user, err := cfg.userStore.Login(r.Context(), cfg.jwtSecret, params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect login or password", err)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}
