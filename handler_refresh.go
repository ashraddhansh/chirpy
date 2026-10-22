package main

import (
	"net/http"
	"time"

	"github.com/ashraddhansh/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh (w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect token", err)
		return
	}
	
	refreshToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect token", err)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now().UTC()) || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token expired or revoked", nil)
		return
	}

	jwtToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating jwt access token", err)
		return
	}

	type returnToken struct {
    AccessToken    string    `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, returnToken{
		AccessToken: jwtToken,
	})

}
