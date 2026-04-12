.PHONY: build run test test-race lint docker-up docker-down docker-logs mock clean help test-upload

# Variables
BINARY_NAME=streaming-log-consumer
MAIN_PATH=cmd/server/main.go

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## [-a-zA-Z0-9_]+:' Makefile | sed 's/## //g' | awk -F: '{printf "  %-15s %s\n", $$1, $$2}'

## build: Build the Go binary locally
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

## run: Run the application locally (expects local DB)
run:
	go run $(MAIN_PATH)

## test: Run all unit and integration tests
test:
	go test ./...

## test-race: Run tests with data race detection
test-race:
	go test -race ./...

## lint: Run all linters (expects golangci-lint installed)
lint:
	golangci-lint run ./...

## clean: Remove build artifacts and local storage files
clean:
	rm -rf bin/
	rm -rf storage/*
