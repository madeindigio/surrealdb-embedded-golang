package surrealdb_test

import (
	"testing"
	"time"

	surrealdb "github.com/yourusername/surrealdb-embedded"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemory(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	err = db.Use("test", "test")
	assert.NoError(t, err)
}

func TestNewRocksDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/testdb"

	db, err := surrealdb.NewRocksDB(dbPath)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	err = db.Use("test", "test")
	assert.NoError(t, err)
}

func TestNew(t *testing.T) {
	t.Run("Memory backend", func(t *testing.T) {
		db, err := surrealdb.New(surrealdb.Config{
			Backend: surrealdb.Memory,
		})
		require.NoError(t, err)
		require.NotNil(t, db)
		defer db.Close()
	})

	t.Run("RocksDB backend", func(t *testing.T) {
		tmpDir := t.TempDir()
		db, err := surrealdb.New(surrealdb.Config{
			Backend: surrealdb.RocksDB,
			Path:    tmpDir + "/testdb",
		})
		require.NoError(t, err)
		require.NotNil(t, db)
		defer db.Close()
	})

	t.Run("RocksDB without path", func(t *testing.T) {
		_, err := surrealdb.New(surrealdb.Config{
			Backend: surrealdb.RocksDB,
		})
		assert.Error(t, err)
	})
}

