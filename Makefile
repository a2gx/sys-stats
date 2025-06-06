BIN := ./bin/daemon
APP := ./cmd/daemon

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" $(APP)

run: build
	$(BIN) run -d -n 3 -m 7

stop: build
	$(BIN) stop

version: build
	$(BIN) -v

help: build
	$(BIN) -h

test:
	go test -race ./...

lint-install-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1.6

lint: lint-install-deps
	golangci-lint run ./...

up:
	docker-compose up --build
down:
	docker-compose down
logs:
	docker-compose logs -f

.PHONY: build run logs stop version help test lint up logs down