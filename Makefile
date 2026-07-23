# KAVACH Makefile

.PHONY: build run dev test css css-watch setup clean

# Build Go binary
build: css
	go build -o server.exe .

# Run server
run:
	./server.exe

# Development mode (Go)
dev:
	go run .

# Run tests
test:
	go test ./...

# Build Tailwind CSS (minified)
css:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/tailwind.css --minify

# Watch Tailwind CSS for changes
css-watch:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/tailwind.css --watch

# Download HTMX (version-locked)
htmx:
	curl -o static/js/htmx.min.js https://unpkg.com/htmx.org@1.9.12/dist/htmx.min.js

# Install frontend dependencies
setup:
	npm install
	$(MAKE) htmx
	$(MAKE) css

# Clean build artifacts
clean:
	rm -f server.exe kavach
	rm -f static/css/tailwind.css
