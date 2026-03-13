.PHONY: build docker-build run lint test test-short test-integration clean

BINARY_NAME := map-api
DOCKER_IMAGE := map-api:dev
ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

build:
	go build -o bin/$(BINARY_NAME) ./cmd/map-api

docker-build:
	docker build -t $(DOCKER_IMAGE) $(ROOT)

run: build
	./bin/$(BINARY_NAME)

lint:
	go vet ./...
	staticcheck ./... 2>/dev/null || true

test: test-short test-integration

test-short:
	go test -short -count=1 ./internal/...

test-integration:
	@if [ -n "$$TEST_DATABASE_URL" ]; then go test -count=1 -run Integration ./internal/adapters/repository; fi

test-coverage:
	go test -short -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out

clean:
	rm -rf bin/
