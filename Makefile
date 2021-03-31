
BIN=sparkles
PKG=git.platform.manulife.io/oa-montreal/peppermint-sparkles
BUILD_IMAGE=golang:latest
VERSION=$(shell cat .version)
GOOS?=linux
PACKAGER?=tar

# docker vars
DOCKER_GRADLE_VER=$(shell cat ./docker/peppermint-sparkles-helper/gradle.version)
DOCKER_MAVEN_VER=$(shell cat ./docker/peppermint-sparkles-helper/maven.version)
DOCKER_HELPER_IMG=peppermint-sparkles-helper
DOCKER_ARTIFACTORY=docker.artifactory.platform.manulife.io

M = $(shell printf "\033[34;1m◉\033[0m")

default: clean build ;                                              @ ## defaulting to clean and build

.PHONY: all
all: clean unit-tests test-integration test-all build package

.PHONY: build
build: ; $(info $(M) building ...)                                  @ ## build the binary
	@mkdir -p ./build/bin/
	@GOOS=$(GOOS) go build -ldflags "-X main.version=$(VERSION)" -o ./build/bin/$(BIN) cmd/*.go

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
	@GOOS=$(GOOS) go build -ldflags "-X main.version=$(VERSION)" -o $(GOPATH)/bin/$(BIN) cmd/*.go

.PHONY: unit-tests
unit-tests: ; $(info $(M) running unit tests ...)                   @ ## run only the unit tests
	@go test -v -cover ./...

.PHONY: test-integration
test-integration: ; $(info $(M) running integration tests ...)      @ ## run the integration tests which has deps on Docker
	@go test -v -cover -tags="docker_env integration" ./...

.PHONY: test-all
test-all: unit-tests test-integration ;                             @ ## run all the tests

.PHONY: clean
clean: ; $(info $(M) running clean ...)                             @ ## clean up the old build dir
	@rm -vrf build

.PHONY: docker-helper
docker-helper: docker-helper-gradle docker-helper-maven docker-helper-push

.PHONY: docker-helper-gradle
docker-helper-gradle: ; $(info $(M) running docker helper gradle...) @ ## build the docker helper gradle image
	@cd ./docker/$(DOCKER_HELPER_IMG)/ && \
		docker build -t $(DOCKER_HELPER_IMG):$(DOCKER_GRADLE_VER)-gradle -f gradle.Dockerfile .
		docker tag $(DOCKER_HELPER_IMG):$(DOCKER_GRADLE_VER)-gradle $(DOCKER_ARTIFACTORY)/$(DOCKER_HELPER_IMG):$(DOCKER_GRADLE_VER)-gradle

.PHONY: docker-helper-maven
docker-helper-maven: ; $(info $(M) running docker helper maven...)   @ ## build the docker helper gradle image
	@cd ./docker/$(DOCKER_HELPER_IMG)/ && \
		docker build -t $(DOCKER_HELPER_IMG):$(DOCKER_MAVEN_VER)-maven -f maven.Dockerfile .
		docker tag $(DOCKER_HELPER_IMG):$(DOCKER_MAVEN_VER)-maven $(DOCKER_ARTIFACTORY)/$(DOCKER_HELPER_IMG):$(DOCKER_MAVEN_VER)-maven

.PHONY: docker-helper-push
docker-helper-push: ; $(info $(M) running docker helper push...)   @ ## push the docker helper images to artifactory
	@docker push $(DOCKER_ARTIFACTORY)/$(DOCKER_HELPER_IMG):$(DOCKER_GRADLE_VER)-gradle
	@docker push $(DOCKER_ARTIFACTORY)/$(DOCKER_HELPER_IMG):$(DOCKER_MAVEN_VER)-maven

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

