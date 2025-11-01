use serde::{Deserialize, Serialize};
use serde_json::Value;
use surrealdb::engine::local::{Db, Mem};
use surrealdb::Surreal;

#[derive(Debug, Serialize, Deserialize)]
struct Document {
    id: Option<String>,
    content: String,
    embedding: Vec<f64>,
}

#[tokio::test]
async fn test_vector_schema_creation() {
    println!("\n=== Test 1: Vector Schema Creation ===");

    let db = Surreal::new::<Mem>(()).await.unwrap();
    db.use_ns("test").await.unwrap();
    db.use_db("test").await.unwrap();

    // Create table with vector field and MTREE index
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
        Ok(mut response) => {
            println!("✅ Schema created successfully");
            let _: Vec<Value> = response.take(0).unwrap_or_default();
        }
        Err(e) => {
            println!("❌ Schema creation failed: {:?}", e);
            panic!("Schema creation failed");
        }
    }
}

#[tokio::test]
async fn test_vector_insertion() {
    println!("\n=== Test 2: Vector Insertion ===");

    let db = Surreal::new::<Mem>(()).await.unwrap();
    db.use_ns("test").await.unwrap();
    db.use_db("test").await.unwrap();

    // Create schema
    db.query(
        "
        DEFINE TABLE document SCHEMAFULL;
        DEFINE FIELD content ON document TYPE string;
        DEFINE FIELD embedding ON document TYPE array;
        DEFINE INDEX embedding_idx ON document FIELDS embedding MTREE DIMENSION 3;
    ",
    )
    .await
    .unwrap();

    // Insert document with embedding
    let result = db
        .query(
            "
        CREATE document:doc1 SET content = 'Hello world', embedding = [0.1, 0.2, 0.3];
    ",
        )
        .await;

    match result {
        Ok(mut response) => {
            let docs: Vec<Document> = response.take(0).unwrap_or_default();
            println!("✅ Inserted {} document(s)", docs.len());
            if let Some(doc) = docs.first() {
                println!("   Content: {}", doc.content);
                println!("   Embedding: {:?}", doc.embedding);
            }
        }
        Err(e) => {
            println!("❌ Insertion failed: {:?}", e);
            panic!("Insertion failed");
        }
    }
}

#[tokio::test]
async fn test_knn_search() {
    println!("\n=== Test 3: KNN Search ===");

    let db = Surreal::new::<Mem>(()).await.unwrap();
    db.use_ns("test").await.unwrap();
    db.use_db("test").await.unwrap();

    // Setup schema and data
    db.query(
        "
        DEFINE TABLE document SCHEMAFULL;
        DEFINE FIELD content ON document TYPE string;
        DEFINE FIELD embedding ON document TYPE array;
        DEFINE INDEX embedding_idx ON document FIELDS embedding MTREE DIMENSION 3;

        CREATE document:doc1 SET content = 'Hello world', embedding = [0.1, 0.2, 0.3];
        CREATE document:doc2 SET content = 'Goodbye world', embedding = [0.5, 0.4, 0.3];
        CREATE document:doc3 SET content = 'Another document', embedding = [0.9, 0.8, 0.7];
    ",
    )
    .await
    .unwrap();

    // Perform KNN search
    let result = db
        .query(
            "
        SELECT content, embedding
        FROM document
        WHERE embedding <|3,EUCLIDEAN|> [0.6, 0.3, 0.2]
        LIMIT 3;
    ",
        )
        .await;

    match result {
        Ok(mut response) => {
            let docs: Vec<Document> = response.take(0).unwrap_or_default();
            println!("✅ KNN search returned {} document(s)", docs.len());
            for (i, doc) in docs.iter().enumerate() {
                println!(
                    "   {}. Content: {}, Embedding: {:?}",
                    i + 1,
                    doc.content,
                    doc.embedding
                );
            }
        }
        Err(e) => {
            println!("❌ KNN search failed: {:?}", e);
            panic!("KNN search failed");
        }
    }
}

