run:
	GOFLAGS=-mod=mod go run ./...

build:
	GOFLAGS=-mod=mod go build -o bin/url-shortener ./main.go

docker:
	docker build -t url-shortener:local .

up:
	docker run --rm -p 8080:8080 --env-file .env url-shortener:local

test:
	go test ./... -v
