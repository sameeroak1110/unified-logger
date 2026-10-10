# Server Makefile

APP_NAME = unified-logger-server
PROTO_DIR = proto
PROTO_SRC = $(PROTO_DIR)/message.proto
GO_FILES = main.go

.PHONY: all build run proto clean help

## all: Default target - generates proto and builds the server binary
all: proto build

## proto: Generates Go gRPC code from the .proto file
proto:
	@echo "Generating gRPC proto code..."
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative $(PROTO_SRC)

## build: Compiles the server binary
build: proto
	@echo "Building server binary..."
	go build -o bin/$(APP_NAME) $(GO_FILES)

## run: Runs the server directly using 'go run'
run: proto
	@echo "Starting gRPC server..."
	go run $(GO_FILES)

## clean: Removes built binaries and output test files
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f output.txt

## help: Shows available targets
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':'
