package main

import (
	"encoding/json"
	"net/http"

	"github.com/ashraddhansh/chirpy/internal/auth"
	"github.com/ashraddhansh/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsers (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password	string `json:"password"`
		Email	string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash the password", err)
	}


	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))


}