#[tokio::test]
async fn test_vector_distance_functions() {
    println!("\n=== Test 4: Vector Distance Functions ===");

    let db = Surreal::new::<Mem>(()).await.unwrap();
    db.use_ns("test").await.unwrap();
    db.use_db("test").await.unwrap();

    // Setup schema and data
    db.query(
        "
        DEFINE TABLE document SCHEMAFULL;
        DEFINE FIELD content ON document TYPE string;
        DEFINE FIELD embedding ON document TYPE array;

        CREATE document:doc1 SET content = 'Hello world', embedding = [0.1, 0.2, 0.3];
        CREATE document:doc2 SET content = 'Goodbye world', embedding = [0.5, 0.4, 0.3];
    ",
    )
    .await
    .unwrap();

    // Test Euclidean distance
    println!("\n--- Testing Euclidean Distance ---");
    let result = db
        .query(
            "
        SELECT
            content,
            vector::distance::euclidean(embedding, [0.1, 0.2, 0.3]) AS distance
        FROM document
        ORDER BY distance ASC;
    ",
        )
        .await;

    match result {
        Ok(mut response) => {
            let result_value: Result<Vec<Value>, _> = response.take(0);
            match result_value {
                Ok(docs) => {
                    println!(
                        "✅ Euclidean distance query returned {} result(s)",
                        docs.len()
                    );
                    for doc in docs {
                        println!(
                            "   Result: {}",
                            serde_json::to_string_pretty(&doc).unwrap_or_default()
                        );
                    }
                }
                Err(e) => {
                    println!("⚠️  Failed to parse results: {:?}", e);
                }
            }
        }
        Err(e) => {
            println!("❌ Euclidean distance query failed: {:?}", e);
        }
    }

    // Test Cosine distance
    println!("\n--- Testing Cosine Distance ---");
    let result = db
        .query(
            "
        SELECT
            content,
            vector::distance::cosine(embedding, [0.1, 0.2, 0.3]) AS distance
        FROM document
        ORDER BY distance ASC;
    ",
        )
        .await;

    match result {
        Ok(mut response) => {
            let result_value: Result<Vec<Value>, _> = response.take(0);
            match result_value {
                Ok(docs) => {
                    println!("✅ Cosine distance query returned {} result(s)", docs.len());
                    for doc in docs {
                        println!(
                            "   Result: {}",
                            serde_json::to_string_pretty(&doc).unwrap_or_default()
                        );
                    }
                }
                Err(e) => {
                    println!("⚠️  Failed to parse results: {:?}", e);
                }
            }
        }
        Err(e) => {
            println!("❌ Cosine distance query failed: {:?}", e);
        }
    }

    // Test Manhattan distance
    println!("\n--- Testing Manhattan Distance ---");
    let result = db
        .query(
            "
        SELECT
            content,
            vector::distance::manhattan(embedding, [0.1, 0.2, 0.3]) AS distance
        FROM document
        ORDER BY distance ASC;
    ",
        )
        .await;

    match result {
        Ok(mut response) => {
            let result_value: Result<Vec<Value>, _> = response.take(0);
            match result_value {
                Ok(docs) => {
                    println!(
                        "✅ Manhattan distance query returned {} result(s)",
                        docs.len()
                    );
                    for doc in docs {
                        println!(
                            "   Result: {}",
                            serde_json::to_string_pretty(&doc).unwrap_or_default()
                        );
                    }
                }
                Err(e) => {
                    println!("⚠️  Failed to parse results: {:?}", e);
                }
            }
        }
        Err(e) => {
            println!("❌ Manhattan distance query failed: {:?}", e);
        }
    }
}

#[tokio::test]
async fn test_hnsw_index() {
    println!("\n=== Test 5: HNSW Index ===");

    let db = Surreal::new::<Mem>(()).await.unwrap();
    db.use_ns("test").await.unwrap();
    db.use_db("test").await.unwrap();

    // Create table with HNSW index
    let result = db
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
        Ok(_) => {
            println!("✅ HNSW index created successfully");
        }
        Err(e) => {
            println!("❌ HNSW index creation failed: {:?}", e);
            panic!("HNSW index creation failed");
        }
    }

    // Insert and search
    db.query(
        "
        CREATE document:doc1 SET content = 'Hello world', embedding = [0.1, 0.2, 0.3];
        CREATE document:doc2 SET content = 'Goodbye world', embedding = [0.5, 0.4, 0.3];
        CREATE document:doc3 SET content = 'Another document', embedding = [0.9, 0.8, 0.7];
    ",
    )
    .await
    .unwrap();

    let result = db
        .query(
            "
        SELECT content, embedding
        FROM document
        WHERE embedding <|2,EUCLIDEAN|> [0.6, 0.3, 0.2]
        LIMIT 2;
    ",
        )
        .await;

    match result {
        Ok(mut response) => {
            let docs: Vec<Document> = response.take(0).unwrap_or_default();
            println!("✅ HNSW KNN search returned {} document(s)", docs.len());
            for (i, doc) in docs.iter().enumerate() {
                println!("   {}. Content: {}", i + 1, doc.content);
            }
        }
        Err(e) => {
            println!("❌ HNSW KNN search failed: {:?}", e);
        }
    }
}

