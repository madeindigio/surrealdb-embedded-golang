# Build Instructions - Quick Reference

## Prerequisites Installation

### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install -y build-essential clang libclang-dev llvm-dev pkg-config libssl-dev

# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env
```

### macOS
```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install LLVM
brew install llvm

# Add LLVM to path (add to ~/.zshrc or ~/.bashrc)
export PATH="/opt/homebrew/opt/llvm/bin:$PATH"
export LIBCLANG_PATH="/opt/homebrew/opt/llvm/lib"

# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env
```

### Fedora/RHEL
```bash
sudo dnf install -y gcc gcc-c++ clang clang-devel llvm-devel openssl-devel

# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env
```

## Build Steps

### 1. Build the Rust Library

```bash
cd surrealdb_embedded_rs
cargo build --release
```

**Note**: First build takes 15-30 minutes. Subsequent builds are much faster.

Expected output:
```
   Compiling surrealdb v2.1.x
   Compiling surrealdb_embedded_rs v0.1.0
    Finished release [optimized] target(s) in 20m 15s
```

Verify the build:
```bash
# Check library was created
ls -lh target/release/libsurrealdb_embedded_rs.a

# Check header was generated
ls -lh include/surrealdb_embedded_rs.h
```

### 2. Test the Go Library

```bash
cd ..  # Back to project root
go mod download
go test -v
```

Expected output:
```
=== RUN   TestNewMemory
--- PASS: TestNewMemory (0.01s)
=== RUN   TestCreate
--- PASS: TestCreate (0.02s)
...
PASS
ok      github.com/yourusername/surrealdb-embedded    0.523s
```

### 3. Run Examples

```bash
# Basic example
go run examples/basic/main.go

# Persistent storage example
go run examples/persistent/main.go

# Graph relations example
go run examples/graph/main.go
```

## Common Build Errors & Solutions

### Error: "couldn't find libclang"

**Solution:**
```bash
# Ubuntu/Debian
sudo apt-get install libclang-dev

# macOS
brew install llvm
export LIBCLANG_PATH="$(brew --prefix llvm)/lib"

# Add to your shell profile (~/.bashrc or ~/.zshrc)
echo 'export LIBCLANG_PATH="/opt/homebrew/opt/llvm/lib"' >> ~/.zshrc
```

### Error: "linker could not find -lsurrealdb_embedded_rs"

**Solution:**
The Rust library hasn't been built yet.
```bash
cd surrealdb_embedded_rs
cargo build --release
cd ..
```

### Error: Out of memory during build

**Solution:**
Limit parallel compilation jobs:
```bash
cd surrealdb_embedded_rs
cargo build --release -j 2  # Use only 2 parallel jobs
```

### Error: "undefined reference to OpenSSL"

**Solution:**
Install OpenSSL development libraries:
```bash
# Ubuntu/Debian
sudo apt-get install libssl-dev

# Fedora/RHEL
sudo dnf install openssl-devel

# macOS (usually already installed)
brew install openssl
```

## Platform-Specific Notes

### macOS Apple Silicon (M1/M2/M3)

The build automatically creates a universal binary. No special configuration needed.

### Windows

1. Install [Visual Studio Build Tools](https://visualstudio.microsoft.com/downloads/)
   - Select "Desktop development with C++"
   
2. Install [LLVM for Windows](https://releases.llvm.org/download.html)
   - Add to PATH during installation

3. Build:
   ```cmd
   cd surrealdb_embedded_rs
   cargo build --release
   ```

## Speeding Up Builds

### Use a Faster Linker (Linux)

```bash
# Install mold linker
sudo apt install mold

# Add to ~/.bashrc or ~/.zshrc
export RUSTFLAGS="-C link-arg=-fuse-ld=mold"
```

### Use Cargo Cache

Install cargo-cache to clean old build artifacts:
```bash
cargo install cargo-cache
cargo cache --autoclean
```

### Incremental Compilation

For development (slower runtime, faster compile):
```bash
# Use debug build
cargo build  # instead of cargo build --release

