VERSION := v1.0.0

BIN := confgr
PKG := github.com/elliottpolk

BUILD_IMAGE ?= golang:alpine

BUILD_NUM := $(shell date +'%s')
IMAGE := $(BIN):$(VERSION)-$(BUILD_NUM)

all: container

build: clean build-dirs
	@docker run --rm -it 			\
		-v $(GOPATH):/go 			\
		-w /go/src/$(PKG)/$(BIN)	\
		$(BUILD_IMAGE) /bin/sh -c 	\
		'go test -v ./... && GOOS=$(GOOS) go build -o $(BIN)'
	@mv $(BIN) build/bin

container: build
	@docker build -t $(IMAGE) .
	@echo "\ncontainer: $(IMAGE)"

test:
	@docker run --rm -it 			\
		-v $(GOPATH):/go 			\
		-w /go/src/$(PKG)/$(BIN)	\
		$(BUILD_IMAGE) /bin/sh -c 	\
		'go test -v ./...'

build-dirs:
	@mkdir -p build/bin

clean:
	@rm -rf build/