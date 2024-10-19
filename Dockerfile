# Stage 1: Build the Go binaries
FROM golang:1.18-alpine AS builder

WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate run_migrations.go

# Build the main application
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

# Copy the binaries from the builder stage
COPY --from=builder /app/migrate .
COPY --from=builder /app/main .

# Copy the migration files
COPY database/migrations ./database/migrations

# Expose the application port
EXPOSE 8080

# Run the migration command before starting the server
CMD ["sh", "-c", "./migrate && ./main"]
