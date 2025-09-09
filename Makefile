PROJECT          ?= go-grpc
PROTO_DIR        ?= protos
PROTO_OUT_DIR    ?= $(PROTO_DIR)/gen
SERVER_DIR       ?= server
CLIENT_DIR       ?= client
BIN_DIR          ?= bin

# tools
PROTOC           ?= protoc
PROTOC_GEN_GO    ?= $(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC ?= $(shell which protoc-gen-go-grpc)
GO               ?= go

# generated files
PROTO_FILES      := $(wildcard $(PROTO_DIR)/*.proto)

.PHONY: all proto build build-server build-client run-server run-client test clean lint

all: proto build test

## Generate Go code from .proto files
proto: $(PROTO_FILES)
	@echo "Generating gRPC code..."
	@mkdir -p $(PROTO_OUT_DIR)
	$(PROTOC) \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$^

## Build both server and client binaries
build: build-server build-client

## Build the gRPC server binary
build-server: proto
	@echo "Building server binary..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(PROJECT)-server ./$(SERVER_DIR)

## Build the gRPC client binary
build-client: proto
	@echo "Building client binary..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(PROJECT)-client ./$(CLIENT_DIR)

## Run the gRPC server
run-server: proto
	@echo "Starting gRPC server..."
	$(GO) run ./$(SERVER_DIR)

## Run the gRPC client
run-client: proto
	@echo "Starting gRPC client..."
	$(GO) run ./$(CLIENT_DIR)

## Run all tests (unit + integration)
test:
	@echo "Running tests..."
	$(GO) test ./... -coverprofile=coverage.out

## Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -f $(PROTO_DIR)/gen/*.pb.go
	rm -rf $(BIN_DIR) coverage.out

## Run go vet and staticcheck if installed
lint:
	@echo "Linting..."
	$(GO) vet ./...
	@if command -v staticcheck > /dev/null; then staticcheck ./...; fi
