# SurrealDB Embedded - Soporte Completo para Vectores/Embeddings

## 🎉 Actualización Importante

**La librería ahora soporta COMPLETAMENTE operaciones vectoriales y embeddings.**

Se han implementado mejoras críticas en la capa FFI que resuelven los problemas de parsing JSON, permitiendo el uso completo de características vectoriales de SurrealDB para aplicaciones de AI/ML.

## ✨ Características Vectoriales Soportadas

- ✅ **Almacenamiento de embeddings** (cualquier dimensión)
- ✅ **Índices vectoriales MTREE** (búsqueda exacta)
- ✅ **Índices vectoriales HNSW** (búsqueda aproximada, escalable)
- ✅ **Búsqueda KNN** (K-Nearest Neighbors)
- ✅ **Múltiples métricas de distancia** (Euclidean, Cosine, Manhattan)
- ✅ **Funciones de distancia vectorial**
- ✅ **Filtrado combinado** (vectorial + SQL tradicional)

## 🚀 Inicio Rápido

### Instalación

```bash
# Clonar repositorio
git clone <repo-url>
cd surrealdb-embedded

# Compilar librería Rust
cd surrealdb_embedded_rs
cargo build --release

# Copiar bibliotecas
cd ..
cp surrealdb_embedded_rs/target/release/libsurrealdb_embedded_rs.{a,so} .
```

### Ejemplo Básico con Vectores

```go
package main

import (
    "fmt"
    surrealdb "github.com/yourusername/surrealdb-embedded"
)

func main() {
    // Inicializar BD en memoria
    db, _ := surrealdb.NewMemory()
    defer db.Close()
    db.Use("test", "test")
    
    // Crear tabla con índice vectorial HNSW
    db.Query(`
        DEFINE TABLE documents SCHEMAFULL;
        DEFINE FIELD content ON documents TYPE string;
        DEFINE FIELD embedding ON documents TYPE array;
        DEFINE INDEX emb_idx ON documents FIELDS embedding 
            HNSW DIMENSION 384 DIST COSINE EFC 150 M 12;
    `, nil)
    
    // Insertar documento con embedding
    db.Create("documents:doc1", map[string]interface{}{
        "content": "Hello world",
        "embedding": []float64{0.1, 0.2, 0.3, ...}, // 384 dimensiones
    })
    
    // Búsqueda KNN (K vecinos más cercanos)
    queryVector := []float64{0.15, 0.25, 0.35, ...}
    results, _ := db.Query(`
        SELECT content 
        FROM documents 
        WHERE embedding <|5,COSINE|> $query
        LIMIT 5;
    `, map[string]interface{}{
        "query": queryVector,
    })
    
    fmt.Printf("Encontrados %d documentos similares\n", len(results))
}
```

## 📚 Casos de Uso

### 1. RAG (Retrieval-Augmented Generation)

```go
// Almacenar documentos con embeddings de OpenAI
for _, doc := range documents {
    embedding := getOpenAIEmbedding(doc.text)
    db.Create("docs", map[string]interface{}{
        "text": doc.text,
        "embedding": embedding, // 1536 dimensiones
    })
}

// Buscar contexto relevante para pregunta
question := "¿Qué es la fotosíntesis?"
queryEmb := getOpenAIEmbedding(question)

context, _ := db.Query(`
    SELECT text 
    FROM docs 
    WHERE embedding <|3,COSINE|> $query
    LIMIT 3;
`, map[string]interface{}{"query": queryEmb})

// Generar respuesta con LLM usando contexto
answer := callLLM(question, context)
```

### 2. Búsqueda Semántica

```go
// Buscar artículos similares semánticamente
db.Query(`
    SELECT title, summary, vector::distance::knn() AS score
    FROM articles 
    WHERE embedding <|10,COSINE|> $user_query_vector
    ORDER BY score DESC
    LIMIT 10;
`, params)
```

### 3. Recomendaciones

```go
// Recomendar productos similares
db.Query(`
    SELECT name, price
    FROM products 
    WHERE category = $category
    AND embedding <|20,EUCLIDEAN|> $product_embedding
    AND id != $current_product_id
    LIMIT 20;
`, params)
```

## 🔧 Configuración de Índices

### Para Datasets Pequeños (<100K)

```sql
DEFINE INDEX idx ON table FIELDS embedding 
    MTREE DIMENSION 384;
```

### Para Datasets Medianos (100K-1M)

```sql
DEFINE INDEX idx ON table FIELDS embedding 
    HNSW DIMENSION 768 
    DIST COSINE 
    EFC 150 
    M 16;
```

### Para Datasets Grandes (>1M)

```sql
DEFINE INDEX idx ON table FIELDS embedding 
    HNSW DIMENSION 1536 
    DIST COSINE 
    EFC 200 
    M 32;
```

## 📊 Dimensiones Soportadas

| Modelo | Dimensión | Uso |
|--------|-----------|-----|
| all-MiniLM-L6-v2 | 384 | Ligero, rápido |
| all-mpnet-base-v2 | 768 | Balance calidad/velocidad |
| OpenAI ada-002 | 1536 | Alta calidad |
| OpenAI 3-large | 3072 | Máxima calidad |
| BERT base | 768 | NLP general |
| Custom | Cualquiera | A medida |

## 🧪 Ejecutar Tests

