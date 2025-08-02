package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/betine97/back-project.git/cmd/config/exceptions"
	"github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	redis "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para CryptoInterface
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

// Mock para PersistenceInterfaceDBMaster
type MockDBMaster struct {
	mock.Mock
}

func (m *MockDBMaster) CreateUser(user entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDBMaster) VerifyExist(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockDBMaster) GetUser(email string) *entity.User {
	args := m.Called(email)
	return args.Get(0).(*entity.User)
}

func (m *MockDBMaster) GetTenantByUserID(userID uint) *entity.Tenants {
	args := m.Called(userID)
	return args.Get(0).(*entity.Tenants)
}

// Teste para casos de erro básicos do LoginUserService
func TestService_LoginUserService_UserNotFound(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)

	// User doesn't exist - retorna user vazio
	emptyUser := &entity.User{}
	mockDBMaster.On("GetUser", "nonexistent@example.com").Return(emptyUser)

	service := &Service{
		crypto:   mockCrypto,
		dbmaster: mockDBMaster,
		dbClient: nil,
		redis:    nil,
	}

	request := dtos.UserLogin{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	// Act
	token, err := service.LoginUserService(request)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "Account not found", err.Message)
	assert.Equal(t, 404, err.Code)
	assert.Empty(t, token)

	mockDBMaster.AssertExpectations(t)
}

func TestService_LoginUserService_IncorrectPassword(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)

	validUser := &entity.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword123",
	}

	mockDBMaster.On("GetUser", "john@example.com").Return(validUser)
	mockCrypto.On("CheckPassword", "wrongpassword", "hashedpassword123").Return(false, errors.New("password mismatch"))

	service := &Service{
		crypto:   mockCrypto,
		dbmaster: mockDBMaster,
		dbClient: nil,
		redis:    nil,
	}

	request := dtos.UserLogin{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}

	// Act
	token, err := service.LoginUserService(request)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "The password entered is incorrect", err.Message)
	assert.Equal(t, 401, err.Code)
	assert.Empty(t, token)

	mockCrypto.AssertExpectations(t)
	mockDBMaster.AssertExpectations(t)
}

// Mock para RedisInterface
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

func (m *MockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	cmd := redis.NewStatusCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	}
	return cmd
}

// Testes de sucesso para LoginUserService
func TestService_LoginUserService_Success_TenantInCache(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)
	mockRedis := new(MockRedis)
	mockTokenGen := new(MockTokenGenerator)

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

	// Setup mocks
	mockDBMaster.On("GetUser", "john@example.com").Return(validUser)
	mockCrypto.On("CheckPassword", "password123", "hashedpassword123").Return(true, nil)
	mockRedis.On("Ping", mock.Anything).Return(nil)
	mockRedis.On("Get", mock.Anything, "user:1:db_info").Return(string(tenantJSON), nil)
	mockTokenGen.On("GenerateToken", uint(1)).Return("valid_jwt_token_123", nil)

	service := &Service{
		crypto:   mockCrypto,
		dbmaster: mockDBMaster,
		dbClient: nil,
		redis:    mockRedis,
		tokenGen: mockTokenGen,
	}

	request := dtos.UserLogin{
		Email:    "john@example.com",
		Password: "password123",
	}

	// Act
	token, err := service.LoginUserService(request)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "valid_jwt_token_123", token)
	assert.NotEmpty(t, token)

	// Verify all mocks were called
	mockCrypto.AssertExpectations(t)
	mockDBMaster.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}

func TestService_LoginUserService_Success_TenantFromDatabase(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)
	mockRedis := new(MockRedis)
	mockTokenGen := new(MockTokenGenerator)

	validUser := &entity.User{
		ID:        2,
		FirstName: "Maria",
		LastName:  "Silva",
		Email:     "maria@example.com",
		Password:  "hashedpassword456",
	}

	validTenant := &entity.Tenants{
		ID:          2,
		UserID:      2,
		NomeEmpresa: "Maria's Company",
		DBName:      "maria_db",
		DBUser:      "maria_user",
		DBPassword:  "maria_pass",
		DBHost:      "localhost",
		DBPort:      "3306",
		CreatedAt:   "2024-01-02",
	}

	// Setup mocks - tenant not in cache, found in database
	mockDBMaster.On("GetUser", "maria@example.com").Return(validUser)
	mockCrypto.On("CheckPassword", "password456", "hashedpassword456").Return(true, nil)
	mockRedis.On("Ping", mock.Anything).Return(nil)
	mockRedis.On("Get", mock.Anything, "user:2:db_info").Return("", redis.Nil)
	mockDBMaster.On("GetTenantByUserID", uint(2)).Return(validTenant)
	mockRedis.On("Set", mock.Anything, "user:2:db_info", mock.Anything, 0*time.Second).Return(nil)
	mockTokenGen.On("GenerateToken", uint(2)).Return("valid_jwt_token_456", nil)

	service := &Service{
		crypto:   mockCrypto,
		dbmaster: mockDBMaster,
		dbClient: nil,
		redis:    mockRedis,
		tokenGen: mockTokenGen,
	}

	request := dtos.UserLogin{
		Email:    "maria@example.com",
		Password: "password456",
	}

	// Act
	token, err := service.LoginUserService(request)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "valid_jwt_token_456", token)
	assert.NotEmpty(t, token)

	// Verify all mocks were called
	mockCrypto.AssertExpectations(t)
	mockDBMaster.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}

