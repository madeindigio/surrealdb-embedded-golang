# Newline Serialization Debug Tool

This tool helps debug issues with newline characters (`\n`) in JSON serialization/deserialization during UPDATE and MERGE operations in surrealdb-embedded.

## Purpose

If you're experiencing issues where newlines appear to be lost or corrupted during UPDATE operations, this tool will help you:

1. Verify that the library itself works correctly
2. Identify where in your application the newlines are being lost
3. See byte-level representations to confirm exact values

## Usage

### Prerequisites

1. Build the Rust library:
```bash
cd surrealdb_embedded_rs
cargo build --release
cd ..
```

2. Set library path:
```bash
export LD_LIBRARY_PATH=$PWD/surrealdb_embedded_rs/target/release:$LD_LIBRARY_PATH
```

### Run the Debug Tool

```bash
go run cmd/debug-newlines/main.go
```

## What It Tests

The tool runs 5 comprehensive tests:

### Test 1: CREATE with newlines
- Creates a record with newline characters in multiple fields
- Verifies the data is stored correctly
- Shows: `"Line 1\nLine 2\nLine 3"`

### Test 2: UPDATE with newlines
- Updates a record, replacing all fields with new values containing newlines
- Verifies all fields are updated correctly
- Shows: Multiple fields with various newline patterns

### Test 3: MERGE with newlines
- Partially updates a record (MERGE operation)
- Verifies only specified fields are updated
- Verifies other fields remain unchanged
- Shows: Selective field updates

### Test 4: Parameterized query with newlines
- Uses parameterized queries: `UPDATE ... SET field = $param`
- Verifies parameters containing newlines work correctly
- Shows: Parameter binding with newlines

### Test 5: Edge cases
- Empty string
- Single newline: `"\n"`
- Only newlines: `"\n\n\n"`
- Leading newline: `"\nstarts here"`
- Trailing newline: `"ends here\n"`
- Tabs and newlines: `"tab:\there\nnewline\nhere"`
- Unicode with newlines: `"你好\n世界\n😀"`

For each edge case, the tool shows:
- The value in escaped format: `%q`
- The raw bytes: `[10, 115, 116, ...]`
- Verification that retrieved value matches expected

## Expected Output

If everything works correctly (which it should), you'll see:

```
=== SurrealDB Newline Serialization Debug Tool ===

✓ Database initialized

--- Test 1: CREATE with newlines ---
Creating with data:
  {
    "content": "First paragraph\n\nSecond paragraph with blank line above",
    "description": "Line 1\nLine 2\nLine 3",
    "title": "Test Document"
  }
...
✓ Field 'description' verified: "Line 1\nLine 2\nLine 3"

... (more tests) ...

=== All tests PASSED ✓ ===
```

## Interpreting Results

### If Tests Pass

The surrealdb-embedded library is working correctly. Your issue is likely in:

1. **Data Preparation**: Check how you prepare data before calling `Update()`
   ```go
   // ❌ Wrong - double encoded
   data := `{"text":"line1\\nline2"}`
   
   // ✅ Correct - use Go map
   data := map[string]interface{}{"text": "line1\nline2"}
   ```

2. **Display/Logging**: Check how you display the data
   ```go
   // This shows escaped version
   fmt.Printf("Value: %q\n", value)  // "line1\nline2"
   
   // This shows actual newlines (might look empty in logs)
   fmt.Printf("Value: %s\n", value)  // line1
                                      // line2
   ```

3. **Wrapper Functions**: Check if you have wrapper/middleware code that processes the data

4. **Race Conditions**: Check for concurrent modifications

### If Tests Fail

This would indicate a problem with the library (unlikely based on comprehensive testing), but would help identify:
- Which exact operation fails (CREATE, UPDATE, MERGE, parameterized)
- Which exact pattern fails (simple newline, multiple, edge case)
- The byte-level difference between expected and actual

## Adding Your Own Test Case

To test your specific scenario, modify the `main()` function:

```go
// Add after existing tests
fmt.Println("\n--- Custom Test ---")
myData := map[string]interface{}{
    "field": "your\nvalue\nhere",
}

result, err := db.Update("documents:custom", myData)
if err != nil {
    log.Fatalf("Custom test failed: %v", err)
}

printData(result)
verifyNewlines(result.([]interface{})[0].(map[string]interface{}), "field", "your\nvalue\nhere")
```

## Debug Output Format

The tool uses several output formats for clarity:

1. **JSON with indent**: For human readability
   ```json
   {
     "field": "value\nwith newline"
   }
   ```

2. **Escaped format** (`%q`): Shows escape sequences
   ```
   "value\nwith newline"
   ```

3. **Byte array**: Shows exact byte values
   ```
   [118 97 108 117 101 10 119 105 116 104 32 110 101 119 108 105 110 101]
   ```
   - Byte 10 is newline (`\n`)
   - Byte 32 is space
   - Byte 116 is 't'

## Related Files

- **newline_update_test.go**: Automated test suite (50+ test cases)
- **NEWLINE_SERIALIZATION_TEST_REPORT.md**: Detailed test report
- **surrealdb.go**: Main library code
- **surrealdb_embedded_rs/src/lib.rs**: Rust FFI layer

## Running Automated Tests

Instead of this interactive tool, you can run the automated test suite:

```bash
export LD_LIBRARY_PATH=$PWD/surrealdb_embedded_rs/target/release:$LD_LIBRARY_PATH
go test -v -run TestUpdateWithNewlines
```

This runs 50+ test cases covering all scenarios.

## Common Issues and Solutions

### Issue: "Newlines don't appear in my logs"

**Cause**: Terminal output might render newlines as actual line breaks

**Solution**: Use escaped format:
```go
fmt.Printf("Value: %q\n", value)  // Shows "line1\nline2"
```

### Issue: "JSON shows \\n instead of \n"

**Cause**: Double encoding - JSON is being encoded twice

**Solution**: Don't pre-encode your data
```go
// ❌ Wrong
jsonStr := `{"text":"line1\\nline2"}`
db.Update(id, jsonStr)

// ✅ Correct
data := map[string]interface{}{"text": "line1\nline2"}
db.Update(id, data)
```

### Issue: "UPDATE doesn't seem to change the value"

**Cause**: Might be updating wrong record or checking wrong field

**Solution**: Add verification query immediately after:
```go
result, _ := db.Update("test:id", data)
log.Printf("Update returned: %+v", result)

verify, _ := db.Query("SELECT * FROM test:id", nil)
log.Printf("Verification: %+v", verify)
```

## Support

If this tool shows tests passing but you still have issues:

1. Compare your code to the working examples in this tool
2. Add similar debug output to your application
3. Check the byte representation of your data at each step
4. Look for middleware or wrapper code that might modify the data

The library is proven to work correctly - the issue will be in how it's being used.