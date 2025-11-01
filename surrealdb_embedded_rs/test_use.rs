use surrealdb::Surreal;
use surrealdb::engine::local::Mem;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let db = Surreal::new::<Mem>(()).await?;
    
    // Prueba 1: encadenado
    println!("Test 1: chained");
    match db.use_ns("test").use_db("test").await {
        Ok(_) => println!("✓ Chained works!"),
        Err(e) => println!("✗ Chained failed: {}", e),
    }
    
    // Prueba 2: separado
    println!("\nTest 2: separate");
    match db.use_ns("test2").await {
        Ok(_) => println!("✓ use_ns works!"),
        Err(e) => println!("✗ use_ns failed: {}", e),
    }
    
    match db.use_db("test2").await {
        Ok(_) => println!("✓ use_db works!"),
        Err(e) => println!("✗ use_db failed: {}", e),
    }
    
    // Crear un registro para verificar
    #[derive(serde::Serialize)]
    struct Person {
        name: String,
    }
    
    let created: Option<Person> = db.create("person").content(Person {
        name: "Test".to_string(),
    }).await?;
    
    println!("\n✓ Created record: {:?}", created.is_some());
    
    Ok(())
}