#[tokio::test]
async fn test_vector_similarity_operator() {
    println!("\n=== Test 6: Vector Similarity Operator ===");

    let db = Surreal::new::<Mem>(()).await.unwrap();
    db.use_ns("test").await.unwrap();
    db.use_db("test").await.unwrap();

    // Setup schema and data
    db.query(
        "
        DEFINE TABLE document SCHEMAFULL;
        DEFINE FIELD content ON document TYPE string;
        DEFINE FIELD embedding ON document TYPE array;

        CREATE document:doc1 SET content = 'Hello world', embedding = [0.1, 0.2, 0.3];
        CREATE document:doc2 SET content = 'Goodbye world', embedding = [0.5, 0.4, 0.3];
    ",
    )
    .await
    .unwrap();

    // Try vector similarity with <||> operator
    println!("\n--- Testing <||> Operator ---");
    let result = db
        .query(
            "
        SELECT content, (embedding <||> [0.1, 0.2, 0.3]) AS similarity
        FROM document;
    ",
        )
        .await;

    match result {
        Ok(mut response) => {
            let result_value: Result<Vec<Value>, _> = response.take(0);
            match result_value {
                Ok(docs) => {
                    println!("✅ Similarity operator returned {} result(s)", docs.len());
                    for doc in docs {
                        println!(
                            "   Result: {}",
                            serde_json::to_string_pretty(&doc).unwrap_or_default()
                        );
                    }
                }
                Err(e) => {
                    println!("⚠️  Failed to parse results: {:?}", e);
                }
            }
        }
        Err(e) => {
            println!("❌ Similarity operator query failed: {:?}", e);
        }
    }
}

#[tokio::test]
async fn test_complex_vector_query() {
    println!("\n=== Test 7: Complex Vector Query with Filtering ===");

    let db = Surreal::new::<Mem>(()).await.unwrap();
    db.use_ns("test").await.unwrap();
    db.use_db("test").await.unwrap();

    // Setup schema with additional fields
    db.query("
        DEFINE TABLE document SCHEMAFULL;
        DEFINE FIELD content ON document TYPE string;
        DEFINE FIELD category ON document TYPE string;
        DEFINE FIELD embedding ON document TYPE array;
        DEFINE INDEX embedding_idx ON document FIELDS embedding MTREE DIMENSION 3;

        CREATE document:doc1 SET content = 'Hello world', category = 'greeting', embedding = [0.1, 0.2, 0.3];
        CREATE document:doc2 SET content = 'Goodbye world', category = 'farewell', embedding = [0.5, 0.4, 0.3];
        CREATE document:doc3 SET content = 'Good morning', category = 'greeting', embedding = [0.2, 0.3, 0.4];
    ").await.unwrap();

    // Complex query: KNN search with category filter
    let result = db
        .query(
            "
        SELECT content, category
        FROM document
        WHERE category = 'greeting'
        AND embedding <|2,EUCLIDEAN|> [0.15, 0.25, 0.35]
        LIMIT 2;
    ",
        )
        .await;

    match result {
        Ok(mut response) => {
            let result_value: Result<Vec<Value>, _> = response.take(0);
            match result_value {
                Ok(docs) => {
                    println!("✅ Complex query returned {} result(s)", docs.len());
                    for doc in docs {
                        println!(
                            "   Result: {}",
                            serde_json::to_string_pretty(&doc).unwrap_or_default()
                        );
                    }
                }
                Err(e) => {
                    println!("⚠️  Failed to parse results: {:?}", e);
                }
            }
        }
        Err(e) => {
            println!("❌ Complex query failed: {:?}", e);
        }
    }
}

#[tokio::test]
async fn test_raw_response_format() {
    println!("\n=== Test 8: Raw Response Format Analysis ===");

    let db = Surreal::new::<Mem>(()).await.unwrap();
    db.use_ns("test").await.unwrap();
    db.use_db("test").await.unwrap();

    // Setup
    db.query(
        "
        DEFINE TABLE document SCHEMAFULL;
        DEFINE FIELD content ON document TYPE string;
        DEFINE FIELD embedding ON document TYPE array;

        CREATE document:doc1 SET content = 'Test', embedding = [0.1, 0.2, 0.3];
    ",
    )
    .await
    .unwrap();

    // Get raw response from distance function
    let mut response = db
        .query(
            "
        SELECT
            content,
            vector::distance::euclidean(embedding, [0.1, 0.2, 0.3]) AS distance
        FROM document;
    ",
        )
        .await
        .unwrap();

    println!("\n--- Analyzing Raw Response Structure ---");

    // Try to get as raw Value first
    let raw_value: Result<Vec<Value>, _> = response.take(0);
    match raw_value {
        Ok(values) => {
            println!("✅ Got raw values: {} item(s)", values.len());
            for (i, value) in values.iter().enumerate() {
                println!("\n   Item {}:", i + 1);
                println!(
                    "   JSON: {}",
                    serde_json::to_string_pretty(value).unwrap_or_default()
                );
                println!("   Debug: {:?}", value);
            }
        }
        Err(e) => {
            println!("❌ Failed to get raw values: {:?}", e);
        }
    }
}
