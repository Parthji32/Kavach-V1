# Multi-stage build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod ./

# Copy source code
COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations

# Tidy and download dependencies (after source is copied so it can detect imports)
RUN GOSUMDB=off go mod tidy

# Download dependencies
RUN go mod download

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

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:3000/health || exit 1

# Run
CMD ["./kavach"]
