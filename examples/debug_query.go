package main

import (
	"fmt"
	"log"
	"os"

	surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
	// Create RocksDB database
	dbPath := "/tmp/debug_rocksdb"
	defer os.RemoveAll(dbPath)

	db, err := surrealdb.NewFromURL(fmt.Sprintf("rocksdb://%s", dbPath))
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Use namespace and database
	if err := db.Use("test", "test"); err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}

	fmt.Println("1. Creating product:1...")
	result1, err := db.Create("product:1", map[string]interface{}{
		"name":  "Laptop",
		"price": 999.99,
	})
	if err != nil {
		log.Fatalf("Failed to create: %v", err)
	}
	fmt.Printf("   Created: %+v\n", result1)

	fmt.Println("\n2. Selecting product:1...")
	result2, err := db.Select("product:1")
	if err != nil {
		log.Fatalf("Failed to select: %v", err)
	}
	fmt.Printf("   Selected: %+v\n", result2)

	fmt.Println("\n3. Querying with SELECT * FROM product...")
	result3, err := db.Query("SELECT * FROM product", nil)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("   Query result: %+v (len=%d)\n", result3, len(result3))

	fmt.Println("\n4. Querying with SELECT * FROM type::table('product')...")
	result4, err := db.Query("SELECT * FROM type::table('product')", nil)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("   Query result: %+v (len=%d)\n", result4, len(result4))

	fmt.Println("\n5. Selecting all products using Select()...")
	result5, err := db.Select("product")
	if err != nil {
		log.Fatalf("Failed to select all: %v", err)
	}
	fmt.Printf("   Select result: %+v\n", result5)
}
