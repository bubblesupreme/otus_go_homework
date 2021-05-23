BIN = "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH = $(shell git log --format="%h" -n 1)
LDFLAGS = -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
GOOSE = goose
SQL_DRIVER = postgres
SQL_USER = ${POSTGRES_USER}
SQL_PASSWORD = ${POSTGRES_PASSWORD}
SQL_DB= ${POSTGRES_DB}
SQL_HOST = ${POSTGRES_HOST}
SQL_PORT = ${POSTGRES_PORT}

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) --config ./configs/config.json

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
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

migrate-up:
	$(GOOSE) -v $(SQL_DRIVER) "user=$(SQL_USER) dbname=$(SQL_DB) password=$(SQL_PASSWORD) port=$(SQL_PORT) host=$(SQL_HOST) sslmode=disable" up

migrate-down:
	$(GOOSE) -v $(SQL_DRIVER) "user=$(SQL_USER) dbname=$(SQL_DB) password=$(SQL_PASSWORD) port=$(SQL_PORT) host=$(SQL_HOST) sslmode=disable" down

.PHONY: build run build-img run-img version test lint
