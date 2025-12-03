# SurrealDB Embedded - Análisis del Problema y Solución

## Fecha de Investigación
2025-11-21

## Problema Identificado

El wrapper de Rust para SurrealDB embedded tiene un problema fundamental con la extracción de resultados de queries. 

### Síntoma
Todas las queries devuelven arrays vacíos `[]`, incluso cuando los datos existen en la base de datos.

### Causa Raíz
**Bug en SurrealDB 2.x**: El método `response.take()` no puede deserializar correctamente el tipo `surrealdb::sql::Value` (que es un enum) a tipos JSON estándar como `serde_json::Value` o `HashMap<String, serde_json::Value>`.

Error específico:
```
Db(Serialization("failed to deserialize; expected an enum variant of $surrealdb::private::sql::Value, found { ... }"))
```

Este es un **bug conocido** reportado en:
- Issue #4921: https://github.com/surrealdb/surrealdb/issues/4921
- Issue #4316: https://github.com/surrealdb/surrealdb/issues/4316

Aunque los issues están cerrados como "completados", el problema persiste en SurrealDB 2.3.

## Approaches Intentados (Todos Fallidos)

1. ✗ `response.take::<Vec<Value>>(0)` - Falla con error de deserialización
2. ✗ `response.take::<Option<Value>>(0)` - Falla con error de deserialización  
3. ✗ `response.take::<Vec<serde_json::Value>>(0)` - Falla: "invalid type: enum, expected any valid JSON value"
4. ✗ `response.take::<Vec<HashMap<String, serde_json::Value>>>(0)` - Mismo error
5. ✗ Conversión manual via `serde_json::to_string()` - No funciona porque `take()` ya falla antes
6. ✗ `#[serde(untagged)] enum FlexibleValue` - Falla: "untagged enums do not support enum input"
7. ✗ `struct DynamicRecord with #[serde(flatten)]` - En progreso pero probablemente fallará por la misma razón

## Por Qué Estos Approaches Fallan

El problema fundamental es que `surrealdb::sql::Value` es un enum complejo:

```rust
pub enum Value {
    None,
    Null,
    Bool(bool),
    Number(Number),
    Strand(Strand),
    Duration(Duration),
    Datetime(Datetime),
    Uuid(Uuid),
    Array(Array),
    Object(Object),  // <- Las queries devuelven esto
    // ... más variantes
}
```

Cuando SurrealDB serializa un registro (Object), serde espera deserializar la **variante del enum** (con su tag), pero SurrealDB está devolviendo directamente el **contenido del Object** sin el wrapper del enum.

## Soluciones Posibles

### Solución 1: Usar SurrealDB Remoto (RECOMENDADA) ⭐

La forma más confiable y probada de usar SurrealDB:

```bash
# Iniciar servidor SurrealDB
surreal start --user root --pass root file://./remembrances.db
```

```yaml
# config.yaml
surrealdb-url: "ws://localhost:8000"
surrealdb-namespace: "remembrances"
surrealdb-database: "main"
```

**Ventajas:**
- ✅ Funciona de forma estable y probada
- ✅ Mejor rendimiento para operaciones concurrentes
- ✅ Más fácil de depurar
- ✅ Permite usar herramientas de administración
- ✅ Separación clara entre almacenamiento y aplicación

**Desventajas:**
- ❌ Requiere proceso separado
- ❌ Mayor complejidad de deployment

### Solución 2: Usar Cliente HTTP de SurrealDB

En lugar de embedded, usar el cliente HTTP/WebSocket nativo de Go:

```go
import "github.com/surrealdb/surrealdb.go"

db, err := surrealdb.New("ws://localhost:8000/rpc")
```

### Solución 3: Fix Profundo del Wrapper (Complejo)

Requiere crear tipos custom que coincidan exactamente con la estructura interna de SurrealDB:

```rust
// Definir struct para cada tipo de query
#[derive(Deserialize)]
struct PersonRecord {
    id: Thing,
    name: String,
    age: i64,
}

// Usar en wrapper
let results: Vec<PersonRecord> = response.take(0)?;
```

**Problema**: Imposible para un wrapper genérico que no conoce los tipos de antemano.

### Solución 4: Fork y Patch de SurrealDB

Modificar el código fuente de SurrealDB para arreglar el problema de serialización.

**Complejidad**: Alta, requiere mantener fork actualizado.

## Recomendación Final

**Usar SurrealDB como servidor externo (Solución 1)** es la mejor opción por:

1. Es la forma oficial y más soportada de usar SurrealDB
2. El modo embedded está claramente en desarrollo y tiene bugs
3. Mejor rendimiento y características
4. Más fácil de mantener a largo plazo

El overhead de tener un proceso separado es mínimo comparado con los beneficios de estabilidad y features.

## Archivos Modificados Durante la Investigación

### Rust (surrealdb-embedded)
- `surrealdb_embedded_rs/src/lib.rs` - Múltiples intentos de fix
- `surrealdb_embedded_rs/Cargo.toml` - Upgrade a SurrealDB 2.3

### Go (remembrances-mcp)  
- `internal/kb/kb.go` - Implementación de verificación de timestamps ✅
- `internal/storage/surrealdb_documents.go` - Modificación de GetDocument() ✅

Los cambios en Go son correctos y funcionarán una vez que el problema del wrapper se resuelva.

## Próximos Pasos

1. **Inmediato**: Cambiar a SurrealDB remoto para desbloquear el desarrollo
2. **Corto plazo**: Monitorear issues de SurrealDB para fix oficial
3. **Largo plazo**: Considerar migrar a cliente Go nativo de SurrealDB

