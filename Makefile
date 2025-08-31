BINARY=kafkaesque

.PHONY: build run test fmt lint

build:
	go build -o bin/$(BINARY) ./cmd/kafkaesque

run: build
	./bin/$(BINARY)

test:
	go test ./... -race -count=1

fmt:
	go fmt ./...
	go vet ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; else echo "Install golangci-lint: https://golangci-lint.run/"; fi
