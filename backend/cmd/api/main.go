package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rock-em/rock-em/backend/internal/problems"
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("connected to PostgreSQL")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := db.Ping(r.Context()); err != nil {
			http.Error(w, `{"status":"unhealthy"}`, http.StatusServiceUnavailable)
			return
		}

		if err := json.NewEncoder(w).Encode(healthResponse{
			Status: "ok",
		}); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	})

	problemHandler := problems.NewHandler(db)

	mux.HandleFunc("GET /api/problems", problemHandler.List)

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	log.Printf("Rock Em API listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
