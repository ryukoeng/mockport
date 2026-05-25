.PHONY: test vet build run docker-build

test:
	go test ./...

vet:
	go vet ./...

build:
	CGO_ENABLED=0 go build -o bin/mockport ./cmd/mockport

run:
	go run ./cmd/mockport run --config mockport.yml

docker-build:
	docker build -t mockport:local -f docker/Dockerfile .
