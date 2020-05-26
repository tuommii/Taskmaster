BIN_DIR = bin

all: run

build-client:
	go build -o $(BIN_DIR)/client cmd/client/*.go

build-server:
	go build -o $(BIN_DIR)/server cmd/server/*.go

build: build-server build-client

run: build-client
	./bin/client

client: build-client
	./bin/client

server: build-server
	./bin/server

test:
	go test ./...

# cross-compile:
# 	GOOS=linux GOARCH=amd64 go build
