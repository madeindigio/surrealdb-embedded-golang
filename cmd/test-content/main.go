package main

import (
	"fmt"
	"log"

	surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
	db, err := surrealdb.NewMemory()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	err = db.Use("test", "test")
	if err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}

	// Create
	result, err := db.Create("person", map[string]interface{}{
		"name": "John",
		"age":  30,
	})
	
	fmt.Printf("Create result type: %T\n", result)
	fmt.Printf("Create result: %+v\n", result)
	fmt.Printf("Create error: %v\n", err)
	
	if arr, ok := result.([]interface{}); ok {
		if len(arr) == 0 {
			fmt.Println("❌ FAIL: Create returned empty array")
		} else {
			fmt.Println("✅ PASS: Create returned data")
		}
	}
	
	// Select
	result, err = db.Select("person")
	fmt.Printf("\nSelect result type: %T\n", result)
	fmt.Printf("Select result: %+v\n", result)
	fmt.Printf("Select error: %v\n", err)
	
	if arr, ok := result.([]interface{}); ok {
		if len(arr) == 0 {
			fmt.Println("❌ FAIL: Select returned empty array")
		} else {
			fmt.Println("✅ PASS: Select returned data")
		}
	}
}
