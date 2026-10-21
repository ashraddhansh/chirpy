package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"
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

func handlerValidateChirp (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body	string `json:"body"`
	}
	
	type returnVals struct {
		CleanedBody	string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	chirpLength := utf8.RuneCount([]byte(params.Body))
	

	if chirpLength > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: replaceProfaneWords(params.Body)})
}
