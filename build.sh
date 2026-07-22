#!/bin/sh
set -e

echo "Building KAVACH..."

# Build with proper CGO settings
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64

go mod download
go build -o kavach ./cmd/server

echo "Build complete!"
