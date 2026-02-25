.PHONY: build test test-cover lint fmt tidy

build:
	go build ./...

test:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

fmt:
	go fmt ./...

tidy:
	go mod tidy
