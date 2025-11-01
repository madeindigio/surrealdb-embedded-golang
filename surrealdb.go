package surrealdb

/*
#cgo CFLAGS: -I${SRCDIR}/surrealdb_embedded_rs/include
#cgo linux LDFLAGS: -L${SRCDIR}/surrealdb_embedded_rs/target/release -lsurrealdb_embedded_rs -lpthread -ldl -lm
#cgo darwin LDFLAGS: -L${SRCDIR}/surrealdb_embedded_rs/target/release -lsurrealdb_embedded_rs -framework Security -framework CoreFoundation
#cgo windows LDFLAGS: -L${SRCDIR}/surrealdb_embedded_rs/target/release -lsurrealdb_embedded_rs -lws2_32 -luserenv -lbcrypt

#include <stdlib.h>

// Error codes
#define SURREAL_OK 0
#define SURREAL_ERR_NULL_PTR -1
#define SURREAL_ERR_INVALID_HANDLE -2
#define SURREAL_ERR_INIT_FAILED -3
#define SURREAL_ERR_QUERY_FAILED -4
#define SURREAL_ERR_USE_FAILED -5

// Function declarations
extern int surreal_init_mem();
extern int surreal_init_rocksdb(const char* path);
extern int surreal_use(int handle, const char* ns, const char* db);
extern char* surreal_query(int handle, const char* query);
extern char* surreal_query_with_params(int handle, const char* query, const char* params);
extern char* surreal_select(int handle, const char* resource);
extern char* surreal_create(int handle, const char* resource, const char* data);
extern char* surreal_update(int handle, const char* resource, const char* data);
extern char* surreal_merge(int handle, const char* resource, const char* data);
extern char* surreal_delete(int handle, const char* resource);
extern char* surreal_insert(int handle, const char* table, const char* data);
extern char* surreal_upsert(int handle, const char* resource, const char* data);
extern char* surreal_version(int handle);
extern int surreal_close(int handle);
extern void surreal_free_string(char* s);
*/
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"
)

// Error definitions
var (
	ErrNullPtr       = errors.New("null pointer error")
	ErrInvalidHandle = errors.New("invalid database handle")
	ErrInitFailed    = errors.New("database initialization failed")
	ErrQueryFailed   = errors.New("query execution failed")
	ErrUseFailed     = errors.New("use namespace/database failed")
)

// DB represents an embedded SurrealDB instance
type DB struct {
	handle int
}

// BackendType represents the storage backend type
type BackendType int

const (
	// Memory backend stores data in memory only
	Memory BackendType = iota
	// RocksDB backend persists data to disk using RocksDB
	RocksDB
)

// Config holds database configuration
type Config struct {
	Backend BackendType
	Path    string // Only used for RocksDB backend
}

// NewMemory creates a new in-memory embedded SurrealDB instance
func NewMemory() (*DB, error) {
	handle := int(C.surreal_init_mem())
	if handle < 0 {
		return nil, handleError(handle)
	}
	return &DB{handle: handle}, nil
}

// NewRocksDB creates a new embedded SurrealDB instance with RocksDB backend
func NewRocksDB(path string) (*DB, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	handle := int(C.surreal_init_rocksdb(cPath))
	if handle < 0 {
		return nil, handleError(handle)
	}
	return &DB{handle: handle}, nil
}

// New creates a new embedded SurrealDB instance with the specified configuration
func New(config Config) (*DB, error) {
	switch config.Backend {
	case Memory:
		return NewMemory()
	case RocksDB:
		if config.Path == "" {
			return nil, errors.New("path is required for RocksDB backend")
		}
		return NewRocksDB(config.Path)
	default:
		return nil, errors.New("unknown backend type")
	}
}

// Use selects a namespace and database to use
func (db *DB) Use(namespace, database string) error {
	cNs := C.CString(namespace)
	cDb := C.CString(database)
	defer C.free(unsafe.Pointer(cNs))
	defer C.free(unsafe.Pointer(cDb))

	result := int(C.surreal_use(C.int(db.handle), cNs, cDb))
	if result != 0 {
		return handleError(result)
	}
	return nil
}

// Query executes a SurrealQL query and returns the result
func (db *DB) Query(query string, vars map[string]interface{}) ([]interface{}, error) {
	cQuery := C.CString(query)
	defer C.free(unsafe.Pointer(cQuery))

	var result *C.char
	if vars != nil && len(vars) > 0 {
		varsJSON, err := json.Marshal(vars)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal query variables: %w", err)
		}
		cVars := C.CString(string(varsJSON))
		defer C.free(unsafe.Pointer(cVars))
		result = C.surreal_query_with_params(C.int(db.handle), cQuery, cVars)
	} else {
		result = C.surreal_query(C.int(db.handle), cQuery)
	}

	if result == nil {
		return nil, ErrQueryFailed
	}
	defer C.surreal_free_string(result)

	jsonStr := C.GoString(result)

	// Check for error in response
	var errorResp map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &errorResp); err == nil {
		if errMsg, ok := errorResp["error"]; ok {
			return nil, fmt.Errorf("query error: %v", errMsg)
		}
	}

	var data []interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return data, nil
}

