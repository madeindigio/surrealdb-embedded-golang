package surrealdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestKnowledgeBasePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/kb_test"

	// Create database and add document
	db1, err := NewFromURL(fmt.Sprintf("surrealkv://%s", dbPath))
	require.NoError(t, err)

	err = db1.Use("test", "test")
	require.NoError(t, err)

	// Test 1: CREATE with CONTENT (exactly like SaveDocument)
	fmt.Println("\n=== Test 1: CREATE with CONTENT ===")
	result, err := db1.Query(`
		CREATE knowledge_base CONTENT {
			file_path: $file_path,
			content: $content,
			embedding: $embedding
		}
	`, map[string]interface{}{
		"file_path": "/test/file1.md",
		"content":   "Test content",
		"embedding": []float64{0.1, 0.2, 0.3},
	})
	require.NoError(t, err)
	fmt.Printf("CREATE result: %+v\n", result)
	fmt.Printf("CREATE result length: %d\n", len(result))

	// Test 2: SELECT to verify immediately
	fmt.Println("\n=== Test 2: SELECT to verify ===")
	result, err = db1.Query("SELECT * FROM knowledge_base WHERE file_path = $file_path", map[string]interface{}{
		"file_path": "/test/file1.md",
	})
	require.NoError(t, err)
	fmt.Printf("SELECT result: %+v\n", result)
	fmt.Printf("SELECT result length: %d\n", len(result))
	if len(result) > 0 {
		fmt.Printf("SELECT result type: %T\n", result[0])
		if arr, ok := result[0].([]interface{}); ok {
			fmt.Printf("SELECT array length: %d\n", len(arr))
			if len(arr) > 0 {
				fmt.Printf("SELECT first item: %+v\n", arr[0])
			} else {
				fmt.Println("WARNING: SELECT returned empty array!")
			}
		}
	}

	// Close the database
	err = db1.Close()
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	// Test 3: Reopen and verify persistence
	fmt.Println("\n=== Test 3: Reopen and verify persistence ===")
	db2, err := NewFromURL(fmt.Sprintf("surrealkv://%s", dbPath))
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		db2, err = NewFromURL(fmt.Sprintf("surrealkv://%s", dbPath))
	}
	require.NoError(t, err)
	defer db2.Close()

	err = db2.Use("test", "test")
	require.NoError(t, err)

	// SELECT after reopen
	result, err = db2.Query("SELECT * FROM knowledge_base", nil)
	require.NoError(t, err)
	fmt.Printf("SELECT after reopen result: %+v\n", result)
	fmt.Printf("SELECT after reopen length: %d\n", len(result))
	
	if len(result) == 0 {
		fmt.Println("ERROR: No results after reopen - data did not persist!")
		t.Fatal("Data did not persist across reopens")
	}
	
	if arr, ok := result[0].([]interface{}); ok {
		fmt.Printf("Array length after reopen: %d\n", len(arr))
		if len(arr) == 0 {
			fmt.Println("ERROR: Empty array after reopen - data did not persist!")
			t.Fatal("Data did not persist - empty result")
		}
		fmt.Printf("First item after reopen: %+v\n", arr[0])
		fmt.Println("SUCCESS: Data persisted correctly!")
	}
}
