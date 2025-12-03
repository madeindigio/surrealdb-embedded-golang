package surrealdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateWithNewlines(t *testing.T) {
	// Initialize database
	db, err := NewMemory()
	require.NoError(t, err, "Failed to initialize database")
	defer db.Close()

	// Use test namespace and database
	err = db.Use("test", "test")
	require.NoError(t, err, "Failed to use namespace/database")

	// Test 1: Create a record with newlines in the content
	t.Run("Create with newlines", func(t *testing.T) {
		data := map[string]interface{}{
			"title":       "Test Document",
			"description": "Line 1\nLine 2\nLine 3",
			"multiline":   "First line\n\nThird line with blank line above",
			"code":        "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}",
		}

		result, err := db.Create("documents:test1", data)
		require.NoError(t, err, "Failed to create record")
		require.NotNil(t, result, "Result should not be nil")

		t.Logf("Create result: %+v", result)

		// Verify the created record
		results, err := db.Query("SELECT * FROM documents:test1", nil)
		require.NoError(t, err, "Failed to query record")
		require.Len(t, results, 1, "Should return 1 result")

		record := results[0].(map[string]interface{})
		assert.Equal(t, "Test Document", record["title"], "Title should match")
		assert.Equal(t, "Line 1\nLine 2\nLine 3", record["description"], "Description with newlines should match")
		assert.Equal(t, "First line\n\nThird line with blank line above", record["multiline"], "Multiline with blank line should match")
		assert.Equal(t, "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}", record["code"], "Code with newlines should match")
	})

	// Test 2: Update with newlines - this is where the bug occurs
	t.Run("Update with newlines", func(t *testing.T) {
		// First verify current state
		results, err := db.Query("SELECT * FROM documents:test1", nil)
		require.NoError(t, err, "Failed to query before update")
		require.Len(t, results, 1, "Should have 1 record")

		beforeUpdate := results[0].(map[string]interface{})
		t.Logf("Before update: %+v", beforeUpdate)

		// Update with new content containing newlines
		updateData := map[string]interface{}{
			"title":       "Updated Document",
			"description": "Updated Line 1\nUpdated Line 2\nUpdated Line 3",
			"multiline":   "New first line\n\nNew third line",
			"code":        "fn main() {\n\tprintln!(\"Rust\");\n}",
			"newField":    "Additional\nfield\nwith\nnewlines",
		}

		result, err := db.Update("documents:test1", updateData)
		require.NoError(t, err, "Failed to update record")
		require.NotNil(t, result, "Update result should not be nil")

		t.Logf("Update result: %+v", result)

		// Query again to verify the update
		results, err = db.Query("SELECT * FROM documents:test1", nil)
		require.NoError(t, err, "Failed to query after update")
		require.Len(t, results, 1, "Should still have 1 record")

		updated := results[0].(map[string]interface{})
		t.Logf("After update: %+v", updated)

		// These assertions should pass but might fail if the bug exists
		assert.Equal(t, "Updated Document", updated["title"], "Title should be updated")
		assert.Equal(t, "Updated Line 1\nUpdated Line 2\nUpdated Line 3", updated["description"], "Description should be updated with newlines preserved")
		assert.Equal(t, "New first line\n\nNew third line", updated["multiline"], "Multiline should be updated with blank lines preserved")
		assert.Equal(t, "fn main() {\n\tprintln!(\"Rust\");\n}", updated["code"], "Code should be updated with newlines preserved")
		assert.Equal(t, "Additional\nfield\nwith\nnewlines", updated["newField"], "New field with newlines should be added")
	})

	// Test 3: Merge with newlines
	t.Run("Merge with newlines", func(t *testing.T) {
		// Merge partial data with newlines
		mergeData := map[string]interface{}{
			"description": "Merged description\nwith newlines\nand tabs\tlike this",
			"notes":       "New notes field\n\nWith blank lines\n\nMultiple times",
		}

		result, err := db.Merge("documents:test1", mergeData)
		require.NoError(t, err, "Failed to merge record")
		require.NotNil(t, result, "Merge result should not be nil")

		t.Logf("Merge result: %+v", result)

		// Verify merge
		results, err := db.Query("SELECT * FROM documents:test1", nil)
		require.NoError(t, err, "Failed to query after merge")
		require.Len(t, results, 1, "Should still have 1 record")

		merged := results[0].(map[string]interface{})
		t.Logf("After merge: %+v", merged)

		assert.Equal(t, "Merged description\nwith newlines\nand tabs\tlike this", merged["description"], "Merged description should preserve newlines and tabs")
		assert.Equal(t, "New notes field\n\nWith blank lines\n\nMultiple times", merged["notes"], "Merged notes should preserve blank lines")
		// Other fields should remain
		assert.Equal(t, "Updated Document", merged["title"], "Title should remain from update")
	})

	// Test 4: Query with newlines in WHERE clause
	t.Run("Query with newlines in parameters", func(t *testing.T) {
		// Create another record for comparison
		data := map[string]interface{}{
			"title":   "Another Document",
			"content": "Simple content\nwith newlines",
		}
		_, err := db.Create("documents:test2", data)
		require.NoError(t, err, "Failed to create second record")

		// Query with newline in parameter
		results, err := db.Query("SELECT * FROM documents WHERE content = $content", map[string]interface{}{
			"content": "Simple content\nwith newlines",
		})
		require.NoError(t, err, "Failed to query with newline parameter")
		require.Len(t, results, 1, "Should find 1 record matching newline content")

		found := results[0].(map[string]interface{})
		assert.Equal(t, "Another Document", found["title"], "Should find the correct document")
		assert.Equal(t, "Simple content\nwith newlines", found["content"], "Content with newlines should match")
	})

	// Test 5: Special characters combinations
	t.Run("Special characters with newlines", func(t *testing.T) {
		specialData := map[string]interface{}{
			"json_string":   `{"key": "value\nwith newline"}`,
			"quote_newline": "Quote: \"Hello\nWorld\"",
			"backslash":     "Path: C:\\Users\\test\nNext line",
			"unicode":       "Unicode: 你好\n世界\nEmoji: 😀\n🎉",
			"tabs_newlines": "Tab:\there\nNewline:\nhere",
		}

		_, err := db.Create("documents:special", specialData)
		require.NoError(t, err, "Failed to create record with special chars")

		// Update with similar special characters
		updateSpecial := map[string]interface{}{
			"json_string":   `{"updated": "new value\nwith newline"}`,
			"quote_newline": "Updated: \"Hi\nThere\"",
			"backslash":     "New path: D:\\Data\\test\nAnother line",
			"unicode":       "更新: こんにちは\n世界\n🌟\n✨",
			"tabs_newlines": "Updated tab:\there\nNew line:\nhere",
		}

		result, err := db.Update("documents:special", updateSpecial)
		require.NoError(t, err, "Failed to update record with special chars")
		require.NotNil(t, result, "Update result should not be nil")

		// Verify
		results, err := db.Query("SELECT * FROM documents:special", nil)
		require.NoError(t, err, "Failed to query special chars record")
		require.Len(t, results, 1, "Should have 1 record")

		updated := results[0].(map[string]interface{})
		assert.Equal(t, `{"updated": "new value\nwith newline"}`, updated["json_string"], "JSON string with newlines should be preserved")
		assert.Equal(t, "Updated: \"Hi\nThere\"", updated["quote_newline"], "Quotes and newlines should be preserved")
		assert.Equal(t, "New path: D:\\Data\\test\nAnother line", updated["backslash"], "Backslashes and newlines should be preserved")
		assert.Equal(t, "更新: こんにちは\n世界\n🌟\n✨", updated["unicode"], "Unicode and newlines should be preserved")
		assert.Equal(t, "Updated tab:\there\nNew line:\nhere", updated["tabs_newlines"], "Tabs and newlines should be preserved")
	})

	// Test 6: Edge case - Empty string vs newline
	t.Run("Empty string vs newline edge cases", func(t *testing.T) {
		edgeData := map[string]interface{}{
			"empty":          "",
			"single_newline": "\n",
			"spaces_newline": "   \n   ",
			"mixed":          "",
		}

		_, err := db.Create("documents:edge", edgeData)
		require.NoError(t, err, "Failed to create edge case record")

		// Update empty field with newline
		updateEdge := map[string]interface{}{
			"empty":          "now has\nvalue",
			"single_newline": "replaced\nnewline",
			"spaces_newline": "\n\n\n",
			"mixed":          "was empty\nnow filled",
		}

		result, err := db.Update("documents:edge", updateEdge)
		require.NoError(t, err, "Failed to update edge case")
		require.NotNil(t, result, "Result should not be nil")

		// Verify
		results, err := db.Query("SELECT * FROM documents:edge", nil)
		require.NoError(t, err, "Failed to query edge case")
		require.Len(t, results, 1, "Should have 1 record")

		updated := results[0].(map[string]interface{})
		assert.Equal(t, "now has\nvalue", updated["empty"], "Empty field updated with newline")
		assert.Equal(t, "replaced\nnewline", updated["single_newline"], "Single newline replaced")
		assert.Equal(t, "\n\n\n", updated["spaces_newline"], "Multiple newlines preserved")
		assert.Equal(t, "was empty\nnow filled", updated["mixed"], "Mixed case preserved")
	})

	// Test 7: Raw JSON string with escaped newlines
	t.Run("Raw JSON with escaped newlines", func(t *testing.T) {
		// Simulate what might happen if JSON is double-encoded
		rawData := map[string]interface{}{
			"normal": "test\nvalue",
			// This simulates if someone passes already JSON-escaped string
			"escaped": `line1\nline2`,
		}

		_, err := db.Create("documents:rawjson", rawData)
		require.NoError(t, err, "Failed to create raw JSON record")

		// Update with literal backslash-n vs actual newline
		updateRaw := map[string]interface{}{
			"normal":  "updated\nwith\nnewlines",
			"escaped": `updated\\nwith\\nescaped`,
		}

		result, err := db.Update("documents:rawjson", updateRaw)
		require.NoError(t, err, "Failed to update raw JSON")
		require.NotNil(t, result, "Result should not be nil")

		// Verify
		results, err := db.Query("SELECT * FROM documents:rawjson", nil)
		require.NoError(t, err, "Failed to query raw JSON")
		require.Len(t, results, 1, "Should have 1 record")

		updated := results[0].(map[string]interface{})
		t.Logf("Raw JSON result: %+v", updated)
		assert.Equal(t, "updated\nwith\nnewlines", updated["normal"], "Normal newlines preserved")
	})
}

