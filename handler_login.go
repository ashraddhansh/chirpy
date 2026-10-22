package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ashraddhansh/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password	string `json:"password"`
		Email	string `json:"email"`
		ExpiresInSeconds	int `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.db.GetUser(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	
	correct, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	if !correct {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	const defaultExpirationTime = 3600

	if params.ExpiresInSeconds < 0 {
    	respondWithError(w, http.StatusBadRequest, "Invalid expiration time", nil)
    	return
	}

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > defaultExpirationTime {
		params.ExpiresInSeconds = defaultExpirationTime

	}

	expiresIn := time.Duration(params.ExpiresInSeconds) * time.Second

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create token", err)
		return
	}

	type loginResponse struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email     string    `json:"email"`
    Token     string    `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	})

}
