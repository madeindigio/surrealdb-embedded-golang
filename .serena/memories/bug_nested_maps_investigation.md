# Investigation: Nested Maps Bug in surrealdb-embedded

## Summary

Investigated the reported bug where nested maps passed as query parameters were supposedly not being serialized correctly. **Finding: The surrealdb-embedded library works correctly - the bug must be elsewhere in the application stack.**

## Tests Performed

Created comprehensive test cases in `test_debug/test_nested_map_bug.go` that test:

1. **Simple nested metadata map** (3 fields: last_modified, source, total_size)
2. **Complex nested metadata** (5 fields with mixed types including nested maps)
3. **Single-key nested map** (edge case that could trigger unwrapping bugs)

### Results

**All tests PASSED ✓**

- Nested maps are correctly serialized from Go to Rust via JSON
- Parameters are correctly bound to SurrealDB queries
- Nested maps are correctly stored in the database
- Nested maps are correctly retrieved and deserialized back to Go

## How It Works

The flow is:

1. **Go → Rust**: Parameters marshaled to JSON (line 157 in `surrealdb.go`)
2. **Rust Parsing**: JSON deserialized to `HashMap<String, serde_json::Value>`
3. **SurrealDB Binding**: Values passed to `.bind((key, value))`
4. **Storage**: SurrealDB stores as `Object(Object({...}))` wrapped format
5. **Retrieval**: Results wrapped in `Array(Array([Object(Object({...}))]))`
6. **Unwrapping**: `unwrap_surrealdb_tagged()` recursively removes SurrealDB type wrappers
7. **Rust → Go**: Final JSON serialized and returned

## The Unwrapping Logic

The `unwrap_surrealdb_tagged()` function correctly handles:

- `{"Array": [...]}` → `[...]`
- `{"Object": {...}}` → `{...}`
- `{"Strand": "str"}` → `"str"`
- `{"Number": {"Int": 42}}` → `42`
- `{"Number": {"Float": 3.14}}` → `3.14`
- `{"Bool": true}` → `true`
- `{"Thing": {...}}` → Record ID object
- Regular objects with multiple keys → Recursively unwrap values

The logic correctly distinguishes between SurrealDB type wrappers (single-key objects with known keys) and user data (regular objects).

## Conclusion

The surrealdb-embedded library correctly handles nested maps. If the remembrances-mcp application is seeing empty metadata, the issue must be in:

1. **The queryEmbedded wrapper** - Type conversions or map copying
2. **Application logic** - Metadata being modified after query
3. **Different code path** - Perhaps using a different method (Create vs Query)
4. **Timing issue** - Race condition or map reference being cleared

## Recommendation

Check the `remembrances-mcp` application's `queryEmbedded()` function and how it processes the results from `embeddedDB.Query()`. The bug is NOT in the surrealdb-embedded library itself.

## Files Modified

- `surrealdb_embedded_rs/src/lib.rs` - Cleaned up debug logging (no functional changes)
- `test_debug/test_nested_map_bug.go` - Comprehensive test suite (new file)
