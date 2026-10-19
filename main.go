package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlerReadiness)
	
	var apiCfg apiConfig

	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)

	mux.Handle("/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix(
				"/app",
				http.FileServer(http.Dir(filepathRoot)),
			),
		),
	)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n",filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
