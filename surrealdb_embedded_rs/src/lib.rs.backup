use std::collections::HashMap;
use std::ffi::{CStr, CString};
use std::os::raw::c_char;
use std::sync::{Arc, Mutex, OnceLock};
use surrealdb::engine::local::{Db, Mem, RocksDb};
use surrealdb::sql::Value;
use surrealdb::Surreal;

// Global database instance storage - unified type
static DB_INSTANCES: OnceLock<Mutex<HashMap<i32, Arc<Surreal<Db>>>>> = OnceLock::new();
static NEXT_HANDLE: Mutex<i32> = Mutex::new(1);

// Global Tokio runtime - must be shared across all operations
static RUNTIME: OnceLock<tokio::runtime::Runtime> = OnceLock::new();

// Get or create the global runtime
fn get_runtime() -> &'static tokio::runtime::Runtime {
    RUNTIME.get_or_init(|| tokio::runtime::Runtime::new().unwrap())
}

// Helper to get DB instances map
fn get_db_instances() -> &'static Mutex<HashMap<i32, Arc<Surreal<Db>>>> {
    DB_INSTANCES.get_or_init(|| Mutex::new(HashMap::new()))
}

// Error codes
pub const SURREAL_OK: i32 = 0;
pub const SURREAL_ERR_NULL_PTR: i32 = -1;
pub const SURREAL_ERR_INVALID_HANDLE: i32 = -2;
pub const SURREAL_ERR_INIT_FAILED: i32 = -3;
pub const SURREAL_ERR_QUERY_FAILED: i32 = -4;
pub const SURREAL_ERR_USE_FAILED: i32 = -5;

/// Initialize an embedded SurrealDB instance with memory backend
/// Returns a handle (positive integer) on success, or negative error code on failure
#[no_mangle]
pub extern "C" fn surreal_init_mem() -> i32 {
    let rt = get_runtime();

    match rt.block_on(async { Surreal::new::<Mem>(()).await }) {
        Ok(db) => {
            let mut next_handle = NEXT_HANDLE.lock().unwrap();
            let handle = *next_handle;
            *next_handle += 1;

            let mut instances = get_db_instances().lock().unwrap();
            instances.insert(handle, Arc::new(db));
            handle
        }
        Err(_) => SURREAL_ERR_INIT_FAILED,
    }
}

/// Initialize an embedded SurrealDB instance with RocksDB backend
/// path: The file path where the RocksDB data will be stored
/// Returns a handle (positive integer) on success, or negative error code on failure
#[no_mangle]
pub extern "C" fn surreal_init_rocksdb(path_ptr: *const c_char) -> i32 {
    if path_ptr.is_null() {
        return SURREAL_ERR_NULL_PTR;
    }

    let path = unsafe { CStr::from_ptr(path_ptr).to_string_lossy().into_owned() };
    let rt = get_runtime();

    match rt.block_on(async { Surreal::new::<RocksDb>(&path).await }) {
        Ok(db) => {
            let mut next_handle = NEXT_HANDLE.lock().unwrap();
            let handle = *next_handle;
            *next_handle += 1;

            let mut instances = get_db_instances().lock().unwrap();
            instances.insert(handle, Arc::new(db));
            handle
        }
        Err(_) => SURREAL_ERR_INIT_FAILED,
    }
}

/// Select namespace and database to use
/// handle: Database handle from init functions
/// ns: Namespace name
/// db: Database name
/// Returns SURREAL_OK on success, or negative error code on failure
#[no_mangle]
pub extern "C" fn surreal_use(handle: i32, ns_ptr: *const c_char, db_ptr: *const c_char) -> i32 {
    if ns_ptr.is_null() || db_ptr.is_null() {
        return SURREAL_ERR_NULL_PTR;
    }

    let instances = get_db_instances().lock().unwrap();
    let db_instance = match instances.get(&handle) {
        Some(db) => db.clone(),
        None => return SURREAL_ERR_INVALID_HANDLE,
    };
    drop(instances);

    let ns = unsafe { CStr::from_ptr(ns_ptr).to_string_lossy().into_owned() };
    let db_name = unsafe { CStr::from_ptr(db_ptr).to_string_lossy().into_owned() };

    let rt = get_runtime();
    let result = rt.block_on(async {
        db_instance.use_ns(&ns).await?;
        db_instance.use_db(&db_name).await
    });

    match result {
        Ok(_) => SURREAL_OK,
        Err(e) => {
            eprintln!("Error in surreal_use: {:?}", e);
            SURREAL_ERR_USE_FAILED
        }
    }
}

