package problems

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Problem struct {
	ID         int64  `json:"id"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
}

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(
		r.Context(),
		`
			SELECT id, slug, title, difficulty
			FROM problems
			ORDER BY id
		`,
	)
	if err != nil {
		log.Printf("failed to query problems: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	problems := make([]Problem, 0)

	for rows.Next() {
		var problem Problem

		if err := rows.Scan(
			&problem.ID,
			&problem.Slug,
			&problem.Title,
			&problem.Difficulty,
		); err != nil {
			log.Printf("failed to scan problem: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		problems = append(problems, problem)
	}

	if err := rows.Err(); err != nil {
		log.Printf("failed while reading problems: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(problems); err != nil {
		log.Printf("failed to encode problems: %v", err)
	}
}
