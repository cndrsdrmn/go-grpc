# Go gRPC CRUD Example

A simple **Go CRUD project** using:

- [gRPC](https://grpc.io/docs/languages/go/) for RPC APIs
- [GORM](https://gorm.io/) ORM with SQLite
- Repository + Service architecture
- Unit tests with [testify](https://github.com/stretchr/testify)
- Makefile for proto generation, builds, and tests

---

## Requirements

- [Go 1.24+](https://go.dev/doc/install)
- [Protobuf Compailer](https://protobuf.dev/installation/)
- [Go Protobuf Plugins](https://protobuf.dev/getting-started/gotutorial/#compiling-protocol-buffers)

---

## Usage

1. Generate gRPC code

   ```shell
   make proto
   ```

2. Build server and client

   ```shell
   make build
   ```

3. Run server

   ```shell
   make run-server
   ```

4. Run client

   ```shell
   make run-client
   ```

   Example output:

   ```text
   Created: id:1 name:"Alice" email:"alice@example.com"
   Fetched: id:1 name:"Alice" email:"alice@example.com"
   Updated: id:1 name:"Alice Updated" email:"alice.new@example.com"
   List Users: [id:1 name:"Alice Updated" email:"alice.new@example.com"]
   Deleted: true
   ```

5. Run tests

   ```shell
   make test
   ```

6. Clean build artifacts

   ```shell
   make clean
   ```

## License

This project is open-sourced software licensed under The MIT License (MIT). See the [LICENSE](./LICENSE) file for more details.
