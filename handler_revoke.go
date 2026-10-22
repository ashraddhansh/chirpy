package main

import (
	"net/http"

	"github.com/ashraddhansh/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke (w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect token", err)
		return
	}

	err = cfg.db.RevokeToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
