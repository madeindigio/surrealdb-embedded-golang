In the folder `surrealdb_embedded_rs` must to implement a C static library with Rust and surrealdb sdk library. Because Golang Surreal SDK library does not implement the feature of embed surrealdb features into the library embedded, using rocksdb or memory, this folder will implement the alternative as static library for this implementation using Rust and surreadldb that it has a working implementation, and later build with this a Golang library for using this static library.

Search in internet what are the features of golang surrealdb sdk methods and implement into this library, use as reference the existing code in the folder. I want to implement the option for memory database and rocksdb database for the embedded surrealdb library.

Later implement the golang library that use this static library with cgo and the tests for this golang library.

[]: #
[]: # ## 📚 References
[]: #
[]: # - [SurrealDB Embedded Documentation](https://surrealdb.com/docs/embedded)
[]: # - [Rust FFI Guide](https://doc.rust-lang.org/nomicon/ffi.html)
[]: # - [Golang cgo Documentation](https://golang.org/cmd/cgo/)
[]: #
