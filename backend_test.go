package surrealdb

import (
	"fmt"
	"os"
	"testing"
)

// TestRocksDBBackend tests the RocksDB backend functionality
func TestRocksDBBackend(t *testing.T) {
	// Setup: create temporary directory for test database
	dbPath := "/tmp/test_rocksdb_backend"
	defer os.RemoveAll(dbPath)

	// Create database instance
	db, err := NewFromURL(fmt.Sprintf("rocksdb://%s", dbPath))
	if err != nil {
		t.Fatalf("Failed to create RocksDB database: %v", err)
	}
	defer db.Close()

	// Use namespace and database
	if err := db.Use("test", "test"); err != nil {
		t.Fatalf("Failed to use namespace/database: %v", err)
	}

	// Insert test data
	testProducts := []map[string]interface{}{
		{
			"id":    "product:1",
			"name":  "Laptop",
			"price": 999.99,
			"stock": 10,
		},
		{
			"id":    "product:2",
			"name":  "Mouse",
			"price": 29.99,
			"stock": 50,
		},
		{
			"id":    "product:3",
			"name":  "Keyboard",
			"price": 79.99,
			"stock": 30,
		},
	}

	t.Log("Inserting products...")
	for _, product := range testProducts {
		id := product["id"].(string)
		delete(product, "id") // Remove id from data map

		_, err := db.Create(id, product)
		if err != nil {
			t.Fatalf("Failed to create product %s: %v", id, err)
		}
		t.Logf("  ✓ Created: %s", id)
	}

	// Query all products
	t.Log("Querying all products...")
	results, err := db.Query("SELECT * FROM product ORDER BY price", nil)
	if err != nil {
		t.Fatalf("Failed to query products: %v", err)
	}

	// Verify results
	t.Logf("Query returned: %+v (len=%d)", results, len(results))
	if len(results) == 0 {
		t.Fatal("Query returned no results")
	}

	// Check that we got all products
	if len(results) != len(testProducts) {
		t.Fatalf("Expected %d products, got %d", len(testProducts), len(results))
	}

	t.Logf("✓ Found %d products", len(results))

	// Query with WHERE clause
	t.Log("Querying products with price > 50...")
	expensiveProducts, err := db.Query("SELECT * FROM product WHERE price > 50", nil)
	if err != nil {
		t.Fatalf("Failed to query expensive products: %v", err)
	}

	if len(expensiveProducts) != 2 {
		t.Fatalf("Expected 2 expensive products, got %d", len(expensiveProducts))
	}

	t.Logf("✓ Found %d expensive products", len(expensiveProducts))

	// Query with parameters
	t.Log("Querying products with parameterized query...")
	params := map[string]interface{}{
		"min_stock": 25,
	}
	inStockProducts, err := db.Query("SELECT * FROM product WHERE stock >= $min_stock", params)
	if err != nil {
		t.Fatalf("Failed to query products in stock: %v", err)
	}

	if len(inStockProducts) != 2 {
		t.Fatalf("Expected 2 products in stock, got %d", len(inStockProducts))
	}

	t.Logf("✓ Found %d products with stock >= 25", len(inStockProducts))

	// Select specific product
	t.Log("Selecting specific product...")
	product, err := db.Select("product:1")
	if err != nil {
		t.Fatalf("Failed to select product:1: %v", err)
	}

	productMap, ok := product.(map[string]interface{})
	if !ok {
		t.Fatal("Expected product to be a map")
	}

	if productMap["name"] != "Laptop" {
		t.Fatalf("Expected product name 'Laptop', got '%v'", productMap["name"])
	}

	t.Logf("✓ Selected product: %v", productMap["name"])

	t.Log("✅ RocksDB backend test passed!")
}

