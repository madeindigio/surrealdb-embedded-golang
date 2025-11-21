import re

with open('lib.rs', 'r') as f:
    content = f.read()

# Add surrealdb::Value import
content = content.replace(
    'use surrealdb::sql::Value;',
    'use surrealdb::sql::Value as SqlValue;\nuse surrealdb::Value;'
)

# Fix first surreal_query function (around line 141)
old_pattern1 = r'''let json_result = match response\.take::<Vec<Value>>\(0\) \{
                Ok\(values\) => \{
                    match serde_json::to_string\(&values\) \{
                        Ok\(json\) => normalize_surrealdb_json\(&json\),
                        Err\(_\) => "\[\]"\.to_string\(\),
                    \}
                \}
                Err\(_\) => "\[\]"\.to_string\(\)
            \};'''

new_pattern1 = '''let json_result = match response.take::<Value>(0) {
                Ok(value) => {
                    match serde_json::to_string(&value) {
                        Ok(json) => normalize_surrealdb_json(&json),
                        Err(e) => {
                            eprintln!("Serialization error: {:?}", e);
                            "[]".to_string()
                        }
                    }
                }
                Err(e) => {
                    eprintln!("Take error: {:?}", e);
                    "[]".to_string()
                }
            };'''

content = re.sub(old_pattern1, new_pattern1, content, flags=re.DOTALL)

with open('lib.rs', 'w') as f:
    f.write(content)

print("Fixed!")
