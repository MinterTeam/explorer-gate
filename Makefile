APP ?= gate
VERSION ?= $(strip $(shell cat VERSION))
GOOS ?= linux
GLIDE ?= $(shell /usr/bin/env glide 2> /dev/null)
SRC = ./

COMMIT = $(shell git rev-parse --short HEAD)
BRANCH = $(strip $(shell git rev-parse --abbrev-ref HEAD))
CHANGES = $(shell git rev-list --count ${COMMIT})
BUILDED ?= $(shell date -u '+%Y-%m-%dT%H:%M:%S')
GOLDFLAGS = "-X main.Version=$(VERSION) -X main.GitCommit=$(COMMIT) -X main.BuildedDate=$(BUILDED)"
BINARY = builds/$(APP)
DOCKER_TAG = latest
SERVER ?= explorer.minter.network

build: clean
	GOOS=${GOOS} go build -ldflags $(GOLDFLAGS) -o $(BINARY)

clean:
	@rm -f $(BINARY)

vendor:
ifndef GLIDE
	@curl https://glide.sh/get | sh
endif
	glide install -v

install: vendor build
	@cp -f $(BINARY) /usr/local/bin
	@cp -f explorer-api.service /etc/systemd/system

update:
	glide update
	$(BINARY) update

deploy: build
	@echo "--> Running latest binary"
	@ssh root@$(SERVER) systemctl stop explorer-api
	@scp $(BINARY) root@$(SERVER):/srv/minter/explorer-api
	@ssh root@$(SERVER) 'systemctl daemon-reload && systemctl restart explorer-api'

run: build
	@echo "--> Running latest binary"
	@$(BINARY)

test:
	@echo "--> Running tests"
	#go tool vet $(strip $(shell glide novendor))
	go test $(strip $(shell glide novendor))

.PHONY: build clean vendor update install test
