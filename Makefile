.PHONY: build run debug clean fmt help

# Show help
help:
	@echo "Available commands:"
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
