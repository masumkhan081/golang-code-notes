# Go Code Notes

A comprehensive, topic-by-topic Go (Golang) interview and learning reference. Each folder covers a core Go concept, with runnable code examples, edge cases, and documentation.

---

## 🚀 How to Run Examples

```sh
go run path/to/file.go
```
Example:
```sh
go run 01_basics/data_types/data_types.go
```

---

## 📚 Directory Structure & Topics

Numbered in study order — prerequisites come before topics that depend on them.

### 🟢 Foundations (1-7)
| Folder | Purpose |
|---|---|
| 01_basics | Data types, constants, variables, import/export |
| 02_control_flow | if, for, switch, labels |
| 03_functions | Parameters, variadic, named return, closures |
| 04_collections | Arrays, slices, maps |
| 05_structs_methods_interfaces | Structs, methods, interfaces, embedding |
| 06_packages_modules | go mod, local imports, versioning |
| 07_errors_handling | error type, wrapping, custom errors |

### 🟡 Core Toolkit (8-13)
| Folder | Purpose |
|---|---|
| 08_files_io | Read/write files, CSV/JSON/Logs |
| 09_logging_config_env | Structured logging, config, env vars |
| 10_testing | Unit tests, benchmarks, table-driven, fuzzing |
| 11_concurrency | Goroutines, channels, sync primitives, patterns (ordered easy→hard) |
| 12_context | Cancellation, timeouts, request-scoped values |
| 13_generics | Type parameters, constraints, stdlib `slices`/`maps`/`cmp` |

### 🔴 Networking & Services (14-19)
| Folder | Purpose |
|---|---|
| 14_json_http_apis | RESTful APIs with `net/http` |
| 15_rest_routers | Routers (chi) |
| 16_database_sql_postgres | `database/sql`, sqlx, pgx |
| 17_caching_redis | Redis for sessions, caching |
| 18_queue_kafka_rabbitmq | Kafka / RabbitMQ, async tasks |
| 19_grpc | Proto definitions, gRPC server/client |

### ⚫ Production / Ops (20-24)
| Folder | Purpose |
|---|---|
| 20_project_layout_clean_arch | Clean project structure (after you've felt the pain) |
| 21_security_jwt_oauth2 | JWT, OAuth2, securing APIs |
| 22_docker_containerization | Dockerfiles, multi-stage builds |
| 23_devops_ci_cd | CI/CD, build pipelines |
| 24_advanced_concepts | Escape analysis, nil interface trap, slice traps, defer/panic/recover |

### 💼 Interview & Portfolio (25-26)
| Folder | Purpose |
|---|---|
| 25_interview_questions | Topic-grouped Q&A (~70 questions) |
| 26_projects | Mini-projects for portfolios |

---

## 🛠 Suggested Projects (in `/26_projects`)
| Project | Concepts Covered |
|---|---|
| expense-tracker | REST API, JWT, CRUD, PostgreSQL |
| concurrent-scraper | Goroutines, channels |
| auth-service | JWT, OAuth2, Redis, PostgreSQL |
| job-queue-service | Kafka/RabbitMQ, worker pool |
| grpc-user-service | Proto files, gRPC server/client |
| stock-price-watcher | Websockets, live updates |
| go-booking-system | Clean architecture, Docker, PostgreSQL |

---

## 🧭 Study Tips
- Start with foundational topics and master them before moving to advanced ones.
- Use the code examples for hands-on practice and interview prep.
- Build at least 2 real-world projects for portfolio strength.

Happy coding and good luck with your Go interviews!

