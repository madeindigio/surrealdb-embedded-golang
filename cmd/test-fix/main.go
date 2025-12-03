package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
	fmt.Println("Testing SurrealDB Query Fix")
	fmt.Println("=============================\n")

	// Create temporary directory for test database
	tmpDir := "/tmp/surrealdb_fix_test"
	os.RemoveAll(tmpDir) // Clean up any previous test
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	// Test 1: Create database and insert data
	fmt.Println("Test 1: Creating database with RocksDB backend...")
	db, err := surrealdb.NewRocksDB(tmpDir + "/testdb")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	err = db.Use("test", "test")
	if err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}

	// Test 2: Define a table schema
	fmt.Println("Test 2: Defining table schema...")
	result, err := db.Query("DEFINE TABLE test_table SCHEMAFULL", nil)
	if err != nil {
		log.Fatalf("Failed to define table: %v", err)
	}
	fmt.Printf("  Define table result: %v\n", result)

	// Test 3: Define a field
	fmt.Println("Test 3: Defining field...")
	result, err = db.Query("DEFINE FIELD name ON test_table TYPE string", nil)
	if err != nil {
		log.Fatalf("Failed to define field: %v", err)
	}
	fmt.Printf("  Define field result: %v\n", result)

	// Test 4: Insert a record
	fmt.Println("Test 4: Creating test record...")
	result, err = db.Query("CREATE test_table CONTENT { name: 'test' }", nil)
	if err != nil {
		log.Fatalf("Failed to create record: %v", err)
	}
	fmt.Printf("  Create result: %v\n", result)
	
	if len(result) == 0 {
		log.Fatal("  ❌ FAILED: Create returned empty array!")
	}
	fmt.Println("  ✓ Create returned data")

	// Test 5: Select immediately (this was failing before)
	fmt.Println("Test 5: Selecting record (immediate)...")
	result, err = db.Query("SELECT * FROM test_table", nil)
	if err != nil {
		log.Fatalf("Failed to select: %v", err)
	}
	fmt.Printf("  Select result: %v\n", result)
	
	if len(result) == 0 {
		log.Fatal("  ❌ FAILED: Select returned empty array!")
	}
	fmt.Println("  ✓ Select returned data")

	// Test 6: Close and reopen database
	fmt.Println("Test 6: Closing database...")
	err = db.Close()
	if err != nil {
		log.Fatalf("Failed to close database: %v", err)
	}
	fmt.Println("  ✓ Database closed")

	fmt.Println("Test 7: Reopening database...")
	db2, err := surrealdb.NewRocksDB(tmpDir + "/testdb")
	if err != nil {
		log.Fatalf("Failed to reopen database: %v", err)
	}
	defer db2.Close()

	err = db2.Use("test", "test")
	if err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}
	fmt.Println("  ✓ Database reopened")

	// Test 8: Select after reopen (this was the main issue)
	fmt.Println("Test 8: Selecting record after reopen...")
	result, err = db2.Query("SELECT * FROM test_table", nil)
	if err != nil {
		log.Fatalf("Failed to select after reopen: %v", err)
	}
	
	if len(result) == 0 {
		log.Fatal("  ❌ FAILED: Data not persisted! Select returned empty array after reopen!")
	}
	
	// Pretty print the result
	jsonBytes, _ := json.MarshalIndent(result, "  ", "  ")
	fmt.Printf("  Select result: %s\n", string(jsonBytes))
	fmt.Println("  ✓ Data persisted correctly!")

	// Test 9: Insert more records
	fmt.Println("Test 9: Inserting multiple records...")
	for i := 1; i <= 5; i++ {
		query := fmt.Sprintf("CREATE test_table CONTENT { name: 'test%d' }", i)
		_, err = db2.Query(query, nil)
		if err != nil {
			log.Fatalf("Failed to create record %d: %v", i, err)
		}
	}
	fmt.Println("  ✓ Multiple records created")

	// Test 10: Count records
	fmt.Println("Test 10: Counting records...")
	result, err = db2.Query("SELECT count() FROM test_table GROUP ALL", nil)
	if err != nil {
		log.Fatalf("Failed to count: %v", err)
	}
	jsonBytes, _ = json.MarshalIndent(result, "  ", "  ")
	fmt.Printf("  Count result: %s\n", string(jsonBytes))

	// Test 11: Select all records
	fmt.Println("Test 11: Selecting all records...")
	result, err = db2.Query("SELECT * FROM test_table", nil)
	if err != nil {
		log.Fatalf("Failed to select all: %v", err)
	}
	fmt.Printf("  Found %d records\n", len(result))
	
	if len(result) != 6 {
		log.Fatalf("  ❌ FAILED: Expected 6 records, got %d", len(result))
	}
	fmt.Println("  ✓ All records retrieved correctly")

	fmt.Println("\n=============================")
	fmt.Println("✅ All tests passed successfully!")
	fmt.Println("The query extraction fix is working correctly.")
}
