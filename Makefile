SHELL := /bin/bash

build: build-bin build-docker

build-bin:
	go build -o bin/image-previewer .

build-docker:
	docker build -t sinuspower/image-previewer .

run: build-docker
	source docker-compose.env && \
	docker-compose up -d

run-bin:
	source run-bin.env && \
	./bin/image-previewer

down:
	source docker-compose.env && \
	docker-compose down --rmi all -v

start:
	source docker-compose.env && \
	docker-compose start image-previewer

stop:
	source docker-compose.env && \
	docker-compose stop image-previewer

test:
	go test -race ./...

lint:
	golangci-lint run ./...

.PHONY: build