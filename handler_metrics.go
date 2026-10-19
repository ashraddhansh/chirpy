package main

import (
	"net/http"
	"fmt"
)

func (cfg *apiConfig) handlerMetrics (w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", cfg.fileServerHits.Load())

}
