# Build Stage
FROM golang:1.17-alpine AS builder

WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o rider-assignment-system main.go

# Run Stage
FROM alpine:latest

WORKDIR /root/

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/rider-assignment-system .

# Copy the configuration file and migrations
COPY config/config.yaml ./config.yaml
COPY database/migrations ./database/migrations

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./rider-assignment-system"]
