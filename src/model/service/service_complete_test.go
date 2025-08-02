package service

import (
	"errors"
	"testing"

	"github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para PersistenceInterfaceDBClient
type MockDBClient struct {
	mock.Mock
}

func (m *MockDBClient) GetAllFornecedores(userID string) ([]entity.Fornecedores, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Fornecedores), args.Error(1)
}

func (m *MockDBClient) CreateFornecedor(fornecedor entity.Fornecedores, userID string) error {
	args := m.Called(fornecedor, userID)
	return args.Error(0)
}

func (m *MockDBClient) GetFornecedorById(id string, userID string) (*entity.Fornecedores, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Fornecedores), args.Error(1)
}

func (m *MockDBClient) UpdateFornecedor(fornecedor entity.Fornecedores, userID string) error {
	args := m.Called(fornecedor, userID)
	return args.Error(0)
}

func (m *MockDBClient) UpdateFornecedorField(id string, campo string, valor string, userID string) error {
	args := m.Called(id, campo, valor, userID)
	return args.Error(0)
}

func (m *MockDBClient) DeleteFornecedor(id string, userID string) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockDBClient) GetAllProducts(userID string) ([]entity.Produto, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Produto), args.Error(1)
}

func (m *MockDBClient) CreateProduct(product entity.Produto, userID string) error {
	args := m.Called(product, userID)
	return args.Error(0)
}

func (m *MockDBClient) DeleteProduct(id string, userID string) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockDBClient) GetAllPedidos(userID string) ([]entity.Pedido, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Pedido), args.Error(1)
}

// TESTES PARA CreateUserService
func TestService_CreateUserService_Success(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)

	request := dtos.CreateUser{
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@teste.com",
		NomeEmpresa: "Empresa Teste",
		Categoria:   "Tecnologia",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "senha123",
	}

	// User doesn't exist yet
	emptyUser := &entity.User{}
	mockDBMaster.On("GetUser", "joao@teste.com").Return(emptyUser)

	// Password hashing succeeds
	mockCrypto.On("HashPassword", "senha123").Return("hashed_senha123", nil)

	// User creation succeeds
	mockDBMaster.On("CreateUser", mock.AnythingOfType("entity.User")).Return(nil)

	service := &Service{
		crypto:   mockCrypto,
		dbmaster: mockDBMaster,
		dbClient: nil,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	user, err := service.CreateUserService(request)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "João", user.FirstName)
	assert.Equal(t, "Silva", user.LastName)
	assert.Equal(t, "joao@teste.com", user.Email)
	assert.Equal(t, "hashed_senha123", user.Password)

	mockCrypto.AssertExpectations(t)
	mockDBMaster.AssertExpectations(t)
}

func TestService_CreateUserService_EmailAlreadyExists(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)

	request := dtos.CreateUser{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@teste.com",
		Password:  "senha123",
	}

	// User already exists
	existingUser := &entity.User{
		ID:    1,
		Email: "joao@teste.com",
	}
	mockDBMaster.On("GetUser", "joao@teste.com").Return(existingUser)

	service := &Service{
		crypto:   mockCrypto,
		dbmaster: mockDBMaster,
		dbClient: nil,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	user, err := service.CreateUserService(request)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "Email is already associated with an existing account", err.Message)
	assert.Equal(t, 400, err.Code)

	mockDBMaster.AssertExpectations(t)
}

func TestService_CreateUserService_HashPasswordError(t *testing.T) {
	// Arrange
	mockCrypto := new(MockCrypto)
	mockDBMaster := new(MockDBMaster)

	request := dtos.CreateUser{
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@teste.com",
		Password:  "senha123",
	}

	// User doesn't exist
	emptyUser := &entity.User{}
	mockDBMaster.On("GetUser", "joao@teste.com").Return(emptyUser)

	// Password hashing fails
	mockCrypto.On("HashPassword", "senha123").Return("", errors.New("hash error"))

	service := &Service{
		crypto:   mockCrypto,
		dbmaster: mockDBMaster,
		dbClient: nil,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	user, err := service.CreateUserService(request)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, 500, err.Code)

	mockCrypto.AssertExpectations(t)
	mockDBMaster.AssertExpectations(t)
}

