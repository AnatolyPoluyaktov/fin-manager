# Variables
APP_NAME = fin-manager
DOCKER_COMPOSE_FILE = docker-compose.yml
MIGRATIONS_DIR = migrations
# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	go mod tidy
	go build -o $(APP_NAME) .

# Run tests
.PHONY: test
test: start-db apply-migrations
	go test ./...

# Start PostgreSQL using Docker Compose
.PHONY: start-db
start-db:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Stop PostgreSQL using Docker Compose
.PHONY: stop-db
stop-db:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Apply migrations
.PHONY: apply-migrations
apply-migrations:
	go run apply_migrations.go

# Clean up
.PHONY: clean
clean:
	rm -f $(APP_NAME)
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v