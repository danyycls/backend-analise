.PHONY: all build run lint lint-fix fix fmt fmt-check test test-unit test-integration test-cover test-race vet tidy deps generate generate-mocks generate-check clean help

APP_NAME := podp
CMD_DIR  := .
API_DIR  := api
GEN_DIR  := internal/generated

GO        := go
GOFMT     := gofmt
GOLINT    := $(shell which golangci-lint 2>/dev/null || echo "go run github.com/golangci/golangci-lint/cmd/golangci-lint")

all: lint test build

help: ## Exibe ajuda com os comandos disponiveis
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Compila o binario
	$(GO) build -o $(APP_NAME) $(CMD_DIR)

run: ## Executa a aplicacao localmente
	$(GO) run $(CMD_DIR)

run-dev: ## Executa com live reload (necessita: go install github.com/cosmtrek/air@latest)
	air

lint: ## Executa analise estatica (golangci-lint)
	$(GOLINT) run ./...

lint-fix: ## Executa analise estatica com correcoes automaticas
	$(GOLINT) run --fix ./...

fmt: ## Formata o codigo
	$(GO) fmt ./...

fix: ## Formata codigo, organiza modulos e aplica autofix do lint
	gofmt -w $(shell find . -path ./dataCSV -prune -o -name '*.go' -print)
	$(GO) mod tidy
	$(GO) mod verify
	-$(GOLINT) run --fix ./... 2>&1 || true

fmt-check: ## Verifica se o codigo esta formatado
	@files=$$(find . -path ./dataCSV -prune -o -name '*.go' -print); \
	test -z "$$($(GOFMT) -l $$files)" || (echo "Arquivos precisam ser formatados:"; $(GOFMT) -l $$files; exit 1)

test: ## Executa todos os testes (unitarios + integracao)
	$(GO) test $(shell $(GO) list ./... | grep -v /dataCSV) -count=1 -timeout 300s

test-unit: ## Executa apenas testes unitarios (pula integracao com servicos externos)
	$(GO) test -short $(shell $(GO) list ./... | grep -v /dataCSV) -count=1 -timeout 60s

test-integration: ## Executa apenas testes de integracao (testcontainers, servicos externos)
	$(GO) test -count=1 -timeout 300s -run "TestRepositorioFornecedores|TestSucesso|TestFalha" ./internal/esferas-brasileiras/tse/ ./internal/ligacao-politica/handler/

test-cover: ## Executa testes com cobertura
	$(GO) test -short $(shell $(GO) list ./... | grep -v /dataCSV) -count=1 -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html

test-race: ## Executa testes com detector de race conditions
	$(GO) test -short $(shell $(GO) list ./... | grep -v /dataCSV) -race -count=1

vet: ## Executa go vet
	$(GO) vet $(shell $(GO) list ./... | grep -v /dataCSV)

tidy: ## Organiza as dependencias do modulo
	$(GO) mod tidy
	$(GO) mod verify

deps: ## Baixa todas as dependencias
	$(GO) mod download

generate: ## Gera codigo (OpenAPI -> Go com oapi-codegen + mocks)
	@mkdir -p $(GEN_DIR)
	oapi-codegen -config $(API_DIR)/oapi-codegen.yaml $(API_DIR)/openapi.yaml

generate-mocks: ## Gera mocks com mockgen
	$(GO) generate ./internal/shared/clients/opencnpj/...
	$(GO) generate ./internal/shared/clients/tcu/...
	$(GO) generate ./internal/shared/redis/...
	$(GO) generate ./internal/ligacao-politica/usecase/...

generate-check: ## Verifica se o codigo gerado esta atualizado
	@cp -r $(GEN_DIR) $(GEN_DIR).bak
	@$(MAKE) generate
	@diff -r $(GEN_DIR) $(GEN_DIR).bak || (echo "Codigo gerado esta desatualizado. Rode 'make generate'."; rm -rf $(GEN_DIR).bak; exit 1)
	@rm -rf $(GEN_DIR).bak

clean: ## Remove binarios e artefatos de build
	rm -f $(APP_NAME)
	rm -f coverage.out coverage.html
	rm -rf $(GEN_DIR)
