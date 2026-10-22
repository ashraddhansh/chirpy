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
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "empty email or password", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash the password", err)
		return
	}


	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))


}
