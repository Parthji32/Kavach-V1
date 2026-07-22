# Build stage - use golang alpine with full build support
FROM golang:1.22-alpine3.20 AS builder

WORKDIR /app

# Install build dependencies upfront
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    pkgconfig \
    ca-certificates

# Copy source code first
COPY . .

# Download dependencies
RUN go mod download

# Build with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o kavach .

# Runtime stage - minimal Alpine
FROM alpine:3.20

WORKDIR /app

# Install only runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

# Copy binary from builder
COPY --from=builder /app/kavach .

# Copy static assets
COPY migrations ./migrations
COPY static ./static
COPY templates ./templates

EXPOSE 3000

CMD ["./kavach"]
