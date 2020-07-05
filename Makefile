.PHONY: setup install build analyze

all: install setup

install:
	GO111MODULE=off go get -u github.com/rubenv/sql-migrate/...
	git submodule update -i
setup:
	go mod download
build:
	go build -ldflags '-s -w' -o ./src/bin/app ./src

generate/go:
	cd src && rm -rf mock && go generate

generate/proto:
	./gen.sh

run:
	cd src && env GO_ENV=local LOCAL_DB_NAME=sample LOCAL_DB_HOST=localhost LOCAL_DB_PORT=3346 LOCAL_DB_USER=root LOCAL_DB_PASS=root go run .

test:
	env GO_ENV=local LOCAL_DB_NAME=sample LOCAL_DB_HOST=localhost LOCAL_DB_PORT=3346 LOCAL_DB_USER=root LOCAL_DB_PASS=root go test -v ./src/...

