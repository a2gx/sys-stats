BIN := ./bin/daemon
APP := ./cmd/daemon

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" $(APP)

run: build
	$(BIN) run --host localhost --port 50051 --config ./configs/config.yaml

logs: build
	$(BIN) logs

stop: build
	$(BIN) stop

version: build
	$(BIN) --version

help: build
	$(BIN) -h

test:
	go test -race ./...

lint-install-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1.6

lint: lint-install-deps
	golangci-lint run ./...

generate:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative proto/**/*.proto

dc-up:
	docker-compose up --build
dc-down:
	docker-compose down
dc-logs:
	docker-compose logs -f

.PHONY: build run logs stop version help test lint generate dc-up dc-down dc-logs