func TestUpdateWithNewlinesDetailed(t *testing.T) {
	// Initialize database
	db, err := NewMemory()
	require.NoError(t, err, "Failed to initialize database")
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err, "Failed to use namespace/database")

	// Create initial record
	initial := map[string]interface{}{
		"name":  "Initial",
		"value": "original",
	}

	_, err = db.Create("test:newline1", initial)
	require.NoError(t, err, "Failed to create initial record")

	// Test various newline scenarios one by one
	testCases := []struct {
		name     string
		update   map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "Single newline",
			update: map[string]interface{}{
				"name":  "Test 1",
				"value": "line1\nline2",
			},
			expected: map[string]interface{}{
				"name":  "Test 1",
				"value": "line1\nline2",
			},
		},
		{
			name: "Multiple newlines",
			update: map[string]interface{}{
				"name":  "Test 2",
				"value": "line1\nline2\nline3\nline4",
			},
			expected: map[string]interface{}{
				"name":  "Test 2",
				"value": "line1\nline2\nline3\nline4",
			},
		},
		{
			name: "Blank lines",
			update: map[string]interface{}{
				"name":  "Test 3",
				"value": "line1\n\nline3",
			},
			expected: map[string]interface{}{
				"name":  "Test 3",
				"value": "line1\n\nline3",
			},
		},
		{
			name: "Starting with newline",
			update: map[string]interface{}{
				"name":  "Test 4",
				"value": "\nstarts with newline",
			},
			expected: map[string]interface{}{
				"name":  "Test 4",
				"value": "\nstarts with newline",
			},
		},
		{
			name: "Ending with newline",
			update: map[string]interface{}{
				"name":  "Test 5",
				"value": "ends with newline\n",
			},
			expected: map[string]interface{}{
				"name":  "Test 5",
				"value": "ends with newline\n",
			},
		},
		{
			name: "Only newlines",
			update: map[string]interface{}{
				"name":  "Test 6",
				"value": "\n\n\n",
			},
			expected: map[string]interface{}{
				"name":  "Test 6",
				"value": "\n\n\n",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Update
			result, err := db.Update("test:newline1", tc.update)
			require.NoError(t, err, "Failed to update: %s", tc.name)
			require.NotNil(t, result, "Result should not be nil")

			t.Logf("Update result for %s: %+v", tc.name, result)

			// Query back
			results, err := db.Query("SELECT * FROM test:newline1", nil)
			require.NoError(t, err, "Failed to query: %s", tc.name)
			require.Len(t, results, 1, "Should have 1 record")

			actual := results[0].(map[string]interface{})
			t.Logf("Queried result for %s: %+v", tc.name, actual)

			// Verify
			for key, expectedValue := range tc.expected {
				actualValue, ok := actual[key]
				require.True(t, ok, "Field %s should exist", key)

				// Debug output
				if actualValue != expectedValue {
					t.Logf("MISMATCH for %s:", key)
					t.Logf("  Expected: %q (%T)", expectedValue, expectedValue)
					t.Logf("  Actual:   %q (%T)", actualValue, actualValue)

					// Show byte representation
					if expStr, ok := expectedValue.(string); ok {
						if actStr, ok := actualValue.(string); ok {
							t.Logf("  Expected bytes: %v", []byte(expStr))
							t.Logf("  Actual bytes:   %v", []byte(actStr))
						}
					}
				}

				assert.Equal(t, expectedValue, actualValue,
					"Field %s should match for test: %s", key, tc.name)
			}
		})
	}
}

