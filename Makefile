# Root Makefile

SERVER_DIR = internal/grpc-server
PROTO_DIR  = proto
PROTO_SRC  = $(PROTO_DIR)/message.proto
BIN_DIR    = bin

SERVER_BIN_DIR = $(BIN_DIR)/server
SERVER_BIN = $(SERVER_BIN_DIR)/server

.PHONY: all proto build build-server run-server clean help

## all: Default target - generates proto and builds server.
all: proto build-server

## proto: Generates Go gRPC code from the .proto definition
proto:
	@echo "==> Generating gRPC proto code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)

## build-server: Compiles only the server binary
build-server: proto
	@echo "==> Building server binary..."
	@mkdir -p $(SERVER_BIN_DIR)
	go build -race -o $(SERVER_BIN) ./$(SERVER_DIR)

## run-server: Runs the server application using 'go run'
run-server: proto
	@echo "==> Starting gRPC Server..."
	go run ./$(SERVER_DIR)

## clean: Deletes generated binaries, test outputs, and proto generated files
clean:
	@echo "==> Cleaning up project artifacts..."
	rm -rf $(BIN_DIR)
	rm -f output.txt
	rm -f $(PROTO_DIR)/*.pb.go

## help: Displays available targets
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':'
