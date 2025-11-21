package main

import (
	"fmt"
	"log"

	surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
	fmt.Println("Simple SurrealDB Test")
	fmt.Println("=====================\n")

	db, err := surrealdb.NewMemory()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	err = db.Use("test", "test")
	if err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}

	// Test 1: Create using the Create method
	fmt.Println("Test 1: Create using db.Create()")
	createResult, err := db.Create("person", map[string]interface{}{
		"name": "John",
		"age":  30,
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Result: %+v\n", createResult)
	}

	// Test 2: Query SELECT
	fmt.Println("\nTest 2: SELECT query")
	queryResult, err := db.Query("SELECT * FROM person", nil)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Result length: %d\n", len(queryResult))
		fmt.Printf("   Result: %+v\n", queryResult)
	}

	// Test 3: Query CREATE
	fmt.Println("\nTest 3: CREATE query")
	createQueryResult, err := db.Query("CREATE person:jane SET name = 'Jane', age = 25", nil)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Result length: %d\n", len(createQueryResult))
		fmt.Printf("   Result: %+v\n", createQueryResult)
	}

	// Test 4: SELECT again
	fmt.Println("\nTest 4: SELECT after CREATE query")
	queryResult2, err := db.Query("SELECT * FROM person", nil)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Result length: %d\n", len(queryResult2))
		fmt.Printf("   Result: %+v\n", queryResult2)
	}
}
