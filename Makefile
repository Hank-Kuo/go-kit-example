VERSION ?= $(shell git describe --tags --always)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
IMAGE_NAME ?= hank-go-kit-example

# ==============================================================================

local:
	@export MODE=dev;\
	@echo Starting local docker compose
	docker-compose -f docker-compose.local.yaml up 

call:
	@echo Starting calling $$APP service
	@bash ./cmd/$$APP/call.sh

pb-complier:
	@echo Compiling $$APP proto...
	@protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. ./pb/$$APP/$$APP.proto
	
run:
	@echo Starting $$APP service
	go run ./cmd/$$APP/main.go

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

linter:
	@echo Starting linters
	golangci-lint run ./...


