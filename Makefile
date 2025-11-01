.PHONY: all build build-rust build-go test clean install help examples deps check-deps

# Default target
all: check-deps build test

# Help target
help:
	@echo "SurrealDB Embedded - Build Targets"
	@echo "==================================="
	@echo ""
	@echo "Setup:"
	@echo "  deps        - Check and install system dependencies"
	@echo "  check-deps  - Check if all required tools are installed"
	@echo ""
	@echo "Building:"
	@echo "  build       - Build both Rust library and Go wrapper"
	@echo "  build-rust  - Build only the Rust FFI library"
	@echo "  build-go    - Build only the Go wrapper"
	@echo ""
	@echo "Testing:"
	@echo "  test        - Run all Go tests"
	@echo "  test-v      - Run tests with verbose output"
	@echo "  bench       - Run benchmarks"
	@echo ""
	@echo "Examples:"
	@echo "  examples    - Run all examples"
	@echo "  ex-basic    - Run basic example"
	@echo "  ex-persist  - Run persistent storage example"
	@echo "  ex-graph    - Run graph relations example"
	@echo ""
	@echo "Maintenance:"
	@echo "  clean       - Clean all build artifacts"
	@echo "  clean-rust  - Clean Rust build artifacts"
	@echo "  clean-go    - Clean Go build artifacts"
	@echo ""
	@echo "Development:"
	@echo "  dev         - Build in debug mode (faster compilation)"
	@echo "  fmt         - Format all code"
	@echo "  lint        - Run linters"

# Check dependencies
check-deps:
	@echo "Checking dependencies..."
	@command -v rustc >/dev/null 2>&1 || { echo "ERROR: Rust is not installed. Run: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"; exit 1; }
	@command -v cargo >/dev/null 2>&1 || { echo "ERROR: Cargo is not installed."; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "ERROR: Go is not installed. Visit: https://golang.org/dl/"; exit 1; }
	@command -v clang >/dev/null 2>&1 || { echo "ERROR: Clang is not installed. See SETUP.md"; exit 1; }
	@echo "✓ All required tools are installed"
	@echo "  Rust:  $$(rustc --version)"
	@echo "  Cargo: $$(cargo --version)"
	@echo "  Go:    $$(go version)"
	@echo "  Clang: $$(clang --version | head -n1)"

# Install dependencies (Ubuntu/Debian)
deps:
	@echo "Installing system dependencies (Ubuntu/Debian)..."
	@if [ -f /etc/debian_version ]; then \
		sudo apt-get update && \
		sudo apt-get install -y build-essential clang libclang-dev llvm-dev pkg-config libssl-dev; \
	elif [ -f /etc/redhat-release ]; then \
		sudo dnf install -y gcc gcc-c++ clang clang-devel llvm-devel openssl-devel; \
	elif [ "$$(uname)" = "Darwin" ]; then \
		brew install llvm; \
	else \
		echo "Unknown OS. Please install dependencies manually. See SETUP.md"; \
		exit 1; \
	fi
	@echo "✓ System dependencies installed"

# Build everything
build: build-rust build-go
	@echo "✓ Build complete!"

# Build Rust library
build-rust:
	@echo "Building Rust FFI library (this may take 15-30 minutes on first build)..."
	@cd surrealdb_embedded_rs && cargo build --release
	@echo "✓ Rust library built successfully"
	@ls -lh surrealdb_embedded_rs/target/release/libsurrealdb_embedded_rs.* 2>/dev/null || true

# Build Go wrapper
build-go:
	@echo "Building Go wrapper..."
	@go mod download
	@go build -v
	@echo "✓ Go wrapper ready"

# Development build (debug mode - faster compilation)
dev:
	@echo "Building in debug mode..."
	@cd surrealdb_embedded_rs && cargo build
	@echo "✓ Debug build complete"
	@echo "Note: Update CGo flags to use target/debug for testing"

# Run all tests
test: build
	@echo "Running tests..."
	@go test -v

# Run tests with verbose output
test-v: build
	@go test -v -run .

# Run benchmarks
bench: build
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem

# Run all examples
examples: ex-basic ex-persist ex-graph

# Run basic example
ex-basic: build
	@echo "Running basic example..."
	@go run examples/basic/main.go

# Run persistent storage example
ex-persist: build
	@echo "Running persistent storage example..."
	@go run examples/persistent/main.go

# Run graph relations example
ex-graph: build
	@echo "Running graph relations example..."
	@go run examples/graph/main.go

# Clean everything
clean: clean-rust clean-go
	@echo "✓ All build artifacts cleaned"

# Clean Rust build artifacts
clean-rust:
	@echo "Cleaning Rust build artifacts..."
	@cd surrealdb_embedded_rs && cargo clean
	@rm -rf surrealdb_embedded_rs/target
	@rm -f surrealdb_embedded_rs/include/*.h
	@echo "✓ Rust artifacts cleaned"

# Clean Go build artifacts
clean-go:
	@echo "Cleaning Go build artifacts..."
	@go clean -cache -testcache -modcache
	@rm -rf examples/*/data examples/*/*.db
	@rm -f *.test
	@echo "✓ Go artifacts cleaned"

# Format code
fmt:
	@echo "Formatting Rust code..."
	@cd surrealdb_embedded_rs && cargo fmt
	@echo "Formatting Go code..."
	@go fmt ./...
	@echo "✓ Code formatted"

# Run linters
lint:
	@echo "Linting Rust code..."
	@cd surrealdb_embedded_rs && cargo clippy -- -D warnings
	@echo "Linting Go code..."
	@go vet ./...
	@echo "✓ Linting complete"

# Install locally (for development)
install: build
	@echo "Installing locally..."
	@go install
	@echo "✓ Installed"

# Rebuild everything from scratch
rebuild: clean build

# Quick test (build and run basic test)
quick: build
	@echo "Quick test..."
	@go test -run TestNewMemory -v

# Generate header file
header:
	@echo "Generating C header..."
	@cd surrealdb_embedded_rs && cargo build --release
	@echo "✓ Header generated at surrealdb_embedded_rs/include/"

# Show build info
info:
	@echo "Build Information"
	@echo "================="
	@echo "Rust library:  surrealdb_embedded_rs/target/release/libsurrealdb_embedded_rs.*"
	@echo "C header:      surrealdb_embedded_rs/include/surrealdb_embedded_rs.h"
	@echo ""
	@echo "Library size:"
	@ls -lh surrealdb_embedded_rs/target/release/libsurrealdb_embedded_rs.* 2>/dev/null || echo "  Not built yet"
	@echo ""
	@echo "Go module:     $$(go list -m)"
	@echo "Go version:    $$(go version)"
	@echo "Rust version:  $$(rustc --version)"
