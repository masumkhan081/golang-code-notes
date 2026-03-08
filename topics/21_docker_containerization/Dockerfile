# Start with a small base image
# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./14_project_layout_clean_arch/cmd/server/main.go

# Stage 2: Create a minimal image for running the app
FROM alpine:latest

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/myapp .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
