package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/ashraddhansh/chirpy/internal/database"
)

func replaceProfaneWords(body string) string {
    profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
    words := strings.Split(body, " ")

    for i, word := range words {
        for _, profane := range profaneWords {
            if strings.EqualFold(word, profane) {
                words[i] = "****"
                break
            }
        }
    }

    return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerChirps (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body	string `json:"body"`
		UserID	uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	chirpLength := utf8.RuneCountInString(params.Body)
	

	if chirpLength > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	chirp, err := cfg.db.CreateChirp(r.Context(),database.CreateChirpParams{
		Body: replaceProfaneWords(params.Body),
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseChirpToChirp(chirp))
}
