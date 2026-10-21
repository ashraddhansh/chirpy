package main

import (
	"net/http"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirpsByID (w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
	}
	chirp, err := cfg.db.GetChirp(r.Context(), id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}


	respondWithJSON(w, http.StatusOK, databaseChirpToChirp(chirp))
}
