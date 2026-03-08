// production_http_handler.go
// Starts an HTTP server on :8080. This is a long-running server — it does not exit.
// Run: go run production_http_handler.go
// Then: curl -H "Authorization: Bearer dev-token" http://localhost:8080/tasks
// Stop: Ctrl+C
package main

import (
    "context"
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "strings"
    "time"
)

type createTaskRequest struct {
    Title string `json:"title"`
}

type taskResponse struct {
    ID    string `json:"id"`
    Title string `json:"title"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(v); err != nil {
        log.Printf("encode response: %v", err)
    }
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}

func createTask(ctx context.Context, title string) (taskResponse, error) {
    select {
    case <-ctx.Done():
        return taskResponse{}, ctx.Err()
    case <-time.After(100 * time.Millisecond):
        if strings.TrimSpace(title) == "" {
            return taskResponse{}, errors.New("title is required")
        }
        return taskResponse{ID: "task_123", Title: title}, nil
    }
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        writeError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    defer r.Body.Close()

    var req createTaskRequest
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()

    if err := dec.Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON body")
        return
    }

    task, err := createTask(r.Context(), req.Title)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
            writeError(w, http.StatusGatewayTimeout, "request canceled or timed out")
            return
        }
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }

    writeJSON(w, http.StatusCreated, task)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/tasks", taskHandler)

    server := &http.Server{
        Addr:              ":8080",
        Handler:           mux,
        ReadHeaderTimeout: 2 * time.Second,
        ReadTimeout:       5 * time.Second,
        WriteTimeout:      5 * time.Second,
        IdleTimeout:       30 * time.Second,
    }

    log.Println("listening on :8080")
    log.Fatal(server.ListenAndServe())
}
