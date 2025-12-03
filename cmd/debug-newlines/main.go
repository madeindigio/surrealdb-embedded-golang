package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	surrealdb "github.com/madeindigio/surrealdb-embedded-golang"
)

func main() {
	fmt.Println("=== SurrealDB Newline Serialization Debug Tool ===\n")

	// Initialize database
	db, err := surrealdb.NewMemory()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	err = db.Use("test", "test")
	if err != nil {
		log.Fatalf("Failed to use namespace/database: %v", err)
	}

	fmt.Println("✓ Database initialized\n")

	// Test 1: Create with newlines
	fmt.Println("--- Test 1: CREATE with newlines ---")
	createData := map[string]interface{}{
		"title":       "Test Document",
		"description": "Line 1\nLine 2\nLine 3",
		"content":     "First paragraph\n\nSecond paragraph with blank line above",
	}

	fmt.Printf("Creating with data:\n")
	printData(createData)

	result, err := db.Create("documents:test1", createData)
	if err != nil {
		log.Fatalf("CREATE failed: %v", err)
	}
	fmt.Printf("\nCreate result:\n")
	printData(result)

	// Verify
	results, err := db.Query("SELECT * FROM documents:test1", nil)
	if err != nil {
		log.Fatalf("SELECT failed: %v", err)
	}
	fmt.Printf("\nSELECT result:\n")
	printData(results[0])
	verifyNewlines(results[0].(map[string]interface{}), "description", "Line 1\nLine 2\nLine 3")

	// Test 2: UPDATE with newlines
	fmt.Println("\n--- Test 2: UPDATE with newlines ---")
	updateData := map[string]interface{}{
		"title":       "Updated Document",
		"description": "Updated Line 1\nUpdated Line 2\nUpdated Line 3",
		"content":     "New first paragraph\n\nNew second paragraph",
		"extra":       "Additional field\nwith\nmultiple\nnewlines",
	}

	fmt.Printf("Updating with data:\n")
	printData(updateData)

	result, err = db.Update("documents:test1", updateData)
	if err != nil {
		log.Fatalf("UPDATE failed: %v", err)
	}
	fmt.Printf("\nUpdate result:\n")
	printData(result)

	// Verify
	results, err = db.Query("SELECT * FROM documents:test1", nil)
	if err != nil {
		log.Fatalf("SELECT after UPDATE failed: %v", err)
	}
	fmt.Printf("\nSELECT after UPDATE:\n")
	updated := results[0].(map[string]interface{})
	printData(updated)

	// Detailed verification
	verifyNewlines(updated, "description", "Updated Line 1\nUpdated Line 2\nUpdated Line 3")
	verifyNewlines(updated, "content", "New first paragraph\n\nNew second paragraph")
	verifyNewlines(updated, "extra", "Additional field\nwith\nmultiple\nnewlines")

	// Test 3: MERGE with newlines
	fmt.Println("\n--- Test 3: MERGE with newlines ---")
	mergeData := map[string]interface{}{
		"description": "Merged description\nwith newlines",
		"notes":       "New field\nadded via\nMERGE",
	}

	fmt.Printf("Merging with data:\n")
	printData(mergeData)

	result, err = db.Merge("documents:test1", mergeData)
	if err != nil {
		log.Fatalf("MERGE failed: %v", err)
	}
	fmt.Printf("\nMerge result:\n")
	printData(result)

	// Verify
	results, err = db.Query("SELECT * FROM documents:test1", nil)
	if err != nil {
		log.Fatalf("SELECT after MERGE failed: %v", err)
	}
	fmt.Printf("\nSELECT after MERGE:\n")
	merged := results[0].(map[string]interface{})
	printData(merged)

	verifyNewlines(merged, "description", "Merged description\nwith newlines")
	verifyNewlines(merged, "notes", "New field\nadded via\nMERGE")
	verifyNewlines(merged, "title", "Updated Document") // Should be unchanged

	// Test 4: Parameterized query with newlines
	fmt.Println("\n--- Test 4: Parameterized query with newlines ---")
	params := map[string]interface{}{
		"newdesc": "Parameterized value\nwith newlines\nand more",
	}

	fmt.Printf("Executing query with params:\n")
	printData(params)

	results, err = db.Query("UPDATE documents:test1 SET description = $newdesc", params)
	if err != nil {
		log.Fatalf("Parameterized UPDATE failed: %v", err)
	}
	fmt.Printf("\nParameterized update result:\n")
	printData(results)

	// Verify
	results, err = db.Query("SELECT * FROM documents:test1", nil)
	if err != nil {
		log.Fatalf("SELECT after parameterized UPDATE failed: %v", err)
	}
	fmt.Printf("\nSELECT after parameterized UPDATE:\n")
	paramResult := results[0].(map[string]interface{})
	printData(paramResult)

	verifyNewlines(paramResult, "description", "Parameterized value\nwith newlines\nand more")

	// Test 5: Edge cases
	fmt.Println("\n--- Test 5: Edge cases ---")
	edgeCases := []struct {
		name     string
		value    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single newline", "\n", "\n"},
		{"Only newlines", "\n\n\n", "\n\n\n"},
		{"Leading newline", "\nstarts here", "\nstarts here"},
		{"Trailing newline", "ends here\n", "ends here\n"},
		{"Tabs and newlines", "tab:\there\nnewline\nhere", "tab:\there\nnewline\nhere"},
		{"Unicode with newlines", "你好\n世界\n😀", "你好\n世界\n😀"},
	}

	for i, tc := range edgeCases {
		fmt.Printf("\nEdge case %d: %s\n", i+1, tc.name)
		fmt.Printf("Value (escaped): %q\n", tc.value)
		fmt.Printf("Value bytes: %v\n", []byte(tc.value))

		updateData := map[string]interface{}{
			"test_field": tc.value,
		}

		_, err := db.Update("documents:test1", updateData)
		if err != nil {
			log.Fatalf("UPDATE failed for edge case '%s': %v", tc.name, err)
		}

		results, err := db.Query("SELECT * FROM documents:test1", nil)
		if err != nil {
			log.Fatalf("SELECT failed for edge case '%s': %v", tc.name, err)
		}

		retrieved := results[0].(map[string]interface{})["test_field"]
		fmt.Printf("Retrieved (escaped): %q\n", retrieved)
		if str, ok := retrieved.(string); ok {
			fmt.Printf("Retrieved bytes: %v\n", []byte(str))
		}

		if retrieved != tc.expected {
			fmt.Printf("❌ MISMATCH! Expected: %q, Got: %q\n", tc.expected, retrieved)
			os.Exit(1)
		} else {
			fmt.Printf("✓ Matches expected value\n")
		}
	}

	fmt.Println("\n=== All tests PASSED ✓ ===")
	fmt.Println("\nIf you're seeing issues in your application but these tests pass,")
	fmt.Println("the problem is likely in how you're:")
	fmt.Println("  1. Preparing the data before passing to Update()")
	fmt.Println("  2. Reading/displaying the data after retrieval")
	fmt.Println("  3. Using a wrapper or middleware that modifies the data")
	fmt.Println("\nAdd similar debug output to your code to identify where the newlines are lost.")
}

func printData(data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "  ", "  ")
	if err != nil {
		fmt.Printf("  Error marshaling: %v\n", err)
		fmt.Printf("  Raw: %+v\n", data)
		return
	}
	fmt.Printf("  %s\n", string(jsonBytes))
}

func verifyNewlines(record map[string]interface{}, field string, expected string) {
	actual, ok := record[field]
	if !ok {
		fmt.Printf("❌ Field '%s' not found in record\n", field)
		os.Exit(1)
	}

	actualStr, ok := actual.(string)
	if !ok {
		fmt.Printf("❌ Field '%s' is not a string, got: %T\n", field, actual)
		os.Exit(1)
	}

	if actualStr != expected {
		fmt.Printf("\n❌ VERIFICATION FAILED for field '%s'\n", field)
		fmt.Printf("Expected: %q\n", expected)
		fmt.Printf("Actual:   %q\n", actualStr)
		fmt.Printf("Expected bytes: %v\n", []byte(expected))
		fmt.Printf("Actual bytes:   %v\n", []byte(actualStr))
		os.Exit(1)
	}

	fmt.Printf("✓ Field '%s' verified: %q\n", field, actualStr)
}
