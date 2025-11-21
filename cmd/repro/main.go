package main

import (
	"encoding/json"
	"fmt"

	"time"

	surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
	fmt.Println("Starting repro...")
	db, err := surrealdb.NewMemory()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Use("test", "test")
	if err != nil {
		panic(err)
	}

	// Create a record with complex types
	now := time.Now().UTC()
	data := map[string]interface{}{
		"name":      "Complex User",
		"created_at": now.Format(time.RFC3339),
		"score":     123.456,
		"active":    true,
		"tags":      []string{"a", "b", "c"},
		"metadata": map[string]interface{}{
			"key": "value",
		},
	}

	created, err := db.Create("person:complex", data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created: %+v\n", created)

	// Select it back
	result, err := db.Select("person:complex")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected: %+v\n", result)

	// Check if we can marshal it back to JSON
	bytes, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("JSON: %s\n", string(bytes))
}
