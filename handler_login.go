package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ashraddhansh/chirpy/internal/auth"
	"github.com/ashraddhansh/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password	string `json:"password"`
		Email	string `json:"email"`
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
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	refreshExpiresAt := time.Now().UTC().Add(60 * 24 * time.Hour)

	dbRefreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: refreshExpiresAt,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't store token", err)
		return
	}

	type loginResponse struct {
    ID             uuid.UUID `json:"id"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    Email          string    `json:"email"`
	IsChirpyRed    bool       `json:"is_chirpy_red"`
    AccessToken    string    `json:"token"`
    RefreshToken   string    `json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
		AccessToken: accessToken,
		RefreshToken: dbRefreshToken.Token,
	})

}
