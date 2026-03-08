package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidAge   = errors.New("age must be positive")
)

// User represents a domain entity.
// In Clean Architecture, this has no dependencies on other layers.
type User struct {
	ID    string
	Name  string
	Email string
	Age   int
}

// NewUser is a factory to create a user with validation.
func NewUser(id, name, email string, age int) (*User, error) {
	if age <= 0 {
		return nil, ErrInvalidAge
	}
	return &User{
		ID:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}, nil
}
