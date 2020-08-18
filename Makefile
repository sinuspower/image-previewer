ifneq (,$(wildcard ./.env))
    include .env
endif

build:
	go build -o bin/image-previewer .

run:
	./bin/image-previewer

test:
	go test -race ./...

lint:
	golangci-lint run ./...

.PHONY: build