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


test:
	env GO_ENV=test TEST_PSQL_DBNAME=sample TEST_PSQL_HOST=localhost TEST_PSQL_PORT=5411 TEST_PSQL_USER=sample_user TEST_PSQL_PASS=sample2020 go test ./src/...

