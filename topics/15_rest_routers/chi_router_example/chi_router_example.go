// chi_router_example.go
// Demonstrates routing with github.com/go-chi/chi/v5 — the modern stdlib-compatible router.
// chi is the current industry standard for Go APIs (Gorilla Mux is archived).
//
// To run:
//   go mod init example && go get github.com/go-chi/chi/v5 && go run chi_router_example.go
package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ---- domain types ----

type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// ---- in-memory store ----

var tasks = map[string]Task{
	"1": {ID: "1", Title: "Buy groceries"},
	"2": {ID: "2", Title: "Write Go notes"},
}

// ---- handlers ----

func listTasks(w http.ResponseWriter, r *http.Request) {
	out := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, t)
	}
	respondJSON(w, http.StatusOK, out)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // chi URL parameter extraction
	t, ok := tasks[id]
	if !ok {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}
	respondJSON(w, http.StatusOK, t)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil || req.Title == "" {
		respondError(w, http.StatusBadRequest, "title is required")
		return
	}
	id := "3" // static for demo; use uuid in real code
	t := Task{ID: id, Title: req.Title}
	tasks[id] = t
	respondJSON(w, http.StatusCreated, t)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := tasks[id]; !ok {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusNoContent)
}

// ---- response helpers ----

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

// ---- router setup ----

func newRouter(logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	// chi built-in middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)    // structured request log
	r.Use(middleware.Recoverer) // recover from panics

	// RESTful resource grouping
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", listTasks)
		r.Post("/", createTask)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", getTask)
			r.Delete("/", deleteTask)
		})
	})

	return r
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           newRouter(logger),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	logger.Info("chi router listening", "addr", ":8080")
	log.Fatal(server.ListenAndServe())
}
