#!/bin/bash
set -e

echo "Current directory: $(pwd)"
echo "Contents:"
ls -la

echo ""
echo "Building KAVACH..."
go build -v -o kavach ./cmd/server

echo "Build complete!"
ls -la kavach
