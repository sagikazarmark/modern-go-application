# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

# Project variables
PACKAGE = $(shell echo $${PWD\#\#*src/})
BINARY_NAME = $(shell basename $$PWD)
DOCKER_IMAGE = $(shell echo ${PACKAGE} | cut -d '/' -f 2,3)
OPENAPI_DESCRIPTOR = swagger.yaml

# Build variables
BUILD_DIR = build
BUILD_PACKAGE = ${PACKAGE}/cmd
VERSION ?= $(shell git rev-parse --abbrev-ref HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
LDFLAGS += -X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE}
export CGO_ENABLED ?= 0
export GOOS = $(shell go env GOOS)
ifeq (${VERBOSE}, 1)
	GOARGS += -v
endif

# Docker variables
DOCKER_TAG ?= ${VERSION}

# Dependency versions
DEP_VERSION = 0.5.0
GOLANGCI_VERSION = 1.10.2
OPENAPI_GENERATOR_VERSION = 3.3.0

GOLANG_VERSION = 1.11

.PHONY: up
up: vendor start .env .env.test ## Set up the development environment

.PHONY: down
down: clean ## Destroy the development environment
	docker-compose down
	rm -rf .docker/

.PHONY: reset
reset: down up ## Reset the development environment

.PHONY: clean
clean: ## Clean the working area and the project
	rm -rf bin/ ${BUILD_DIR}/ vendor/

docker-compose.override.yml:
	cp docker-compose.override.yml.dist docker-compose.override.yml

.PHONY: start
start: docker-compose.override.yml ## Start docker development environment
	docker-compose up -d

.PHONY: stop
stop: ## Stop docker development environment
	docker-compose stop

bin/dep: bin/dep-${DEP_VERSION}
	@ln -sf dep-${DEP_VERSION} bin/dep
bin/dep-${DEP_VERSION}:
	@mkdir -p bin
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | INSTALL_DIRECTORY=bin DEP_RELEASE_TAG=v${DEP_VERSION} sh
	@mv bin/dep $@

.PHONY: vendor
vendor: bin/dep ## Install dependencies
	bin/dep ensure -v -vendor-only

.env:
	cp .env.dist .env

.env.test:
	cp .env.dist .env.test

.PHONY: run
run: GOTAGS += dev
run: build .env ## Build and execute a binary
	${BUILD_DIR}/${BINARY_NAME}-${GOOS} ${ARGS}

.PHONY: debug
debug: GOTAGS += dev
debug: build-debug bin/dlv .env ## Build and execute a binary with remote debugging enabled
	bin/dlv --listen=127.0.0.1:40000 --headless=true --api-version=2 --log exec -- ${BUILD_DIR}/${BINARY_NAME}-debug-${GOOS} ${ARGS}

bin/dlv:
	@mkdir -p bin
	GOBIN=${PWD}/bin go get -u github.com/derekparker/delve/cmd/dlv

.PHONY: build
build: GOARGS += -tags "${GOTAGS}" -ldflags "${LDFLAGS}"
build: ## Build a binary
ifeq (${VERBOSE}, 1)
	go env
endif
ifneq (${IGNORE_GOLANG_VERSION_REQ}, 1)
	@printf "${GOLANG_VERSION}\n$$(go version | awk '{sub(/^go/, "", $$3);print $$3}')" | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -g | head -1 | grep -q -E "^${GOLANG_VERSION}$$" || (printf "Required Go version is ${GOLANG_VERSION}\nInstalled: `go version`" && exit 1)
endif

	@$(eval GENERATED_BINARY_NAME = ${BINARY_NAME})
	@$(if $(strip ${BINARY_NAME_SUFFIX}),$(eval GENERATED_BINARY_NAME = ${BINARY_NAME}-$(subst $(eval) ,-,$(strip ${BINARY_NAME_SUFFIX}))),)
	go build ${GOARGS} -o ${BUILD_DIR}/${GENERATED_BINARY_NAME} ${BUILD_PACKAGE}

.PHONY: build-release
build-release: LDFLAGS += -w
build-release: build ## Build a binary without debug information

.PHONY: build-debug
build-debug: GOARGS += -gcflags "all=-N -l"
build-debug: BINARY_NAME_SUFFIX += debug
build-debug: build ## Build a binary with remote debugging capabilities

.PHONY: docker
docker: export GOOS = linux
docker: BINARY_NAME_SUFFIX += docker
docker: build-release ## Build a Docker image
	docker build --build-arg BUILD_DIR=${BUILD_DIR} --build-arg BINARY_NAME=${GENERATED_BINARY_NAME} -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
ifeq (${DOCKER_LATEST}, 1)
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
endif

.PHONY: docker-debug
docker-debug: export GOOS = linux
docker-debug: BINARY_NAME_SUFFIX += docker
docker-debug: build-debug ## Build a Docker image with remote debugging capabilities
	docker build --build-arg BUILD_DIR=${BUILD_DIR} --build-arg BINARY_NAME=${GENERATED_BINARY_NAME} -t ${DOCKER_IMAGE}:${DOCKER_TAG}-debug -f Dockerfile.debug .

.PHONY: check
check: test lint ## Run tests and linters

.PHONY: test
test: GOTAGS ?= unit integration acceptance
test: GOARGS += -tags "${GOTAGS}"
test: ## Run all tests
	go test ${GOARGS} ./...

.PHONY: test-%
test-%: ## Run a specific test suite
	@${MAKE} VERBOSE=0 GOTAGS=$* test

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
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v${OPENAPI_GENERATOR_VERSION} validate --recommend -i /local/${OPENAPI_DESCRIPTOR}

.PHONY: generate-api
generate-api: ## Generate server stubs from the OpenAPI descriptor
	rm -rf .gen/openapi
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v${OPENAPI_GENERATOR_VERSION} generate \
	--additional-properties packageName=api \
	--additional-properties withGoCodegenComment=true \
	-i /local/${OPENAPI_DESCRIPTOR} \
	-g go-server \
	-o /local/.gen/openapi

release-%: ## Release a new version
	@sed -e "s/^## \[Unreleased\]$$/## [Unreleased]\\"$$'\n'"\\"$$'\n'"\\"$$'\n'"## [$*] - $$(date +%Y-%m-%d)/g" CHANGELOG.md > CHANGELOG.md.new
	@mv CHANGELOG.md.new CHANGELOG.md

	@sed -e "s|^\[Unreleased\]: \(.*\)HEAD$$|[Unreleased]: https://${PACKAGE}/compare/v$*...HEAD\\"$$'\n'"[$*]: \1v$*|g" CHANGELOG.md > CHANGELOG.md.new
	@mv CHANGELOG.md.new CHANGELOG.md

ifeq (${TAG}, 1)
	git add CHANGELOG.md
	git commit -s -S -m 'Prepare release v$*'
	git tag -s -m 'Release v$*' v$*
endif

	@echo "Version updated to $*!"
	@echo
	@echo "Review the changes made by this script then execute the following:"
ifneq (${TAG}, 1)
	@echo
	@echo "git add CHANGELOG.md && git commit -S -m 'Prepare release v$*' && git tag -s -m 'Release v$*' v$*"
	@echo
	@echo "Finally, push the changes:"
endif
	@echo
	@echo "git push; git push --tags"

.PHONY: patch
patch: ## Release a new patch version
	@${MAKE} release-$(shell git describe --abbrev=0 --tags | sed 's/^v//' | awk -F'[ .]' '{print $$1"."$$2"."$$3+1}')

.PHONY: minor
minor: ## Release a new minor version
	@${MAKE} release-$(shell git describe --abbrev=0 --tags | sed 's/^v//' | awk -F'[ .]' '{print $$1"."$$2+1".0"}')

.PHONY: major
major: ## Release a new major version
	@${MAKE} release-$(shell git describe --abbrev=0 --tags | sed 's/^v//' | awk -F'[ .]' '{print $$1+1".0.0"}')

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
