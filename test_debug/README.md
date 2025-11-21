# Nested Map Bug Test Suite

This test suite verifies that the surrealdb-embedded library correctly handles nested maps in query parameters.

## Background

A bug was reported where nested `map[string]interface{}` values passed as query parameters were supposedly not being serialized/deserialized correctly, resulting in empty maps in the database.

## Test Results

**All tests PASS ✓** - The surrealdb-embedded library correctly handles nested maps.

## Running the Tests

```bash
# From the surrealdb-embedded directory
cd test_debug
LD_LIBRARY_PATH=../surrealdb_embedded_rs/target/release:$LD_LIBRARY_PATH go run test_nested_map_bug.go
```

To see only the test results without debug output:

```bash
LD_LIBRARY_PATH=../surrealdb_embedded_rs/target/release:$LD_LIBRARY_PATH go run test_nested_map_bug.go 2>&1 | grep -E "(^===|^Test|^---|PASS|FAIL|metadata field:)"
```

## Test Cases

### Test 1: Simple Nested Metadata Map

Tests the exact pattern from the bug report:

```go
metadata := map[string]interface{}{
    "last_modified": "2025-11-18T12:58:01+01:00",
    "source":        "watcher",
    "total_size":    8958,
}
```

**Result**: ✓ PASS - Metadata is correctly stored and retrieved with all 3 fields.

### Test 2: Complex Nested Metadata

Tests a more complex nested structure with multiple data types:

```go
complexMeta := map[string]interface{}{
    "string_field": "test",
    "int_field":    42,
    "float_field":  3.14,
    "bool_field":   true,
    "nested": map[string]interface{}{
        "inner": "value",
    },
}
```

**Result**: ✓ PASS - All 5 fields including nested map are correctly preserved.

### Test 3: Single-Key Nested Map

Tests an edge case where a nested map has only one key (which could trigger incorrect unwrapping logic):

```go
singleKeyMeta := map[string]interface{}{
    "only_field": "single value",
}
```

**Result**: ✓ PASS - Single-key map is correctly distinguished from SurrealDB type wrappers.

## Conclusion

The surrealdb-embedded library works correctly. If you're experiencing issues with empty nested maps in your application:

1. **Check your application's wrapper functions** - The issue may be in how you're processing results
2. **Verify type assertions** - Ensure you're not losing data during type conversions
3. **Check for map reference issues** - Ensure maps aren't being modified after queries
4. **Use different code paths** - Try using `db.Query()` directly vs wrapped helpers

## Files

- `test_nested_map_bug.go` - Test suite
- `go.mod` - Go module with local replace directive for surrealdb-embedded

## Building

The tests use the compiled Rust library from `../surrealdb_embedded_rs/target/release`. If you've made changes to the Rust code:

```bash
cd ../surrealdb_embedded_rs
cargo build --release
cd ../test_debug
```

Then run the tests again.