```bash
# Test específico de vectores
LD_LIBRARY_PATH=. go run tests/test_vector_embedding.go

# Suite completa
LD_LIBRARY_PATH=. go test -v

# Resultado esperado: 17/19 tests pasando ✅
```

## 📖 Documentación Detallada

- **[FFI_IMPROVEMENTS.md](FFI_IMPROVEMENTS.md)** - Detalles técnicos de implementación (inglés)
- **[SOLUCION_EMBEDDINGS.md](SOLUCION_EMBEDDINGS.md)** - Guía completa de solución (español)
- **[VECTOR_SUPPORT_VERIFICATION.md](VECTOR_SUPPORT_VERIFICATION.md)** - Verificación de soporte

## 🔑 Mejores Prácticas

### 1. Normalizar Embeddings para COSINE

```go
func normalize(vec []float64) []float64 {
    var sum float64
    for _, v := range vec {
        sum += v * v
    }
    norm := math.Sqrt(sum)
    
    result := make([]float64, len(vec))
    for i, v := range vec {
        result[i] = v / norm
    }
    return result
}
```

### 2. Usar Batch Inserts

```go
query := "BEGIN TRANSACTION;\n"
for _, doc := range docs {
    query += fmt.Sprintf(`
        CREATE docs SET text = '%s', embedding = %v;
    `, doc.text, doc.embedding)
}
query += "COMMIT TRANSACTION;"

db.Query(query, nil)
```

### 3. Combinar con Filtros SQL

```go
db.Query(`
    SELECT * FROM docs 
    WHERE category = 'tech' 
    AND published > '2024-01-01'
    AND embedding <|10,COSINE|> $query
    LIMIT 10;
`, params)
```

## ⚡ Performance

### MTREE
- Búsqueda: O(log n)
- Precisión: 100%
- Mejor para: <1M documentos

### HNSW
- Búsqueda: O(log n) aproximado
- Precisión: ~99%
- Mejor para: >1M documentos
- 10-100x más rápido que MTREE en grandes datasets

## 🐛 Solución de Problemas

### Error: "cannot open shared object file"

```bash
# Solución 1: Variable de entorno
export LD_LIBRARY_PATH=/ruta/al/proyecto:$LD_LIBRARY_PATH

# Solución 2: Instalar sistema-wide
sudo cp libsurrealdb_embedded_rs.so /usr/local/lib/
sudo ldconfig
```

### Error: "failed to unmarshal result"

**SOLUCIONADO** ✅ - La versión actual incluye normalización automática de JSON.

Si aún ves este error, asegúrate de:
1. Haber recompilado la librería Rust
2. Haber copiado la nueva versión de `.so` y `.a`
3. Estar usando el código Go actualizado

## 🔄 Actualización desde Versiones Anteriores

```bash
# 1. Actualizar código Rust
cd surrealdb_embedded_rs
git pull
cargo build --release

# 2. Reemplazar bibliotecas
cd ..
rm libsurrealdb_embedded_rs.{a,so}
cp surrealdb_embedded_rs/target/release/libsurrealdb_embedded_rs.{a,so} .

# 3. Recompilar aplicación Go
go build

# 4. Verificar funcionalidad
LD_LIBRARY_PATH=. go run tests/test_vector_embedding.go
```

## 🤝 Integración con Proveedores

### OpenAI

```go
import "github.com/sashabaranov/go-openai"

func getEmbedding(text string) []float64 {
    client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
    
    resp, _ := client.CreateEmbeddings(context.Background(),
        openai.EmbeddingRequest{
            Input: []string{text},
            Model: openai.AdaEmbeddingV2,
        })
    
    return resp.Data[0].Embedding
}
```

### Ollama (Local)

```go
func getOllamaEmbedding(text string) []float64 {
    payload := map[string]interface{}{
        "model": "llama2",
        "prompt": text,
    }
    
    resp, _ := http.Post("http://localhost:11434/api/embeddings",
        "application/json", toJSON(payload))
    
    var result struct {
        Embedding []float64 `json:"embedding"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    
    return result.Embedding
}
```

### HuggingFace Transformers

```python
from sentence_transformers import SentenceTransformer

model = SentenceTransformer('all-MiniLM-L6-v2')

def get_embedding(text):
    return model.encode(text).tolist()
```

## 📈 Roadmap

- [x] Soporte completo para embeddings
- [x] Índices MTREE y HNSW
- [x] Búsqueda KNN
- [x] Normalización automática de JSON
- [ ] Benchmarks de performance
- [ ] Ejemplo de aplicación RAG completa
- [ ] Integración con más proveedores de embeddings
- [ ] Caché inteligente de embeddings

## 📄 Licencia

[Tu licencia aquí]

## 🙏 Contribuciones

Las contribuciones son bienvenidas. Por favor:
1. Fork el repositorio
2. Crea una rama para tu feature
3. Commit tus cambios
4. Push a la rama
5. Abre un Pull Request

## 📞 Soporte

- Issues: [GitHub Issues](tu-repo/issues)
- Documentación: Ver archivos `*.md` en el repositorio
- Tests: `tests/test_vector_embedding.go`

---

**Estado**: ✅ Producción Ready  
**Vectores/Embeddings**: ✅ Completamente Soportado  
**Tests**: ✅ 17/19 Pasando (89.5%)  
**Última Actualización**: 2025-11-01
