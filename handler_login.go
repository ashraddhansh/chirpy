package main

import (
	"encoding/json"
	"net/http"

	"github.com/ashraddhansh/chirpy/internal/auth"
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
		respondWithError(w, http.StatusBadGateway, "Couldn't decode parameters", err)
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
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))

}
