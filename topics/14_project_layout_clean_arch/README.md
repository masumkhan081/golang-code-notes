# Project Layout & Clean Architecture

This folder demonstrates a standard Go project layout following **Clean Architecture** principles.

## Structure

```
14_project_layout_clean_arch/
├── cmd/
│   └── server/
│       └── main.go         # Entry point, wires everything together
├── internal/
│   ├── domain/             # Core business entities (No dependencies)
│   │   └── user.go
│   ├── repository/         # Data access interfaces & implementations
│   │   └── user_repo.go
│   └── service/            # Business logic (Depends on domain & repository)
│       └── user_service.go
└── README.md
```

## Key Principles

1.  **Dependency Rule**: Dependencies point _inwards_. `cmd` depends on `service`, `service` depends on `domain` and `repository`. `domain` depends on nothing.
2.  **Internal Package**: Code in `internal/` cannot be imported by external projects, enforcing encapsulation.
3.  **Dependency Injection**: `main.go` injects the repository into the service.

## How to Run

From the project root:

```sh
go run 14_project_layout_clean_arch/cmd/server/main.go
```
