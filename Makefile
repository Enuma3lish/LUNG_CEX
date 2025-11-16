.PHONY: help setup db-up db-down backend frontend docker-up docker-down clean

help:
	@echo "LUNG CEX - Virtual Trading Platform"
	@echo ""
	@echo "Available commands:"
	@echo "  make setup        - Setup project (install dependencies)"
	@echo "  make docker-up    - Start PostgreSQL and Redis with Docker"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make backend      - Run backend server"
	@echo "  make frontend     - Run frontend dev server"
	@echo "  make clean        - Clean build artifacts"

setup:
	@echo "Installing backend dependencies..."
	cd backend && go mod download && go mod tidy
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Setup complete!"

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	sleep 5
	@echo "Docker containers are running!"

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

backend:
	@echo "Starting backend server..."
	cd backend && go run cmd/api/main.go

frontend:
	@echo "Starting frontend dev server..."
	cd frontend && npm run dev

clean:
	@echo "Cleaning build artifacts..."
	cd backend && go clean
	cd frontend && rm -rf dist node_modules/.vite
	@echo "Clean complete!"

dev: docker-up
	@echo "Starting development environment..."
	@echo "Backend will run on http://localhost:8080"
	@echo "Frontend will run on http://localhost:3000"
