package repository

import (
	"context"
	"github.com/masumkhan081/golang-code-notes/topics/14_project_layout_clean_arch/internal/domain"
	"sync"
)

// UserRepository defines the interface for user storage.
// This belongs in the 'internal' layer but is often defined near the domain or in a 'ports' package.
type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

// InMemoryUserRepo is a simple implementation of UserRepository.
type InMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepo) Save(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}
