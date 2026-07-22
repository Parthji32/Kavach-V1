# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy mod files
COPY go.mod go.sum ./

# Download dependencies  
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o kavach ./cmd/server

# Runtime
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/kavach .
COPY migrations ./migrations
COPY static ./static
COPY templates ./templates

EXPOSE 3000

CMD ["./kavach"]
