BIN_DIR = bin

all: run

client:
	go build -o $(BIN_DIR)/client cmd/client/client.go

server:
	go build -o $(BIN_DIR)/server cmd/server/server.go

build: server client

run: client
	./bin/client

test:
	go test ./...

# cross-compile:
# 	GOOS=linux GOARCH=amd64 go build
