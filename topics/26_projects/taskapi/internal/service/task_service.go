package service

import (
	"context"
	"fmt"

	"github.com/masumkhan081/taskapi/internal/domain"
	"github.com/masumkhan081/taskapi/internal/repo"
)

// TaskService holds business logic, decoupled from transport.
type TaskService struct {
	repo repo.TaskRepo
}

func NewTaskService(r repo.TaskRepo) *TaskService {
	return &TaskService{repo: r}
}

func (s *TaskService) GetTask(ctx context.Context, id string) (domain.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) ListTasks(ctx context.Context) ([]domain.Task, error) {
	return s.repo.List(ctx)
}

func (s *TaskService) CreateTask(ctx context.Context, title string) (domain.Task, error) {
	if title == "" {
		return domain.Task{}, fmt.Errorf("title is required")
	}
	t := domain.Task{
		ID:    fmt.Sprintf("task-%d", nextID()),
		Title: title,
		Done:  false,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return domain.Task{}, err
	}
	return t, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

var idCounter int

func nextID() int {
	idCounter++
	return idCounter
}
