package main

import (
    "fmt"
    "log"
    surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
    fmt.Println("Testing Vector/Embeddings Support")
    fmt.Println("==================================")
    
    db, err := surrealdb.NewMemory()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    err = db.Use("test", "test")
    if err != nil {
        log.Fatal(err)
    }
    
    // Test 1: Create table with vector field
    fmt.Println("\n1. Creating schema with vector field...")
    _, err = db.Query(`
        DEFINE TABLE document SCHEMAFULL;
        DEFINE FIELD content ON document TYPE string;
        DEFINE FIELD embedding ON document TYPE array;
        DEFINE INDEX embedding_idx ON document FIELDS embedding MTREE DIMENSION 3;
    `, nil)
    if err != nil {
        fmt.Printf("   ✗ Schema creation failed: %v\n", err)
    } else {
        fmt.Println("   ✓ Schema created successfully")
    }
    
    // Test 2: Insert document with embedding vector
    fmt.Println("\n2. Inserting document with embedding vector...")
    result, err := db.Create("document", map[string]interface{}{
        "content": "Hello world",
        "embedding": []float64{0.1, 0.2, 0.3},
    })
    if err != nil {
        fmt.Printf("   ✗ Insert failed: %v\n", err)
    } else {
        fmt.Printf("   ✓ Document created: %+v\n", result)
    }
    
    // Test 3: Vector similarity search
    fmt.Println("\n3. Testing vector similarity search...")
    results, err := db.Query(`
        SELECT * FROM document 
        WHERE embedding <|3|> [0.15, 0.25, 0.35]
        ORDER BY embedding <|3|> [0.15, 0.25, 0.35] ASC
        LIMIT 5
    `, nil)
    if err != nil {
        fmt.Printf("   ✗ Vector search failed: %v\n", err)
    } else {
        fmt.Printf("   ✓ Vector search results: %+v\n", results)
    }
    
    // Test 4: Insert multiple documents with embeddings
    fmt.Println("\n4. Inserting multiple documents...")
    docs := []map[string]interface{}{
        {"content": "AI and machine learning", "embedding": []float64{0.8, 0.2, 0.1}},
        {"content": "Deep learning models", "embedding": []float64{0.7, 0.3, 0.2}},
        {"content": "Natural language processing", "embedding": []float64{0.6, 0.4, 0.1}},
    }
    
    for _, doc := range docs {
        _, err := db.Create("document", doc)
        if err != nil {
            fmt.Printf("   ✗ Failed to insert: %v\n", err)
        } else {
            fmt.Printf("   ✓ Inserted: %s\n", doc["content"])
        }
    }
    
    // Test 5: Vector distance calculation
    fmt.Println("\n5. Testing vector distance functions...")
    results, err = db.Query(`
        SELECT 
            content,
            vector::distance::euclidean(embedding, [0.5, 0.3, 0.2]) AS euclidean_dist,
            vector::distance::cosine(embedding, [0.5, 0.3, 0.2]) AS cosine_dist
        FROM document
        ORDER BY euclidean_dist ASC
    `, nil)
    if err != nil {
        fmt.Printf("   ✗ Distance calculation failed: %v\n", err)
    } else {
        fmt.Printf("   ✓ Distance results: %+v\n", results)
    }
    
    // Test 6: KNN search
    fmt.Println("\n6. Testing KNN search...")
    results, err = db.Query(`
        SELECT content, embedding 
        FROM document 
        WHERE embedding <|3,EUCLIDEAN|> [0.6, 0.3, 0.2]
        LIMIT 3
    `, nil)
    if err != nil {
        fmt.Printf("   ✗ KNN search failed: %v\n", err)
    } else {
        fmt.Printf("   ✓ KNN results: %+v\n", results)
    }
    
    fmt.Println("\n=== Vector Support Test Complete ===")
}
