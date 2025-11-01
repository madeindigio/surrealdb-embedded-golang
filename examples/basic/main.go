package main

import (
	"fmt"
	"log"

	surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
	fmt.Println("SurrealDB Embedded - Basic Example")
	fmt.Println("===================================")

	// Create an in-memory database
	db, err := surrealdb.NewMemory()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Select namespace and database
	if err := db.Use("test", "test"); err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}

	// Create some people
	people := []struct {
		ID   string
		Name string
		Age  int
	}{
		{"person:alice", "Alice Johnson", 28},
		{"person:bob", "Bob Smith", 35},
		{"person:charlie", "Charlie Brown", 42},
	}

	fmt.Println("\n1. Creating records...")
	for _, p := range people {
		result, err := db.Create(p.ID, map[string]interface{}{
			"name": p.Name,
			"age":  p.Age,
		})
		if err != nil {
			log.Fatalf("Failed to create person: %v", err)
		}
		fmt.Printf("   Created: %v\n", result)
	}

	// Select all people
	fmt.Println("\n2. Selecting all records...")
	allPeople, err := db.Select("person")
	if err != nil {
		log.Fatalf("Failed to select people: %v", err)
	}
	fmt.Printf("   All people: %v\n", allPeople)

	// Query with filter
	fmt.Println("\n3. Querying with filter (age > 30)...")
	results, err := db.Query(
		"SELECT * FROM person WHERE age > $minAge ORDER BY age",
		map[string]interface{}{"minAge": 30},
	)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("   Results: %v\n", results)

	// Update a record
	fmt.Println("\n4. Updating a record...")
	updated, err := db.Update("person:alice", map[string]interface{}{
		"name": "Alice Cooper",
		"age":  29,
	})
	if err != nil {
		log.Fatalf("Failed to update: %v", err)
	}
	fmt.Printf("   Updated: %v\n", updated)

	// Merge partial data
	fmt.Println("\n5. Merging partial data...")
	merged, err := db.Merge("person:bob", map[string]interface{}{
		"email": "bob@example.com",
	})
	if err != nil {
		log.Fatalf("Failed to merge: %v", err)
	}
	fmt.Printf("   Merged: %v\n", merged)

	// Delete a record
	fmt.Println("\n6. Deleting a record...")
	deleted, err := db.Delete("person:charlie")
	if err != nil {
		log.Fatalf("Failed to delete: %v", err)
	}
	fmt.Printf("   Deleted: %v\n", deleted)

	// Get version
	fmt.Println("\n7. Getting version info...")
	version, err := db.Version()
	if err != nil {
		log.Fatalf("Failed to get version: %v", err)
	}
	fmt.Printf("   SurrealDB Version: %v\n", version)

	fmt.Println("\n✓ All operations completed successfully!")
}
