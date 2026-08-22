# golang-cli-template — development tasks
GO ?= go
BINARY ?= golang-cli-template

.PHONY: all fmt vet tidy lint test test-race cover build clean run

all: fmt vet lint test

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

lint:
	golangci-lint run ./...

test:
	$(GO) test -count=1 ./...

test-race:
	$(GO) test -race -count=1 ./...

cover:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out

build:
	$(GO) build -o bin/$(BINARY) ./cmd/app

run:
	$(GO) run ./cmd/app

clean:
	rm -f coverage.out
	rm -rf bin/
