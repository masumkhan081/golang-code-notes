package repo

import (
	"context"
	"sync"

	"github.com/masumkhan081/taskapi/internal/domain"
)

// TaskRepo defines the storage contract.
type TaskRepo interface {
	GetByID(ctx context.Context, id string) (domain.Task, error)
	List(ctx context.Context) ([]domain.Task, error)
	Create(ctx context.Context, t domain.Task) error
	Delete(ctx context.Context, id string) error
}

// InMemoryTaskRepo is an in-memory implementation of TaskRepo.
type InMemoryTaskRepo struct {
	mu    sync.RWMutex
	store map[string]domain.Task
}

func NewInMemoryTaskRepo() *InMemoryTaskRepo {
	return &InMemoryTaskRepo{
		store: map[string]domain.Task{
			"1": {ID: "1", Title: "Buy groceries", Done: false},
			"2": {ID: "2", Title: "Write Go notes", Done: true},
		},
	}
}

func (r *InMemoryTaskRepo) GetByID(_ context.Context, id string) (domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.store[id]
	if !ok {
		return domain.Task{}, domain.ErrNotFound
	}
	return t, nil
}

func (r *InMemoryTaskRepo) List(_ context.Context) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Task, 0, len(r.store))
	for _, t := range r.store {
		out = append(out, t)
	}
	return out, nil
}

func (r *InMemoryTaskRepo) Create(_ context.Context, t domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[t.ID] = t
	return nil
}

func (r *InMemoryTaskRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.store, id)
	return nil
}
