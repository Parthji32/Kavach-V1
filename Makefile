build:
	go build -o kavach ./cmd/server

run:
	./kavach

dev:
	go run ./cmd/server/main.go

test:
	go test ./...

.PHONY: build run dev test
