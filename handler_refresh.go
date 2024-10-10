package main

import (
	"net/http"
	"time"

	"github.com/tquid/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	userID, err := cfg.userStore.GetUserIDFromValidRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	jwtStr, err := auth.MakeJWT(userID, cfg.jwtSecret, jwtExpireSeconds*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Response{
		Token: jwtStr,
	})
}
