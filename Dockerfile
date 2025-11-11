# ---------- Build Stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build tools and goose
RUN apk add --no-cache bash git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

# Install goose globally
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy project files
COPY . .

# Build binaries
RUN go build -o medistream ./cmd/api/main.go
RUN go build -o worker ./cmd/worker/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.19

WORKDIR /app
RUN apk add --no-cache bash

# Copy binaries
COPY --from=builder /app/medistream .
COPY --from=builder /app/worker .
RUN chmod +x /app/medistream /app/worker

# Copy .env and migrations
COPY --from=builder /app/.env .env
COPY --from=builder /app/migrations ./migrations

# Copy goose binary from builder
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy migration script
COPY --from=builder /app/scripts/migrate.sh ./scripts/migrate.sh
RUN chmod +x ./scripts/migrate.sh

# Expose backend port
EXPOSE 8080

# Default command — runs migrations then starts app
CMD ./scripts/migrate.sh up && ./medistream
