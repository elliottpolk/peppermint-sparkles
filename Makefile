
BIN=sparkles
PKG=git.platform.manulife.io/oa-montreal/peppermint-sparkles
BUILD_IMAGE=golang:latest
VERSION=`cat .version`
GOOS?=linux
PACKAGER?=tar

M = $(shell printf "\033[34;1m◉\033[0m")

default: clean build ;                                              @ ## defaulting to clean and build

.PHONY: all
all: clean unit-tests test-integration test-all build package

.PHONY: build
build: ; $(info $(M) building ...)                                  @ ## build the binary
	@mkdir -p ./build/bin/
	@GOOS=$(GOOS) go build -ldflags "-X main.version=$(VERSION)" -o ./build/bin/$(BIN)

.PHONY: package
package: ; $(info $(M) packaging ...)                               @ ## package up the binary for distribution to Artifactory or PCF
ifeq ($(PACKAGER),zip)
	@cd ./build/bin/ && zip $(BIN).zip $(shell ls -A ./build/bin) && cd -
else
	@cd ./build/bin/ && tar jcvf $(BIN).tar.bz2 $(shell ls -A ./build/bin) && cd -
endif

.PHONY: distro
distro: ;                                          					@ ## build and package in a distro dir for each OS
	@printf "\033[34;1m◉\033[0m cleaning up ...\n" \
		&& rm -vrf dist; mkdir dist
	@printf "\033[34;1m◉\033[0m building for Linux ...\n" \
		&& GOOS=linux $(MAKE) clean build package \
		&& mv ./build/bin/$(BIN).tar.bz2 dist/peppermint-sparkles-v$(VERSION).linux.tar.bz2
	@printf "\033[34;1m◉\033[0m building for macOS ...\n" \
		&& GOOS=darwin $(MAKE) clean build package \
		&& mv ./build/bin/$(BIN).tar.bz2 dist/peppermint-sparkles-v$(VERSION).macos.tar.bz2
	@printf "\033[34;1m◉\033[0m building for Windows ...\n" \
		&& GOOS=windows $(MAKE) clean build \
		&& $(MAKE) package && mv ./build/bin/$(BIN).tar.bz2 dist/peppermint-sparkles-v$(VERSION).windows.tar.bz2 \
		&& PACKAGER=zip $(MAKE) package && mv ./build/bin/$(BIN).zip dist/peppermint-sparkles-v$(VERSION).windows.zip
	@$(MAKE) clean

.PHONY: install
install: ; $(info $(M) installing locally...)                       @ ## install the binary locally
	@GOOS=$(GOOS) go build -ldflags "-X main.version=$(VERSION)" -o $(GOPATH)/bin/$(BIN)

.PHONY: unit-tests
unit-tests: ; $(info $(M) running unit tests ...)                   @ ## run only the unit tests
	@go test -v -cover ./... && go test -v -cover ./...

.PHONY: test-integration
test-integration: ; $(info $(M) running integration tests ...)      @ ## run the integration tests which has deps on Docker
	@go test -v -cover -tags="docker_env integration" ./...

.PHONY: test-all
test-all: unit-tests test-integration ;                             @ ## run all the tests

.PHONY: clean
clean: ; $(info $(M) running clean ...)                             @ ## clean up the old build dir
	@rm -vrf build

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

