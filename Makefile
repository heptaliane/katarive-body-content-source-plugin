.PHONY: build test format

build:
	go build .

test:
	go test ./...

format:
	go fmt ./...
