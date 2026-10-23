package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ashraddhansh/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolka (w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "invalid API key", err)
	}

	type PolkaWebhookData struct {
    UserID uuid.UUID `json:"user_id"`
	}

	type parameters struct {
    Event string          `json:"event"`
    Data  PolkaWebhookData `json:"data"`
	}


	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	user, err := cfg.db.UpgradeUserToRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user", err)
		return
	}
	_ = user
	w.WriteHeader(http.StatusNoContent)


}
