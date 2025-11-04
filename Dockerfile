# Build Stage - supports any platform
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# These args are automatically provided by Docker BuildKit
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the binaries for the target platform
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o medistream ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o worker ./cmd/worker/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /app/medistream .
COPY --from=builder /app/worker .
RUN chmod +x /app/medistream /app/worker

# Copy the .env file from builder stage
COPY --from=builder /app/.env .env

# Copy the migrations folder
COPY --from=builder /app/migrations ./migrations

# Expose the backend port
EXPOSE 8080

# Run the backend
CMD ["./medistream"]