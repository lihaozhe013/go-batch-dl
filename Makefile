.PHONY: build run lint test clean

BIN_DIR=bin
APP_NAME=gobatchdl


default: test

test: clean
	uv run tester/tester.py

build:
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)/

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BIN_DIR)
	rm -rf test_downloads
	rm -rf test_server_root
