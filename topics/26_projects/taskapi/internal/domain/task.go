package domain

import "errors"

// ErrNotFound is returned when a task does not exist.
var ErrNotFound = errors.New("task not found")

// Task is the core domain entity.
type Task struct {
	ID    string
	Title string
	Done  bool
}
