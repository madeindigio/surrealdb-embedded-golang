# SurrealDB Embedded - Multiple Backend Support

This library now supports multiple storage backends for SurrealDB embedded instances.

## Supported Backends

### 1. Memory (In-Memory)
- **URL Scheme**: `memory`
- **Use Case**: Testing, temporary data, fast operations
- **Persistence**: No (data is lost when the instance is closed)

### 2. RocksDB
- **URL Scheme**: `rocksdb://path` or `rocksdb:path`
- **Use Case**: Production, high-performance persistent storage
- **Persistence**: Yes (data is stored on disk)
- **Legacy**: Also supports `file://path` (deprecated)

### 3. SurrealKV
- **URL Scheme**: `surrealkv://path` or `surrealkv:path`
- **Use Case**: Modern alternative to RocksDB, native SurrealDB storage engine
- **Persistence**: Yes (data is stored on disk)
- **Note**: This is SurrealDB's own key-value store implementation

## Usage Examples

### Method 1: Using URL strings (Recommended)

```go
import surrealdb "github.com/yourusername/surrealdb-embedded"

// Memory backend
db, err := surrealdb.NewFromURL("memory")

// RocksDB backend
db, err := surrealdb.NewFromURL("rocksdb:///var/data/mydb")

// SurrealKV backend
db, err := surrealdb.NewFromURL("surrealkv:///var/data/mydb")
```

### Method 2: Using specific constructors

```go
// Memory backend
db, err := surrealdb.NewMemory()

// RocksDB backend
db, err := surrealdb.NewRocksDB("/var/data/mydb")

// SurrealKV backend
db, err := surrealdb.NewSurrealKV("/var/data/mydb")
```

### Method 3: Using Config struct

```go
// Memory backend
db, err := surrealdb.New(surrealdb.Config{
    Backend: surrealdb.Memory,
})

// RocksDB backend
db, err := surrealdb.New(surrealdb.Config{
    Backend: surrealdb.RocksDB,
    Path:    "/var/data/mydb",
})

// SurrealKV backend
db, err := surrealdb.New(surrealdb.Config{
    Backend: surrealdb.SurrealKV,
    Path:    "/var/data/mydb",
})
```

## Complete Example

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
        "name": "John Doe",
        "age":  30,
        "email": "john@example.com",
    }
    
    result, err := db.Create("person:john", person)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Created: %+v\n", result)
    
    // Query records
    people, err := db.Query("SELECT * FROM person", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("People: %+v\n", people)
}
```

## Performance Considerations

### Memory
- **Pros**: Fastest, no I/O overhead
- **Cons**: No persistence, limited by RAM

### RocksDB
- **Pros**: Battle-tested, good performance, compression
- **Cons**: C++ dependency, larger binary size

### SurrealKV
- **Pros**: Pure Rust, optimized for SurrealDB, smaller binary
- **Cons**: Newer, less battle-tested than RocksDB

## Migration Between Backends

To migrate data between backends, use the export/import functionality:

```go
// Export from one backend
sourceDB, _ := surrealdb.NewRocksDB("./old_data")
sourceDB.Use("myapp", "production")
// ... export logic ...

// Import to another backend
targetDB, _ := surrealdb.NewSurrealKV("./new_data")
targetDB.Use("myapp", "production")
// ... import logic ...
```

## Building

The library requires the following features to be enabled:
- `kv-mem`: In-memory storage
- `kv-rocksdb`: RocksDB storage (requires C++ compiler and RocksDB libraries)
- `kv-surrealkv`: SurrealKV storage

These are all enabled by default in the `Cargo.toml`.

### Build Requirements

- **RocksDB**: Requires `clang`, `cmake`, and standard build tools
- **SurrealKV**: Pure Rust, no additional dependencies

## Troubleshooting

### Error: "Unsupported database URL"
Make sure you're using one of the supported URL schemes: `memory`, `rocksdb://`, `surrealkv://`

### Error: "path is required for X backend"
File-based backends (RocksDB, SurrealKV) require a path to be specified.

### Error initializing database
Check that:
1. The path exists and is writable (for file-based backends)
2. You have sufficient permissions
3. The path is not already in use by another process

## Future Backends

The architecture supports adding more backends in the future:
- TiKV (distributed)
- FoundationDB
- SurrealCS (cloud storage)

To add support for these, update:
1. `Cargo.toml` - Add feature flags
2. `lib.rs` (Rust) - Add URL parsing and initialization
3. `surrealdb.go` (Go) - Add constants and helper functions
