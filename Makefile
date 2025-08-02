# Makefile para o projeto Back Project

.PHONY: test test-unit test-integration build run clean help

# Variáveis
GO_CMD=go
BUILD_DIR=build
BINARY_NAME=server
MAIN_PATH=cmd/server/main.go

# Comandos de teste
test: ## Executa todos os testes
	$(GO_CMD) test ./... -v

test-unit: ## Executa apenas testes unitários
	$(GO_CMD) test -tags=test ./src/model/service -v

test-coverage: ## Executa testes com cobertura
	$(GO_CMD) test -tags=test ./... -coverprofile=coverage.out
	$(GO_CMD) tool cover -html=coverage.out -o coverage.html

test-watch: ## Executa testes em modo watch (requer entr)
	find . -name "*.go" | entr -c make test-unit

# Comandos de build
build: ## Compila a aplicação
	mkdir -p $(BUILD_DIR)
	$(GO_CMD) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-linux: ## Compila para Linux
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO_CMD) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)

# Comandos de execução
run: ## Executa a aplicação
	$(GO_CMD) run $(MAIN_PATH)

run-dev: ## Executa em modo desenvolvimento com hot reload (requer air)
	air

# Comandos de limpeza
clean: ## Remove arquivos de build
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Comandos de qualidade
lint: ## Executa linter
	golangci-lint run

fmt: ## Formata o código
	$(GO_CMD) fmt ./...

vet: ## Executa go vet
	$(GO_CMD) vet ./...

# Comandos de dependências
deps: ## Instala dependências
	$(GO_CMD) mod download
	$(GO_CMD) mod tidy

deps-update: ## Atualiza dependências
	$(GO_CMD) get -u ./...
	$(GO_CMD) mod tidy

# Comandos de Docker
docker-build: ## Constrói imagem Docker
	docker build -t back-project .

docker-run: ## Executa container Docker
	docker run -p 8080:8080 back-project

# Comandos de banco
db-migrate: ## Executa migrações do banco
	$(GO_CMD) run scripts/migrate.go

db-seed: ## Popula banco com dados de teste
	$(GO_CMD) run scripts/seed.go

# Help
help: ## Mostra esta ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Default target
.DEFAULT_GOAL := help