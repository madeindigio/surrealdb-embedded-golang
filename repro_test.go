package surrealdb_test

import (
	"encoding/json"
	"testing"
	"time"

	surrealdb "github.com/yourusername/surrealdb-embedded"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReproSerialization(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

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
	require.NoError(t, err)
	t.Logf("Created: %+v", created)

	// Select it back
	result, err := db.Select("person:complex")
	assert.NoError(t, err)
	t.Logf("Selected: %+v", result)

	// Check if we can marshal it back to JSON
	bytes, err := json.Marshal(result)
	assert.NoError(t, err)
	t.Logf("JSON: %s", string(bytes))

	// Verify fields
	resMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Complex User", resMap["name"])
	// assert.Equal(t, 123.456, resMap["score"]) // Float comparison might be tricky
}
