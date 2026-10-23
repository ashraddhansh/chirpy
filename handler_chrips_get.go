package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerGetChirps (w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("author_id")

	sortMethod := r.URL.Query().Get("sort")
	if sortMethod != "asc" && sortMethod != "desc" && sortMethod != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid sort query parameter", nil)
		return
	}

	if id != "" {
		authorID, err := uuid.Parse(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id", err)
			return
		}
		dbChirpsUser, err := cfg.db.GetChirpsByUserID(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps of author from database", err)
			return
		}
		chirps := make([]Chirp, 0, len(dbChirpsUser))
		for _, dbChirp := range dbChirpsUser {
			chirps = append(chirps, databaseChirpToChirp(dbChirp))
		}
		if sortMethod == "desc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			})
			respondWithJSON(w, http.StatusOK, chirps)
			return

		}
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}
	
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps from database", err)
		return
	}

	chirps := make([]Chirp, 0, len(dbChirps))
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, databaseChirpToChirp(dbChirp))
	}

	if sortMethod == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
	})
}
	respondWithJSON(w, http.StatusOK, chirps)
}
