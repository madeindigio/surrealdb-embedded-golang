package main

import (
	"fmt"
	"log"

	surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
	fmt.Println("SurrealDB Embedded - Graph Relations Example")
	fmt.Println("============================================")

	// Create an in-memory database
	db, err := surrealdb.NewMemory()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	if err := db.Use("test", "test"); err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}

	// Create people
	fmt.Println("\n1. Creating people...")
	people := []struct {
		ID   string
		Name string
		Role string
	}{
		{"person:alice", "Alice", "Developer"},
		{"person:bob", "Bob", "Designer"},
		{"person:charlie", "Charlie", "Manager"},
		{"person:diana", "Diana", "Developer"},
	}

	for _, p := range people {
		_, err := db.Create(p.ID, map[string]interface{}{
			"name": p.Name,
			"role": p.Role,
		})
		if err != nil {
			log.Fatalf("Failed to create person: %v", err)
		}
		fmt.Printf("   ✓ Created: %s (%s)\n", p.Name, p.Role)
	}

	// Create relationships
	fmt.Println("\n2. Creating relationships...")
	relationships := []struct {
		From string
		To   string
		Type string
	}{
		{"person:alice", "person:bob", "knows"},
		{"person:alice", "person:charlie", "reports_to"},
		{"person:bob", "person:charlie", "reports_to"},
		{"person:diana", "person:charlie", "reports_to"},
		{"person:bob", "person:alice", "works_with"},
		{"person:diana", "person:alice", "works_with"},
	}

	for _, rel := range relationships {
		query := fmt.Sprintf("RELATE %s->%s->%s", rel.From, rel.Type, rel.To)
		_, err := db.Query(query, nil)
		if err != nil {
			log.Fatalf("Failed to create relationship: %v", err)
		}
		fmt.Printf("   ✓ %s -> %s -> %s\n", rel.From, rel.Type, rel.To)
	}

	// Query: Who does Alice know?
	fmt.Println("\n3. Who does Alice know?")
	results, err := db.Query(`
		SELECT ->knows->person.name as knows
		FROM person:alice
	`, nil)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("   Result: %v\n", results)

	// Query: Who reports to Charlie?
	fmt.Println("\n4. Who reports to Charlie?")
	results, err = db.Query(`
		SELECT <-reports_to<-person.* as team_members
		FROM person:charlie
	`, nil)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("   Result: %v\n", results)

	// Query: Who works with Alice?
	fmt.Println("\n5. Who works with Alice?")
	results, err = db.Query(`
		SELECT ->works_with->person.name as colleagues
		FROM person:alice
	`, nil)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("   Result: %v\n", results)

	// Complex graph traversal
	fmt.Println("\n6. Complex query: Alice's network (2 degrees)")
	results, err = db.Query(`
		SELECT
			name,
			->knows->person.name as direct_contacts,
			->knows->person->knows->person.name as extended_network
		FROM person:alice
	`, nil)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("   Result: %v\n", results)

	// Count relationships
	fmt.Println("\n7. Relationship statistics")
	results, err = db.Query(`
		SELECT
			count() as total_people
		FROM person;

		SELECT
			count() as total_relationships
		FROM knows;

		SELECT
			count() as reporting_relationships
		FROM reports_to;
	`, nil)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("   Statistics: %v\n", results)

	fmt.Println("\n✓ Graph operations completed successfully!")
}
