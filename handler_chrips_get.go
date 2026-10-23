package main

import (
	"net/http"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerGetChirps (w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("author_id")

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

	respondWithJSON(w, http.StatusOK, chirps)
}
