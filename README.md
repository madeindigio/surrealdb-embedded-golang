# SurrealDB Embedded - Go Wrapper with Multi-Backend Support

A Go wrapper for embedded SurrealDB with support for multiple storage backends.

## Features

- 🚀 **Multiple Storage Backends**:
  - **Memory**: Fast in-memory storage for testing and temporary data
  - **RocksDB**: High-performance persistent storage
  - **SurrealKV**: Native SurrealDB storage engine (fastest persistent option)

- 🔧 **Flexible Initialization**: Three ways to create database instances
- 🔄 **Full CRUD Operations**: Create, Read, Update, Delete
- 📊 **SurrealQL Support**: Execute complex queries with parameters
- 🔐 **Type-Safe**: Strongly typed Go API
- ⚡ **High Performance**: Optimized native Rust implementation

## Installation

### Prerequisites

- Go 1.16 or higher
- Rust toolchain (for building the native library)
- For RocksDB: `clang`, `cmake`, and standard build tools

### Build

```bash
# Clone the repository
git clone <your-repo-url>
cd surrealdb-embedded

# Build the Rust library
cd surrealdb_embedded_rs
cargo build --release
cd ..
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
    // Create a SurrealKV database
    db, err := surrealdb.NewFromURL("surrealkv://./data/mydb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Use namespace and database
    if err := db.Use("myapp", "production"); err != nil {
        log.Fatal(err)
    }
    
    // Create a record
    person := map[string]interface{}{
        "name": "Alice",
        "age":  30,
        "email": "alice@example.com",
    }
    
    result, err := db.Create("person:alice", person)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Created: %+v\n", result)
    
    // Query records
    people, err := db.Query("SELECT * FROM person WHERE age > 25", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("People: %+v\n", people)
}
```

## Usage

### Method 1: URL-based Initialization (Recommended)

```go
// Memory backend
db, err := surrealdb.NewFromURL("memory")

// RocksDB backend
db, err := surrealdb.NewFromURL("rocksdb://./data/mydb")

// SurrealKV backend
db, err := surrealdb.NewFromURL("surrealkv://./data/mydb")
```

### Method 2: Specific Constructors

```go
// Memory
db, err := surrealdb.NewMemory()

// RocksDB
db, err := surrealdb.NewRocksDB("./data/mydb")

// SurrealKV
db, err := surrealdb.NewSurrealKV("./data/mydb")
```

### Method 3: Config-based

```go
// Memory
db, err := surrealdb.New(surrealdb.Config{
    Backend: surrealdb.Memory,
})

// RocksDB
db, err := surrealdb.New(surrealdb.Config{
    Backend: surrealdb.RocksDB,
    Path:    "./data/mydb",
})

// SurrealKV
db, err := surrealdb.New(surrealdb.Config{
    Backend: surrealdb.SurrealKV,
    Path:    "./data/mydb",
})
```

## API Reference

### Database Operations

#### Use Namespace and Database
```go
err := db.Use("namespace", "database")
```

#### Create Records
```go
data := map[string]interface{}{
    "name": "John",
    "age": 25,
}
result, err := db.Create("person:john", data)
```

#### Select Records
```go
// Select specific record
person, err := db.Select("person:john")

// Select all records from table
people, err := db.Select("person")
```

#### Update Records
```go
// Replace entire record
data := map[string]interface{}{
    "name": "John Doe",
    "age": 26,
}
result, err := db.Update("person:john", data)
```

#### Merge/Patch Records
```go
// Partial update
updates := map[string]interface{}{
    "age": 27,
}
result, err := db.Merge("person:john", updates)
```

#### Delete Records
```go
// Delete specific record
result, err := db.Delete("person:john")

// Delete all records from table
result, err := db.Delete("person")
```

#### Insert Records
```go
data := map[string]interface{}{
    "id": "person:jane",
    "name": "Jane",
    "age": 28,
}
result, err := db.Insert("person", data)
```

#### Upsert Records
```go
data := map[string]interface{}{
    "name": "Bob",
    "age": 30,
}
result, err := db.Upsert("person:bob", data)
```

### Query Operations

