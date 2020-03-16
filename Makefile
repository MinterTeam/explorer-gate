include .env
export

APP ?= gate
VERSION ?= $(strip $(shell cat VERSION))
GOOS ?= linux
SRC = ./

COMMIT = $(shell git rev-parse --short HEAD)
BRANCH = $(strip $(shell git rev-parse --abbrev-ref HEAD))
CHANGES = $(shell git rev-list --count ${COMMIT})
BUILDED ?= $(shell date -u '+%Y-%m-%dT%H:%M:%S')
BUILD_FLAGS = "-X main.Version=$(VERSION) -X main.GitCommit=$(COMMIT) -X main.BuildedDate=$(BUILDED)"
BUILD_TAGS?=minter-gate
DOCKER_TAG = latest
SERVER ?= gate-api.minter.network

all: get_deps build

### Tools & dependencies ####
#Run this from CI
get_deps:
	@echo "--> Running dep"
	@go mod download

### Build ###################
build: clean
	GOOS=${GOOS} GOARCH=amd64 go build -ldflags $(BUILD_FLAGS) -o ./builds/${GOOS}/$(APP)

install:
	GOOS=${GOOS} GOARCH=amd64 go install -ldflags $(BUILD_FLAGS)

clean:
	@rm -f $(BINARY)

fmt:
	@go fmt ./...

docker_build:
	docker build -f .docker/gate/Dockerfile --build-arg DOKCER_GO_VER=$(DOKCER_GO_VER) .

.PHONY: get_deps build install clean fmt all
