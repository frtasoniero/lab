# Makefile

# Binary name
BINARY_NAME=pwh-auth

# Docker variables
DOCKER_IMAGE=pwh-auth-go
DOCKER_CONTAINER=pwh-auth-app
DOCKER_COMPOSE=docker-compose.yml

# Default Go entrypoint
CMD=go run ./cmd/server

# MongoDB URI for local run
export MONGO_URI=mongodb://localhost:27017
export PORT=5005

# === Commands ===

.PHONY: all build run docker docker-up docker-down clean tidy fmt

## Run the app locally
run:
	@echo "Running the app locally..."
	$(CMD)

## Build the Go binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) ./cmd/server

## Format code
fmt:
	go fmt ./...

## Clean up binaries
clean:
	rm -f $(BINARY_NAME)

## Tidy Go dependencies
tidy:
	go mod tidy

## Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

## Run with Docker Compose
docker-up:
	@echo "Starting Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE) up --build -d

## Stop Docker Compose
docker-down:
	@echo "Stopping Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE) down -d

## Run only the app service (requires db already running)
docker-app-up:
	docker-compose -f $(DOCKER_COMPOSE) up --build app -d

## Run only the database (MongoDB)
docker-db-up:
	docker-compose -f $(DOCKER_COMPOSE) up mongo -d

# Seed the test database with mock users
mock-db:
	@echo "Populating auth_demo_test with mock users..."
	go run ./cmd/seed/main.go