import re

with open('lib.rs', 'r') as f:
    content = f.read()

# Add detailed logging before take
old = '''    let rt = get_runtime();
    let result = rt.block_on(async { db_instance.query(&query).await });

    match result {
        Ok(mut response) => {
            // Extract as single SurrealDB Value and convert to JSON
            let json_result = match response.take::<Vec<Value>>(0) {'''

new = '''    let rt = get_runtime();
    let result = rt.block_on(async { db_instance.query(&query).await });

    match result {
        Ok(mut response) => {
            eprintln!("DEBUG: Query executed successfully");
            eprintln!("DEBUG: Number of statements: {}", response.num_statements());
            
            // Extract as single SurrealDB Value and convert to JSON
            let json_result = match response.take::<Vec<Value>>(0) {'''

content = content.replace(old, new)

with open('lib.rs', 'w') as f:
    f.write(content)

print("Added logging!")
