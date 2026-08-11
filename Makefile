APP_NAME=user-service
BINARY_NAME=api

.PHONY: run build test clean fmt vet

run:
	go run ./cmd/api

build:
	mkdir -p bin
	go build -o bin/$(BINARY_NAME) ./cmd/api

test:
	go test -v ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin