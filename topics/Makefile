.PHONY: help test build tidy \
        run-escape-analysis run-slice-capacity run-interface-trap \
        run-range-copy run-defer-panic run-worker-pool \
        run-semaphore run-mutex-cache run-context-pitfalls \
        run-di-example run-slog run-rate-limiter \
        run-http run-middleware run-chi run-graceful-shutdown \
        run-taskapi

help:
	@echo ""
	@echo "  make test                   run all tests"
	@echo "  make build                  build all packages"
	@echo "  make tidy                   go mod tidy"
	@echo ""
	@echo "  --- standalone demos ---"
	@echo "  make run-escape-analysis"
	@echo "  make run-slice-capacity"
	@echo "  make run-interface-trap"
	@echo "  make run-range-copy"
	@echo "  make run-defer-panic"
	@echo "  make run-worker-pool"
	@echo "  make run-semaphore"
	@echo "  make run-mutex-cache"
	@echo "  make run-context-pitfalls"
	@echo "  make run-di-example"
	@echo "  make run-slog"
	@echo "  make run-rate-limiter        (uses golang.org/x/time/rate)"
	@echo ""
	@echo "  --- servers (Ctrl+C to stop) ---"
	@echo "  make run-http                production_http_handler on :8080"
	@echo "  make run-middleware          middleware_chain on :8080"
	@echo "  make run-chi                 chi router on :8080"
	@echo "  make run-graceful-shutdown   server with graceful drain on :8080"
	@echo ""
	@echo "  --- full app (needs DATABASE_URL) ---"
	@echo "  make run-taskapi"
	@echo ""

test:
	go test -v ./10_testing/...

build:
	go build ./...

tidy:
	go mod tidy

# --- advanced concepts ---
run-escape-analysis:
	go run ./99_advanced_concepts/escape_analysis/escape_analysis.go

run-slice-capacity:
	go run ./99_advanced_concepts/slice_growth_and_capacity/slice_growth_and_capacity.go

run-interface-trap:
	go run ./99_advanced_concepts/interface_internal_representation/interface_internal_representation.go

run-range-copy:
	go run ./99_advanced_concepts/range_value_copy_traps/range_value_copy_traps.go

run-defer-panic:
	go run ./99_advanced_concepts/defer_panic_recover/defer_panic_recover.go

# --- concurrency ---
run-worker-pool:
	go run ./17_go_routines_advanced_patterns/worker_pool/worker_pool.go

run-semaphore:
	go run ./17_go_routines_advanced_patterns/semaphore_bounded_concurrency/semaphore_bounded_concurrency.go

run-mutex-cache:
	go run ./17_go_routines_advanced_patterns/mutex_protected_cache/mutex_protected_cache.go

run-rate-limiter:
	go run ./17_go_routines_advanced_patterns/rate_limiter/rate_limiter.go

# --- context ---
run-context-pitfalls:
	go run ./18_contexts_cancellation_timeout/context_misuse_pitfalls/context_misuse_pitfalls.go

# --- architecture ---
run-di-example:
	go run ./14_project_layout_clean_arch/dependency_injection_example/dependency_injection_example.go

# --- logging ---
run-slog:
	go run ./24_logging_config_env/slog_example/slog_example.go

# --- servers ---
run-http:
	@echo "Server on :8080  — curl -H 'Authorization: Bearer dev-token' http://localhost:8080/tasks"
	go run ./12_json_http_apis/production_http_handler/production_http_handler.go

run-middleware:
	@echo "Server on :8080  — curl -H 'Authorization: Bearer dev-token' http://localhost:8080/tasks"
	go run ./12_json_http_apis/middleware_chain/middleware_chain.go

run-chi:
	@echo "Server on :8080  — curl -H 'X-Employee-ID: emp-1' http://localhost:8080/me/tasks"
	go run ./15_rest_routers/chi_router_example/chi_router_example.go

run-graceful-shutdown:
	@echo "Server on :8080  — curl http://localhost:8080/health    Ctrl+C to trigger graceful drain"
	go run ./18_contexts_cancellation_timeout/graceful_shutdown.go

# --- full app ---
run-taskapi:
	@echo "Requires: export DATABASE_URL=..."
	cd 26_projects/taskapi && go run ./cmd/api/main.go
