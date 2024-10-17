# Build Stage
FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o rider-assignment-system main.go

# Run Stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/rider-assignment-system .
COPY config/config.yaml .
COPY postman_collection.json .

CMD ["./rider-assignment-system"]
