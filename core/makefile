build:
	go build -o ./.bin/main cmd/main/main.go

run: build
	./.bin/main

IMAGE_NAME = gses-api
VOLUME_PATH = $(shell pwd)/local

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE_NAME) .

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 -v $(VOLUME_PATH):/root/local --env-file .env $(IMAGE_NAME)

test-unit:
	go test -v ./...

test-integration:
	go test -v ./... --tags=integration