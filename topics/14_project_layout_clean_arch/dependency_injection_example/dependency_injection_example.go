package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type TaskRepository interface {
	ListByEmployeeID(ctx context.Context, employeeID string) ([]Task, error)
}

type InMemoryTaskRepository struct {
	data map[string][]Task
}

func (r *InMemoryTaskRepository) ListByEmployeeID(ctx context.Context, employeeID string) ([]Task, error) {
	return r.data[employeeID], nil
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) ListMyTasks(ctx context.Context, employeeID string) ([]Task, error) {
	return s.repo.ListByEmployeeID(ctx, employeeID)
}

type Server struct {
	taskService *TaskService
}

func NewServer(taskService *TaskService) *Server {
	return &Server{taskService: taskService}
}

func (s *Server) handleMyTasks(w http.ResponseWriter, r *http.Request) {
	employeeID := r.Header.Get("X-Employee-ID")
	if employeeID == "" {
		http.Error(w, "missing employee id", http.StatusUnauthorized)
		return
	}

	tasks, err := s.taskService.ListMyTasks(r.Context(), employeeID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
	}
}

func main() {
	repo := &InMemoryTaskRepository{
		data: map[string][]Task{
			"emp-1": {
				{ID: "t1", Title: "Prepare report"},
				{ID: "t2", Title: "Review PR"},
			},
		},
	}

	taskService := NewTaskService(repo)
	server := NewServer(taskService)

	mux := http.NewServeMux()
	mux.HandleFunc("/me/tasks", server.handleMyTasks)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
