// http_handler_test.go
// Run with: go test -v ./10_testing/
// NOT runnable with go run — this is a test file.
package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

type createUserRequest struct {
    Name string `json:"name"`
}

type createUserResponse struct {
    Name string `json:"name"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req createUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Name) == "" {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(createUserResponse{Name: req.Name})
}

func TestCreateUserHandler(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name       string
        method     string
        body       string
        wantStatus int
    }{
        {
            name:       "created",
            method:     http.MethodPost,
            body:       `{"name":"alice"}`,
            wantStatus: http.StatusCreated,
        },
        {
            name:       "bad json",
            method:     http.MethodPost,
            body:       `{`,
            wantStatus: http.StatusBadRequest,
        },
        {
            name:       "empty name",
            method:     http.MethodPost,
            body:       `{"name":"   "}`,
            wantStatus: http.StatusBadRequest,
        },
        {
            name:       "wrong method",
            method:     http.MethodGet,
            body:       ``,
            wantStatus: http.StatusMethodNotAllowed,
        },
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            req := httptest.NewRequest(tt.method, "/users", strings.NewReader(tt.body))
            req.Header.Set("Content-Type", "application/json")
            rec := httptest.NewRecorder()

            createUserHandler(rec, req)

            if rec.Code != tt.wantStatus {
                t.Fatalf("got status %d want %d body=%s", rec.Code, tt.wantStatus, rec.Body.String())
            }
        })
    }
}