// TestUpdateJSONRoundtrip tests the complete JSON serialization round trip
func TestUpdateJSONRoundtrip(t *testing.T) {
	db, err := NewMemory()
	require.NoError(t, err, "Failed to initialize database")
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err, "Failed to use namespace/database")

	// Create initial record
	initial := map[string]interface{}{
		"text": "initial value",
	}

	_, err = db.Create("test:roundtrip", initial)
	require.NoError(t, err, "Failed to create record")

	// Test various problematic strings
	testCases := []struct {
		name  string
		value string
	}{
		{"simple_newline", "line1\nline2"},
		{"multiple_newlines", "a\nb\nc\nd"},
		{"blank_lines", "line1\n\nline3"},
		{"leading_newline", "\nstarts here"},
		{"trailing_newline", "ends here\n"},
		{"only_newlines", "\n\n\n"},
		{"mixed_whitespace", "tab\there\nnewline\nhere"},
		{"json_like", `{"key": "value\nwith newline"}`},
		{"quotes_newlines", `"quoted\nvalue"`},
		{"backslashes", `path\to\file\nwith newline`},
		{"unicode_newlines", "你好\n世界\n😀"},
		{"carriage_return", "line1\r\nline2"},
		{"all_escapes", "tab:\t newline:\n quote:\" backslash:\\ end"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Update with the test value
			updateData := map[string]interface{}{
				"text": tc.value,
			}

			updateResult, err := db.Update("test:roundtrip", updateData)
			require.NoError(t, err, "Failed to update for case: %s", tc.name)
			require.NotNil(t, updateResult, "Update result should not be nil")

			// Query back
			results, err := db.Query("SELECT * FROM test:roundtrip", nil)
			require.NoError(t, err, "Failed to query for case: %s", tc.name)
			require.Len(t, results, 1, "Should have 1 record")

			record := results[0].(map[string]interface{})
			actual := record["text"]

			if actual != tc.value {
				t.Errorf("Value mismatch for %s:", tc.name)
				t.Errorf("  Expected: %q", tc.value)
				t.Errorf("  Actual:   %q", actual)
				t.Errorf("  Expected bytes: %v", []byte(tc.value))
				if str, ok := actual.(string); ok {
					t.Errorf("  Actual bytes:   %v", []byte(str))
				}
			}

			assert.Equal(t, tc.value, actual, "Value should match exactly for: %s", tc.name)
		})
	}
}

