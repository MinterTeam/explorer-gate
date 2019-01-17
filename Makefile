APP ?= gate
VERSION ?= $(strip $(shell cat VERSION))
GOOS ?= linux
SRC = ./

COMMIT = $(shell git rev-parse --short HEAD)
BRANCH = $(strip $(shell git rev-parse --abbrev-ref HEAD))
CHANGES = $(shell git rev-list --count ${COMMIT})
BUILDED ?= $(shell date -u '+%Y-%m-%dT%H:%M:%S')
GOLDFLAGS = "-X main.Version=$(VERSION) -X main.GitCommit=$(COMMIT) -X main.BuildedDate=$(BUILDED)"
BINARY = builds/$(APP)
DOCKER_TAG = latest
SERVER ?= gate.minter.network

build: clean
	GOOS=${GOOS} go build -ldflags $(GOLDFLAGS) -o $(BINARY)

clean:
	@rm -f $(BINARY)

test:
	@echo "--> Running tests"
	go test -v ${SRC}

fmt:
	@go fmt ./...

.PHONY: build clean fmt test
