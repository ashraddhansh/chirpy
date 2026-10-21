package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ashraddhansh/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Print(err)
	}

	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		db: dbQueries,
	}

	const port = "8080"
	const filepathRoot = "."
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)


	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n",filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
