// graceful_shutdown.go
// Starts an HTTP server on :8080 with SIGINT/SIGTERM graceful shutdown.
// Run: go run graceful_shutdown.go
// Then: curl http://localhost:8080/health
// Stop: Ctrl+C  (triggers graceful drain)
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func slowHandler(w http.ResponseWriter, r *http.Request) {
    select {
    case <-r.Context().Done():
        http.Error(w, "request canceled", http.StatusRequestTimeout)
    case <-time.After(2 * time.Second):
        w.Write([]byte("ok"))
    }
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/slow", slowHandler)

    server := &http.Server{
        Addr:              ":8080",
        Handler:           mux,
        ReadHeaderTimeout: 2 * time.Second,
    }

    go func() {
        log.Println("server listening on :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %v", err)
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    log.Println("shutting down")
    if err := server.Shutdown(ctx); err != nil {
        log.Printf("shutdown error: %v", err)
    }
}
