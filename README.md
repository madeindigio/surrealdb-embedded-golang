# SurrealDB Embedded for Go

An embedded SurrealDB implementation for Go using a Rust FFI layer. This library provides full SurrealDB functionality with both in-memory and RocksDB persistent storage backends, without requiring a separate SurrealDB server.

## Features

- **Embedded Database**: No separate server required
- **Multiple Backends**:
  - **Memory**: Fast in-memory storage for development and testing
  - **RocksDB**: Persistent disk-based storage for production
- **Full SurrealQL Support**: Execute any SurrealQL query
- **CRUD Operations**: Create, Read, Update, Delete operations
- **Type-Safe**: Strongly-typed Go API
- **Zero Network Overhead**: Direct FFI calls to Rust

## Why This Library?

The official Go SDK for SurrealDB (`github.com/surrealdb/surrealdb.go`) does not support embedded mode. It only works with remote SurrealDB servers via WebSocket or HTTP. This library fills that gap by providing embedded SurrealDB functionality using the Rust SDK, which has full embedded support.

## Embedded Mode Limitations

**Important**: This library provides **embedded mode only**. The following features are NOT available because they require a SurrealDB server:

❌ **Not Available in Embedded Mode:**
- Authentication (SignIn, SignUp, JWT tokens)
- Live Queries (WebSocket-based real-time updates)
- Remote connections
- User management
- Access control and permissions

✅ **Available in Embedded Mode:**
- Full SurrealQL queries
- CRUD operations
- Transactions
- Graph relations
- Schema definitions
- Indexes
- Data persistence (with RocksDB)

## Architecture

```
┌─────────────────┐
│   Go Application│
└────────┬────────┘
         │
         │ CGo FFI
         ▼
┌─────────────────┐
│  Rust C Library │
│  (Static .a)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  SurrealDB Rust │
│      SDK        │
└────────┬────────┘
         │
    ┌────┴─────┐
    ▼          ▼
┌────────┐  ┌─────────┐
│ Memory │  │ RocksDB │
└────────┘  └─────────┘
```

## Installation

### Prerequisites

1. **Rust toolchain** (for building the static library):
   ```bash
   curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
   ```

2. **System dependencies**:
   ```bash
   # Ubuntu/Debian
   sudo apt-get install -y build-essential clang libclang-dev llvm-dev pkg-config libssl-dev
   
   # macOS
   brew install llvm
   export LIBCLANG_PATH="$(brew --prefix llvm)/lib"
   ```

3. **Go 1.21+**

### Build the Rust Library

```bash
cd surrealdb_embedded_rs
cargo build --release
```

This will create:
- Static library: `target/release/libsurrealdb_embedded_rs.a`
- C header: `include/surrealdb_embedded_rs.h`

### Install the Go Module

```bash
go get github.com/yourusername/surrealdb-embedded
```

## Quick Start

### In-Memory Database

```go
package main

import (
    "fmt"
    "log"
    
    surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
    // Create an in-memory database
    db, err := surrealdb.NewMemory()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Select namespace and database
    if err := db.Use("test", "test"); err != nil {
        log.Fatal(err)
    }
    
    // Create a record
    person := map[string]interface{}{
        "name": "John Doe",
        "age":  30,
        "email": "john@example.com",
    }
    
    result, err := db.Create("person", person)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Created: %+v\n", result)
    
    // Query records
    results, err := db.Query("SELECT * FROM person WHERE age > $age", 
        map[string]interface{}{"age": 25})
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Results: %+v\n", results)
}
```

### Persistent RocksDB Database

```go
package main

import (
    "log"
    
    surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
    // Create a persistent database
    db, err := surrealdb.NewRocksDB("./data/mydb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Use namespace and database
    if err := db.Use("app", "production"); err != nil {
        log.Fatal(err)
    }
    
    // Your data persists across restarts
    // ...
}
```

## API Reference

### Database Operations

#### `NewMemory() (*DB, error)`
Create a new in-memory database instance.

#### `NewRocksDB(path string) (*DB, error)`
Create a new RocksDB-backed database instance with persistent storage.

#### `New(config Config) (*DB, error)`
Create a new database instance with custom configuration.

#### `Use(namespace, database string) error`
Select a namespace and database to use.

#### `Close() error`
Close the database and free resources.

### CRUD Operations

#### `Create(resource string, data interface{}) (interface{}, error)`
Create a new record.

```go
result, err := db.Create("person", map[string]interface{}{
    "name": "Alice",
    "age":  28,
})
```

#### `Select(resource string) (interface{}, error)`
Select all records from a table or a specific record.

```go
// Select all
allPeople, err := db.Select("person")

// Select specific record
person, err := db.Select("person:alice")
```

#### `Update(resource string, data interface{}) (interface{}, error)`
Replace an entire record.

```go
result, err := db.Update("person:alice", map[string]interface{}{
    "name": "Alice Smith",
    "age":  29,
})
```

#### `Merge(resource string, data interface{}) (interface{}, error)`
Partially update a record.

```go
result, err := db.Merge("person:alice", map[string]interface{}{
    "email": "alice@example.com",
})
```

