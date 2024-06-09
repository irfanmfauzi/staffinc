
# Simple Makefile for a Go project

# Build the application
all: build

templ-generate:
	@echo "Generating templates..."
	@templ generate

build: templ-generate
	@echo "Building..."
	@go build -o main cmd/monolith/main.go
	@echo "Finish Build"

# Run the application
run:
	@go run cmd/monolith/main.go


swagger-doc:
	@swag init -g cmd/monolith/main.go 

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test -v -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out
	@go tool cover -func coverage.out
	@echo "Finish testing"

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

mock:
	mockery --dir ./internals/repository --all --output ./internals/repository/mocks


.PHONY: all build generate run test clean
