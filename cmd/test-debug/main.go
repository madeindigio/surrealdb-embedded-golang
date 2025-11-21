package main

import (
	"fmt"
	"log"
	"os"

	surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
	fmt.Println("Debugging SurrealDB Query Response")
	fmt.Println("===================================\n")

	tmpDir := "/tmp/surrealdb_debug_test"
	os.RemoveAll(tmpDir)
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	db, err := surrealdb.NewRocksDB(tmpDir + "/testdb")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	err = db.Use("test", "test")
	if err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}

	// Test different query types
	queries := []string{
		"CREATE person:john SET name = 'John', age = 30",
		"SELECT * FROM person",
		"SELECT * FROM person WHERE name = 'John'",
		"INFO FOR DB",
	}

	for i, query := range queries {
		fmt.Printf("\n%d. Query: %s\n", i+1, query)
		result, err := db.Query(query, nil)
		if err != nil {
			fmt.Printf("   Error: %v\n", err)
		} else {
			fmt.Printf("   Result type: %T\n", result)
			fmt.Printf("   Result length: %d\n", len(result))
			fmt.Printf("   Result: %+v\n", result)
		}
	}
}
