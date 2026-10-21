package main

import "net/http"

func (cfg *apiConfig) handlerReset (w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "You are not in dev environment", nil)
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error deleting users", err)
		return
	}

	type returnVals struct {
		Message	string `json:"message"`
	}

	respondWithJSON(w, http.StatusOK, returnVals{Message: "Users deleted successufully"})

}