// TestUpdateWithParametersNewlines tests UPDATE with parameterized queries
func TestUpdateWithParametersNewlines(t *testing.T) {
	db, err := NewMemory()
	require.NoError(t, err, "Failed to initialize database")
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err, "Failed to use namespace/database")

	// Create initial record
	_, err = db.Query("CREATE test:param SET value = 'initial'", nil)
	require.NoError(t, err, "Failed to create record")

	// Update using parameterized query with newlines
	params := map[string]interface{}{
		"newvalue": "line1\nline2\nline3",
	}

	results, err := db.Query("UPDATE test:param SET value = $newvalue", params)
	require.NoError(t, err, "Failed to update with parameters")
	require.NotNil(t, results, "Results should not be nil")

	t.Logf("Parameterized update result: %+v", results)

	// Verify
	results, err = db.Query("SELECT * FROM test:param", nil)
	require.NoError(t, err, "Failed to query after param update")
	require.Len(t, results, 1, "Should have 1 record")

	record := results[0].(map[string]interface{})
	assert.Equal(t, "line1\nline2\nline3", record["value"], "Parameterized newlines should be preserved")
}

// TestUpdateDirectQueryInjection tests if newlines cause query parsing issues
func TestUpdateDirectQueryInjection(t *testing.T) {
	db, err := NewMemory()
	require.NoError(t, err, "Failed to initialize database")
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err, "Failed to use namespace/database")

	// Create initial record
	_, err = db.Create("test:inject", map[string]interface{}{
		"field1": "initial",
		"field2": "value",
	})
	require.NoError(t, err, "Failed to create record")

	// Test UPDATE with newlines that could break query parsing
	// This is the critical test case - if the JSON is interpolated directly
	// into the query string without proper handling, newlines could break it
	problemData := map[string]interface{}{
		"field1": "Line 1\nLine 2",
		"field2": "Another\nvalue\nwith\nmultiple\nlines",
		// This could be particularly problematic:
		"field3": "Value with\n-- comment\nand more",
		"field4": "{\n  \"nested\": \"json\n  with newlines\"\n}",
	}

	result, err := db.Update("test:inject", problemData)
	require.NoError(t, err, "Failed to update with problematic newlines")
	require.NotNil(t, result, "Result should not be nil")

	t.Logf("Update result: %+v", result)

	// Verify all fields were updated correctly
	results, err := db.Query("SELECT * FROM test:inject", nil)
	require.NoError(t, err, "Failed to query after update")
	require.Len(t, results, 1, "Should have 1 record")

	updated := results[0].(map[string]interface{})

	// These assertions will fail if newlines break the query
	assert.Equal(t, "Line 1\nLine 2", updated["field1"], "field1 should preserve newlines")
	assert.Equal(t, "Another\nvalue\nwith\nmultiple\nlines", updated["field2"], "field2 should preserve multiple newlines")
	assert.Equal(t, "Value with\n-- comment\nand more", updated["field3"], "field3 should preserve comment-like content")
	assert.Equal(t, "{\n  \"nested\": \"json\n  with newlines\"\n}", updated["field4"], "field4 should preserve JSON-like content with newlines")
}

