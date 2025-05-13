BIN := ./bin/daemon
APP := ./cmd/daemon

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" $(APP)

start:
	go run $(APP) run

run: build
	$(BIN) run

run-detect: build
	#$(BIN) run -n 5 -m 15
	#$(BIN) run --detach
	$(BIN) run -d

logs: build
	$(BIN) logs -n 5 -m 15

stop: build
	$(BIN) stop

version: build
	#$(BIN) --version
	$(BIN) -v

help: build
	#$(BIN) --help
	$(BIN) -h

test:
	go test -race ./internal/...
	go test -race ./pkg/...

lint-install-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1

lint: lint-install-deps
	golangci-lint run ./cmd/...
	golangci-lint run ./internal/...
	golangci-lint run ./pkg/...

dc-up:
	docker-compose up --build -d

dc-logs:
	docker-compose logs -f

dc-down:
	docker-compose down

.PHONY: build run run-detect logs stop version help test lint dc-up dc-logs dc-down