func TestService_LoginUserService_Error_RedisNotReachable(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)
	mockRedis := new(MockRedis)
	mockTokenGen := new(MockTokenGenerator)

	validUser := &entity.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword123",
	}

	// Setup mocks - Redis not reachable
	mockDBMaster.On("GetUser", "john@example.com").Return(validUser)
	mockCrypto.On("CheckPassword", "password123", "hashedpassword123").Return(true, nil)
	mockRedis.On("Ping", mock.Anything).Return(errors.New("redis connection failed"))

	service := &Service{
		crypto:   mockCrypto,
		dbmaster: mockDBMaster,
		dbClient: nil,
		redis:    mockRedis,
		tokenGen: mockTokenGen,
	}

	request := dtos.UserLogin{
		Email:    "john@example.com",
		Password: "password123",
	}

	// Act
	token, err := service.LoginUserService(request)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "Redis is not reachable", err.Message)
	assert.Equal(t, 500, err.Code)
	assert.Empty(t, token)

	// Verify mocks were called
	mockCrypto.AssertExpectations(t)
	mockDBMaster.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

func TestService_LoginUserService_Error_TenantNotFound(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)
	mockRedis := new(MockRedis)
	mockTokenGen := new(MockTokenGenerator)

	validUser := &entity.User{
		ID:        3,
		FirstName: "Carlos",
		LastName:  "Santos",
		Email:     "carlos@example.com",
		Password:  "hashedpassword789",
	}

	emptyTenant := &entity.Tenants{} // Tenant não encontrado

	// Setup mocks - tenant not found
	mockDBMaster.On("GetUser", "carlos@example.com").Return(validUser)
	mockCrypto.On("CheckPassword", "password789", "hashedpassword789").Return(true, nil)
	mockRedis.On("Ping", mock.Anything).Return(nil)
	mockRedis.On("Get", mock.Anything, "user:3:db_info").Return("", redis.Nil)
	mockDBMaster.On("GetTenantByUserID", uint(3)).Return(emptyTenant)

	service := &Service{
		crypto:   mockCrypto,
		dbmaster: mockDBMaster,
		dbClient: nil,
		redis:    mockRedis,
		tokenGen: mockTokenGen,
	}

	request := dtos.UserLogin{
		Email:    "carlos@example.com",
		Password: "password789",
	}

	// Act
	token, err := service.LoginUserService(request)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "Tenant not found", err.Message)
	assert.Equal(t, 500, err.Code)
	assert.Empty(t, token)

	// Verify mocks were called
	mockCrypto.AssertExpectations(t)
	mockDBMaster.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

func TestService_LoginUserService_Error_TokenGenerationFails(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)
	mockRedis := new(MockRedis)
	mockTokenGen := new(MockTokenGenerator)

	validUser := &entity.User{
		ID:        4,
		FirstName: "Ana",
		LastName:  "Costa",
		Email:     "ana@example.com",
		Password:  "hashedpassword000",
	}

	validTenant := &entity.Tenants{
		ID:          4,
		UserID:      4,
		NomeEmpresa: "Ana's Company",
		DBName:      "ana_db",
		DBUser:      "ana_user",
		DBPassword:  "ana_pass",
		DBHost:      "localhost",
		DBPort:      "3306",
		CreatedAt:   "2024-01-04",
	}

	tenantJSON, _ := json.Marshal(validTenant)

	// Setup mocks - token generation fails
	mockDBMaster.On("GetUser", "ana@example.com").Return(validUser)
	mockCrypto.On("CheckPassword", "password000", "hashedpassword000").Return(true, nil)
	mockRedis.On("Ping", mock.Anything).Return(nil)
	mockRedis.On("Get", mock.Anything, "user:4:db_info").Return(string(tenantJSON), nil)
	mockTokenGen.On("GenerateToken", uint(4)).Return("", exceptions.NewInternalServerError("Token generation failed"))

	service := &Service{
		crypto:   mockCrypto,
		dbmaster: mockDBMaster,
		dbClient: nil,
		redis:    mockRedis,
		tokenGen: mockTokenGen,
	}

	request := dtos.UserLogin{
		Email:    "ana@example.com",
		Password: "password000",
	}

	// Act
	token, err := service.LoginUserService(request)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "Token generation failed", err.Message)
	assert.Equal(t, 500, err.Code)
	assert.Empty(t, token)

	// Verify mocks were called
	mockCrypto.AssertExpectations(t)
	mockDBMaster.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}

// Teste simples para verificar se a função buildUserEntity funciona corretamente
func TestBuildUserEntity(t *testing.T) {
	// Arrange
	request := dtos.CreateUser{
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john@example.com",
		NomeEmpresa: "Test Company",
		Categoria:   "Tech",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "password123",
	}
	hashedPassword := "hashedpassword123"

	// Act
	user := buildUserEntity(request, hashedPassword)

	// Assert
	assert.Equal(t, request.FirstName, user.FirstName)
	assert.Equal(t, request.LastName, user.LastName)
	assert.Equal(t, request.Email, user.Email)
	assert.Equal(t, request.NomeEmpresa, user.NomeEmpresa)
	assert.Equal(t, request.Categoria, user.Categoria)
	assert.Equal(t, request.Segmento, user.Segmento)
	assert.Equal(t, request.City, user.City)
	assert.Equal(t, request.State, user.State)
	assert.Equal(t, hashedPassword, user.Password)
}
