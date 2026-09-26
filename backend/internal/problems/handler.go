package problems

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Problem struct {
	ID         int64  `json:"id"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
}

type ProblemDetail struct {
	ID          int64  `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	var problem ProblemDetail

	err := h.db.QueryRow(
		r.Context(),
		`
			SELECT id, slug, title, description, difficulty
			FROM problems
			WHERE slug = $1
		`,
		slug,
	).Scan(
		&problem.ID,
		&problem.Slug,
		&problem.Title,
		&problem.Description,
		&problem.Difficulty,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(
			w,
			`{"error":"problem not found"}`,
			http.StatusNotFound,
		)
		return
	}

	if err != nil {
		log.Printf("failed to query problem %q: %v", slug, err)
		http.Error(
			w,
			`{"error":"internal server error"}`,
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(problem); err != nil {
		log.Printf("failed to encode problem: %v", err)
	}
}
