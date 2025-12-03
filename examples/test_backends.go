package main

import (
	"fmt"
	"log"
	"os"

	surrealdb "github.com/yourusername/surrealdb-embedded"
)

func testBackend(name string, db *surrealdb.DB, cleanup func()) {
	fmt.Printf("\n=== Testing %s backend ===\n", name)
	defer func() {
		if cleanup != nil {
			cleanup()
		}
	}()

	// Use namespace and database
	if err := db.Use("test", "test"); err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}

	// Create a record
	data := map[string]interface{}{
		"name": fmt.Sprintf("Test User (%s)", name),
		"age":  25,
	}

	result, err := db.Create("person:test", data)
	if err != nil {
		log.Fatalf("Failed to create record: %v", err)
	}
	fmt.Printf("Created: %+v\n", result)

	// Select the record
	selected, err := db.Select("person:test")
	if err != nil {
		log.Fatalf("Failed to select record: %v", err)
	}
	fmt.Printf("Selected: %+v\n", selected)

	// Query records
	records, err := db.Query("SELECT * FROM person", nil)
	if err != nil {
		log.Fatalf("Failed to query records: %v", err)
	}
	fmt.Printf("Query result: %+v\n", records)

	// Delete the record
	deleted, err := db.Delete("person:test")
	if err != nil {
		log.Fatalf("Failed to delete record: %v", err)
	}
	fmt.Printf("Deleted: %+v\n", deleted)

	// Close the database
	if err := db.Close(); err != nil {
		log.Fatalf("Failed to close database: %v", err)
	}

	fmt.Printf("✓ %s backend test passed!\n", name)
}

func main() {
	fmt.Println("Testing SurrealDB Embedded Backends")

	// Test 1: Memory backend
	fmt.Println("\n--- Test 1: Memory Backend (using NewMemory) ---")
	memDB, err := surrealdb.NewMemory()
	if err != nil {
		log.Fatalf("Failed to create memory database: %v", err)
	}
	testBackend("Memory", memDB, nil)

	// Test 2: Memory backend using NewFromURL
	fmt.Println("\n--- Test 2: Memory Backend (using NewFromURL) ---")
	memDB2, err := surrealdb.NewFromURL("memory")
	if err != nil {
		log.Fatalf("Failed to create memory database from URL: %v", err)
	}
	testBackend("Memory (URL)", memDB2, nil)

	// Test 3: RocksDB backend
	fmt.Println("\n--- Test 3: RocksDB Backend ---")
	rocksPath := "/tmp/test_rocksdb"
	rocksDB, err := surrealdb.NewRocksDB(rocksPath)
	if err != nil {
		log.Fatalf("Failed to create RocksDB database: %v", err)
	}
	testBackend("RocksDB", rocksDB, func() {
		os.RemoveAll(rocksPath)
	})

	// Test 4: RocksDB backend using NewFromURL
	fmt.Println("\n--- Test 4: RocksDB Backend (using NewFromURL) ---")
	rocksPath2 := "/tmp/test_rocksdb2"
	rocksDB2, err := surrealdb.NewFromURL(fmt.Sprintf("rocksdb://%s", rocksPath2))
	if err != nil {
		log.Fatalf("Failed to create RocksDB database from URL: %v", err)
	}
	testBackend("RocksDB (URL)", rocksDB2, func() {
		os.RemoveAll(rocksPath2)
	})

	// Test 5: SurrealKV backend
	fmt.Println("\n--- Test 5: SurrealKV Backend ---")
	skvPath := "/tmp/test_surrealkv"
	skvDB, err := surrealdb.NewSurrealKV(skvPath)
	if err != nil {
		log.Fatalf("Failed to create SurrealKV database: %v", err)
	}
	testBackend("SurrealKV", skvDB, func() {
		os.RemoveAll(skvPath)
	})

	// Test 6: SurrealKV backend using NewFromURL
	fmt.Println("\n--- Test 6: SurrealKV Backend (using NewFromURL) ---")
	skvPath2 := "/tmp/test_surrealkv2"
	skvDB2, err := surrealdb.NewFromURL(fmt.Sprintf("surrealkv://%s", skvPath2))
	if err != nil {
		log.Fatalf("Failed to create SurrealKV database from URL: %v", err)
	}
	testBackend("SurrealKV (URL)", skvDB2, func() {
		os.RemoveAll(skvPath2)
	})

	// Test 7: Using Config with different backends
	fmt.Println("\n--- Test 7: Using Config ---")

	// Memory via Config
	memConfigDB, err := surrealdb.New(surrealdb.Config{Backend: surrealdb.Memory})
	if err != nil {
		log.Fatalf("Failed to create memory database via Config: %v", err)
	}
	testBackend("Memory (Config)", memConfigDB, nil)

	// RocksDB via Config
	rocksPath3 := "/tmp/test_rocksdb_config"
	rocksConfigDB, err := surrealdb.New(surrealdb.Config{
		Backend: surrealdb.RocksDB,
		Path:    rocksPath3,
	})
	if err != nil {
		log.Fatalf("Failed to create RocksDB database via Config: %v", err)
	}
	testBackend("RocksDB (Config)", rocksConfigDB, func() {
		os.RemoveAll(rocksPath3)
	})

	// SurrealKV via Config
	skvPath3 := "/tmp/test_surrealkv_config"
	skvConfigDB, err := surrealdb.New(surrealdb.Config{
		Backend: surrealdb.SurrealKV,
		Path:    skvPath3,
	})
	if err != nil {
		log.Fatalf("Failed to create SurrealKV database via Config: %v", err)
	}
	testBackend("SurrealKV (Config)", skvConfigDB, func() {
		os.RemoveAll(skvPath3)
	})

	fmt.Println("\n✓✓✓ All backend tests passed! ✓✓✓")
}