// TestSurrealKVBackend tests the SurrealKV backend functionality
func TestSurrealKVBackend(t *testing.T) {
	// Setup: create temporary directory for test database
	dbPath := "/tmp/test_surrealkv_backend"
	defer os.RemoveAll(dbPath)

	// Create database instance
	db, err := NewFromURL(fmt.Sprintf("surrealkv://%s", dbPath))
	if err != nil {
		t.Fatalf("Failed to create SurrealKV database: %v", err)
	}
	defer db.Close()

	// Use namespace and database
	if err := db.Use("test", "test"); err != nil {
		t.Fatalf("Failed to use namespace/database: %v", err)
	}

	// Insert test data
	testUsers := []map[string]interface{}{
		{
			"id":     "user:alice",
			"name":   "Alice Smith",
			"email":  "alice@example.com",
			"age":    30,
			"active": true,
		},
		{
			"id":     "user:bob",
			"name":   "Bob Johnson",
			"email":  "bob@example.com",
			"age":    25,
			"active": true,
		},
		{
			"id":     "user:charlie",
			"name":   "Charlie Brown",
			"email":  "charlie@example.com",
			"age":    35,
			"active": false,
		},
		{
			"id":     "user:diana",
			"name":   "Diana Prince",
			"email":  "diana@example.com",
			"age":    28,
			"active": true,
		},
	}

	t.Log("Inserting users...")
	for _, user := range testUsers {
		id := user["id"].(string)
		delete(user, "id") // Remove id from data map

		_, err := db.Create(id, user)
		if err != nil {
			t.Fatalf("Failed to create user %s: %v", id, err)
		}
		t.Logf("  ✓ Created: %s", id)
	}

	// Query all users
	t.Log("Querying all users...")
	results, err := db.Query("SELECT * FROM user ORDER BY age", nil)
	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}

	// Verify results
	if len(results) == 0 {
		t.Fatal("Query returned no results")
	}

	// Check that we got all users
	if len(results) != len(testUsers) {
		t.Fatalf("Expected %d users, got %d", len(testUsers), len(results))
	}

	t.Logf("✓ Found %d users", len(results))

	// Query active users only
	t.Log("Querying active users...")
	activeUsers, err := db.Query("SELECT * FROM user WHERE active = true", nil)
	if err != nil {
		t.Fatalf("Failed to query active users: %v", err)
	}

	if len(activeUsers) != 3 {
		t.Fatalf("Expected 3 active users, got %d", len(activeUsers))
	}

	t.Logf("✓ Found %d active users", len(activeUsers))

	// Query with age range using parameters
	t.Log("Querying users in age range...")
	params := map[string]interface{}{
		"min_age": 26,
		"max_age": 32,
	}
	ageRangeUsers, err := db.Query("SELECT * FROM user WHERE age >= $min_age AND age <= $max_age", params)
	if err != nil {
		t.Fatalf("Failed to query users in age range: %v", err)
	}

	if len(ageRangeUsers) != 2 {
		t.Fatalf("Expected 2 users in age range, got %d", len(ageRangeUsers))
	}

	t.Logf("✓ Found %d users in age range 26-32", len(ageRangeUsers))

	// Count query
	t.Log("Counting total users...")
	countResult, err := db.Query("SELECT count() AS total FROM user GROUP ALL", nil)
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	if len(countResult) == 0 {
		t.Fatal("Count query returned no results")
	}

	t.Logf("✓ Total users: %v", countResult)

	// Update a user
	t.Log("Updating user:bob...")
	_, err = db.Merge("user:bob", map[string]interface{}{
		"age": 26,
	})
	if err != nil {
		t.Fatalf("Failed to update user:bob: %v", err)
	}

	// Verify update
	updatedBob, err := db.Select("user:bob")
	if err != nil {
		t.Fatalf("Failed to select updated user:bob: %v", err)
	}

	bobMap, ok := updatedBob.(map[string]interface{})
	if !ok {
		t.Fatal("Expected user to be a map")
	}

	// Age might be returned as float64 from JSON
	bobAge := int(bobMap["age"].(float64))
	if bobAge != 26 {
		t.Fatalf("Expected Bob's age to be 26, got %d", bobAge)
	}

	t.Logf("✓ Updated Bob's age to %d", bobAge)

	// Delete a user
	t.Log("Deleting user:charlie...")
	_, err = db.Delete("user:charlie")
	if err != nil {
		t.Fatalf("Failed to delete user:charlie: %v", err)
	}

	// Verify deletion
	remainingUsers, err := db.Query("SELECT * FROM user", nil)
	if err != nil {
		t.Fatalf("Failed to query remaining users: %v", err)
	}

	if len(remainingUsers) != 3 {
		t.Fatalf("Expected 3 remaining users, got %d", len(remainingUsers))
	}

	t.Logf("✓ Remaining users: %d", len(remainingUsers))

	t.Log("✅ SurrealKV backend test passed!")
}

// TestBackendComparison tests both backends with the same data
func TestBackendComparison(t *testing.T) {
	testData := []map[string]interface{}{
		{
			"id":    "item:1",
			"title": "First Item",
			"value": 100,
		},
		{
			"id":    "item:2",
			"title": "Second Item",
			"value": 200,
		},
		{
			"id":    "item:3",
			"title": "Third Item",
			"value": 300,
		},
	}

	backends := []struct {
		name string
		url  string
		path string
	}{
		{"RocksDB", "rocksdb://%s", "/tmp/test_comparison_rocksdb"},
		{"SurrealKV", "surrealkv://%s", "/tmp/test_comparison_surrealkv"},
	}

	for _, backend := range backends {
		t.Run(backend.name, func(t *testing.T) {
			defer os.RemoveAll(backend.path)

			// Create database
			db, err := NewFromURL(fmt.Sprintf(backend.url, backend.path))
			if err != nil {
				t.Fatalf("[%s] Failed to create database: %v", backend.name, err)
			}
			defer db.Close()

			// Use namespace and database
			if err := db.Use("test", "test"); err != nil {
				t.Fatalf("[%s] Failed to use namespace/database: %v", backend.name, err)
			}

			// Insert data
			for _, item := range testData {
				id := item["id"].(string)
				itemCopy := make(map[string]interface{})
				for k, v := range item {
					if k != "id" {
						itemCopy[k] = v
					}
				}

				_, err := db.Create(id, itemCopy)
				if err != nil {
					t.Fatalf("[%s] Failed to create item %s: %v", backend.name, id, err)
				}
			}

			// Query and verify
			results, err := db.Query("SELECT * FROM item ORDER BY value", nil)
			if err != nil {
				t.Fatalf("[%s] Failed to query items: %v", backend.name, err)
			}

			if len(results) != len(testData) {
				t.Fatalf("[%s] Expected %d items, got %d", backend.name, len(testData), len(results))
			}

			// Query with aggregation
			sumResult, err := db.Query("SELECT math::sum(value) AS total FROM item GROUP ALL", nil)
			if err != nil {
				t.Fatalf("[%s] Failed to calculate sum: %v", backend.name, err)
			}

			if len(sumResult) == 0 {
				t.Fatalf("[%s] Sum query returned no results", backend.name)
			}

			t.Logf("[%s] ✓ Test passed with %d items", backend.name, len(results))
		})
	}

	t.Log("✅ Backend comparison test passed!")
}
