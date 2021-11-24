include main.mk

# Project variables
OPENAPI_DESCRIPTOR_DIR = api/openapi

# Dependency versions
MGA_VERSION = 0.2.0

.PHONY: up
up: start config.toml ## Set up the development environment

.PHONY: down
down: clear ## Destroy the development environment
	docker-compose down --volumes --remove-orphans --rmi local
	rm -rf var/docker/volumes/*

.PHONY: reset
reset: down up ## Reset the development environment

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
	sed 's/production/development/g; s/debug = false/debug = true/g; s/shutdownTimeout = "15s"/shutdownTimeout = "0s"/g; s/format = "json"/format = "logfmt"/g; s/level = "info"/level = "debug"/g; s/addr = ":10000"/addr = "127.0.0.1:10000"/g; s/httpAddr = ":8000"/httpAddr = "127.0.0.1:8000"/g; s/grpcAddr = ":8001"/grpcAddr = "127.0.0.1:8001"/g' config.toml.dist > config.toml

bin/entc:
	@mkdir -p bin
	go build -o bin/entc github.com/facebook/ent/cmd/entc

bin/mga: bin/mga-${MGA_VERSION}
	@ln -sf mga-${MGA_VERSION} bin/mga
bin/mga-${MGA_VERSION}:
	@mkdir -p bin
	curl -sfL https://git.io/mgatool | bash -s v${MGA_VERSION}
	@mv bin/mga $@

.PHONY: generate
generate: bin/mga bin/entc ## Generate code
	go generate -x ./...
	mga generate kit endpoint ./internal/app/mga/todo/...
	mga generate event handler --output subpkg:suffix=gen ./internal/app/mga/todo/...
	mga generate event dispatcher --output subpkg:suffix=gen ./internal/app/mga/todo/...
	entc generate ./internal/app/mga/todo/todoadapter/ent/schema
