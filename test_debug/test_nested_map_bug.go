package main

import (
	"fmt"
	"log"

	embedded "github.com/yourusername/surrealdb-embedded"
)

func main() {
	fmt.Println("=== Testing Nested Map Bug ===\n")

	// Create temporary database
	db, err := embedded.NewFromURL("memory")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Use("test", "test")
	if err != nil {
		log.Fatal(err)
	}

	// Test Case 1: Simple nested map (like in bug report)
	fmt.Println("Test 1: Simple nested metadata map")
	fmt.Println("-----------------------------------")

	metadata := map[string]interface{}{
		"last_modified": "2025-11-18T12:58:01+01:00",
		"source":        "watcher",
		"total_size":    8958,
	}

	params := map[string]interface{}{
		"name":     "test1",
		"metadata": metadata,
	}

	fmt.Printf("Input params: %+v\n", params)
	fmt.Printf("Input metadata: %+v\n\n", metadata)

	result, err := db.Query(`
		CREATE test_table CONTENT {
			name: $name,
			metadata: $metadata
		}
	`, params)

	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Create result: %+v\n", result)
		if len(result) > 0 {
			if resMap, ok := result[0].(map[string]interface{}); ok {
				fmt.Printf("Created metadata field: %+v\n", resMap["metadata"])
				if meta, ok := resMap["metadata"].(map[string]interface{}); ok {
					fmt.Printf("  - last_modified: %v\n", meta["last_modified"])
					fmt.Printf("  - source: %v\n", meta["source"])
					fmt.Printf("  - total_size: %v\n", meta["total_size"])

					if len(meta) == 0 {
						fmt.Println("❌ FAIL: metadata is empty!")
					} else {
						fmt.Println("✓ PASS: metadata has values")
					}
				} else {
					fmt.Printf("❌ FAIL: metadata is not a map: %T\n", resMap["metadata"])
				}
			}
		}
	}

	// Retrieve and verify
	fmt.Println("\nRetrieving record...")
	result, err = db.Query("SELECT * FROM test_table WHERE name = $name",
		map[string]interface{}{"name": "test1"})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Select result: %+v\n", result)
		if len(result) > 0 {
			if resMap, ok := result[0].(map[string]interface{}); ok {
				fmt.Printf("Retrieved metadata field: %+v\n", resMap["metadata"])
				if meta, ok := resMap["metadata"].(map[string]interface{}); ok {
					if len(meta) == 0 {
						fmt.Println("❌ FAIL: retrieved metadata is empty!")
					} else {
						fmt.Println("✓ PASS: retrieved metadata has values")
					}
				}
			}
		}
	}

	// Test Case 2: Nested map with different types
	fmt.Println("\n\nTest 2: Complex nested metadata")
	fmt.Println("--------------------------------")

	complexMeta := map[string]interface{}{
		"string_field": "test",
		"int_field":    42,
		"float_field":  3.14,
		"bool_field":   true,
		"nested": map[string]interface{}{
			"inner": "value",
		},
	}

	params2 := map[string]interface{}{
		"name":     "test2",
		"metadata": complexMeta,
	}

	fmt.Printf("Input complex metadata: %+v\n\n", complexMeta)

	result, err = db.Query(`
		CREATE test_table CONTENT {
			name: $name,
			metadata: $metadata
		}
	`, params2)

	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Create result: %+v\n", result)
		if len(result) > 0 {
			if resMap, ok := result[0].(map[string]interface{}); ok {
				fmt.Printf("Created metadata field: %+v\n", resMap["metadata"])
				if meta, ok := resMap["metadata"].(map[string]interface{}); ok {
					if len(meta) == 0 {
						fmt.Println("❌ FAIL: complex metadata is empty!")
					} else {
						fmt.Printf("✓ PASS: complex metadata has %d fields\n", len(meta))
					}
				}
			}
		}
	}

	// Test Case 3: Single-key nested map (edge case that might trigger the bug)
	fmt.Println("\n\nTest 3: Single-key nested map")
	fmt.Println("------------------------------")

	singleKeyMeta := map[string]interface{}{
		"only_field": "single value",
	}

	params3 := map[string]interface{}{
		"name":     "test3",
		"metadata": singleKeyMeta,
	}

	fmt.Printf("Input single-key metadata: %+v\n\n", singleKeyMeta)

	result, err = db.Query(`
		CREATE test_table CONTENT {
			name: $name,
			metadata: $metadata
		}
	`, params3)

	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Create result: %+v\n", result)
		if len(result) > 0 {
			if resMap, ok := result[0].(map[string]interface{}); ok {
				fmt.Printf("Created metadata field: %+v\n", resMap["metadata"])
				if meta, ok := resMap["metadata"].(map[string]interface{}); ok {
					if len(meta) == 0 {
						fmt.Println("❌ FAIL: single-key metadata is empty!")
					} else if val, ok := meta["only_field"]; ok {
						fmt.Printf("✓ PASS: single-key metadata preserved: only_field=%v\n", val)
					} else {
						fmt.Println("❌ FAIL: only_field not found in metadata")
					}
				} else {
					fmt.Printf("❌ FAIL: metadata is not a map: %T = %v\n", resMap["metadata"], resMap["metadata"])
				}
			}
		}
	}

	fmt.Println("\n=== Test Complete ===")
}
