package main

import (
	"fmt"
	"log"
	"os"

	surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
	// Create a SurrealKV database instance
	db, err := surrealdb.NewFromURL("surrealkv://./data/example.skv")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()
	defer os.RemoveAll("./data") // Cleanup

	// Select namespace and database
	if err := db.Use("example", "main"); err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}

	// Create some users
	users := []map[string]interface{}{
		{"name": "Alice", "age": 30, "email": "alice@example.com"},
		{"name": "Bob", "age": 25, "email": "bob@example.com"},
		{"name": "Charlie", "age": 35, "email": "charlie@example.com"},
	}

	fmt.Println("Creating users...")
	for i, user := range users {
		result, err := db.Create(fmt.Sprintf("user:%d", i+1), user)
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Printf("  ✓ Created: %+v\n", result)
	}

	// Query all users
	fmt.Println("\nQuerying all users...")
	allUsers, err := db.Query("SELECT * FROM user ORDER BY age", nil)
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}
	fmt.Printf("Result: %+v\n", allUsers)

	// Query with parameters
	fmt.Println("\nQuerying users older than 28...")
	params := map[string]interface{}{
		"min_age": 28,
	}
	oldUsers, err := db.Query("SELECT * FROM user WHERE age > $min_age", params)
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}
	fmt.Printf("Result: %+v\n", oldUsers)

	// Update a user
	fmt.Println("\nUpdating Bob's age...")
	updated, err := db.Merge("user:2", map[string]interface{}{"age": 26})
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
	fmt.Printf("Updated: %+v\n", updated)

	// Select specific user
	fmt.Println("\nSelecting user:2...")
	bob, err := db.Select("user:2")
	if err != nil {
		log.Fatalf("Failed to select user: %v", err)
	}
	fmt.Printf("Bob: %+v\n", bob)

	fmt.Println("\n✅ Example completed successfully!")
}
