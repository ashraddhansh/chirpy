package main

import "github.com/ashraddhansh/chirpy/internal/database"

func databaseUserToUser(dbUser database.User) User {
    return User{
        ID:        dbUser.ID,
        CreatedAt: dbUser.CreatedAt,
        UpdatedAt: dbUser.UpdatedAt,
        Email:     dbUser.Email,
    }
}


func databaseChirpToChirp(dbUser database.Chirp) Chirp {
    return Chirp{
        ID:        dbUser.ID,
        CreatedAt: dbUser.CreatedAt,
        UpdatedAt: dbUser.UpdatedAt,
		Body: dbUser.Body,
		UserID: dbUser.UserID,

    }
}
