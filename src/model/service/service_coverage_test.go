package service

import (
	"errors"
	"testing"

	dtos "github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// =============================================================================
// TESTES PARA AUMENTAR COBERTURA - CENÁRIOS NÃO COBERTOS
// =============================================================================

func TestService_CreateUserService_DatabaseError(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)

	// User doesn't exist
	emptyUser := &entity.User{}
	mockDBMaster.On("GetUser", "joao@teste.com").Return(emptyUser)

	// Hash password succeeds
	mockCrypto.On("HashPassword", "senha123!").Return("hashedpassword123", nil)

	// Database creation fails
	mockDBMaster.On("CreateUser", mock.AnythingOfType("entity.User")).Return(errors.New("database connection failed"))

	service := &Service{
		crypto:   mockCrypto,
		dbmaster: mockDBMaster,
		dbClient: nil,
		redis:    nil,
	}

	request := dtos.CreateUser{
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@teste.com",
		NomeEmpresa: "Empresa Teste",
		Categoria:   "Tech",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "senha123!",
	}

	// Act
	user, err := service.CreateUserService(request)

	// Assert
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, 500, err.Code)

	mockCrypto.AssertExpectations(t)
	mockDBMaster.AssertExpectations(t)
}

func TestService_GetAllFornecedoresService_DatabaseError(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	mockDBClient.On("GetAllFornecedores", "1").Return(nil, errors.New("database error"))

	service := &Service{
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
	}

	// Act
	result, err := service.GetAllFornecedoresService("1")

	// Assert
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Error retrieving fornecedores", err.Message)
	assert.Equal(t, 500, err.Code)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestService_CreateFornecedorService_DatabaseError(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	tenant := &entity.Tenants{
		ID:          1,
		UserID:      1,
		NomeEmpresa: "Test Company",
		DBName:      "test_db",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(tenant)
	mockDBClient.On("CreateFornecedor", mock.AnythingOfType("entity.Fornecedores"), mock.AnythingOfType("*context.valueCtx")).Return(errors.New("database error"))

	service := &Service{
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
	}

	request := dtos.CreateFornecedorRequest{
		Nome:     "Novo Fornecedor",
		Email:    "novo@fornecedor.com",
		Telefone: "11999999999",
		Status:   "Ativo",
	}

	// Act
	success, err := service.CreateFornecedorService("1", request)

	// Assert
	assert.False(t, success)
	assert.NotNil(t, err)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, 500, err.Code)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestService_ChangeStatusFornecedorService_GetFornecedorError(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	tenant := &entity.Tenants{
		ID:          1,
		UserID:      1,
		NomeEmpresa: "Test Company",
		DBName:      "test_db",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(tenant)
	mockDBClient.On("GetFornecedorById", "1", mock.AnythingOfType("*context.valueCtx")).Return(nil, errors.New("fornecedor not found"))

	service := &Service{
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
	}

	// Act
	success, err := service.ChangeStatusFornecedorService("1", "1")

	// Assert
	assert.False(t, success)
	assert.NotNil(t, err)
	assert.Equal(t, "Error retrieving fornecedor", err.Message)
	assert.Equal(t, 500, err.Code)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestService_ChangeStatusFornecedorService_UpdateError(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	tenant := &entity.Tenants{
		ID:          1,
		UserID:      1,
		NomeEmpresa: "Test Company",
		DBName:      "test_db",
	}

	fornecedor := &entity.Fornecedores{
		ID:     1,
		Nome:   "Fornecedor Teste",
		Status: "Ativo",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(tenant)
	mockDBClient.On("GetFornecedorById", "1", mock.AnythingOfType("*context.valueCtx")).Return(fornecedor, nil)
	mockDBClient.On("UpdateFornecedor", mock.AnythingOfType("entity.Fornecedores"), mock.AnythingOfType("*context.valueCtx")).Return(errors.New("update failed"))

	service := &Service{
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
	}

	// Act
	success, err := service.ChangeStatusFornecedorService("1", "1")

	// Assert
	assert.False(t, success)
	assert.NotNil(t, err)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, 500, err.Code)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestService_UpdateFornecedorFieldService_DatabaseError(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	tenant := &entity.Tenants{
		ID:          1,
		UserID:      1,
		NomeEmpresa: "Test Company",
		DBName:      "test_db",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(tenant)
	mockDBClient.On("UpdateFornecedorField", "1", "nome", "Novo Nome", mock.AnythingOfType("*context.valueCtx")).Return(errors.New("update failed"))

	service := &Service{
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
	}

	// Act
	success, err := service.UpdateFornecedorFieldService("1", "1", "nome", "Novo Nome")

	// Assert
	assert.False(t, success)
	assert.NotNil(t, err)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, 500, err.Code)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestService_DeleteFornecedorService_DatabaseError(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	tenant := &entity.Tenants{
		ID:          1,
		UserID:      1,
		NomeEmpresa: "Test Company",
		DBName:      "test_db",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(tenant)
	mockDBClient.On("DeleteFornecedor", "1", mock.AnythingOfType("*context.valueCtx")).Return(errors.New("delete failed"))

	service := &Service{
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
	}

	// Act
	success, err := service.DeleteFornecedorService("1", "1")

	// Assert
	assert.False(t, success)
	assert.NotNil(t, err)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, 500, err.Code)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestService_GetAllProductsService_DatabaseError(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	tenant := &entity.Tenants{
		ID:          1,
		UserID:      1,
		NomeEmpresa: "Test Company",
		DBName:      "test_db",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(tenant)
	mockDBClient.On("GetAllProducts", mock.AnythingOfType("*context.valueCtx")).Return(nil, errors.New("database error"))

	service := &Service{
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
	}

	// Act
	result, err := service.GetAllProductsService("1")

	// Assert
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Error retrieving products", err.Message)
	assert.Equal(t, 500, err.Code)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestService_CreateProductService_DatabaseError(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	tenant := &entity.Tenants{
		ID:          1,
		UserID:      1,
		NomeEmpresa: "Test Company",
		DBName:      "test_db",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(tenant)
	mockDBClient.On("CreateProduct", mock.AnythingOfType("entity.Produto"), mock.AnythingOfType("*context.valueCtx")).Return(errors.New("database error"))

	service := &Service{
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
	}

	request := dtos.CreateProductRequest{
		NomeProduto: "Novo Produto",
		SKU:         "SKU001",
		PrecoVenda:  99.99,
		Status:      "Ativo",
	}

	// Act
	success, err := service.CreateProductService("1", request)

	// Assert
	assert.False(t, success)
	assert.NotNil(t, err)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, 500, err.Code)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestService_DeleteProductService_DatabaseError(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	tenant := &entity.Tenants{
		ID:          1,
		UserID:      1,
		NomeEmpresa: "Test Company",
		DBName:      "test_db",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(tenant)
	mockDBClient.On("DeleteProduct", "1", mock.AnythingOfType("*context.valueCtx")).Return(errors.New("delete failed"))

	service := &Service{
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
	}

	// Act
	success, err := service.DeleteProductService("1", "1")

	// Assert
	assert.False(t, success)
	assert.NotNil(t, err)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, 500, err.Code)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestService_GetAllPedidosService_DatabaseError(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	tenant := &entity.Tenants{
		ID:          1,
		UserID:      1,
		NomeEmpresa: "Test Company",
		DBName:      "test_db",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(tenant)
	mockDBClient.On("GetAllPedidos", mock.AnythingOfType("*context.valueCtx")).Return(nil, errors.New("database error"))

	service := &Service{
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
	}

	// Act
	result, err := service.GetAllPedidosService("1")

	// Assert
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Error retrieving pedidos", err.Message)
	assert.Equal(t, 500, err.Code)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

// =============================================================================
// TESTES DE EDGE CASES ADICIONAIS
// =============================================================================

func TestService_ChangeStatusFornecedorService_StatusToggle(t *testing.T) {
	tests := []struct {
		name           string
		currentStatus  string
		expectedStatus string
	}{
		{"Ativo to Inativo", "Ativo", "Inativo"},
		{"Inativo to Ativo", "Inativo", "Ativo"},
		{"Empty to Ativo", "", "Ativo"},
		{"Random to Ativo", "Pendente", "Ativo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDBMaster := new(MockDBMaster)
			mockDBClient := new(MockDBClient)

			tenant := &entity.Tenants{
				ID:          1,
				UserID:      1,
				NomeEmpresa: "Test Company",
				DBName:      "test_db",
			}

			fornecedor := &entity.Fornecedores{
				ID:     1,
				Nome:   "Fornecedor Teste",
				Status: tt.currentStatus,
			}

			mockDBMaster.On("GetTenantByUserID", uint(1)).Return(tenant)
			mockDBClient.On("GetFornecedorById", "1", mock.AnythingOfType("*context.valueCtx")).Return(fornecedor, nil)
			mockDBClient.On("UpdateFornecedor", mock.MatchedBy(func(f entity.Fornecedores) bool {
				return f.Status == tt.expectedStatus
			}), mock.AnythingOfType("*context.valueCtx")).Return(nil)

			service := &Service{
				dbmaster: mockDBMaster,
				dbClient: mockDBClient,
			}

			// Act
			success, err := service.ChangeStatusFornecedorService("1", "1")

			// Assert
			assert.True(t, success)
			assert.Nil(t, err)

			mockDBMaster.AssertExpectations(t)
			mockDBClient.AssertExpectations(t)
		})
	}
}

func TestService_NewServiceInstance(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)
	mockRedis := new(MockRedis)
	mockTokenGen := new(MockTokenGenerator)

	// Act
	service := NewServiceInstance(mockCrypto, mockDBMaster, mockDBClient, mockRedis, mockTokenGen)

	// Assert
	assert.NotNil(t, service)
	assert.IsType(t, &Service{}, service)

	// Verificar se as dependências foram injetadas corretamente
	serviceImpl := service.(*Service)
	assert.Equal(t, mockCrypto, serviceImpl.crypto)
	assert.Equal(t, mockDBMaster, serviceImpl.dbmaster)
	assert.Equal(t, mockDBClient, serviceImpl.dbClient)
	assert.Equal(t, mockRedis, serviceImpl.redis)
	assert.Equal(t, mockTokenGen, serviceImpl.tokenGen)
}