// Select retrieves all records from a table or a specific record
func (db *DB) Select(resource string) (interface{}, error) {
	cResource := C.CString(resource)
	defer C.free(unsafe.Pointer(cResource))

	result := C.surreal_select(C.int(db.handle), cResource)
	if result == nil {
		return nil, ErrQueryFailed
	}
	defer C.surreal_free_string(result)

	return parseResult(C.GoString(result))
}

// Create creates a new record in the database
func (db *DB) Create(resource string, data interface{}) (interface{}, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	cResource := C.CString(resource)
	cData := C.CString(string(dataJSON))
	defer C.free(unsafe.Pointer(cResource))
	defer C.free(unsafe.Pointer(cData))

	result := C.surreal_create(C.int(db.handle), cResource, cData)
	if result == nil {
		return nil, ErrQueryFailed
	}
	defer C.surreal_free_string(result)

	return parseResult(C.GoString(result))
}

// Update replaces a record in the database
func (db *DB) Update(resource string, data interface{}) (interface{}, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	cResource := C.CString(resource)
	cData := C.CString(string(dataJSON))
	defer C.free(unsafe.Pointer(cResource))
	defer C.free(unsafe.Pointer(cData))

	result := C.surreal_update(C.int(db.handle), cResource, cData)
	if result == nil {
		return nil, ErrQueryFailed
	}
	defer C.surreal_free_string(result)

	return parseResult(C.GoString(result))
}

// Merge partially updates a record in the database
func (db *DB) Merge(resource string, data interface{}) (interface{}, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	cResource := C.CString(resource)
	cData := C.CString(string(dataJSON))
	defer C.free(unsafe.Pointer(cResource))
	defer C.free(unsafe.Pointer(cData))

	result := C.surreal_merge(C.int(db.handle), cResource, cData)
	if result == nil {
		return nil, ErrQueryFailed
	}
	defer C.surreal_free_string(result)

	return parseResult(C.GoString(result))
}

// Delete removes records from the database
func (db *DB) Delete(resource string) (interface{}, error) {
	cResource := C.CString(resource)
	defer C.free(unsafe.Pointer(cResource))

	result := C.surreal_delete(C.int(db.handle), cResource)
	if result == nil {
		return nil, ErrQueryFailed
	}
	defer C.surreal_free_string(result)

	return parseResult(C.GoString(result))
}

// Insert creates one or more records with specified IDs
func (db *DB) Insert(table string, data interface{}) (interface{}, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	cTable := C.CString(table)
	cData := C.CString(string(dataJSON))
	defer C.free(unsafe.Pointer(cTable))
	defer C.free(unsafe.Pointer(cData))

	result := C.surreal_insert(C.int(db.handle), cTable, cData)
	if result == nil {
		return nil, ErrQueryFailed
	}
	defer C.surreal_free_string(result)

	return parseResult(C.GoString(result))
}

// Upsert creates or updates a record
func (db *DB) Upsert(resource string, data interface{}) (interface{}, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	cResource := C.CString(resource)
	cData := C.CString(string(dataJSON))
	defer C.free(unsafe.Pointer(cResource))
	defer C.free(unsafe.Pointer(cData))

	result := C.surreal_upsert(C.int(db.handle), cResource, cData)
	if result == nil {
		return nil, ErrQueryFailed
	}
	defer C.surreal_free_string(result)

	return parseResult(C.GoString(result))
}

// Version returns the SurrealDB version information
func (db *DB) Version() (map[string]interface{}, error) {
	result := C.surreal_version(C.int(db.handle))
	if result == nil {
		return nil, errors.New("failed to get version")
	}
	defer C.surreal_free_string(result)

	var version map[string]interface{}
	if err := json.Unmarshal([]byte(C.GoString(result)), &version); err != nil {
		return nil, fmt.Errorf("failed to unmarshal version: %w", err)
	}

	return version, nil
}

// Close closes the database connection and frees resources
func (db *DB) Close() error {
	result := int(C.surreal_close(C.int(db.handle)))
	if result != 0 {
		return handleError(result)
	}
	return nil
}

// Helper functions

func handleError(code int) error {
	switch code {
	case -1:
		return ErrNullPtr
	case -2:
		return ErrInvalidHandle
	case -3:
		return ErrInitFailed
	case -4:
		return ErrQueryFailed
	case -5:
		return ErrUseFailed
	default:
		return fmt.Errorf("unknown error code: %d", code)
	}
}

func parseResult(jsonStr string) (interface{}, error) {
	// Check for error in response
	var errorResp map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &errorResp); err == nil {
		if errMsg, ok := errorResp["error"]; ok {
			return nil, fmt.Errorf("database error: %v", errMsg)
		}
	}

	var result interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return result, nil
}
