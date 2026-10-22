package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/ashraddhansh/chirpy/internal/auth"
	"github.com/ashraddhansh/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirpsByID (w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "access token malformed or missing", err)
		return
	}

	userUUID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "access token not valid", err)
		return
	}

	chirpID := r.PathValue("chirpID")
	parsedChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), parsedChirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp", err)
		return
	}

	if chirp.UserID != userUUID {
		respondWithError(w, http.StatusForbidden, "unauthorized user trying to delete chirp", nil)
		return

	}

	rowsAffected, err := cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID: parsedChirpID,
		UserID: userUUID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp", err)
		return
	}
	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "chirp not found", nil)
		return

	}
	w.WriteHeader(http.StatusNoContent)

}
