# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

OS = $(shell uname)

# Project variables
BUILD_PACKAGE ?= ./cmd/modern-go-application
BINARY_NAME ?= modern-go-application
DOCKER_IMAGE = sagikazarmark/modern-go-application
OPENAPI_DESCRIPTOR_DIR = api/openapi

# Build variables
BUILD_DIR ?= build
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
LDFLAGS += -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}
export CGO_ENABLED ?= 0
ifeq (${VERBOSE}, 1)
ifeq ($(filter -v,${GOARGS}),)
	GOARGS += -v
endif
TEST_FORMAT = short-verbose
endif

# Docker variables
DOCKER_TAG ?= ${VERSION}

# Dependency versions
GOTESTSUM_VERSION = 0.3.3
GOLANGCI_VERSION = 1.15.0
OPENAPI_GENERATOR_VERSION = 3.3.4
PROTOLOCK_VERSION = 0.10.0
PROTOLOCK_BUILD_DATE = 20190101T225741Z
GOBIN_VERSION = 0.0.4
PROTOC_GEN_GO_VERSION = 1.2.0
PROTOTOOL_VERSION = 1.3.0

GOLANG_VERSION = 1.12

# Add the ability to override some variables
# Use with care
-include override.mk

.PHONY: up
up: start config.toml ## Set up the development environment

