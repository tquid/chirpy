package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0\n"))
	rows, err := cfg.userStore.DeleteAllUsers(r.Context())
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Unable to delete all users: %v\n", err)))
		return
	}
	w.Write([]byte(fmt.Sprintf("Deleted %d users\n", rows)))
}
