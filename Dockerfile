#Build Stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download



# COPY the entire project
COPY . .

# Build the binary

RUN go build -o medistream cmd/main.go

# Run stage

FROM alpine:latest
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/medistream .


# Copy the .env file from builder stage
COPY --from=builder /app/.env .env

# Copy the migrations folder
COPY --from=builder /app/migrations ./migrations

# Expose the backend port
EXPOSE 8080

# Run the backend
CMD ["./medistream"]

