BINARY   := sailor
PKG      := ./...
TESTPKGS := ./internal/...

.PHONY: build test test-cover cover-html lint

build:
	go build -v -o $(BINARY) main.go

test:
	go test $(TESTPKGS) -v

test-cover:
	go test $(TESTPKGS) -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

cover-html: test-cover
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	go vet $(PKG)
