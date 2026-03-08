# Running the examples

This repo contains three types of content. Not every file works the same way.

---

## 1. Standalone `go run` examples

These are self-contained and print output immediately.

```bash
go run 99_advanced_concepts/escape_analysis/escape_analysis.go
go run 99_advanced_concepts/slice_growth_and_capacity/slice_growth_and_capacity.go
go run 99_advanced_concepts/interface_internal_representation/interface_internal_representation.go
go run 99_advanced_concepts/range_value_copy_traps/range_value_copy_traps.go
go run 99_advanced_concepts/defer_panic_recover/defer_panic_recover.go
go run 17_go_routines_advanced_patterns/worker_pool/worker_pool.go
go run 17_go_routines_advanced_patterns/semaphore_bounded_concurrency/semaphore_bounded_concurrency.go
go run 17_go_routines_advanced_patterns/mutex_protected_cache/mutex_protected_cache.go
go run 17_go_routines_advanced_patterns/fan_out_in/fan_out_in.go
go run 17_go_routines_advanced_patterns/pipeline/pipeline.go
go run 18_contexts_cancellation_timeout/context_misuse_pitfalls/context_misuse_pitfalls.go
go run 14_project_layout_clean_arch/dependency_injection_example/dependency_injection_example.go
go run 24_logging_config_env/slog_example/slog_example.go
go run 11_generics/generics.go
go run 09_files_io/files_io.go
```

Or use `make run-<topic>` shortcuts from this directory.

---

## 2. `go test` examples

These are test files. Run them with `go test`, not `go run`.

```bash
go test ./10_testing/...
```

Or individually:

```bash
go test -v -run TestAdd ./10_testing/
go test -v -run TestDiv ./10_testing/
go test -v -run TestDivSubtest ./10_testing/
go test -bench=. ./10_testing/
```

---

## 3. Server examples (long-running)

These start an HTTP server on `:8080` and do not exit on their own.

**Start the server:**
```bash
go run 12_json_http_apis/production_http_handler/production_http_handler.go
go run 12_json_http_apis/middleware_chain/middleware_chain.go
go run 15_rest_routers/chi_router_example/chi_router_example.go
go run 18_contexts_cancellation_timeout/graceful_shutdown.go
```

**Then hit an endpoint in another terminal:**
```bash
# middleware_chain and production_http_handler
curl -H "Authorization: Bearer dev-token" http://localhost:8080/tasks

# chi router
curl -H "X-Employee-ID: emp-1" http://localhost:8080/me/tasks

# health
curl http://localhost:8080/health
```

**Stop with:** `Ctrl+C`

---

## 4. Rate limiter (requires `golang.org/x/time/rate`)

Already in `go.mod`. Just run:

```bash
go run 17_go_routines_advanced_patterns/rate_limiter/rate_limiter.go
```

---

## 5. Examples requiring environment variables / external services

### Database examples — require a running Postgres instance

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/mydb?sslmode=disable"
go run 13_database_sql_postgres/database_sql_postgres.go
go run 13_database_sql_postgres/query_timeout_and_tx/query_timeout_and_tx.go
```

### Config example — requires DATABASE_URL

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/mydb?sslmode=disable"
go run 24_logging_config_env/env_config/env_config.go
```

---

## 6. Full app: taskapi

A multi-file service with its own `go.mod`. Run from its own directory:

```bash
cd 26_projects/taskapi
export DATABASE_URL="postgres://user:pass@localhost:5432/mydb?sslmode=disable"
go run ./cmd/api/main.go
```

The internal files (`internal/handler`, `internal/service`, etc.) are **not standalone** — they are part of the app and must be run via `cmd/api/main.go`.

---

## 7. Internal package files (not standalone runnable)

These are package files, not executables:

- `14_project_layout_clean_arch/internal/...`
- `26_projects/taskapi/internal/...`

They are wired together in `cmd/server/main.go` and `cmd/api/main.go` respectively.
