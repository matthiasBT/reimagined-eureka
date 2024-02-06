.PHONY: all windows macos linux clean

BINARY_NAME=client
BUILD_DIR=./bin

all: windows macos linux

windows:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/client/main.go

macos:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64 ./cmd/client/main.go
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64 ./cmd/client/main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/client/main.go

clean:
	rm -rf $(BUILD_DIR)/*
