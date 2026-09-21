package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(healthResponse{
		Status: "ok",
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", healthHandler)

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	log.Printf("Rock Em API listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
