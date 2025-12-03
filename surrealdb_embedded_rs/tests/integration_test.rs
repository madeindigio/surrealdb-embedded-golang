use surrealdb_embedded_rs::*;
use std::ffi::{CStr, CString};

#[test]
fn test_create_and_select() {
    unsafe {
        // Initialize memory DB
        let handle = surreal_init_memory();
        assert!(handle > 0, "Failed to initialize DB");

        // Use namespace and database
        let ns = CString::new("test").unwrap();
        let db = CString::new("test").unwrap();
        let result = surreal_use(handle, ns.as_ptr(), db.as_ptr());
        assert_eq!(result, 0, "Failed to use namespace/database");

        // Create a record
        let query = CString::new("CREATE person SET name = 'John', age = 30").unwrap();
        let result_ptr = surreal_query(handle, query.as_ptr());
        assert!(!result_ptr.is_null(), "Query returned NULL");

        let result_str = CStr::from_ptr(result_ptr).to_string_lossy();
        println!("Create result: {}", result_str);
        surreal_free_string(result_ptr);
        
        assert_ne!(result_str, "[]", "Create should return data");

        // Select records
        let query = CString::new("SELECT * FROM person").unwrap();
        let result_ptr = surreal_query(handle, query.as_ptr());
        assert!(!result_ptr.is_null(), "Query returned NULL");

        let result_str = CStr::from_ptr(result_ptr).to_string_lossy();
        println!("Select result: {}", result_str);
        surreal_free_string(result_ptr);

        assert_ne!(result_str, "[]", "Select should return data");
        assert!(result_str.contains("John"), "Result should contain 'John'");

        // Close
        surreal_close(handle);
    }
}
