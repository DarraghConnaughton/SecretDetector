# Makefile
.PHONY: build test coverage
# Output binary name with timestamp
BINARY_NAME = secretdetector

build:
	make clean
	mkdir releases
	go build -o ./releases/$(BINARY_NAME) ./cmd

clean:
	@if [ -d ./releases/ ]; then rm -rf ./releases/; fi

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./... | tee coverage.report
	go tool cover -html=coverage.out -o coverage.html
