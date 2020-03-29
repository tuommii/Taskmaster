BIN_DIR = bin

all: run

build:
	go build -o $(BIN_DIR)/taskmaster cmd/taskmaster/main.go

run: build
	./bin/taskmaster

deploy:
	# scp -P 3001 run.sh yoda@46.101.105.101:/home/yoda/nolife/run.sh
	# scp -P 3001 users.json yoda@46.101.105.101:/home/yoda/nolife/users.json
	# scp -P 3001 bin/crawler yoda@46.101.105.101:/home/yoda/nolife/bin/crawler
	# scp -P 3001 bin/parser yoda@46.101.105.101:/home/yoda/nolife/bin/parser
	# scp -P 3001 html/template.html yoda@46.101.105.101:/home/yoda/nolife/html/template.html

test:
	go test ./...

# cross-compile:
# 	GOOS=linux GOARCH=amd64 go build
