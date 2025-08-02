package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/betine97/back-project.git/cmd/config/exceptions"
	"github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks melhorados usando as interfaces
type MockCryptoImproved struct {
	mock.Mock
}

func (m *MockCryptoImproved) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockCryptoImproved) CheckPassword(password, hashedPassword string) (bool, error) {
	args := m.Called(password, hashedPassword)
	return args.Bool(0), args.Error(1)
}

type MockDBMasterImproved struct {
	mock.Mock
}

func (m *MockDBMasterImproved) CreateUser(user entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDBMasterImproved) VerifyExist(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockDBMasterImproved) GetUser(email string) *entity.User {
	args := m.Called(email)
	return args.Get(0).(*entity.User)
}

func (m *MockDBMasterImproved) GetTenantByUserID(userID uint) *entity.Tenants {
	args := m.Called(userID)
	return args.Get(0).(*entity.Tenants)
}

type MockRedisImproved struct {
	mock.Mock
}

func (m *MockRedisImproved) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	cmd := redis.NewStatusCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	}
	return cmd
}

func (m *MockRedisImproved) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx, "get", key)
	if args.Error(1) != nil {
		cmd.SetErr(args.Error(1))
	} else {
		cmd.SetVal(args.String(0))
	}
	return cmd
}

func (m *MockRedisImproved) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	cmd := redis.NewStatusCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	}
	return cmd
}

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateToken(userID uint) (string, *exceptions.RestErr) {
	args := m.Called(userID)
	if args.Get(1) == nil {
		return args.String(0), nil
	}
	return args.String(0), args.Get(1).(*exceptions.RestErr)
}

func TestService_LoginUserService_Complete(t *testing.T) {
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
		setupMocks    func(*MockCryptoImproved, *MockDBMasterImproved, *MockRedisImproved, *MockTokenGenerator)
		expectedError *exceptions.RestErr
		expectError   bool
	}{
		{
			name: "Success - Complete flow with token generation",
			request: dtos.UserLogin{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func(crypto *MockCryptoImproved, dbMaster *MockDBMasterImproved, redisClient *MockRedisImproved, tokenGen *MockTokenGenerator) {
				// User exists
				dbMaster.On("GetUser", "john@example.com").Return(validUser)

				// Password is correct
				crypto.On("CheckPassword", "password123", "hashedpassword123").Return(true, nil)

				// Redis is reachable
				redisClient.On("Ping", mock.Anything).Return(nil)

				// Tenant found in cache
				redisClient.On("Get", mock.Anything, "user:1:db_info").Return(string(tenantJSON), nil)

				// Token generation successful
				tokenGen.On("GenerateToken", uint(1)).Return("valid_jwt_token", nil)
			},
			expectError: false,
		},
		{
			name: "Error - Token generation fails",
			request: dtos.UserLogin{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func(crypto *MockCryptoImproved, dbMaster *MockDBMasterImproved, redisClient *MockRedisImproved, tokenGen *MockTokenGenerator) {
				// User exists
				dbMaster.On("GetUser", "john@example.com").Return(validUser)

				// Password is correct
				crypto.On("CheckPassword", "password123", "hashedpassword123").Return(true, nil)

				// Redis is reachable
				redisClient.On("Ping", mock.Anything).Return(nil)

				// Tenant found in cache
				redisClient.On("Get", mock.Anything, "user:1:db_info").Return(string(tenantJSON), nil)

				// Token generation fails
				tokenGen.On("GenerateToken", uint(1)).Return("", exceptions.NewInternalServerError("Token generation failed"))
			},
			expectedError: exceptions.NewInternalServerError("Token generation failed"),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - Setup mocks
			mockCrypto := new(MockCryptoImproved)
			mockDBMaster := new(MockDBMasterImproved)
			mockRedis := new(MockRedisImproved)
			mockTokenGen := new(MockTokenGenerator)

			// Setup specific mocks for this test case
			tt.setupMocks(mockCrypto, mockDBMaster, mockRedis, mockTokenGen)

			// Create service instance with mocks
			service := &Service{
				crypto:   mockCrypto,
				dbmaster: mockDBMaster,
				dbClient: nil,
				redis:    mockRedis,
				tokenGen: mockTokenGen,
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
				assert.NotEmpty(t, token)
			}

			// Verify all mocks were called as expected
			mockCrypto.AssertExpectations(t)
			mockDBMaster.AssertExpectations(t)
			mockRedis.AssertExpectations(t)
			mockTokenGen.AssertExpectations(t)
		})
	}
}
