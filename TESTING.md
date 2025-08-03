# Guia de Testes

Este documento explica como executar os testes do projeto.

## Arquivos de Configuração para Testes

### `.env.test`
Contém variáveis de ambiente específicas para testes. Este arquivo é versionado no Git e contém apenas valores de teste (não sensíveis).

### Scripts de Teste

#### Windows
```bash
# Executar via script batch
test.bat

# Ou manualmente
cmd /c test.bat
```

#### Linux/Mac
```bash
# Dar permissão de execução
chmod +x test.sh

# Executar
./test.sh
```

## Executando Testes Específicos

### Testes de Configuração
```bash
# Todos os testes de configuração
go test -v ./cmd/config

# Testes específicos
go test -v ./cmd/config -run "TestConfig_Structure"
go test -v ./cmd/config -run "TestConfig_DefaultValues"
go test -v ./cmd/config -run "TestGetEnvWithDefault"
```

### Testes com Cobertura
```bash
# Gerar relatório de cobertura
go test -cover ./cmd/config

# Relatório detalhado em HTML
go test -coverprofile=coverage.out ./cmd/config
go tool cover -html=coverage.out

# Relatório de cobertura por função
go tool cover -func coverage.out
```

## Cobertura Atual: 91.1%

### Cobertura por Função:
- `NewConfig()` - **100.0%**
- `getEnvWithDefault()` - **100.0%**
- `getEnvOrFail()` - **50.0%** (limitado por log.Fatalf)
- `init()` - **85.7%** (limitado por tratamento de erros)
- `NewDatabaseConnection()` - **100.0%**
- `ConnectionDBClients()` - **100.0%**

## Estrutura dos Testes

### Funções Principais Testadas

1. **`testGetEnvWithDefault`**
   - Testa a função auxiliar que busca variáveis de ambiente com fallback
   - Cenários: variável existe, não existe, está vazia

2. **`TestConfig_Structure`**
   - Verifica se a struct Config armazena valores corretamente
   - Testa todos os campos da configuração

3. **`TestConfig_DefaultValues`**
   - Verifica se valores padrão são aplicados quando variáveis de ambiente não existem
   - Testa o comportamento de fallback

### Outros Testes Incluídos

- Validação de configuração de banco de dados
- Validação de configuração JWT
- Validação de porta do servidor web
- Testes com caracteres especiais e Unicode
- Testes de edge cases (valores muito longos, vazios, etc.)

## Variáveis de Ambiente Necessárias

### Obrigatórias
- `DB_HOST`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `JWT_SECRET`

### Opcionais (com valores padrão)
- `DB_DRIVER` (padrão: mysql)
- `DB_PORT` (padrão: 3306)
- `WEB_SERVER_PORT` (padrão: 8080)
- `JWT_EXPIRES_IN` (padrão: 30)
- `CORS_ORIGINS` (padrão: localhost URLs)
- `REDIS_ADDR` (padrão: localhost:6379)

## Troubleshooting

### Erro: "Required environment variable X is not set"
Certifique-se de que as variáveis obrigatórias estão definidas no `.env.test` ou use os scripts fornecidos.

### Aviso: "Error loading .env file"
Este aviso é normal em ambiente de teste quando não há arquivo `.env` na raiz do projeto.