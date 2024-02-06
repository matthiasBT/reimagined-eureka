.PHONY: all windows macos linux server-macos clean

BINARY_NAME_CLIENT=client
BINARY_NAME_SERVER=server
BUILD_DIR=./bin

all: windows macos linux

windows:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME_CLIENT)-windows-amd64.exe ./cmd/client/main.go
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME_SERVER)-windows-amd64.exe ./cmd/server/main.go

macos:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME_CLIENT)-macos-amd64 ./cmd/client/main.go
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME_CLIENT)-macos-arm64 ./cmd/client/main.go
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME_SERVER)-macos-amd64 ./cmd/server/main.go
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME_SERVER)-macos-arm64 ./cmd/server/main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME_CLIENT)-linux-amd64 ./cmd/client/main.go
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME_SERVER)-linux-amd64 ./cmd/server/main.go

clean:
	rm -rf $(BUILD_DIR)/*
