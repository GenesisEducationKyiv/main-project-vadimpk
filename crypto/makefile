build:
	go build -o ./.bin/app cmd/main.go

run: build
	./.bin/app

IMAGE_NAME = gses-crypto-service

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE_NAME) .

.PHONY: docker-run
docker-run:
	docker run -p 8081:8081 --env-file .env $(IMAGE_NAME)

test-functional:
	go test -v ./... --tags=functional