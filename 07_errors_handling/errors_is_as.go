package main

import (
    "errors"
    "fmt"
)

var ErrNotFound = errors.New("not found")

type ValidationError struct {
    Field string
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Msg)
}

func loadUser(id int) error {
    if id == 0 {
        return fmt.Errorf("load user: %w", ErrNotFound)
    }
    if id < 0 {
        return fmt.Errorf("load user: %w", &ValidationError{
            Field: "id",
            Msg:   "must be positive",
        })
    }
    return nil
}

func main() {
    for _, id := range []int{0, -1, 42} {
        err := loadUser(id)
        if err == nil {
            fmt.Println("ok:", id)
            continue
        }

        fmt.Println("error:", err)

        if errors.Is(err, ErrNotFound) {
            fmt.Println("matched sentinel ErrNotFound")
        }

        var vErr *ValidationError
        if errors.As(err, &vErr) {
            fmt.Printf("matched ValidationError field=%s msg=%s\n", vErr.Field, vErr.Msg)
        }
    }
}
