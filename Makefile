#!make
include .env

default: help

.PHONY: help
help:   ## show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

.PHONY: dependencies
dependencies: ## install dependencies
	go mod tidy
	go mod verify
	go mod vendor

.PHONY: test
test: ## run tests
	go test ./internal/... -coverprofile coverage.out
	go tool cover -func coverage.out | grep ^total:

.PHONY: tools
tools: ## get tools
	git config core.hooksPath .githooks
	go get -u github.com/golang/mock/gomock
	go get -u github.com/golang/mock/mockgen

.PHONY: generate-mocks
generate-mocks: ## generate mocks
	mockgen -package=mock_converters -source internal/pkg/converters/converter.go -destination=internal/pkg/converters/mocks/converter_mock.go	go mod verify

.PHONY: generate-api-schema
generate-api-schema: ## generate api schema using swagger
	./scripts/bundle-api.sh
	./scripts/generate-code.sh

.PHONY: run-local
run-local: ## run the application locally
	go run cmd/app/main.go

# .PHONY: migrate
# migrate: ## wraps golang-migrate. Use with arguments such as 'up', 'down 2', 'version' etc. run 'migrate help for more info'
# 	migrate -database '$(DB_URI)' -path ./migrations $(RUN_ARGS)

# .PHONY: migrate-create
# migrate-create: ## creates migrations with one argument for a suffix
# 	migrate -database '$(DB_URI)' create -dir migrations -ext .sql $(RUN_ARGS)

