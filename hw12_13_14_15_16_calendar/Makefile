BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"
DOCKER_IMG_CALENDAR="calendar:develop"
DOCKER_IMG_SCHEDULER="scheduler:develop"
DOCKER_IMG_SENDER="sender:develop"
DB_USER ?= postgres
DB_PASSWORD ?= password
DB_DSN="host=localhost port=5432 user=$(DB_USER) password=$(DB_PASSWORD) dbname=calendar sslmode=disable"
MIGRATIONS_DIR="./migrations"
PROTO_DIR="./api"
PROTO_FILE="calendar.proto"
GRPC_OUT_DIR="./internal/server/grpc/pb"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN_CALENDAR) -config ./configs/calendar_config.yaml
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.yaml
	$(BIN_SENDER) -config ./configs/sender_config.yaml

build-img:
	docker build --build-arg=LDFLAGS="$(LDFLAGS)" -t $(DOCKER_IMG_CALENDAR) -f ./build/calendar/Dockerfile .
	docker build --build-arg=LDFLAGS="$(LDFLAGS)" -t $(DOCKER_IMG_SCHEDULER) -f ./build/scheduler/Dockerfile .
	docker build --build-arg=LDFLAGS="$(LDFLAGS)" -t $(DOCKER_IMG_CALENDAR) -f ./build/sender/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race -count=10 ./internal/logger/... ./internal/storage/... ./internal/storage/memory/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...

install-goose:
	(which goose > /dev/null) || go install github.com/pressly/goose/v3/cmd/goose@v3.24.3

migrate: install-goose
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

install-protoc-deps:
	(which protoc > /dev/null) || (echo "Please install protoc: https://grpc.io/docs/protoc-installation/" && exit 1)
	(which protoc-gen-go > /dev/null) || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	(which protoc-gen-go-grpc > /dev/null) || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate: install-protoc-deps
	@echo "Generating gRPC code from $(PROTO_DIR)/$(PROTO_FILE)..."
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(GRPC_OUT_DIR) \
		--go-grpc_out=$(GRPC_OUT_DIR) \
		$(PROTO_DIR)/$(PROTO_FILE)
	@echo "gRPC code generated in $(GRPC_OUT_DIR)"

up:
	docker compose -f './deployments/docker/docker-compose.yaml' up -d --build 

down:
	docker compose -f './deployments/docker/docker-compose.yaml' down 

integration-tests:
	docker compose -f './deployments/docker/docker-compose.yaml' -f './deployments/docker/docker-compose.test.yaml' up --build --exit-code-from integration-tests-otus ; \
	EXIT_CODE=$$?; \
	docker compose -f './deployments/docker/docker-compose.yaml' -f './deployments/docker/docker-compose.test.yaml' down; \
	docker rmi -f 1aeef68f-calendar:develop 1aeef68f-scheduler:develop 1aeef68f-sender:develop 1aeef68f-integration-tests:develop 2>/dev/null || true; \
	docker system prune -f; \
	echo "integration-tests-otus exited with code: $$EXIT_CODE" ; \
	exit $$EXIT_CODE

.PHONY: build run build-img run-img version test lint install-goose migrate install-protoc-deps generate up down integration-tests