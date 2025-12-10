.PHONY: build run test clean swagger wire docker-up docker-down install-tools

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/server cmd/server/main.go

# Run the application
run:
	@echo "Running application..."
	@go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf docs/

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs

# Generate Wire dependency injection code
wire:
	@echo "Generating Wire code..."
	@cd internal/wire && wire

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/google/wire/cmd/wire@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

# Start Docker services (MySQL and Redis)
docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d

# Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down

# Start all (docker services + app)
start: docker-up
	@echo "Waiting for services to be ready..."
	@sleep 5
	@make run

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
