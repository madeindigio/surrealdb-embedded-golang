# Setup Guide

This guide will help you set up the development environment and build the SurrealDB embedded library.

## Prerequisites

### 1. Install Rust

Install the Rust toolchain using rustup:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env
```

Verify installation:
```bash
rustc --version
cargo --version
```

### 2. Install System Dependencies

#### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install -y \
    build-essential \
    clang \
    libclang-dev \
    llvm-dev \
    pkg-config \
    libssl-dev
```

#### Fedora/RHEL/CentOS
```bash
sudo dnf install -y \
    gcc \
    gcc-c++ \
    clang \
    clang-devel \
    llvm-devel \
    openssl-devel
```

#### macOS
```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install LLVM (includes clang)
brew install llvm
```

#### Windows
1. Install [Visual Studio Build Tools](https://visualstudio.microsoft.com/downloads/)
2. Install [LLVM](https://releases.llvm.org/download.html)
3. Add LLVM to PATH

### 3. Install Go

Download and install Go 1.21 or later from [golang.org](https://golang.org/dl/)

Verify installation:
```bash
go version
```

## Building the Project

### Step 1: Build the Rust Static Library

Navigate to the Rust project directory:

```bash
cd surrealdb_embedded_rs
```

Build the release version:

```bash
cargo build --release
```

This will:
- Download and compile all Rust dependencies (this may take 10-20 minutes on first build)
- Generate the static library at `target/release/libsurrealdb_embedded_rs.a` (Linux/macOS) or `surrealdb_embedded_rs.lib` (Windows)
- Generate the C header file at `include/surrealdb_embedded_rs.h`

### Step 2: Verify the Build

Check that the files were created:

```bash
# Linux/macOS
ls -lh target/release/libsurrealdb_embedded_rs.a
ls -lh include/surrealdb_embedded_rs.h

# Windows
dir target\release\surrealdb_embedded_rs.lib
dir include\surrealdb_embedded_rs.h
```

### Step 3: Build the Go Library

Navigate back to the project root:

```bash
cd ..
```

Download Go dependencies:

```bash
go mod download
```

### Step 4: Run Tests

Run the Go tests:

```bash
go test -v
```

If the tests pass, the integration is working correctly!

## Troubleshooting

### Error: "couldn't find libclang"

This means clang development libraries are not installed or not in the system path.

**Linux:**
```bash
sudo apt-get install libclang-dev
# or
export LIBCLANG_PATH=/usr/lib/llvm-14/lib  # Adjust version as needed
```

**macOS:**
```bash
brew install llvm
export LIBCLANG_PATH="$(brew --prefix llvm)/lib"
```

**Windows:**
```cmd
set LIBCLANG_PATH=C:\Program Files\LLVM\bin
```

### Error: "linker could not find -lsurrealdb_embedded_rs"

The Rust library hasn't been built yet. Run:

```bash
cd surrealdb_embedded_rs
cargo build --release
cd ..
```

### Error: "undefined reference" on Linux

You may need to install additional libraries:

```bash
sudo apt-get install -y libssl-dev pkg-config
```

### Slow Rust Build Times

First builds can take 15-30 minutes because RocksDB and other dependencies need to compile from source.

To speed up subsequent builds:

```bash
# Use the faster linker
cargo install -f cargo-binutils
rustup component add llvm-tools-preview

# Or use mold (Linux only)
sudo apt install mold
export RUSTFLAGS="-C link-arg=-fuse-ld=mold"
```

### Memory Issues During Build

If cargo runs out of memory:

```bash
# Build with fewer parallel jobs
cargo build --release -j 2

# Or limit memory per job
export CARGO_BUILD_JOBS=2
```

## Platform-Specific Notes

### Linux

The library should work on any modern Linux distribution. Tested on:
- Ubuntu 20.04+
- Debian 11+
- Fedora 36+
- CentOS 8+

### macOS

Works on:
- macOS 11 (Big Sur) or later
- Both Intel (x86_64) and Apple Silicon (aarch64)

For Apple Silicon, the build automatically creates a universal binary.

### Windows

Requires:
- Windows 10 or later
- Visual Studio 2019 Build Tools or later
- LLVM/Clang for Windows

Known limitations:
- Some file paths may need to use backslashes
- Administrator privileges might be needed for some operations

## Development Workflow

### Rebuilding After Changes

After modifying Rust code:

```bash
cd surrealdb_embedded_rs
cargo build --release
cd ..
go test -v
```

After modifying Go code:

```bash
go test -v
```

### Using Different Build Modes

Debug build (faster compilation, slower runtime):

```bash
cargo build  # Creates target/debug/libsurrealdb_embedded_rs.a
```

Update CGo flags in Go code to use debug build:

```go
#cgo LDFLAGS: -L${SRCDIR}/surrealdb_embedded_rs/target/debug -lsurrealdb_embedded_rs
```

### Cross-Compilation

To build for a different platform:

```bash
# Install target
rustup target add x86_64-unknown-linux-gnu

# Build for target
cargo build --release --target x86_64-unknown-linux-gnu
```

Common targets:
- `x86_64-unknown-linux-gnu` - Linux x64
- `aarch64-unknown-linux-gnu` - Linux ARM64
- `x86_64-apple-darwin` - macOS Intel
- `aarch64-apple-darwin` - macOS Apple Silicon
- `x86_64-pc-windows-msvc` - Windows x64

## Next Steps

After successful setup:

1. Run the examples:
   ```bash
   go run examples/basic/main.go
   go run examples/persistent/main.go
   go run examples/graph/main.go
   ```

2. Read the [README.md](README.md) for API documentation

3. Check out the [tests](surrealdb_test.go) for more usage examples

## Getting Help

If you encounter issues:

1. Check this troubleshooting guide
2. Review the build logs for specific error messages
3. Open an issue on GitHub with:
   - Your operating system and version
   - Rust version (`rustc --version`)
   - Go version (`go version`)
   - Full error output
   - Steps to reproduce

## Performance Tips

### For Development

Use debug builds for faster compilation:
```bash
cargo build  # Instead of cargo build --release
```

### For Production

1. Use release builds with optimizations:
   ```bash
   cargo build --release
   ```

2. Consider enabling CPU-specific optimizations:
   ```bash
   RUSTFLAGS="-C target-cpu=native" cargo build --release
   ```

3. Use thin LTO for smaller binaries:
   ```bash
   # Add to Cargo.toml
   [profile.release]
   lto = "thin"
   ```

## Size Optimization

The static library can be quite large (50-100 MB). To reduce size:

```toml
# In Cargo.toml
[profile.release]
opt-level = "z"     # Optimize for size
lto = true          # Enable Link Time Optimization
codegen-units = 1   # Better optimization
strip = true        # Strip symbols
```

Then rebuild:
```bash
cargo build --release
```

This can reduce the binary size by 30-50%.