#### `Delete(resource string) (interface{}, error)`
Delete records.

```go
result, err := db.Delete("person:alice")
```

#### `Insert(table string, data interface{}) (interface{}, error)`
Insert one or more records with specified IDs.

```go
result, err := db.Insert("person", []map[string]interface{}{
    {"id": "person:bob", "name": "Bob"},
    {"id": "person:carol", "name": "Carol"},
})
```

#### `Upsert(resource string, data interface{}) (interface{}, error)`
Create or update a record.

```go
result, err := db.Upsert("person:dave", map[string]interface{}{
    "name": "Dave",
    "age":  35,
})
```

### Query Operations

#### `Query(query string, vars map[string]interface{}) ([]interface{}, error)`
Execute a SurrealQL query with optional variables.

```go
// Simple query
results, err := db.Query("SELECT * FROM person", nil)

// Query with variables
results, err := db.Query(
    "SELECT * FROM person WHERE age > $minAge AND age < $maxAge",
    map[string]interface{}{
        "minAge": 25,
        "maxAge": 50,
    },
)

// Complex query with transactions
results, err := db.Query(`
    BEGIN TRANSACTION;
    CREATE person:john SET name = 'John', age = 30;
    CREATE person:jane SET name = 'Jane', age = 28;
    COMMIT TRANSACTION;
`, nil)
```

### Utility

#### `Version() (map[string]interface{}, error)`
Get SurrealDB version information.

```go
version, err := db.Version()
fmt.Printf("Version: %s\n", version["version"])
```

## Advanced Examples

### Relations and Graph Queries

```go
// Create nodes
db.Create("person:alice", map[string]interface{}{"name": "Alice"})
db.Create("person:bob", map[string]interface{}{"name": "Bob"})

// Create relationship
db.Query("RELATE person:alice->knows->person:bob SET since = '2020-01-01'", nil)

// Query graph
results, err := db.Query("SELECT *, ->knows->person.* as friends FROM person:alice", nil)
```

### Transactions

```go
results, err := db.Query(`
    BEGIN TRANSACTION;
    
    LET $alice = CREATE person SET name = 'Alice', balance = 100;
    LET $bob = CREATE person SET name = 'Bob', balance = 100;
    
    UPDATE $alice SET balance = balance - 50;
    UPDATE $bob SET balance = balance + 50;
    
    COMMIT TRANSACTION;
`, nil)
```

### Schema Definitions

```go
db.Query(`
    DEFINE TABLE person SCHEMAFULL;
    
    DEFINE FIELD name ON person TYPE string;
    DEFINE FIELD age ON person TYPE int;
    DEFINE FIELD email ON person TYPE string
        ASSERT string::is::email($value);
    
    DEFINE INDEX idx_email ON person COLUMNS email UNIQUE;
`, nil)
```

## Testing

Run the tests:

```bash
# Run all tests
go test -v

# Run specific test
go test -v -run TestCreate

# Run benchmarks
go test -bench=. -benchmem
```

## Error Handling

All functions return standard Go errors. Check errors appropriately:

```go
result, err := db.Create("person", data)
if err != nil {
    if errors.Is(err, surrealdb.ErrQueryFailed) {
        // Handle query failure
    }
    log.Fatal(err)
}
```

## Thread Safety

The library is thread-safe. Multiple goroutines can safely use the same `*DB` instance concurrently.

## Limitations

1. **No Authentication**: Embedded mode does not support user authentication
2. **No Live Queries**: WebSocket-based live queries not available
3. **No Remote Connections**: This library is for embedded use only
4. **Platform-Specific Builds**: The Rust library must be compiled for your target platform

## Comparison with Official Go SDK

| Feature | This Library | Official Go SDK |
|---------|-------------|-----------------|
| Embedded Mode | ✅ Yes | ❌ No |
| Memory Backend | ✅ Yes | ❌ No |
| RocksDB Backend | ✅ Yes | ❌ No |
| Remote Connection | ❌ No | ✅ Yes |
| WebSocket | ❌ No | ✅ Yes |
| HTTP | ❌ No | ✅ Yes |
| Authentication | ❌ No | ✅ Yes |
| Live Queries | ❌ No | ✅ Yes |
| SurrealQL | ✅ Yes | ✅ Yes |
| Transactions | ✅ Yes | ✅ Yes |
| Graph Queries | ✅ Yes | ✅ Yes |

Use this library when you need embedded SurrealDB. Use the official SDK when you need to connect to remote SurrealDB servers.

## Contributing

Contributions are welcome! Please ensure all tests pass before submitting a PR.

## License

Apache 2.0 License - see LICENSE file for details.

## Acknowledgments

- [SurrealDB](https://surrealdb.com/) - The amazing multi-model database
- [SurrealDB Rust SDK](https://github.com/surrealdb/surrealdb) - The Rust implementation we wrap

## Support

For issues and questions:
- Open an issue on GitHub
- Check the [SurrealDB documentation](https://surrealdb.com/docs)
- Join the [SurrealDB Discord](https://discord.gg/surrealdb)
