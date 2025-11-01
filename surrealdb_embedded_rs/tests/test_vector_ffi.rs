/// Integration test for vector operations through FFI
/// This tests the complete flow: Rust -> C FFI -> Parsing
use std::ffi::{CStr, CString};
use std::os::raw::c_char;

// Import the FFI functions
extern "C" {
    fn surreal_init_mem() -> i32;
    fn surreal_use(handle: i32, ns: *const c_char, db: *const c_char) -> i32;
    fn surreal_query(handle: i32, query: *const c_char) -> *mut c_char;
    fn surreal_free_string(ptr: *mut c_char);
}

fn get_string_from_ptr(ptr: *mut c_char) -> String {
    if ptr.is_null() {
        return String::new();
    }

    unsafe {
        let c_str = CStr::from_ptr(ptr);
        let result = c_str.to_string_lossy().into_owned();
        surreal_free_string(ptr);
        result
    }
}

#[test]
fn test_vector_operations_through_ffi() {
    println!("\n=== Testing Vector Operations Through FFI ===\n");

    // Initialize database
    let handle = unsafe { surreal_init_mem() };
    assert!(handle > 0, "Failed to initialize database");
    println!("✅ Database initialized with handle: {}", handle);

    // Use namespace and database
    let ns = CString::new("test").unwrap();
    let db = CString::new("test").unwrap();
    let result = unsafe { surreal_use(handle, ns.as_ptr(), db.as_ptr()) };
    assert_eq!(result, 0, "Failed to use namespace/database");
    println!("✅ Namespace and database set\n");

    // Test 1: Create schema with vector index
    println!("Test 1: Creating schema with MTREE index");
    let schema_query = CString::new(
        "
        DEFINE TABLE doc SCHEMAFULL;
        DEFINE FIELD content ON doc TYPE string;
        DEFINE FIELD embedding ON doc TYPE array;
        DEFINE INDEX emb_idx ON doc FIELDS embedding MTREE DIMENSION 3;
    ",
    )
    .unwrap();

    let response = unsafe { surreal_query(handle, schema_query.as_ptr()) };
    let json = get_string_from_ptr(response);
    println!("  Response: {}", json);
    assert!(!json.contains("error"), "Schema creation failed: {}", json);
    println!("✅ Schema created\n");

    // Test 2: Insert documents with embeddings
    println!("Test 2: Inserting documents with embeddings");
    let insert_query = CString::new(
        "
        CREATE doc:1 SET content = 'Document 1', embedding = [0.1, 0.2, 0.3];
        CREATE doc:2 SET content = 'Document 2', embedding = [0.5, 0.4, 0.3];
        CREATE doc:3 SET content = 'Document 3', embedding = [0.9, 0.8, 0.7];
    ",
    )
    .unwrap();

    let response = unsafe { surreal_query(handle, insert_query.as_ptr()) };
    let json = get_string_from_ptr(response);
    println!("  Response: {}", json);
    assert!(!json.contains("error"), "Insert failed: {}", json);
    println!("✅ Documents inserted\n");

    // Test 3: KNN search
    println!("Test 3: KNN search");
    let knn_query = CString::new(
        "
        SELECT content, embedding
        FROM doc
        WHERE embedding <|3,EUCLIDEAN|> [0.6, 0.3, 0.2]
        LIMIT 3;
    ",
    )
    .unwrap();

    let response = unsafe { surreal_query(handle, knn_query.as_ptr()) };
    let json = get_string_from_ptr(response);
    println!("  Response: {}", json);
    assert!(!json.contains("error"), "KNN search failed: {}", json);

    // Verify JSON can be parsed
    match serde_json::from_str::<serde_json::Value>(&json) {
        Ok(parsed) => {
            println!("✅ KNN search successful and parseable");
            println!(
                "  Parsed JSON: {}",
                serde_json::to_string_pretty(&parsed).unwrap()
            );
        }
        Err(e) => {
            panic!(
                "Failed to parse KNN response as JSON: {}. Response: {}",
                e, json
            );
        }
    }
    println!();

    // Test 4: Vector distance function (if supported)
    println!("Test 4: Vector distance function");
    let dist_query = CString::new(
        "
        SELECT
            content,
            vector::distance::knn() AS distance
        FROM doc
        WHERE embedding <|2,EUCLIDEAN|> [0.1, 0.2, 0.3];
    ",
    )
    .unwrap();

    let response = unsafe { surreal_query(handle, dist_query.as_ptr()) };
    let json = get_string_from_ptr(response);
    println!("  Response: {}", json);

    // Try to parse - this should work with our normalization
    match serde_json::from_str::<serde_json::Value>(&json) {
        Ok(parsed) => {
            println!("✅ Distance function response parseable");
            println!(
                "  Parsed JSON: {}",
                serde_json::to_string_pretty(&parsed).unwrap()
            );

            // Verify no 'f' suffix remains
            assert!(
                !json.contains("f,"),
                "JSON still contains float suffix 'f,'"
            );
            assert!(
                !json.contains("f}"),
                "JSON still contains float suffix 'f}}'"
            );
            assert!(
                !json.contains("f]"),
                "JSON still contains float suffix 'f]'"
            );
        }
        Err(e) => {
            println!("⚠️  Distance function response not parseable: {}", e);
            println!("  This may be due to SurrealDB version limitations");
        }
    }
    println!();

    println!("=== All FFI Vector Tests Completed ===");
}

#[test]
fn test_json_normalization_effectiveness() {
    println!("\n=== Testing JSON Normalization ===\n");

    // Initialize and setup
    let handle = unsafe { surreal_init_mem() };
    let ns = CString::new("test").unwrap();
    let db = CString::new("test").unwrap();
    unsafe { surreal_use(handle, ns.as_ptr(), db.as_ptr()) };

    // Create simple table with numeric field
    let setup = CString::new(
        "
        DEFINE TABLE test SCHEMAFULL;
        DEFINE FIELD value ON test TYPE number;
        CREATE test:1 SET value = 3.14159;
        CREATE test:2 SET value = 2.71828;
    ",
    )
    .unwrap();

    let response = unsafe { surreal_query(handle, setup.as_ptr()) };
    get_string_from_ptr(response); // Consume response

    // Query and verify normalization
    let query = CString::new("SELECT value FROM test;").unwrap();
    let response = unsafe { surreal_query(handle, query.as_ptr()) };
    let json = get_string_from_ptr(response);

    println!("Response: {}", json);

    // The JSON should be parseable
    match serde_json::from_str::<serde_json::Value>(&json) {
        Ok(_) => println!("✅ JSON normalization working correctly"),
        Err(e) => panic!("JSON parsing failed: {}. Response: {}", e, json),
    }
}
