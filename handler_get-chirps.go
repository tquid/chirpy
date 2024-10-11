package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	// Check if we have a parameter to limit chirps to one author
	authorId := r.URL.Query().Get("author_id")
	sortParam := strings.ToLower(GetSortOrder(r.URL.Query().Get("sort")))
	var sort SortDirection
	switch sortParam {
	case "desc":
		sort = SortDesc
	case "asc", "":
		sort = SortAsc
	default:
		errStr := fmt.Sprintf("invalid sort parameter '%s'", sortParam)
		respondWithError(w, http.StatusBadRequest, errStr, errors.New(errStr))
	}
	if authorId != "" {
		id, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author ID", err)
			return
		}
		chirps, err := cfg.chirpStore.GetChirpsByUserId(r.Context(), id, sort)
		if err != nil {
			switch err.(type) {
			case *NotFoundError:
				respondWithError(w, http.StatusNotFound, "author not found", err)
				return
			case *DatabaseError:
				respondWithError(w, http.StatusInternalServerError, "database error", err)
				return
			default:
				respondWithError(w, http.StatusInternalServerError, "unknown error", err)
				return
			}
		}
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}
	chirps, err := cfg.chirpStore.GetChirps(r.Context(), sort)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to get chirps", err)
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
