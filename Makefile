BIN := "./bin/resizer"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd

run:
	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.prod.yaml -p resizer-prod up -d

test:
	go test -race -count 100 ./internal/... ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.45.2

lint: install-lint-deps
	golangci-lint run ./...

integration-tests:
	set -e; \
	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml -p integration_test_resizer up --build -d; \
	test_status_code=0 ;\
	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml run integration-tests go test -v -mod=readonly -tags integration ./tests/integration || test_status_code=$$? ;\
	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml down ;\
	exit $$test_status_code ;

.PHONY: build run test lint integration-tests
