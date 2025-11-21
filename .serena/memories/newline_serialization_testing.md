# Newline Serialization Testing - Investigation Results

## Date
2024-11-21

## Issue Reported
User reported problems with JSON serialization/deserialization when values contain newlines (`\n`) during UPDATE operations. The claim was that new values weren't being updated correctly.

## Investigation Approach

Created comprehensive test suite to verify all aspects of newline handling:

### Test Files Created
1. **newline_update_test.go** - Comprehensive Go test suite
   - 6 test functions
   - 50+ individual test cases
   - Covers: CREATE, UPDATE, MERGE, parameterized queries, edge cases

2. **cmd/debug-newlines/main.go** - Interactive debugging tool
   - Demonstrates working behavior
   - Shows byte-level representation
   - Tests edge cases interactively

3. **NEWLINE_SERIALIZATION_TEST_REPORT.md** - Detailed test report

## Test Results

**ALL TESTS PASSED ✅**

### Test Coverage
- ✅ Simple newlines
- ✅ Multiple newlines
- ✅ Blank lines (consecutive newlines)
- ✅ Leading newlines
- ✅ Trailing newlines
- ✅ Only newlines
- ✅ Mixed whitespace (tabs + newlines)
- ✅ JSON-like strings with newlines
- ✅ Quotes + newlines
- ✅ Backslashes + newlines
- ✅ Unicode + newlines
- ✅ Carriage return + newline
- ✅ All escape characters combined
- ✅ Parameterized queries
- ✅ MERGE operations
- ✅ UPDATE operations
- ✅ Potential SQL-injection-like patterns

### Verified Operations
- `db.Create()` - ✅ Works correctly
- `db.Update()` - ✅ Works correctly
- `db.Merge()` - ✅ Works correctly
- `db.Query()` with parameters - ✅ Works correctly

## How It Works

### Data Flow for UPDATE
1. **Go → JSON**: `json.Marshal()` escapes newlines to `\n`
2. **Go → Rust (FFI)**: CString safely passes JSON
3. **Rust Query Building**: `format!("UPDATE {} CONTENT {}", resource, data)`
4. **SurrealDB Parsing**: Correctly interprets `\n` in JSON
5. **Storage**: Stores actual newline character (byte 10)
6. **Retrieval**: Reverse process preserves newlines

### Example Verification
```go
// Input
db.Update("test:1", map[string]interface{}{
    "text": "line1\nline2\nline3",
})

// Stored as bytes: [108, 105, 110, 101, 49, 10, 108, 105, 110, 101, 50, 10, ...]
//                   l    i    n    e    1    \n   l    i    n    e    2    \n

// Retrieved correctly
result["text"] == "line1\nline2\nline3" // ✅ True
```

## Conclusion

**NO BUG FOUND IN SURREALDB-EMBEDDED**

The library correctly handles newlines in all tested scenarios. If the user is experiencing issues, the problem must be in:

1. **Application layer** - Wrapper code, middleware
2. **Data preparation** - Double encoding, pre-escaped strings
3. **Display/logging** - Terminal not interpreting newlines
4. **Race conditions** - Concurrent modifications
5. **Different code path** - Using different methods

## Recommendations for Debugging User's Issue

### 1. Add Debug Logging
```go
log.Printf("Before update: %q", data)
result, err := db.Update(id, data)
log.Printf("After update: %q", result)

check, _ := db.Query("SELECT * FROM ...", nil)
log.Printf("Verification: %q", check[0])
```

### 2. Check Byte Representation
```go
value := result["field"].(string)
log.Printf("Bytes: %v", []byte(value))
// Should show [10] for newline character
```

### 3. Run Test Suite
```bash
export LD_LIBRARY_PATH=$PWD/surrealdb_embedded_rs/target/release:$LD_LIBRARY_PATH
go test -v -run TestUpdateWithNewlines
```

### 4. Run Debug Tool
```bash
go run cmd/debug-newlines/main.go
```

### 5. Check for Double Encoding
Bad (double-encoded):
```go
data := `{"text":"line1\\nline2"}` // Wrong - literal backslash-n
```

Good:
```go
data := map[string]interface{}{"text": "line1\nline2"} // Correct
```

## Files
- `newline_update_test.go` - Test suite
- `cmd/debug-newlines/main.go` - Debug tool
- `NEWLINE_SERIALIZATION_TEST_REPORT.md` - Full report

## Test Execution
```bash
# All tests pass
PASS: TestUpdateWithNewlines (0.03s)
PASS: TestUpdateWithNewlinesDetailed (0.00s)
PASS: TestUpdateJSONRoundtrip (0.01s)
PASS: TestUpdateWithParametersNewlines (0.00s)
PASS: TestUpdateDirectQueryInjection (0.00s)
PASS: TestUpdateMergeWithNewlines (0.00s)

# Debug tool output
=== All tests PASSED ✓ ===
```

## Key Insight

The surrealdb-embedded library passes JSON strings directly into SurrealDB queries like:
```rust
let query = format!("UPDATE {} CONTENT {}", resource, data);
```

This is SAFE because:
- The `data` is already valid JSON with properly escaped newlines (`\n`)
- SurrealDB parser correctly interprets the JSON
- No additional escaping is needed or wanted

The newlines are preserved through the entire round trip:
Go string → JSON → Rust → SurrealDB → Rust → JSON → Go string
