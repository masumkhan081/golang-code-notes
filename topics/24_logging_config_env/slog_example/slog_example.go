// slog_example.go
// Demonstrates structured logging with log/slog (Go 1.21+).
package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {
	// ---- text handler (human-readable) ----
	textLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(textLogger)

	slog.Info("server started", "addr", ":8080", "env", "production")
	slog.Debug("debug detail", "goroutines", 4)
	slog.Warn("slow query", "duration_ms", 312, "query", "SELECT * FROM users")
	slog.Error("db connection failed", "err", "timeout after 2s")

	// ---- JSON handler (for log aggregators like Datadog, Loki) ----
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	jsonLogger.Info("request complete",
		"method", "POST",
		"path", "/tasks",
		"status", 201,
		"duration_ms", 45,
		"request_id", "abc123",
	)

	// ---- logger with persistent fields (With) ----
	// Good for request-scoped or service-scoped logging
	requestLogger := jsonLogger.With(
		"request_id", "req-789",
		"user_id", "usr-42",
	)
	requestLogger.Info("task created", "task_id", "task-001")
	requestLogger.Warn("permission check skipped", "reason", "admin user")

	// ---- context-aware logging ----
	// Attach logger to context so handlers can pull it without passing it explicitly
	ctx := context.Background()
	ctx = contextWithLogger(ctx, requestLogger)

	logger := loggerFromContext(ctx)
	logger.Info("handler called from context logger")

	// ---- log levels summary ----
	// slog.LevelDebug  = -4  (verbose internal state)
	// slog.LevelInfo   =  0  (normal operations)
	// slog.LevelWarn   =  4  (unexpected but recoverable)
	// slog.LevelError  =  8  (failures that need attention)
}

// contextKey is unexported to avoid collisions
type contextKey string

const loggerKey contextKey = "logger"

func contextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func loggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
