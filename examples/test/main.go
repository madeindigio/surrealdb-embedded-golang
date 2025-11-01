package main

import (
    "fmt"
    "log"
    surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
    fmt.Println("Creating database...")
    db, err := surrealdb.NewMemory()
    if err != nil {
        log.Fatalf("Failed to create database: %v", err)
    }
    defer db.Close()
    fmt.Println("✓ Database created successfully!")
    
    fmt.Println("\nUsing namespace and database...")
    err = db.Use("test", "test")
    if err != nil {
        log.Fatalf("Failed to use namespace/database: %v", err)
    }
    fmt.Println("✓ Namespace/database selected successfully!")
    
    fmt.Println("\nCreating a record...")
    result, err := db.Create("person", map[string]interface{}{
        "name": "Test User",
        "age":  25,
    })
    if err != nil {
        log.Fatalf("Failed to create record: %v", err)
    }
    fmt.Printf("✓ Record created: %+v\n", result)
    
    fmt.Println("\n✓ All tests passed!")
}
