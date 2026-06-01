.PHONY: test vet build run docker-build

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/albert-einshutoin/mockport/internal/cli.Version=$(VERSION)

test:
	go test ./...

vet:
	go vet ./...

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/mockport ./cmd/mockport

run:
	go run ./cmd/mockport run --config mockport.yml

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t mockport:local -f docker/Dockerfile .
