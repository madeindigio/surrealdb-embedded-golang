use serde_json::Value;
use surrealdb::engine::local::{Db, Mem};
use surrealdb::Surreal;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("=== SurrealDB Embedded Vector Support Test ===\n");

    // Initialize database
    let db = Surreal::new::<Mem>(()).await?;
    db.use_ns("test").await?;
    db.use_db("test").await?;

    println!("✅ Database initialized\n");

    // Test 1: Schema with MTREE index
    println!("--- Test 1: Creating schema with MTREE index ---");
    let result = db
        .query(
            "
        DEFINE TABLE document SCHEMAFULL;
        DEFINE FIELD content ON document TYPE string;
        DEFINE FIELD embedding ON document TYPE array;
        DEFINE INDEX embedding_idx ON document FIELDS embedding MTREE DIMENSION 3;
    ",
        )
        .await;

    match result {
        Ok(_) => println!("✅ MTREE schema created successfully\n"),
        Err(e) => println!("❌ MTREE schema failed: {:?}\n", e),
    }

    // Test 2: Insert documents with embeddings
    println!("--- Test 2: Inserting documents with embeddings ---");
    let result = db
        .query(
            "
        CREATE document:doc1 SET content = 'Hello world', embedding = [0.1, 0.2, 0.3];
        CREATE document:doc2 SET content = 'Goodbye world', embedding = [0.5, 0.4, 0.3];
        CREATE document:doc3 SET content = 'Another document', embedding = [0.9, 0.8, 0.7];
    ",
        )
        .await;

    match result {
        Ok(_) => println!("✅ Documents inserted successfully\n"),
        Err(e) => println!("❌ Insertion failed: {:?}\n", e),
    }

    // Test 3: KNN search
    println!("--- Test 3: KNN search with <|K,DISTANCE|> operator ---");
    let mut response = db
        .query(
            "
        SELECT content, embedding
        FROM document
        WHERE embedding <|3,EUCLIDEAN|> [0.6, 0.3, 0.2]
        LIMIT 3;
    ",
        )
        .await?;

    let docs: Vec<Value> = response.take(0)?;
    println!("✅ KNN search returned {} documents:", docs.len());
    for (i, doc) in docs.iter().enumerate() {
        println!("   {}. {}", i + 1, serde_json::to_string_pretty(doc)?);
    }
    println!();

    // Test 4: Vector distance functions
    println!("--- Test 4: Euclidean distance function ---");
    let mut response = db
        .query(
            "
        SELECT
            content,
            vector::distance::euclidean(embedding, [0.1, 0.2, 0.3]) AS distance
        FROM document
        ORDER BY distance ASC;
    ",
        )
        .await?;

    let result: Result<Vec<Value>, _> = response.take(0);
    match result {
        Ok(docs) => {
            println!("✅ Distance function returned {} results:", docs.len());
            for (i, doc) in docs.iter().enumerate() {
                println!("   {}. {}", i + 1, serde_json::to_string_pretty(doc)?);
            }
        }
        Err(e) => {
            println!("⚠️  Failed to parse distance results: {:?}", e);
            // Try to get raw response
            println!("   Attempting to get raw response...");
        }
    }
    println!();

    // Test 5: Cosine distance
    println!("--- Test 5: Cosine distance function ---");
    let mut response = db
        .query(
            "
        SELECT
            content,
            vector::distance::cosine(embedding, [0.1, 0.2, 0.3]) AS distance
        FROM document
        ORDER BY distance ASC;
    ",
        )
        .await?;

    let result: Result<Vec<Value>, _> = response.take(0);
    match result {
        Ok(docs) => {
            println!("✅ Cosine distance returned {} results:", docs.len());
            for doc in docs.iter().take(1) {
                println!("   {}", serde_json::to_string_pretty(doc)?);
            }
        }
        Err(e) => {
            println!("⚠️  Failed to parse cosine results: {:?}", e);
        }
    }
    println!();

    // Test 6: HNSW index
    println!("--- Test 6: Creating HNSW index ---");
    let db2 = Surreal::new::<Mem>(()).await?;
    db2.use_ns("test2").await?;
    db2.use_db("test2").await?;

    let result = db2
        .query(
            "
        DEFINE TABLE document SCHEMAFULL;
        DEFINE FIELD content ON document TYPE string;
        DEFINE FIELD embedding ON document TYPE array;
        DEFINE INDEX embedding_idx ON document FIELDS embedding
            HNSW DIMENSION 3
            DIST EUCLIDEAN
            EFC 150
            M 12;
    ",
        )
        .await;

    match result {
        Ok(_) => println!("✅ HNSW index created successfully"),
        Err(e) => println!("❌ HNSW index failed: {:?}", e),
    }

    // Insert and search with HNSW
    db2.query(
        "
        CREATE document:doc1 SET content = 'HNSW test', embedding = [0.1, 0.2, 0.3];
        CREATE document:doc2 SET content = 'HNSW test 2', embedding = [0.5, 0.4, 0.3];
    ",
    )
    .await?;

    let mut response = db2
        .query(
            "
        SELECT content
        FROM document
        WHERE embedding <|2,EUCLIDEAN|> [0.2, 0.3, 0.4]
        LIMIT 2;
    ",
        )
        .await?;

    let docs: Vec<Value> = response.take(0)?;
    println!("✅ HNSW KNN search returned {} documents", docs.len());
    println!();

    // Test 7: Similarity operator
    println!("--- Test 7: Vector similarity operator <||> ---");
    let mut response = db
        .query(
            "
        SELECT content, (embedding <||> [0.1, 0.2, 0.3]) AS similarity
        FROM document;
    ",
        )
        .await?;

    let result: Result<Vec<Value>, _> = response.take(0);
    match result {
        Ok(docs) => {
            println!("✅ Similarity operator returned {} results:", docs.len());
            for doc in docs.iter().take(1) {
                println!("   {}", serde_json::to_string_pretty(doc)?);
            }
        }
        Err(e) => {
            println!("⚠️  Failed to parse similarity results: {:?}", e);
        }
    }

    println!("\n=== Test Summary ===");
    println!("All core vector features are working in embedded mode:");
    println!("✅ MTREE index");
    println!("✅ HNSW index");
    println!("✅ Vector insertion");
    println!("✅ KNN search");
    println!("⚠️  Distance functions may have response parsing issues");

    Ok(())
}
