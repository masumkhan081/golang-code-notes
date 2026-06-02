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

### 🟢 Foundational (Must Know First)
| Folder | Purpose |
|---|---|
| 01_basics | Data types, constants, variables, import/export |
| 02_control_flow | if, for, switch, labels |
| 03_functions | Parameters, variadic, named return, closures |
| 04_arrays_slices_maps | Core collections used everywhere |
| 05_structs_methods_interfaces | Structs, methods, interfaces, embedding |
| 06_packages_modules | go mod, local imports, versioning |
| 07_errors_handling | error type, wrapping, custom errors |

### 🟡 Intermediate (Production-Ready)
| Folder | Purpose |
|---|---|
| 08_concurrency | Goroutines, channels, sync primitives, patterns (ordered easy→hard) |
| 09_files_io | Read/write files, CSV/JSON/Logs |
| 10_testing | Unit tests, benchmarks, table-driven tests |
| 11_generics | Go 1.18+ generics |
| 12_json_http_apis | RESTful APIs with net/http |
| 13_database_sql_postgres | database/sql, sqlx, gorm usage |

### 🔴 Advanced / Real-World Engineering
| Folder | Purpose |
|---|---|
| 14_project_layout_clean_arch | Clean project structure |
| 15_gorilla_mux_rest_api | REST API with Gorilla Mux |
| 16_grpc | Proto definitions, gRPC servers/clients |
| 18_contexts_cancellation_timeout | Context, cancellation, timeout |
| 19_caching_redis | Redis for sessions, caching |
| 20_queue_kafka_rabbitmq | Kafka/RabbitMQ, async tasks |
| 21_docker_containerization | Dockerfiles, multi-stage builds |
| 22_devops_ci_cd | CI/CD, build pipelines |
| 23_security_jwt_oauth2 | JWT, OAuth2, securing APIs |
| 24_logging_config_env | Logging, config, env, logrus/zap |

### 💼 Job Focused & Final Touch
| Folder | Purpose |
|---|---|
| 25_interview_questions | Real-world Go interview Q&A |
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

