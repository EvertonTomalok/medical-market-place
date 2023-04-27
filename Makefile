ENV_VARS =$(shell grep -v "^\#" .env | xargs)

envvars:
	@echo "export $(ENV_VARS)"

swag:
	swag init --parseDependency  --parseInternal --parseDepth 1

build:
	go build -o exe main.go

PHONY: up
up:
	docker-compose up

PHONY: up-build
up-build:
	docker-compose up --build

PHONY: down
down:
	docker-compose down

.PHONY: pre-run
pre-run: envvars swag

.PHONY: server-dev
server-dev: pre-run
	go run main.go server

.PHONY: server
server: pre-run build
	./exe server

.PHONY: setup 
setup: 
	go mod tidy
	go install github.com/golang/mock/mockgen@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/rakyll/gotest@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

.PHONY: lint 
lint: 
	gosec ./...
	golangci-lint run -v

.PHONY: setup_dev
setup_dev: 
	go install github.com/go-delve/delve/cmd/dlv@latest

.PHONY: format
format: gofmt goimports

.PHONY: gofmt
gofmt: 
	@echo "Formatting go files..."
	@go fmt ./...

.PHONY: goimports
goimports: 
	@echo "Organizing the imports"
	@goimports -w $(shell \
					find . -not \( \( -name .git -o -name .go \) -prune \) \
					-name '*.go')

.PHONY: test
test: 
	@echo "Running tasks..."
	@go test -short -parallel 30 -failfast -coverprofile=coverage.out.tmp ./internal/app/... ./internal/infra/postgres/entities/... ./internal/infra/postgres/repositories/...  ./pkg/...
	@cat coverage.out.tmp | grep -v "_mock.go" | grep -v "_gen.go" > coverage.out
	@rm -rf coverage.out.tmp
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: test-ci
test-ci: setup test lint

.PHONY: generate
generate: 
	@echo "Generating code..."
	@go generate ./...