# Update Go CGo flags to use debug:
# Change: -L${SRCDIR}/surrealdb_embedded_rs/target/release
# To:     -L${SRCDIR}/surrealdb_embedded_rs/target/debug
```

## Verifying Installation

### Check Versions

```bash
# Rust
rustc --version
# Should show: rustc 1.7x.x or later

# Go
go version
# Should show: go version go1.21.x or later

# Clang
clang --version
# Should show: clang version 10.x or later
```

### Run All Tests

```bash
# Run all Go tests
go test -v

# Run with race detector
go test -race -v

# Run benchmarks
go test -bench=. -benchmem
```

## Clean Build

If you need to start fresh:

```bash
# Clean Rust build
cd surrealdb_embedded_rs
cargo clean
cargo build --release
cd ..

# Clean Go cache
go clean -cache -testcache -modcache

# Rebuild
go test -v
```

## Cross-Compilation

### Build for Linux from macOS

```bash
# Install target
rustup target add x86_64-unknown-linux-gnu

# Install cross-compilation tools
brew install FiloSottile/musl-cross/musl-cross

# Build
cd surrealdb_embedded_rs
cargo build --release --target x86_64-unknown-linux-gnu
```

### Build for ARM64

```bash
rustup target add aarch64-unknown-linux-gnu
cd surrealdb_embedded_rs
cargo build --release --target aarch64-unknown-linux-gnu
```

## Optimizing Binary Size

Add to `surrealdb_embedded_rs/Cargo.toml`:

```toml
[profile.release]
opt-level = "z"     # Optimize for size
lto = true          # Link-time optimization
codegen-units = 1   # Better optimization
strip = true        # Strip debug symbols
```

Then rebuild:
```bash
cargo build --release
```

This can reduce binary size by 30-50%.

## Docker Build (Alternative)

If you prefer using Docker:

```dockerfile
# Dockerfile (create this file)
FROM rust:1.75 as rust-builder

RUN apt-get update && apt-get install -y \
    clang libclang-dev llvm-dev pkg-config libssl-dev

WORKDIR /build
COPY surrealdb_embedded_rs ./surrealdb_embedded_rs
RUN cd surrealdb_embedded_rs && cargo build --release

FROM golang:1.21

RUN apt-get update && apt-get install -y libssl-dev
COPY --from=rust-builder /build/surrealdb_embedded_rs/target/release/libsurrealdb_embedded_rs.a /usr/local/lib/
COPY --from=rust-builder /build/surrealdb_embedded_rs/include/surrealdb_embedded_rs.h /usr/local/include/

WORKDIR /app
COPY . .
RUN go build -o myapp
CMD ["./myapp"]
```

Build with Docker:
```bash
docker build -t surrealdb-embedded .
```

## Success Checklist

- [ ] Rust installed and working (`rustc --version`)
- [ ] Go installed and working (`go version`)
- [ ] Clang/LLVM installed (`clang --version`)
- [ ] Rust library built successfully
- [ ] Library file exists: `surrealdb_embedded_rs/target/release/libsurrealdb_embedded_rs.a`
- [ ] Header file exists: `surrealdb_embedded_rs/include/surrealdb_embedded_rs.h`
- [ ] Go tests pass: `go test -v`
- [ ] Examples run: `go run examples/basic/main.go`

## Getting Help

If you're still stuck:

1. Check the full [SETUP.md](SETUP.md) for detailed troubleshooting
2. Review build logs for specific errors
3. Search for the error message online
4. Open an issue on GitHub with:
   - Your OS and version
   - Output of `rustc --version` and `go version`
   - Full error output
   - Steps you've tried

## Quick Commands Reference

```bash
# Initial setup (once)
cd surrealdb_embedded_rs && cargo build --release && cd ..

# Run tests
go test -v

# Run example
go run examples/basic/main.go

# Rebuild after Rust changes
cd surrealdb_embedded_rs && cargo build --release && cd ..

# Clean everything
cd surrealdb_embedded_rs && cargo clean && cd .. && go clean -cache

# Benchmark
go test -bench=. -benchmem
```

---

**Estimated time for first-time setup:** 30-45 minutes  
**Estimated time for rebuilds:** 2-5 minutes
