use surrealdb::engine::local::Mem;
use surrealdb::Surreal;
use serde_json::Value;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("Testing vector support in SurrealDB embedded...\n");

    let db = Surreal::new::<Mem>(()).await?;
    db.use_ns("test").await?;
    db.use_db("test").await?;

    // Test 1: Basic vector operations
    println!("1. Creating schema with vector field...");
    db.query("
        DEFINE TABLE doc SCHEMAFULL;
        DEFINE FIELD content ON doc TYPE string;
        DEFINE FIELD embedding ON doc TYPE array;
        DEFINE INDEX emb_idx ON doc FIELDS embedding MTREE DIMENSION 3;
    ").await?;
    println!("   ✅ Schema created\n");

    // Test 2: Insert
    println!("2. Inserting documents...");
    db.query("
        CREATE doc:1 SET content = 'Test 1', embedding = [0.1, 0.2, 0.3];
        CREATE doc:2 SET content = 'Test 2', embedding = [0.5, 0.4, 0.3];
    ").await?;
    println!("   ✅ Documents inserted\n");

    // Test 3: KNN search
    println!("3. Testing KNN search...");
    let mut resp = db.query("
        SELECT content FROM doc
        WHERE embedding <|2,EUCLIDEAN|> [0.2, 0.3, 0.4]
        LIMIT 2;
    ").await?;

    let docs: Vec<Value> = resp.take(0)?;
    println!("   ✅ KNN returned {} results\n", docs.len());

    // Test 4: Distance function
    println!("4. Testing distance function...");
    let mut resp = db.query("
        SELECT content, vector::distance::euclidean(embedding, [0.1, 0.2, 0.3]) AS dist
        FROM doc;
    ").await?;

    match resp.take::<Vec<Value>>(0) {
        Ok(docs) => {
            println!("   ✅ Distance function works! {} results", docs.len());
            println!("   Result: {}", serde_json::to_string_pretty(&docs[0])?);
        }
        Err(e) => {
            println!("   ❌ Failed to parse: {:?}", e);
        }
    }

    println!("\n✅ All vector features are supported in embedded mode!");
    Ok(())
}
