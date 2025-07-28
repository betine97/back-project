# Design Document

## Overview

Este documento define a arquitetura e padrões de design para o projeto back-project em Go. O sistema segue uma arquitetura MVC (Model-View-Controller) com injeção de dependências, utilizando Go Fiber como framework web e GORM para persistência de dados.

## Architecture

### Estrutura de Camadas

```
┌─────────────────────────────────────┐
│            HTTP Layer               │
│         (Go Fiber Routes)           │
└─────────────────────────────────────┘
                    │
┌─────────────────────────────────────┐
│          Controller Layer           │
│     (src/controller/*.go)           │
└─────────────────────────────────────┘
                    │
┌─────────────────────────────────────┐
│           Service Layer             │
│    (src/model/service/*.go)         │
└─────────────────────────────────────┘
                    │
┌─────────────────────────────────────┐
│         Persistence Layer           │
│   (src/model/persistence/*.go)      │
└─────────────────────────────────────┘
                    │
┌─────────────────────────────────────┐
│           Database Layer            │
│        (SQLite/MySQL)               │
└─────────────────────────────────────┘
```

### Padrão de Injeção de Dependências

O sistema utiliza injeção de dependências manual através da função `initDependencies()` em `main.go`:

```go
func initDependencies(database *gorm.DB) controller.ControllerInterface {
    cryptoService := &crypto.Crypto{}
    persistence := persistence.NewDBConnection(database)
    service := service.NewServiceInstance(cryptoService, persistence)
    return controller.NewControllerInstance(service)
}
```

## Components and Interfaces

### Controller Layer
- **Responsabilidade**: Receber requisições HTTP, validar entrada, chamar serviços e retornar respostas
- **Padrão**: Implementar `ControllerInterface`
- **Localização**: `src/controller/`
- **Exemplo**: `controller.go` com métodos como `CreateUser`, `LoginUser`

### Service Layer
- **Responsabilidade**: Lógica de negócio, validações, transformações de dados
- **Padrão**: Implementar interfaces específicas para cada domínio
- **Localização**: `src/model/service/`
- **Dependências**: Crypto services, Persistence layer

### Model Layer
- **Entidades**: Estruturas de dados do domínio (`src/model/entity/`)
- **DTOs**: Objetos de transferência de dados (`src/model/dtos/`)
- **Persistência**: Acesso a dados (`src/model/persistence/`)

### Configuration Layer
- **Localização**: `cmd/config/`
- **Responsabilidade**: Configuração de banco, variáveis de ambiente
- **Padrão**: Utilizar `godotenv` para carregar configurações

## Data Models

### Padrão de Entidades

```go
type CreateUser struct {
    ID         string    `json:"id" gorm:"primaryKey"`
    First_Name string    `json:"first_name" validate:"required"`
    Last_Name  string    `json:"last_name" validate:"required"`
    Email      string    `json:"email" validate:"required,email"`
    CepBR      string    `json:"cep_br"`
    Country    string    `json:"country"`
    City       string    `json:"city"`
    Address    string    `json:"address"`
    Password   string    `json:"password" validate:"required,min=6"`
    CreateAt   time.Time `json:"create_at"`
}
```

### Padrão de DTOs

```go
type CreateUser struct {
    First_Name string `json:"first_name" validate:"required,min=2,max=50"`
    Last_Name  string `json:"last_name" validate:"required,min=2,max=50"`
    Email      string `json:"email" validate:"required,email"`
    CepBR      string `json:"cep_br" validate:"omitempty,len=8"`
    Country    string `json:"country" validate:"required,min=2,max=50"`
    City       string `json:"city" validate:"required,min=2,max=100"`
    Address    string `json:"address" validate:"required,min=5,max=200"`
    Password   string `json:"password" validate:"required,min=6,max=100"`
}
```

## Error Handling

### Padrão de Erros Estruturados

```go
type RestErr struct {
    Message string   `json:"message"`
    Err     string   `json:"error"`
    Code    int      `json:"code"`
    Causes  []Causes `json:"causes,omitempty"`
}
```

### Tipos de Erro Padronizados
- `NewBadRequestError()` - Erro 400
- `NewUnauthorizedRequestError()` - Erro 401
- `NewBadRequestValidationError()` - Erro 400 com detalhes de validação

## Testing Strategy

### Estrutura de Testes
- **Unit Tests**: Para cada service e controller
- **Integration Tests**: Para endpoints completos
- **Localização**: Arquivos `*_test.go` ao lado do código testado

### Padrão de Testes
```go
func TestServiceMethod(t *testing.T) {
    // Arrange
    mockPersistence := &MockPersistence{}
    service := NewServiceInstance(mockCrypto, mockPersistence)
    
    // Act
    result, err := service.Method(input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

## Security Patterns

### Autenticação JWT
- Utilizar `golang-jwt/jwt/v4`
- Secret configurável via `JWT_SECRET`
- Expiração configurável via `JWT_EXPIRESIN`

### Hash de Senhas
- Utilizar `golang.org/x/crypto/bcrypt`
- Implementar através do `crypto.Crypto` service

### Validação de Entrada
- Utilizar `go-playground/validator/v10`
- Validações em DTOs e entidades
- Retornar erros estruturados para validações

## Configuration Management

### Variáveis de Ambiente
```env
DB_DRIVER=sqlite|mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=bff-golang
WEB_SERVER_PORT=8080
JWT_SECRET=secret
JWT_EXPIRESIN=30
```

### Padrão de Configuração
- Carregar via `godotenv`
- Validar configurações obrigatórias
- Suportar diferentes ambientes (dev/test/prod)

## Logging Strategy

### Padrão de Logging
- Utilizar `go.uber.org/zap`
- Logs estruturados em JSON
- Níveis: Debug, Info, Warn, Error, Fatal

### Exemplo de Uso
```go
zap.L().Info("User created successfully", 
    zap.String("user_id", user.ID),
    zap.String("email", user.Email))
```

## Database Patterns

### GORM Configuration
- Suporte a SQLite (desenvolvimento) e MySQL (produção)
- Auto-migration para desenvolvimento
- Connection pooling configurado

### Repository Pattern
```go
type PersistenceInterface interface {
    CreateUser(user entity.CreateUser) (*entity.CreateUser, error)
    FindUserByEmail(email string) (*entity.CreateUser, error)
}
```