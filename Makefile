.PHONY: build test lint clean run

BINARY := promptkit
PKG := ./...

build:
	go build -o bin/$(BINARY) ./cmd/promptkit

test:
	go test -race -count=1 $(PKG)

cover:
	go test -race -coverprofile=coverage.txt $(PKG)
	go tool cover -html=coverage.txt -o coverage.html

lint:
	golangci-lint run

clean:
	rm -rf bin/ coverage.txt coverage.html

run: build
	./bin/$(BINARY)
