SHELL := /bin/bash

build: build-bin build-docker

build-bin:
	go build -o bin/image-previewer .

build-docker:
	docker build -t sinuspower/image-previewer .

build-nginx:
	docker build -t nginx/image-server ./test/integration

run: # build-docker
	source docker-compose.env && \
	docker-compose up -d

run-bin:
	source run-bin.env && \
	./bin/image-previewer

run-nginx:
	IMAGE_SERVER_PORT=8082 docker-compose -f ./test/integration/docker-compose.yml up -d 

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
	go test -v -race ./...

lint:
	golangci-lint run ./...

.PHONY: build test