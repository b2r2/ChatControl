.PHONY: build
build:
	go build -v ./cmd/botapi

.DEFAULT_GOAT := build