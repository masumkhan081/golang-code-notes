// with_value_and_keys.go
// context.WithValue done right.
//
// WithValue carries request-scoped data across API boundaries — request IDs,
// auth tokens, tracing spans. It is NOT a general-purpose dependency injector.
package main

import (
	"context"
	"fmt"
	"net/http"
)

// RULE 1: the key must be an unexported custom type, not a string.
// Why: prevents collisions across packages. A string "userID" from
// package A can clash with "userID" from package B.
type ctxKey int

const (
	keyRequestID ctxKey = iota
	keyUser
)

// RULE 2: export typed accessor functions, never the key itself.
// Callers don't reach into ctx with raw keys.

type User struct {
	ID   string
	Name string
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

func RequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keyRequestID).(string)
	return v, ok
}

func WithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, keyUser, u)
}

func CurrentUser(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(keyUser).(*User)
	return u, ok
}

// ---------------------------------------------------------------------------
// Realistic flow: HTTP middleware seeds the context; handler reads it.
// ---------------------------------------------------------------------------

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = "generated-id"
		}
		ctx := WithRequestID(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	if id, ok := RequestID(r.Context()); ok {
		fmt.Fprintf(w, "request id: %s", id)
	}
}

func main() {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-123")
	ctx = WithUser(ctx, &User{ID: "u1", Name: "Alice"})

	if id, ok := RequestID(ctx); ok {
		fmt.Println("rid:", id)
	}
	if u, ok := CurrentUser(ctx); ok {
		fmt.Println("user:", u.Name)
	}
}

/*
DO
  - Use an unexported custom key type, one per concern.
  - Wrap WithValue/Value in typed helper functions.
  - Pass ctx as the FIRST parameter, always named `ctx`.
  - Use only for request-scoped data that crosses API boundaries.

DON'T
  - Don't put structs full of configuration/dependencies into ctx.
    DI belongs in constructors, not in Value().
  - Don't use plain string/int keys. `go vet` warns about this.
  - Don't store context in a struct field. Pass it through call chains.
  - Don't pass nil context — use context.TODO() while wiring is incomplete,
    context.Background() at program/test roots.

PERFORMANCE
  ctx.Value walks a linked list of parents. O(depth). Fine for the typical
  3-5 deep chain; don't put hot-path values there.
*/
