package main

import (
	"sync/atomic"

	"github.com/ashraddhansh/chirpy/internal/database"
)


type apiConfig struct {
	fileServerHits atomic.Int32
	db *database.Queries
}