// TestUpdateMergeWithNewlines tests MERGE operation with newlines
func TestUpdateMergeWithNewlines(t *testing.T) {
	db, err := NewMemory()
	require.NoError(t, err, "Failed to initialize database")
	defer db.Close()

	err = db.Use("test", "test")
	require.NoError(t, err, "Failed to use namespace/database")

	// Create record with initial values
	_, err = db.Create("test:merge", map[string]interface{}{
		"keep":   "this value",
		"change": "old value",
	})
	require.NoError(t, err, "Failed to create record")

	// Merge with newlines - should only update specified fields
	mergeData := map[string]interface{}{
		"change": "new value\nwith newlines\nin it",
		"add":    "new field\nalso with\nnewlines",
	}

	result, err := db.Merge("test:merge", mergeData)
	require.NoError(t, err, "Failed to merge with newlines")
	require.NotNil(t, result, "Merge result should not be nil")

	// Verify merge
	results, err := db.Query("SELECT * FROM test:merge", nil)
	require.NoError(t, err, "Failed to query after merge")
	require.Len(t, results, 1, "Should have 1 record")

	merged := results[0].(map[string]interface{})
	assert.Equal(t, "this value", merged["keep"], "Unmodified field should remain")
	assert.Equal(t, "new value\nwith newlines\nin it", merged["change"], "Merged field should have newlines")
	assert.Equal(t, "new field\nalso with\nnewlines", merged["add"], "Added field should have newlines")
}
