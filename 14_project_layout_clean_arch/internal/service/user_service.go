package service

import (
	"context"
	"golang-code-notes/14_project_layout_clean_arch/internal/domain"
	"golang-code-notes/14_project_layout_clean_arch/internal/repository"
)

// UserService handles business logic for users.
type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, id, name, email string, age int) error {
	// Business logic: Validate and create
	user, err := domain.NewUser(id, name, email, age)
	if err != nil {
		return err
	}

	// Business logic: Check if user exists (omitted for brevity)

	return s.repo.Save(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}