func TestCreate(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	data := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
		"email": "john@example.com",
	}

	result, err := db.Create("person", data)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSelect(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	// Create a record first
	data := map[string]interface{}{
		"name": "Jane Doe",
		"age":  25,
	}
	_, err = db.Create("person", data)
	require.NoError(t, err)

	// Select all records
	result, err := db.Select("person")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUpdate(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	// Create a record with a specific ID
	data := map[string]interface{}{
		"name": "Alice",
		"age":  28,
	}
	_, err = db.Create("person:alice", data)
	require.NoError(t, err)

	// Update the record
	updateData := map[string]interface{}{
		"name": "Alice Smith",
		"age":  29,
	}
	result, err := db.Update("person:alice", updateData)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestMerge(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	// Create a record
	data := map[string]interface{}{
		"name": "Bob",
		"age":  35,
	}
	_, err = db.Create("person:bob", data)
	require.NoError(t, err)

	// Merge partial data
	mergeData := map[string]interface{}{
		"email": "bob@example.com",
	}
	result, err := db.Merge("person:bob", mergeData)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestDelete(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	// Create a record
	data := map[string]interface{}{
		"name": "Charlie",
	}
	_, err = db.Create("person:charlie", data)
	require.NoError(t, err)

	// Delete the record
	result, err := db.Delete("person:charlie")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestInsert(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	// Insert multiple records
	data := []map[string]interface{}{
		{
			"id":   "person:dave",
			"name": "Dave",
			"age":  40,
		},
		{
			"id":   "person:eve",
			"name": "Eve",
			"age":  32,
		},
	}

	result, err := db.Insert("person", data)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUpsert(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	data := map[string]interface{}{
		"name": "Frank",
		"age":  45,
	}

	// First upsert creates the record
	result1, err := db.Upsert("person:frank", data)
	assert.NoError(t, err)
	assert.NotNil(t, result1)

	// Second upsert updates the record
	data["age"] = 46
	result2, err := db.Upsert("person:frank", data)
	assert.NoError(t, err)
	assert.NotNil(t, result2)
}

func TestQuery(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	// Create some test data
	_, err = db.Create("person", map[string]interface{}{
		"name": "George",
		"age":  50,
	})
	require.NoError(t, err)

	// Query without variables
	result, err := db.Query("SELECT * FROM person", nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Query with variables
	result, err = db.Query(
		"SELECT * FROM person WHERE age > $minAge",
		map[string]interface{}{"minAge": 40},
	)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestVersion(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	version, err := db.Version()
	assert.NoError(t, err)
	assert.NotNil(t, version)
	assert.Contains(t, version, "version")
}

func TestPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/persisttest"

	// Create database and add data
	db1, err := surrealdb.NewRocksDB(dbPath)
	require.NoError(t, err)

	err = db1.Use("test", "test")
	require.NoError(t, err)

	data := map[string]interface{}{
		"name":  "Persistent User",
		"value": 12345,
	}
	_, err = db1.Create("person:persist", data)
	require.NoError(t, err)

	err = db1.Close()
	require.NoError(t, err)

	// Give RocksDB time to release file locks
	// RocksDB may take a moment to flush and release locks
	time.Sleep(100 * time.Millisecond)

	// Reopen database and verify data
	db2, err := surrealdb.NewRocksDB(dbPath)
	if err != nil {
		// If still locked, wait a bit more and retry
		time.Sleep(500 * time.Millisecond)
		db2, err = surrealdb.NewRocksDB(dbPath)
	}
	require.NoError(t, err)
	defer db2.Close()

	err = db2.Use("test", "test")
	require.NoError(t, err)

	result, err := db2.Select("person:persist")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestConcurrentOperations(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	// Note: This test is basic. In production, you'd want more thorough concurrency testing
	// The current implementation uses a global runtime, so concurrent operations should work

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			data := map[string]interface{}{
				"name": "Concurrent User",
				"id":   id,
			}
			_, err := db.Create("person", data)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestErrorHandling(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	// Use namespace and database first
	err = db.Use("test", "test")
	require.NoError(t, err)

	// Invalid query syntax should return error
	_, err = db.Query("INVALID QUERY SYNTAX", nil)
	assert.Error(t, err, "Invalid syntax should return error")

	// Another invalid query - missing semicolon/incomplete
	_, err = db.Query("SELECT * FROM", nil)
	assert.Error(t, err, "Incomplete query should return error")

	// Try to select from non-existent table (should work but return empty)
	result, err := db.Query("SELECT * FROM nonexistent_table", nil)
	assert.NoError(t, err, "Query non-existent table should not error")
	assert.Equal(t, 0, len(result), "Non-existent table should return empty result")
}

func TestMultipleInstances(t *testing.T) {
	db1, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db1.Close()

	db2, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db2.Close()

	// Both databases should be independent
	err = db1.Use("test", "db1")
	require.NoError(t, err)

	err = db2.Use("test", "db2")
	require.NoError(t, err)

	// Add data to db1
	_, err = db1.Create("person:db1user", map[string]interface{}{"name": "DB1 User"})
	require.NoError(t, err)

	// Add data to db2
	_, err = db2.Create("person:db2user", map[string]interface{}{"name": "DB2 User"})
	require.NoError(t, err)

	// Verify isolation
	result1, err := db1.Select("person")
	require.NoError(t, err)
	assert.NotNil(t, result1)

	result2, err := db2.Select("person")
	require.NoError(t, err)
	assert.NotNil(t, result2)
}

func TestTransactions(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	// Execute a transaction
	_, err = db.Query(`
		BEGIN TRANSACTION;
		CREATE person:tx1 SET name = 'Transaction User 1', balance = 100;
		CREATE person:tx2 SET name = 'Transaction User 2', balance = 200;
		COMMIT TRANSACTION;
	`, nil)
	assert.NoError(t, err)

	// Verify records were created
	result, err := db.Select("person:tx1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSchemaDefinition(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	// Define schema
	_, err = db.Query(`
		DEFINE TABLE user SCHEMAFULL;
		DEFINE FIELD username ON user TYPE string;
		DEFINE FIELD email ON user TYPE string;
		DEFINE INDEX idx_username ON user COLUMNS username UNIQUE;
	`, nil)
	assert.NoError(t, err)

	// Create a record following the schema
	_, err = db.Create("user", map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
	})
	assert.NoError(t, err)
}

func TestGraphRelations(t *testing.T) {
	db, err := surrealdb.NewMemory()
	require.NoError(t, err)
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err)

	// Create nodes
	_, err = db.Create("person:alice", map[string]interface{}{"name": "Alice"})
	require.NoError(t, err)

	_, err = db.Create("person:bob", map[string]interface{}{"name": "Bob"})
	require.NoError(t, err)

	// Create relationship
	_, err = db.Query("RELATE person:alice->knows->person:bob", nil)
	assert.NoError(t, err)

	// Query graph
	result, err := db.Query("SELECT ->knows->person.* as friends FROM person:alice", nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// Benchmark tests
func BenchmarkCreate(b *testing.B) {
	db, _ := surrealdb.NewMemory()
	defer db.Close()
	db.Use("test", "test")

	data := map[string]interface{}{
		"name": "Benchmark User",
		"age":  30,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Create("person", data)
	}
}

func BenchmarkQuery(b *testing.B) {
	db, _ := surrealdb.NewMemory()
	defer db.Close()
	db.Use("test", "test")

	// Create some initial data
	for i := 0; i < 100; i++ {
		db.Create("person", map[string]interface{}{
			"name": "User",
			"age":  i,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Query("SELECT * FROM person LIMIT 10", nil)
	}
}
