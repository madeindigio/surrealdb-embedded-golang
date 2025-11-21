import re

with open('lib.rs', 'r') as f:
    content = f.read()

# Fix to use Vec<Value> instead of Value
content = content.replace(
    'let json_result = match response.take::<Value>(0) {',
    'let json_result = match response.take::<Vec<Value>>(0) {'
)

content = content.replace(
    '''let json_result = match response.take::<Vec<Value>>(0) {
                Ok(value) => {
                    match serde_json::to_string(&value) {''',
    '''let json_result = match response.take::<Vec<Value>>(0) {
                Ok(values) => {
                    match serde_json::to_string(&values) {'''
)

content = content.replace(
    '                Ok(value) => {',
    '                Ok(values) => {'
)

with open('lib.rs', 'w') as f:
    f.write(content)

print("Fixed to Vec<Value>!")
