# Exemplo de Teste para LoginUserService

## Estrutura do Teste com Table-Driven Tests

Aqui está um exemplo completo de como testar a função `LoginUserService` usando table-driven tests e mocks:

```go
package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/betine97/back-project.git/cmd/config/exceptions"
	"github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks para as dependências
type MockCrypto struct {
	mock.Mock
}

func (m *MockCrypto) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockCrypto) CheckPassword(password, hashedPassword string) (bool, error) {
	args := m.Called(password, hashedPassword)
	return args.Bool(0), args.Error(1)
}

type MockDBMaster struct {
	mock.Mock
}

func (m *MockDBMaster) GetUser(email string) *entity.User {
	args := m.Called(email)
	return args.Get(0).(*entity.User)
}

func (m *MockDBMaster) GetTenantByUserID(userID uint) *entity.Tenants {
	args := m.Called(userID)
	return args.Get(0).(*entity.Tenants)
}

// Interface para Redis (para facilitar mocking)
type RedisInterface interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration interface{}) *redis.StatusCmd
}

type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	cmd := redis.NewStatusCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	}
	return cmd
}

func (m *MockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx, "get", key)
	if args.Error(1) != nil {
		cmd.SetErr(args.Error(1))
	} else {
		cmd.SetVal(args.String(0))
	}
	return cmd
}

func (m *MockRedis) Set(ctx context.Context, key string, value interface{}, expiration interface{}) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	cmd := redis.NewStatusCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	}
	return cmd
}

func TestService_LoginUserService(t *testing.T) {
	// Dados de teste comuns
	validUser := &entity.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword123",
	}

	validTenant := &entity.Tenants{
		ID:          1,
		UserID:      1,
		NomeEmpresa: "Test Company",
		DBName:      "test_db",
		DBUser:      "test_user",
		DBPassword:  "test_pass",
		DBHost:      "localhost",
		DBPort:      "3306",
		CreatedAt:   "2024-01-01",
	}

	tenantJSON, _ := json.Marshal(validTenant)

	tests := []struct {
		name          string
		request       dtos.UserLogin
		setupMocks    func(*MockCrypto, *MockDBMaster, *MockRedis)
		expectedError *exceptions.RestErr
		expectError   bool
	}{
		{
			name: "Error - User not found",
			request: dtos.UserLogin{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMocks: func(crypto *MockCrypto, dbMaster *MockDBMaster, redisClient *MockRedis) {
				// User doesn't exist - retorna user vazio
				emptyUser := &entity.User{}
				dbMaster.On("GetUser", "nonexistent@example.com").Return(emptyUser)
			},
			expectedError: exceptions.NewNotFoundError("Account not found"),
			expectError:   true,
		},
		{
			name: "Error - Incorrect password",
			request: dtos.UserLogin{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			setupMocks: func(crypto *MockCrypto, dbMaster *MockDBMaster, redisClient *MockRedis) {
				// User exists
				dbMaster.On("GetUser", "john@example.com").Return(validUser)
				
				// Password is incorrect
				crypto.On("CheckPassword", "wrongpassword", "hashedpassword123").Return(false, errors.New("password mismatch"))
			},
			expectedError: exceptions.NewUnauthorizedRequestError("The password entered is incorrect"),
			expectError:   true,
		},
		{
			name: "Error - Redis not reachable",
			request: dtos.UserLogin{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func(crypto *MockCrypto, dbMaster *MockDBMaster, redisClient *MockRedis) {
				// User exists
				dbMaster.On("GetUser", "john@example.com").Return(validUser)
				
				// Password is correct
				crypto.On("CheckPassword", "password123", "hashedpassword123").Return(true, nil)
				
				// Redis is not reachable
				redisClient.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(errors.New("redis connection failed"))
			},
			expectedError: exceptions.NewInternalServerError("Redis is not reachable"),
			expectError:   true,
		},
		{
			name: "Success - Tenant found in cache",
			request: dtos.UserLogin{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func(crypto *MockCrypto, dbMaster *MockDBMaster, redisClient *MockRedis) {
				// User exists
				dbMaster.On("GetUser", "john@example.com").Return(validUser)
				
				// Password is correct
				crypto.On("CheckPassword", "password123", "hashedpassword123").Return(true, nil)
				
				// Redis is reachable
				redisClient.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
				
				// Tenant found in cache
				redisClient.On("Get", mock.AnythingOfType("*context.emptyCtx"), "user:1:db_info").Return(string(tenantJSON), nil)
			},
			expectError: false,
		},
		{
			name: "Success - Tenant not in cache, found in database",
			request: dtos.UserLogin{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func(crypto *MockCrypto, dbMaster *MockDBMaster, redisClient *MockRedis) {
				// User exists
				dbMaster.On("GetUser", "john@example.com").Return(validUser)
				
				// Password is correct
				crypto.On("CheckPassword", "password123", "hashedpassword123").Return(true, nil)
				
				// Redis is reachable
				redisClient.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
				
				// Tenant not found in cache
				redisClient.On("Get", mock.AnythingOfType("*context.emptyCtx"), "user:1:db_info").Return("", redis.Nil)
				
				// Tenant found in database
				dbMaster.On("GetTenantByUserID", uint(1)).Return(validTenant)
				
				// Successfully store in Redis
				redisClient.On("Set", mock.AnythingOfType("*context.emptyCtx"), "user:1:db_info", string(tenantJSON), 0).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Error - Tenant not found in database",
			request: dtos.UserLogin{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func(crypto *MockCrypto, dbMaster *MockDBMaster, redisClient *MockRedis) {
				// User exists
				dbMaster.On("GetUser", "john@example.com").Return(validUser)
				
				// Password is correct
				crypto.On("CheckPassword", "password123", "hashedpassword123").Return(true, nil)
				
				// Redis is reachable
				redisClient.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
				
				// Tenant not found in cache
				redisClient.On("Get", mock.AnythingOfType("*context.emptyCtx"), "user:1:db_info").Return("", redis.Nil)
				
				// Tenant not found in database
				emptyTenant := &entity.Tenants{}
				dbMaster.On("GetTenantByUserID", uint(1)).Return(emptyTenant)
			},
			expectedError: exceptions.NewInternalServerError("Tenant not found"),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - Setup mocks
			mockCrypto := new(MockCrypto)
			mockDBMaster := new(MockDBMaster)
			mockRedis := new(MockRedis)

			// Setup specific mocks for this test case
			tt.setupMocks(mockCrypto, mockDBMaster, mockRedis)

			// Create service instance with mocks
			service := &Service{
				crypto:   mockCrypto,
				dbmaster: mockDBMaster,
				dbClient: nil, // Não usado nestes testes
				redis:    mockRedis,
			}

			// Act - Call the function
			token, err := service.LoginUserService(tt.request)

			// Assert - Verify results
			if tt.expectError {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError.Message, err.Message)
				assert.Equal(t, tt.expectedError.Code, err.Code)
				assert.Empty(t, token)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, token) // Para casos de sucesso, apenas verifica se o token não está vazio
			}

			// Verify all mocks were called as expected
			mockCrypto.AssertExpectations(t)
			mockDBMaster.AssertExpectations(t)
			mockRedis.AssertExpectations(t)
		})
	}
}
```

## Vantagens desta Abordagem

1. **Cobertura Completa**: Testa todos os cenários possíveis da função
2. **Isolamento**: Cada teste é independente e não afeta os outros
3. **Mocks Controlados**: Você controla exatamente o que cada dependência retorna
4. **Fácil Manutenção**: Adicionar novos cenários é simples
5. **Documentação**: Os testes servem como documentação do comportamento esperado

## Cenários Testados

- ✅ Usuário não encontrado
- ✅ Senha incorreta  
- ✅ Redis não acessível
- ✅ Tenant encontrado no cache
- ✅ Tenant não no cache, mas encontrado no banco
- ✅ Tenant não encontrado no banco

## Próximos Passos

Para resolver o problema do config que está impedindo os testes de rodar:

1. **Criar uma interface para o Redis** no service
2. **Modificar o service** para aceitar interfaces ao invés de structs concretas
3. **Criar um arquivo de configuração específico para testes**
4. **Usar build tags** para separar código de teste do código de produção

Isso permitirá que você execute os testes sem depender de configurações externas.