.PHONY: down
down: clear ## Destroy the development environment
	docker-compose down
	rm -rf var/docker/volumes/*

.PHONY: reset
reset: down up ## Reset the development environment

.PHONY: clear
clear: ## Clear the working area and the project
	rm -rf bin/

docker-compose.override.yml:
	cp docker-compose.override.yml.dist docker-compose.override.yml

.PHONY: start
start: docker-compose.override.yml ## Start docker development environment
	@ if [ docker-compose.override.yml -ot docker-compose.override.yml.dist ]; then diff -u docker-compose.override.yml docker-compose.override.yml.dist || (echo "!!! The distributed docker-compose.override.yml example changed. Please update your file accordingly (or at least touch it). !!!" && false); fi
	docker-compose up -d

.PHONY: stop
stop: ## Stop docker development environment
	docker-compose stop

config.toml:
	sed 's/production/development/g; s/debug = false/debug = true/g; s/shutdownTimeout = "15s"/shutdownTimeout = "0s"/g; s/format = "json"/format = "logfmt"/g; s/level = "info"/level = "debug"/g; s/addr = ":10000"/addr = "127.0.0.1:10000"/g; s/addr = ":8000"/addr = "127.0.0.1:8000"/g' config.toml.dist > config.toml

.PHONY: run
run: GOTAGS += dev
run: build ## Build and execute a binary
	${BUILD_DIR}/${BINARY_NAME} ${ARGS}

.PHONY: clean
clean: ## Clean builds
	rm -rf ${BUILD_DIR}/

.PHONY: build
build: ## Build a binary
ifeq (${VERBOSE}, 1)
	go env
endif
ifneq (${IGNORE_GOLANG_VERSION_REQ}, 1)
	@printf "${GOLANG_VERSION}\n$$(go version | awk '{sub(/^go/, "", $$3);print $$3}')" | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -g | head -1 | grep -q -E "^${GOLANG_VERSION}$$" || (printf "Required Go version is ${GOLANG_VERSION}\nInstalled: `go version`" && exit 1)
endif

	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${BINARY_NAME} ${BUILD_PACKAGE}

.PHONY: build-release
build-release: ## Build a binary without debug information
	@${MAKE} LDFLAGS="-w ${LDFLAGS}" BUILD_DIR="${BUILD_DIR}/release" build

.PHONY: build-debug
build-debug: ## Build a binary with remote debugging capabilities
	@${MAKE} GOARGS="${GOARGS} -gcflags \"all=-N -l\"" BUILD_DIR="${BUILD_DIR}/debug" build

.PHONY: docker
docker: ## Build a Docker image
ifneq (${DOCKER_PREBUILT}, 1)
	@${MAKE} BINARY_NAME="${BINARY_NAME}-linux-amd64" GOOS=linux GOARCH=amd64 build-release
endif
	docker build --build-arg BUILD_DIR=${BUILD_DIR}/release --build-arg BINARY_NAME=${BINARY_NAME}-linux-amd64 -t ${DOCKER_IMAGE}:${DOCKER_TAG} -f Dockerfile.local .
ifeq (${DOCKER_LATEST}, 1)
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
endif

.PHONY: docker-debug
docker-debug: ## Build a Docker image with remote debugging capabilities
ifneq (${DOCKER_PREBUILT}, 1)
	@${MAKE} BINARY_NAME="${BINARY_NAME}-linux-amd64" GOOS=linux GOARCH=amd64 build-debug
endif
	docker build --build-arg BUILD_DIR=${BUILD_DIR}/debug --build-arg BINARY_NAME=${BINARY_NAME}-linux-amd64 -t ${DOCKER_IMAGE}:${DOCKER_TAG}-debug -f Dockerfile.debug .
ifeq (${DOCKER_LATEST}, 1)
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG}-debug ${DOCKER_IMAGE}:latest-debug
endif

.PHONY: check
check: test-all lint ## Run tests and linters

bin/gotestsum: bin/gotestsum-${GOTESTSUM_VERSION}
	@ln -sf gotestsum-${GOTESTSUM_VERSION} bin/gotestsum
bin/gotestsum-${GOTESTSUM_VERSION}:
	@mkdir -p bin
ifeq (${OS}, Darwin)
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_darwin_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}
endif
ifeq (${OS}, Linux)
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_linux_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}
endif

TEST_PKGS ?= ./...
TEST_REPORT_NAME ?= results.xml
.PHONY: test
test: TEST_REPORT ?= main
test: TEST_FORMAT ?= short
test: SHELL = /bin/bash
test: bin/gotestsum ## Run tests
	@mkdir -p ${BUILD_DIR}/test_results/${TEST_REPORT}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/test_results/${TEST_REPORT}/${TEST_REPORT_NAME} --format ${TEST_FORMAT} -- $(filter-out -v,${GOARGS}) $(if ${TEST_PKGS},${TEST_PKGS},./...)

.PHONY: test-all
test-all: ## Run all tests
	@${MAKE} GOARGS="${GOARGS} -run .\*" TEST_REPORT=all test

.PHONY: test-integration
test-integration: ## Run integration tests
	@${MAKE} GOARGS="${GOARGS} -run ^TestIntegration\$$\$$" TEST_REPORT=integration test

.PHONY: test-functional
test-functional: ## Run functional tests
	@${MAKE} GOARGS="${GOARGS} -run ^TestFunctional\$$\$$" TEST_REPORT=functional test

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ./bin/ v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run

.PHONY: validate-openapi
validate-openapi: ## Validate the OpenAPI descriptor
	for api in `ls ${OPENAPI_DESCRIPTOR_DIR}`; do \
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v${OPENAPI_GENERATOR_VERSION} validate --recommend -i /local/${OPENAPI_DESCRIPTOR_DIR}/$$api/swagger.yaml; \
	done

.PHONY: openapi
openapi: ## Generate client and server stubs from the OpenAPI definition
	rm -rf .gen/${OPENAPI_DESCRIPTOR_DIR}
	for api in `ls ${OPENAPI_DESCRIPTOR_DIR}`; do \
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v${OPENAPI_GENERATOR_VERSION} generate \
	--additional-properties packageName=api \
	--additional-properties withGoCodegenComment=true \
	-i /local/${OPENAPI_DESCRIPTOR_DIR}/$$api/swagger.yaml \
	-g go-server \
	-o /local/.gen/${OPENAPI_DESCRIPTOR_DIR}/$$api; \
	done

bin/protolock: bin/protolock-${PROTOLOCK_VERSION}
	@ln -sf protolock-${PROTOLOCK_VERSION} bin/protolock
bin/protolock-${PROTOLOCK_VERSION}:
	@mkdir -p bin
ifeq (${OS}, Darwin)
	curl -L https://github.com/nilslice/protolock/releases/download/v${PROTOLOCK_VERSION}/protolock.${PROTOLOCK_BUILD_DATE}.darwin-amd64.tgz | tar -zOxf - protolock > ./bin/protolock-${PROTOLOCK_VERSION} && chmod +x ./bin/protolock-${PROTOLOCK_VERSION}
endif
ifeq (${OS}, Linux)
	curl -L https://github.com/nilslice/protolock/releases/download/v${PROTOLOCK_VERSION}/protolock.${PROTOLOCK_BUILD_DATE}.linux-amd64.tgz | tar -zOxf - protolock > ./bin/protolock-${PROTOLOCK_VERSION} && chmod +x ./bin/protolock-${PROTOLOCK_VERSION}
endif

proto.lock: bin/protolock
	bin/protolock init

.PHONY: protolock
protolock: bin/protolock
	bin/protolock status

bin/gobin: bin/gobin-${GOBIN_VERSION}
	@ln -sf gobin-${GOBIN_VERSION} bin/gobin
bin/gobin-${GOBIN_VERSION}:
	@mkdir -p bin
ifeq (${OS}, Darwin)
	curl -L https://github.com/myitcv/gobin/releases/download/v${GOBIN_VERSION}/darwin-amd64 > ./bin/gobin-${GOBIN_VERSION} && chmod +x ./bin/gobin-${GOBIN_VERSION}
endif
ifeq (${OS}, Linux)
	curl -L https://github.com/myitcv/gobin/releases/download/v${GOBIN_VERSION}/linux-amd64 > ./bin/gobin-${GOBIN_VERSION} && chmod +x ./bin/gobin-${GOBIN_VERSION}
endif

bin/protoc-gen-go: bin/protoc-gen-go-${PROTOC_GEN_GO_VERSION}
	@ln -sf protoc-gen-go-${PROTOC_GEN_GO_VERSION} bin/protoc-gen-go
bin/protoc-gen-go-${PROTOC_GEN_GO_VERSION}: bin/gobin
	@mkdir -p bin
	GOBIN=bin/ bin/gobin github.com/golang/protobuf/protoc-gen-go@v1.2.0

bin/prototool: bin/prototool-${PROTOTOOL_VERSION}
	@ln -sf prototool-${PROTOTOOL_VERSION} bin/prototool
bin/prototool-${PROTOTOOL_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/uber/prototool/releases/download/v${PROTOTOOL_VERSION}/prototool-${OS}-x86_64 > ./bin/prototool-${PROTOTOOL_VERSION} && chmod +x ./bin/prototool-${PROTOTOOL_VERSION}

.PHONY: proto
proto: bin/prototool bin/protoc-gen-go protolock ## Generate client and server stubs from the protobuf definition
	bin/prototool $(if ${VERBOSE},--debug,) all

release-%: TAG_PREFIX = v
release-%:
ifneq (${DRY}, 1)
	@sed -e "s/^## \[Unreleased\]$$/## [Unreleased]\\"$$'\n'"\\"$$'\n'"\\"$$'\n'"## [$*] - $$(date +%Y-%m-%d)/g; s|^\[Unreleased\]: \(.*\/compare\/\)\(.*\)...HEAD$$|[Unreleased]: \1${TAG_PREFIX}$*...HEAD\\"$$'\n'"[$*]: \1\2...${TAG_PREFIX}$*|g" CHANGELOG.md > CHANGELOG.md.new
	@mv CHANGELOG.md.new CHANGELOG.md

ifeq (${TAG}, 1)
	git add CHANGELOG.md
	git commit -m 'Prepare release $*'
	git tag -m 'Release $*' ${TAG_PREFIX}$*
ifeq (${PUSH}, 1)
	git push; git push origin ${TAG_PREFIX}$*
endif
endif
endif

	@echo "Version updated to $*!"
ifneq (${PUSH}, 1)
	@echo
	@echo "Review the changes made by this script then execute the following:"
ifneq (${TAG}, 1)
	@echo
	@echo "git add CHANGELOG.md && git commit -m 'Prepare release $*' && git tag -m 'Release $*' ${TAG_PREFIX}$*"
	@echo
	@echo "Finally, push the changes:"
endif
	@echo
	@echo "git push; git push origin ${TAG_PREFIX}$*"
endif

.PHONY: patch
patch: ## Release a new patch version
	@${MAKE} release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[ .]' '{print $$1"."$$2"."$$3+1}')

.PHONY: minor
minor: ## Release a new minor version
	@${MAKE} release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[ .]' '{print $$1"."$$2+1".0"}')

.PHONY: major
major: ## Release a new major version
	@${MAKE} release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[ .]' '{print $$1+1".0.0"}')

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Variable outputting/exporting rules
var-%: ; @echo $($*)
varexport-%: ; @echo $*=$($*)
