# Multi-stage build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy all go files
COPY go.mod go.sum ./

# Download dependencies with GOSUMDB off
RUN GOSUMDB=off go mod download

# Copy source code
COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations

# Build binary
RUN CGO_ENABLED=1 GOOS=linux go build -o kavach ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/kavach .

# Copy migrations and static files
COPY migrations ./migrations
COPY static ./static
COPY templates ./templates

# Expose port
EXPOSE 3000

# Run
CMD ["./kavach"]
