package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/masumkhan081/taskapi/internal/config"
	"github.com/masumkhan081/taskapi/internal/handler"
	"github.com/masumkhan081/taskapi/internal/middleware"
	"github.com/masumkhan081/taskapi/internal/repo"
	"github.com/masumkhan081/taskapi/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 1. Load and validate config
	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", "err", err)
		os.Exit(1)
	}

	// 2. Wire dependencies (repo → service → handler)
	taskRepo := repo.NewInMemoryTaskRepo()
	taskSvc := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskSvc, logger)

	// 3. Build router with middleware chain
	mux := http.NewServeMux()
	taskHandler.RegisterRoutes(mux)

	var root http.Handler = mux
	root = middleware.RequestID(root)
	root = middleware.Logger(logger)(root)

	// 4. Configure server with timeouts
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           root,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	// 5. Start server in a goroutine
	go func() {
		logger.Info("server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// 6. Graceful shutdown on signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutdown signal received, draining connections...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "err", err)
	}
	logger.Info("server stopped cleanly")
}
