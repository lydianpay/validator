# Pinned golangci-lint version — matches the org CI default (lydianpay/ci qlty-backend.yml).
GOLANGCI := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.1

.PHONY: build test lint tidy

build:
	@go build ./...

test:
	@go test ./... -coverprofile coverage.out
	@go tool cover -html=coverage.out

lint:
	$(GOLANGCI) run

tidy:
	go mod tidy
	go mod verify
