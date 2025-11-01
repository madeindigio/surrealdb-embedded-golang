package main

import (
	"fmt"
	"log"
	"os"

	surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
	fmt.Println("SurrealDB Embedded - Persistent Storage Example")
	fmt.Println("================================================")

	dbPath := "./data/myapp.db"

	// Check if database exists
	_, err := os.Stat(dbPath)
	isNew := os.IsNotExist(err)

	// Create a persistent RocksDB database
	db, err := surrealdb.NewRocksDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Select namespace and database
	if err := db.Use("myapp", "production"); err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}

	if isNew {
		fmt.Println("\n✓ Created new database")

		// Initialize schema
		fmt.Println("\n1. Setting up schema...")
		_, err = db.Query(`
			DEFINE TABLE user SCHEMAFULL;
			DEFINE FIELD username ON user TYPE string;
			DEFINE FIELD email ON user TYPE string ASSERT string::is::email($value);
			DEFINE FIELD created_at ON user TYPE datetime DEFAULT time::now();
			DEFINE INDEX idx_username ON user COLUMNS username UNIQUE;
		`, nil)
		if err != nil {
			log.Fatalf("Failed to create schema: %v", err)
		}
		fmt.Println("   ✓ Schema created")

		// Add initial data
		fmt.Println("\n2. Adding initial data...")
		users := []map[string]interface{}{
			{"username": "admin", "email": "admin@example.com"},
			{"username": "user1", "email": "user1@example.com"},
		}

		for _, user := range users {
			_, err := db.Create("user", user)
			if err != nil {
				log.Fatalf("Failed to create user: %v", err)
			}
			fmt.Printf("   ✓ Created user: %s\n", user["username"])
		}
	} else {
		fmt.Println("\n✓ Opened existing database")
	}

	// Query all users
	fmt.Println("\n3. Listing all users...")
	results, err := db.Query("SELECT * FROM user ORDER BY created_at", nil)
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}
	fmt.Printf("   Total users: %d\n", len(results))
	for i, user := range results {
		fmt.Printf("   %d. %v\n", i+1, user)
	}

	// Add a new user on each run
	fmt.Println("\n4. Adding a new user...")
	newUser := map[string]interface{}{
		"username": fmt.Sprintf("user_%d", len(results)+1),
		"email":    fmt.Sprintf("user%d@example.com", len(results)+1),
	}
	created, err := db.Create("user", newUser)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("   ✓ Created: %v\n", created)

	fmt.Println("\n✓ Data persisted to disk at:", dbPath)
	fmt.Println("\n  Run this program again to see the data persist!")
}
