.PHONY: build run debug clean fmt help docker-build docker-up docker-down docker-restart docker-logs docker-shell docker-clean dev

# Show help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Local Development:"
	@echo "  make build        - Build the binary"
	@echo "  make build-debug  - Build with debug symbols"
	@echo "  make run          - Run without building"
	@echo "  make run-debug    - Run with API throttle monitoring (DEBUG=true)"
	@echo "  make debug        - Run with Delve debugger"
	@echo "  make run-verbose  - Run with race detector and verbose output"
	@echo "  make start        - Build and run"
	@echo "  make clean        - Remove binaries"
	@echo "  make fmt          - Format code"
	@echo "  make dev-deps     - Install development dependencies"
	@echo ""
	@echo "Docker Development:"
	@echo "  make dev          - Start Docker with live development setup"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make docker-restart - Rebuild and restart Docker containers"
	@echo "  make docker-logs  - Show Docker container logs (tail)"
	@echo "  make docker-logs-f - Follow Docker container logs"
	@echo "  make docker-shell - Open shell in running container"
	@echo "  make docker-clean - Stop containers and remove images"
	@echo "  make docker-test  - Test API endpoint in Docker"

# Build the binary
build:
	go build -o stargazer-be ./cmd/server

# Build with debug symbols and optimizations disabled
build-debug:
	go build -gcflags="all=-N -l" -o stargazer-be-debug ./cmd/server

# Run without building
run:
	go run ./cmd/server

# Run with API throttle monitoring (DEBUG mode)
run-debug:
	@echo "Running in DEBUG mode with API throttle monitoring..."
	DEBUG=true go run ./cmd/server

# Run with Delve debugger
debug:
	dlv debug ./cmd/server

# Run with verbose logging
run-verbose:
	@echo "Running with detailed output..."
	go run -race ./cmd/server

# Build and run
start: build
	./stargazer-be

# Clean build artifacts
clean:
	rm -f stargazer-be stargazer-be-debug

# Format code
fmt:
	go fmt ./...

# Install development dependencies
dev-deps:
	@echo "Installing Delve debugger..."
	go install github.com/go-delve/delve/cmd/dlv@latest

# ============================================================================
# Docker Development Commands
# ============================================================================

# Start development environment with Docker
dev: docker-down docker-build docker-up
	@echo "✅ Development environment started!"
	@echo "📡 API: http://localhost:8080/launches"
	@echo "❤️  Health: http://localhost:8080/health"
	@echo "📋 Logs: make docker-logs-f"

# Build Docker image
docker-build:
	@echo "🔨 Building Docker image..."
	docker-compose build

# Start Docker containers
docker-up:
	@echo "🚀 Starting Docker containers..."
	docker-compose up -d
	@echo "⏳ Waiting for services to be ready..."
	@sleep 3
	@echo "✅ Containers started!"

# Stop Docker containers
docker-down:
	@echo "🛑 Stopping Docker containers..."
	docker-compose down

# Rebuild and restart Docker containers (fast development workflow)
docker-restart:
	@echo "🔄 Rebuilding and restarting containers..."
	docker-compose build
	docker-compose up -d
	@echo "⏳ Waiting for services to be ready..."
	@sleep 3
	@echo "✅ Containers restarted!"
	@make docker-logs

# Show Docker container logs (last 50 lines)
docker-logs:
	@echo "📋 Container logs (last 50 lines):"
	docker logs rleu-api --tail 50

# Follow Docker container logs in real-time
docker-logs-f:
	@echo "📋 Following container logs (Ctrl+C to exit):"
	docker logs rleu-api -f

# Open shell in running container
docker-shell:
	@echo "🐚 Opening shell in container..."
	docker exec -it rleu-api /bin/sh

# Stop containers and remove images
docker-clean:
	@echo "🧹 Cleaning up Docker resources..."
	docker-compose down
	docker rmi stargazer-be-rleu-api 2>/dev/null || true
	@echo "✅ Cleanup complete!"

# Test API endpoint in Docker
docker-test:
	@echo "🧪 Testing API endpoints..."
	@echo ""
	@echo "Health Check:"
	@curl -s http://localhost:8080/health | python3 -m json.tool || echo "❌ Health check failed"
	@echo ""
	@echo "Launches Endpoint:"
	@curl -s http://localhost:8080/launches | python3 -m json.tool | head -20 || echo "❌ API test failed"
	@echo ""
	@echo "Container Status:"
	@docker ps --filter name=rleu-api --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