/// Execute a SurrealQL query and return result as JSON string
/// handle: Database handle from init functions
/// query: SurrealQL query string
/// Returns JSON result string on success, or NULL on failure
#[no_mangle]
pub extern "C" fn surreal_query(handle: i32, query_ptr: *const c_char) -> *mut c_char {
    if query_ptr.is_null() {
        return std::ptr::null_mut();
    }

    let instances = get_db_instances().lock().unwrap();
    let db_instance = match instances.get(&handle) {
        Some(db) => db.clone(),
        None => return std::ptr::null_mut(),
    };
    drop(instances);

    let query = unsafe { CStr::from_ptr(query_ptr).to_string_lossy().into_owned() };

    let rt = get_runtime();
    let result = rt.block_on(async { db_instance.query(&query).await });

    match result {
        Ok(mut response) => {
            // Try to take all results - use a more robust approach
            let json_result = match response.take::<Vec<Value>>(0) {
                Ok(values) => {
                    // Convert to JSON and normalize any SurrealDB-specific formatting
                    match serde_json::to_string(&values) {
                        Ok(json) => normalize_surrealdb_json(&json),
                        Err(_) => "[]".to_string(),
                    }
                }
                Err(_) => {
                    // If Vec<Value> fails, just return empty array
                    // SurrealDB responses should always be Vec<Value>
                    "[]".to_string()
                }
            };

            CString::new(json_result).unwrap().into_raw()
        }
        Err(e) => {
            let error_json = format!(r#"{{"error": "{}"}}"#, e);
            CString::new(error_json).unwrap().into_raw()
        }
    }
}

/// Normalize SurrealDB JSON format to standard JSON
/// Removes 'f' suffix from float numbers and fixes other SurrealDB-specific formatting
fn normalize_surrealdb_json(json: &str) -> String {
    // Replace float suffix 'f' with nothing (e.g., "0.412f" -> "0.412")
    // Use regex-like pattern matching
    let result = json.to_string();

    // Handle patterns like: 0.123f, 1.0f, etc.
    // We need to be careful to only match float literals, not 'f' in strings
    let mut normalized = String::with_capacity(result.len());
    let mut chars = result.chars().peekable();
    let mut in_string = false;
    let mut escape_next = false;

    while let Some(ch) = chars.next() {
        if escape_next {
            normalized.push(ch);
            escape_next = false;
            continue;
        }

        if ch == '\\' && in_string {
            escape_next = true;
            normalized.push(ch);
            continue;
        }

        if ch == '"' {
            in_string = !in_string;
            normalized.push(ch);
            continue;
        }

        // If we're in a string, just copy the character
        if in_string {
            normalized.push(ch);
            continue;
        }

        // Check if this is a float suffix 'f'
        if ch == 'f' {
            // Look back to see if previous character is a digit or '.'
            if let Some(last_ch) = normalized.chars().last() {
                if last_ch.is_ascii_digit() || last_ch == '.' {
                    // Check if next char is not alphanumeric (to avoid matching 'from', 'for', etc.)
                    if let Some(&next_ch) = chars.peek() {
                        if !next_ch.is_alphanumeric() && next_ch != '_' {
                            // Skip the 'f' suffix
                            continue;
                        }
                    } else {
                        // End of string, skip the 'f'
                        continue;
                    }
                }
            }
        }

        normalized.push(ch);
    }

    normalized
}

/// Execute a SurrealQL query with parameters and return result as JSON string
/// handle: Database handle from init functions
/// query: SurrealQL query string
/// params: JSON string containing query parameters
/// Returns JSON result string on success, or NULL on failure
#[no_mangle]
pub extern "C" fn surreal_query_with_params(
    handle: i32,
    query_ptr: *const c_char,
    params_ptr: *const c_char,
) -> *mut c_char {
    if query_ptr.is_null() {
        return std::ptr::null_mut();
    }

    let instances = get_db_instances().lock().unwrap();
    let db_instance = match instances.get(&handle) {
        Some(db) => db.clone(),
        None => return std::ptr::null_mut(),
    };
    drop(instances);

    let query = unsafe { CStr::from_ptr(query_ptr).to_string_lossy().into_owned() };

    let params: HashMap<String, serde_json::Value> = if params_ptr.is_null() {
        HashMap::new()
    } else {
        let params_str = unsafe { CStr::from_ptr(params_ptr).to_string_lossy().into_owned() };
        serde_json::from_str(&params_str).unwrap_or_default()
    };

    let rt = get_runtime();
    let result = rt.block_on(async {
        let mut q = db_instance.query(&query);

        for (key, value) in params {
            q = q.bind((key, value));
        }

        q.await
    });

    match result {
        Ok(mut response) => {
            // Use the same robust approach as surreal_query
            let json_result = match response.take::<Vec<Value>>(0) {
                Ok(values) => match serde_json::to_string(&values) {
                    Ok(json) => normalize_surrealdb_json(&json),
                    Err(_) => "[]".to_string(),
                },
                Err(_) => {
                    // If Vec<Value> fails, just return empty array
                    // SurrealDB responses should always be Vec<Value>
                    "[]".to_string()
                }
            };

            CString::new(json_result).unwrap().into_raw()
        }
        Err(e) => {
            let error_json = format!(r#"{{"error": "{}"}}"#, e);
            CString::new(error_json).unwrap().into_raw()
        }
    }
}

/// Select all records from a table or a specific record
/// handle: Database handle from init functions
/// resource: Table name or record ID (e.g., "person" or "person:123")
/// Returns JSON result string on success, or NULL on failure
#[no_mangle]
pub extern "C" fn surreal_select(handle: i32, resource_ptr: *const c_char) -> *mut c_char {
    if resource_ptr.is_null() {
        return std::ptr::null_mut();
    }

    let resource = unsafe { CStr::from_ptr(resource_ptr).to_string_lossy().into_owned() };
    let query = format!("SELECT * FROM {}", resource);

    surreal_query(handle, CString::new(query).unwrap().as_ptr())
}

/// Create a record in the database
/// handle: Database handle from init functions
/// resource: Table name or record ID (e.g., "person" or "person:123")
/// data: JSON string containing record data
/// Returns JSON result string on success, or NULL on failure
#[no_mangle]
pub extern "C" fn surreal_create(
    handle: i32,
    resource_ptr: *const c_char,
    data_ptr: *const c_char,
) -> *mut c_char {
    if resource_ptr.is_null() || data_ptr.is_null() {
        return std::ptr::null_mut();
    }

    let resource = unsafe { CStr::from_ptr(resource_ptr).to_string_lossy().into_owned() };
    let data = unsafe { CStr::from_ptr(data_ptr).to_string_lossy().into_owned() };

    let query = format!("CREATE {} CONTENT {}", resource, data);
    surreal_query(handle, CString::new(query).unwrap().as_ptr())
}

/// Update a record in the database (replaces entire content)
/// handle: Database handle from init functions
/// resource: Table name or record ID
/// data: JSON string containing new record data
/// Returns JSON result string on success, or NULL on failure
#[no_mangle]
pub extern "C" fn surreal_update(
    handle: i32,
    resource_ptr: *const c_char,
    data_ptr: *const c_char,
) -> *mut c_char {
    if resource_ptr.is_null() || data_ptr.is_null() {
        return std::ptr::null_mut();
    }

    let resource = unsafe { CStr::from_ptr(resource_ptr).to_string_lossy().into_owned() };
    let data = unsafe { CStr::from_ptr(data_ptr).to_string_lossy().into_owned() };

    let query = format!("UPDATE {} CONTENT {}", resource, data);
    surreal_query(handle, CString::new(query).unwrap().as_ptr())
}

/// Merge data into a record (partial update)
/// handle: Database handle from init functions
/// resource: Table name or record ID
/// data: JSON string containing fields to merge
/// Returns JSON result string on success, or NULL on failure
#[no_mangle]
pub extern "C" fn surreal_merge(
    handle: i32,
    resource_ptr: *const c_char,
    data_ptr: *const c_char,
) -> *mut c_char {
    if resource_ptr.is_null() || data_ptr.is_null() {
        return std::ptr::null_mut();
    }

    let resource = unsafe { CStr::from_ptr(resource_ptr).to_string_lossy().into_owned() };
    let data = unsafe { CStr::from_ptr(data_ptr).to_string_lossy().into_owned() };

    let query = format!("UPDATE {} MERGE {}", resource, data);
    surreal_query(handle, CString::new(query).unwrap().as_ptr())
}

/// Delete records from database
/// handle: Database handle from init functions
/// resource: Table name or record ID
/// Returns JSON result string on success, or NULL on failure
#[no_mangle]
pub extern "C" fn surreal_delete(handle: i32, resource_ptr: *const c_char) -> *mut c_char {
    if resource_ptr.is_null() {
        return std::ptr::null_mut();
    }

    let resource = unsafe { CStr::from_ptr(resource_ptr).to_string_lossy().into_owned() };
    let query = format!("DELETE FROM {}", resource);

    surreal_query(handle, CString::new(query).unwrap().as_ptr())
}

/// Insert one or more records with specified IDs
/// handle: Database handle from init functions
/// table: Table name
/// data: JSON string or array of records to insert
/// Returns JSON result string on success, or NULL on failure
#[no_mangle]
pub extern "C" fn surreal_insert(
    handle: i32,
    table_ptr: *const c_char,
    data_ptr: *const c_char,
) -> *mut c_char {
    if table_ptr.is_null() || data_ptr.is_null() {
        return std::ptr::null_mut();
    }

    let table = unsafe { CStr::from_ptr(table_ptr).to_string_lossy().into_owned() };
    let data = unsafe { CStr::from_ptr(data_ptr).to_string_lossy().into_owned() };

    let query = format!("INSERT INTO {} {}", table, data);
    surreal_query(handle, CString::new(query).unwrap().as_ptr())
}

/// Upsert a record (create or update)
/// handle: Database handle from init functions
/// resource: Table name or record ID
/// data: JSON string containing record data
/// Returns JSON result string on success, or NULL on failure
#[no_mangle]
pub extern "C" fn surreal_upsert(
    handle: i32,
    resource_ptr: *const c_char,
    data_ptr: *const c_char,
) -> *mut c_char {
    if resource_ptr.is_null() || data_ptr.is_null() {
        return std::ptr::null_mut();
    }

    let resource = unsafe { CStr::from_ptr(resource_ptr).to_string_lossy().into_owned() };
    let data = unsafe { CStr::from_ptr(data_ptr).to_string_lossy().into_owned() };

    let query = format!("UPSERT {} CONTENT {}", resource, data);
    surreal_query(handle, CString::new(query).unwrap().as_ptr())
}

/// Get version information
/// handle: Database handle from init functions
/// Returns JSON string with version info, or NULL on failure
#[no_mangle]
pub extern "C" fn surreal_version(handle: i32) -> *mut c_char {
    let instances = get_db_instances().lock().unwrap();
    let db_instance = match instances.get(&handle) {
        Some(db) => db.clone(),
        None => return std::ptr::null_mut(),
    };
    drop(instances);

    let rt = get_runtime();
    let result = rt.block_on(async { db_instance.version().await });

    match result {
        Ok(version) => {
            let version_json = serde_json::json!({
                "version": version.to_string()
            });
            CString::new(version_json.to_string()).unwrap().into_raw()
        }
        Err(_) => std::ptr::null_mut(),
    }
}

/// Close and cleanup a database instance
/// handle: Database handle to close
/// Returns SURREAL_OK on success
#[no_mangle]
pub extern "C" fn surreal_close(handle: i32) -> i32 {
    let mut instances = get_db_instances().lock().unwrap();
    instances.remove(&handle);
    SURREAL_OK
}

/// Free a string returned by the library
/// This must be called for every string returned by surreal_* functions
/// s: Pointer to string to free
#[no_mangle]
pub extern "C" fn surreal_free_string(s: *mut c_char) {
    if !s.is_null() {
        unsafe {
            let _ = CString::from_raw(s);
        }
    }
}