// TESTES PARA GetAllFornecedoresService
func TestService_GetAllFornecedoresService_Success(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	userID := "1"

	fornecedores := []entity.Fornecedores{
		{
			ID:       1,
			Nome:     "Fornecedor 1",
			Email:    "fornecedor1@teste.com",
			Telefone: "11999999999",
			Status:   "Ativo",
		},
		{
			ID:       2,
			Nome:     "Fornecedor 2",
			Email:    "fornecedor2@teste.com",
			Telefone: "11888888888",
			Status:   "Ativo",
		},
	}

	mockDBClient.On("GetAllFornecedores", "1").Return(fornecedores, nil)

	service := &Service{
		crypto:   nil,
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	result, err := service.GetAllFornecedoresService(userID)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 2, len(result.Fornecedores))
	assert.Equal(t, "Fornecedor 1", result.Fornecedores[0].Nome)
	assert.Equal(t, "Fornecedor 2", result.Fornecedores[1].Nome)

	mockDBClient.AssertExpectations(t)
}

func TestService_GetAllFornecedoresService_DBError(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	userID := "1"

	mockDBClient.On("GetAllFornecedores", userID).Return(nil, errors.New("database error"))

	service := &Service{
		crypto:   nil,
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	result, err := service.GetAllFornecedoresService(userID)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Error retrieving fornecedores", err.Message)
	assert.Equal(t, 500, err.Code)

	mockDBClient.AssertExpectations(t)
}

// TESTES PARA CreateFornecedorService
func TestService_CreateFornecedorService_Success(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	userID := "1"
	validTenant := &entity.Tenants{
		ID:     1,
		UserID: 1,
	}

	request := dtos.CreateFornecedorRequest{
		Nome:     "Novo Fornecedor",
		Email:    "novo@fornecedor.com",
		Telefone: "11999999999",
		Cidade:   "São Paulo",
		Estado:   "SP",
		Status:   "Ativo",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(validTenant)
	mockDBClient.On("CreateFornecedor", mock.AnythingOfType("entity.Fornecedores"), "1").Return(nil)

	service := &Service{
		crypto:   nil,
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	success, err := service.CreateFornecedorService(userID, request)

	// Assert
	assert.Nil(t, err)
	assert.True(t, success)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

// TESTES PARA ChangeStatusFornecedorService
func TestService_ChangeStatusFornecedorService_Success(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	userID := "1"
	fornecedorID := "1"
	validTenant := &entity.Tenants{
		ID:     1,
		UserID: 1,
	}

	fornecedor := &entity.Fornecedores{
		ID:     1,
		Nome:   "Fornecedor Teste",
		Status: "Ativo",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(validTenant)
	mockDBClient.On("GetFornecedorById", fornecedorID, "1").Return(fornecedor, nil)
	mockDBClient.On("UpdateFornecedor", mock.AnythingOfType("entity.Fornecedores"), "1").Return(nil)

	service := &Service{
		crypto:   nil,
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	success, err := service.ChangeStatusFornecedorService(userID, fornecedorID)

	// Assert
	assert.Nil(t, err)
	assert.True(t, success)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

// TESTES PARA UpdateFornecedorFieldService
func TestService_UpdateFornecedorFieldService_Success(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	userID := "1"
	fornecedorID := "1"
	campo := "nome"
	valor := "Novo Nome"
	validTenant := &entity.Tenants{
		ID:     1,
		UserID: 1,
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(validTenant)
	mockDBClient.On("UpdateFornecedorField", fornecedorID, campo, valor, "1").Return(nil)

	service := &Service{
		crypto:   nil,
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	success, err := service.UpdateFornecedorFieldService(userID, fornecedorID, campo, valor)

	// Assert
	assert.Nil(t, err)
	assert.True(t, success)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestService_UpdateFornecedorFieldService_InvalidField(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	userID := "1"
	fornecedorID := "1"
	campo := "campo_invalido"
	valor := "Valor"

	service := &Service{
		crypto:   nil,
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	success, err := service.UpdateFornecedorFieldService(userID, fornecedorID, campo, valor)

	// Assert
	assert.NotNil(t, err)
	assert.False(t, success)
	assert.Equal(t, "Invalid field to update", err.Message)
	assert.Equal(t, 400, err.Code)
}

// TESTES PARA DeleteFornecedorService
func TestService_DeleteFornecedorService_Success(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	userID := "1"
	fornecedorID := "1"
	validTenant := &entity.Tenants{
		ID:     1,
		UserID: 1,
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(validTenant)
	mockDBClient.On("DeleteFornecedor", fornecedorID, "1").Return(nil)

	service := &Service{
		crypto:   nil,
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	success, err := service.DeleteFornecedorService(userID, fornecedorID)

	// Assert
	assert.Nil(t, err)
	assert.True(t, success)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

// TESTES PARA GetAllProductsService
func TestService_GetAllProductsService_Success(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	userID := "1"
	validTenant := &entity.Tenants{
		ID:     1,
		UserID: 1,
	}

	products := []entity.Produto{
		{
			IDProduto:   1,
			NomeProduto: "Produto 1",
			SKU:         "SKU001",
			PrecoVenda:  10.50,
			Status:      "Ativo",
		},
		{
			IDProduto:   2,
			NomeProduto: "Produto 2",
			SKU:         "SKU002",
			PrecoVenda:  25.99,
			Status:      "Ativo",
		},
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(validTenant)
	mockDBClient.On("GetAllProducts", "1").Return(products, nil)

	service := &Service{
		crypto:   nil,
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	result, err := service.GetAllProductsService(userID)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 2, len(result.Products))
	assert.Equal(t, "Produto 1", result.Products[0].NomeProduto)
	assert.Equal(t, "Produto 2", result.Products[1].NomeProduto)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

// TESTES PARA CreateProductService
func TestService_CreateProductService_Success(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	userID := "1"
	validTenant := &entity.Tenants{
		ID:     1,
		UserID: 1,
	}

	request := dtos.CreateProductRequest{
		NomeProduto: "Novo Produto",
		SKU:         "SKU003",
		PrecoVenda:  15.99,
		Status:      "Ativo",
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(validTenant)
	mockDBClient.On("CreateProduct", mock.AnythingOfType("entity.Produto"), "1").Return(nil)

	service := &Service{
		crypto:   nil,
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	success, err := service.CreateProductService(userID, request)

	// Assert
	assert.Nil(t, err)
	assert.True(t, success)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

// TESTES PARA DeleteProductService
func TestService_DeleteProductService_Success(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	userID := "1"
	productID := "1"
	validTenant := &entity.Tenants{
		ID:     1,
		UserID: 1,
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(validTenant)
	mockDBClient.On("DeleteProduct", productID, "1").Return(nil)

	service := &Service{
		crypto:   nil,
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	success, err := service.DeleteProductService(userID, productID)

	// Assert
	assert.Nil(t, err)
	assert.True(t, success)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

// TESTES PARA GetAllPedidosService
func TestService_GetAllPedidosService_Success(t *testing.T) {
	// Arrange
	mockDBMaster := new(MockDBMaster)
	mockDBClient := new(MockDBClient)

	userID := "1"
	validTenant := &entity.Tenants{
		ID:     1,
		UserID: 1,
	}

	pedidos := []entity.Pedido{
		{
			IDPedido:     1,
			IDFornecedor: 1,
			DataPedido:   "2024-01-01",
			ValorTotal:   100.50,
			Status:       "Pendente",
		},
		{
			IDPedido:     2,
			IDFornecedor: 2,
			DataPedido:   "2024-01-02",
			ValorTotal:   250.99,
			Status:       "Concluído",
		},
	}

	mockDBMaster.On("GetTenantByUserID", uint(1)).Return(validTenant)
	mockDBClient.On("GetAllPedidos", "1").Return(pedidos, nil)

	service := &Service{
		crypto:   nil,
		dbmaster: mockDBMaster,
		dbClient: mockDBClient,
		redis:    nil,
		tokenGen: nil,
	}

	// Act
	result, err := service.GetAllPedidosService(userID)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 2, len(result.Pedidos))
	assert.Equal(t, 1, result.Pedidos[0].ID)
	assert.Equal(t, 2, result.Pedidos[1].ID)

	mockDBMaster.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}