#### Execute SurrealQL Queries
```go
// Simple query
results, err := db.Query("SELECT * FROM person", nil)

// Query with parameters
params := map[string]interface{}{
    "min_age": 25,
}
results, err := db.Query("SELECT * FROM person WHERE age >= $min_age", params)
```

#### Complex Queries
```go
// With ORDER BY
results, err := db.Query("SELECT * FROM person ORDER BY age DESC", nil)

// With aggregation
results, err := db.Query("SELECT count() AS total FROM person GROUP ALL", nil)

// With relationships
results, err := db.Query("SELECT *, ->knows->person AS friends FROM person:john", nil)
```

### Other Operations

#### Get Version
```go
version, err := db.Version()
```

#### Close Database
```go
err := db.Close()
```

## Running Tests

```bash
# Set library path
export LD_LIBRARY_PATH="$(pwd)/surrealdb_embedded_rs/target/release:$LD_LIBRARY_PATH"

# Run all tests
go test -v

# Run specific test
go test -v -run TestSurrealKVBackend

# Run with timeout
go test -v -timeout 120s
```

## Examples

See the `examples/` directory for complete examples:

- `simple_surrealkv.go` - Basic SurrealKV usage
- `test_backends.go` - Comprehensive backend testing
- `debug_query.go` - Query debugging

To run an example:
```bash
cd examples
LD_LIBRARY_PATH="../surrealdb_embedded_rs/target/release:$LD_LIBRARY_PATH" go run simple_surrealkv.go
```

## Backend Comparison

| Feature | Memory | RocksDB | SurrealKV |
|---------|--------|---------|-----------|
| Persistence | ❌ No | ✅ Yes | ✅ Yes |
| Speed | ⚡ Fastest | 🚀 Fast | 🚀 Very Fast |
| Memory Usage | High | Low | Low |
| Setup Complexity | None | Medium | Low |
| Production Ready | Testing only | ✅ Yes | ✅ Yes |
| Dependencies | None | C++ libs | None |

### When to Use Each Backend

**Memory**:
- Unit tests
- Temporary data
- Development
- Prototyping

**RocksDB**:
- Production applications
- High write throughput
- Battle-tested reliability
- Large datasets

**SurrealKV**:
- Production applications
- Best performance for SurrealDB
- Pure Rust (smaller binary)
- Modern alternative to RocksDB

## Performance Tips

1. **Use SurrealKV for production** - It's optimized for SurrealDB and offers the best performance
2. **Use Memory for testing** - Fast and doesn't require cleanup
3. **Use parameterized queries** - More secure and can be cached
4. **Batch operations** when possible
5. **Close databases** when done to free resources

## Troubleshooting

### Symbol lookup error
```
error: undefined symbol: surreal_init
```
**Solution**: Make sure `LD_LIBRARY_PATH` includes the library directory:
```bash
export LD_LIBRARY_PATH="$(pwd)/surrealdb_embedded_rs/target/release:$LD_LIBRARY_PATH"
```

### Path errors
```
error: path is required for X backend
```
**Solution**: File-based backends (RocksDB, SurrealKV) require a path:
```go
db, err := surrealdb.NewSurrealKV("./data/mydb") // ✅ Correct
db, err := surrealdb.NewSurrealKV("")            // ❌ Error
```

### Compilation errors
```
error: failed to compile surrealdb_embedded_rs
```
**Solution**: Make sure you have all dependencies:
```bash
# For RocksDB support
sudo apt-get install clang cmake librocksdb-dev

# Rebuild
cd surrealdb_embedded_rs
cargo clean
cargo build --release
```

## Documentation

- [BACKENDS.md](BACKENDS.md) - Detailed backend documentation
- [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) - Technical implementation details
- [TEST_RESULTS.md](TEST_RESULTS.md) - Test results and coverage

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Your License Here]

## Acknowledgments

- [SurrealDB](https://surrealdb.com/) - The database powering this library
- [RocksDB](https://rocksdb.org/) - High-performance storage engine
- [SurrealKV](https://github.com/surrealdb/surrealkv) - Native SurrealDB storage

## Support

For issues and questions:
- Open an issue on GitHub
- Check existing documentation
- Review test examples

---

**Status**: ✅ Production Ready  
**Version**: 1.0.0  
**Last Updated**: 2025-11-21
