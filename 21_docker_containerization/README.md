# Docker & Containerization

This folder contains examples of how to containerize a Go application.

## Files

- **Dockerfile**: A multi-stage build optimized for size. It builds the Go binary in a `golang` image and then copies it to a minimal `alpine` image.
- **docker-compose.yml**: Orchestrates the Go application along with a PostgreSQL database.

## How to Build and Run

1.  **Build the Docker Image**:

    ```sh
    docker build -t my-go-app -f 21_docker_containerization/Dockerfile .
    ```

2.  **Run with Docker Compose**:
    ```sh
    docker-compose -f 21_docker_containerization/docker-compose.yml up --build
    ```

## Best Practices

- **Multi-stage builds**: Reduces image size drastically (from ~800MB to ~15MB).
- **Non-root user**: For security (not shown here for simplicity, but recommended).
- **Environment Variables**: Pass config via `ENV` or `.env` files.
