BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
DB_USER ?= user
DB_PASSWORD ?= password
DB_DSN="host=localhost port=5432 user=$(DB_USER) password=$(DB_PASSWORD) dbname=calendar sslmode=disable"
MIGRATIONS_DIR="./migrations"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/calendar_config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race -count=100 ./internal/logger/... ./internal/storage/... ./internal/storage/memory/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...

install-goose:
	(which goose > /dev/null) || go install github.com/pressly/goose/v3/cmd/goose@v3.24.3

migrate: install-goose
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

.PHONY: build run build-img run-img version test lint install-goose migrate