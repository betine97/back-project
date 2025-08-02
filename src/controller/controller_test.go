package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/betine97/back-project.git/cmd/config/exceptions"
	"github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para ServiceInterface
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateUserService(request dtos.CreateUser) (*entity.User, *exceptions.RestErr) {
	args := m.Called(request)
	if args.Get(1) == nil {
		return args.Get(0).(*entity.User), nil
	}
	return nil, args.Get(1).(*exceptions.RestErr)
}

func (m *MockService) LoginUserService(request dtos.UserLogin) (string, *exceptions.RestErr) {
	args := m.Called(request)
	if args.Get(1) == nil {
		return args.String(0), nil
	}
	return "", args.Get(1).(*exceptions.RestErr)
}

func (m *MockService) GetAllFornecedoresService(userID string) (*dtos.FornecedorListResponse, *exceptions.RestErr) {
	args := m.Called(userID)
	if args.Get(1) == nil {
		return args.Get(0).(*dtos.FornecedorListResponse), nil
	}
	return nil, args.Get(1).(*exceptions.RestErr)
}

func (m *MockService) CreateFornecedorService(userID string, request dtos.CreateFornecedorRequest) (bool, *exceptions.RestErr) {
	args := m.Called(userID, request)
	if args.Get(1) == nil {
		return args.Bool(0), nil
	}
	return false, args.Get(1).(*exceptions.RestErr)
}

func (m *MockService) ChangeStatusFornecedorService(userID string, id string) (bool, *exceptions.RestErr) {
	args := m.Called(userID, id)
	if args.Get(1) == nil {
		return args.Bool(0), nil
	}
	return false, args.Get(1).(*exceptions.RestErr)
}

func (m *MockService) UpdateFornecedorFieldService(userID string, id string, campo string, valor string) (bool, *exceptions.RestErr) {
	args := m.Called(userID, id, campo, valor)
	if args.Get(1) == nil {
		return args.Bool(0), nil
	}
	return false, args.Get(1).(*exceptions.RestErr)
}

func (m *MockService) DeleteFornecedorService(userID string, id string) (bool, *exceptions.RestErr) {
	args := m.Called(userID, id)
	if args.Get(1) == nil {
		return args.Bool(0), nil
	}
	return false, args.Get(1).(*exceptions.RestErr)
}

func (m *MockService) GetAllProductsService(userID string) (*dtos.ProductListResponse, *exceptions.RestErr) {
	args := m.Called(userID)
	if args.Get(1) == nil {
		return args.Get(0).(*dtos.ProductListResponse), nil
	}
	return nil, args.Get(1).(*exceptions.RestErr)
}

func (m *MockService) CreateProductService(userID string, request dtos.CreateProductRequest) (bool, *exceptions.RestErr) {
	args := m.Called(userID, request)
	if args.Get(1) == nil {
		return args.Bool(0), nil
	}
	return false, args.Get(1).(*exceptions.RestErr)
}

func (m *MockService) DeleteProductService(userID string, id string) (bool, *exceptions.RestErr) {
	args := m.Called(userID, id)
	if args.Get(1) == nil {
		return args.Bool(0), nil
	}
	return false, args.Get(1).(*exceptions.RestErr)
}

func (m *MockService) GetAllPedidosService(userID string) (*dtos.PedidoListResponse, *exceptions.RestErr) {
	args := m.Called(userID)
	if args.Get(1) == nil {
		return args.Get(0).(*dtos.PedidoListResponse), nil
	}
	return nil, args.Get(1).(*exceptions.RestErr)
}

// Helper function para criar app de teste
func setupTestApp(_ ControllerInterface) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
	return app
}

// Helper function para fazer requests
func makeRequest(app *fiber.App, method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, url, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return app.Test(req)
}

// TESTES DE HEALTH CHECKS
func TestController_HealthCheck(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)
	app.Get("/health", controller.HealthCheck)

	// Act
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "alive", response["status"])
	assert.Contains(t, response, "timestamp")
	assert.Contains(t, response, "version")
	assert.Contains(t, response, "uptime")
	assert.Contains(t, response, "message")
}

// TESTES DE LOGIN
func TestController_LoginUser_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)
	app.Post("/login", controller.LoginUser)

	loginRequest := dtos.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockService.On("LoginUserService", loginRequest).Return("valid_jwt_token", nil)

	// Act
	jsonBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Login successful", response["message"])
	assert.Equal(t, "valid_jwt_token", response["token"])

	mockService.AssertExpectations(t)
}

func TestController_LoginUser_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)
	app.Post("/login", controller.LoginUser)

	// Act - Enviar JSON inválido
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Dados de login inválidos", response["error"])
}

func TestController_LoginUser_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)
	app.Post("/login", controller.LoginUser)

	loginRequest := dtos.UserLogin{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockService.On("LoginUserService", loginRequest).Return("", exceptions.NewUnauthorizedRequestError("Invalid credentials"))

	// Act
	jsonBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Invalid credentials", response["error"])

	mockService.AssertExpectations(t)
}

// TESTES DE CREATE USER
func TestController_CreateUser_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	// Simular middleware que adiciona createUser ao context
	app.Post("/cadastro", func(c *fiber.Ctx) error {
		createUser := dtos.CreateUser{
			FirstName:   "João",
			LastName:    "Silva",
			Email:       "joao@test.com",
			NomeEmpresa: "Test Company",
			Password:    "password123",
		}
		c.Locals("createUser", createUser)
		return controller.CreateUser(c)
	})

	expectedUser := &entity.User{
		ID:          1,
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@test.com",
		NomeEmpresa: "Test Company",
	}

	mockService.On("CreateUserService", mock.AnythingOfType("dtos.CreateUser")).Return(expectedUser, nil)

	// Act
	req := httptest.NewRequest("POST", "/cadastro", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "User created successfully", response["message"])
	assert.Contains(t, response, "usuário")

	mockService.AssertExpectations(t)
}

func TestController_CreateUser_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Post("/cadastro", func(c *fiber.Ctx) error {
		createUser := dtos.CreateUser{
			Email: "existing@test.com",
		}
		c.Locals("createUser", createUser)
		return controller.CreateUser(c)
	})

	mockService.On("CreateUserService", mock.AnythingOfType("dtos.CreateUser")).Return(nil, exceptions.NewBadRequestError("Email already exists"))

	// Act
	req := httptest.NewRequest("POST", "/cadastro", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Contains(t, response["error"], "Email already exists")

	mockService.AssertExpectations(t)
}

// TESTES DE FORNECEDORES
func TestController_GetAllFornecedores_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Get("/fornecedores", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.GetAllFornecedores(c)
	})

	expectedResponse := &dtos.FornecedorListResponse{
		Fornecedores: []dtos.FornecedorResponse{
			{
				ID:   1,
				Nome: "Fornecedor 1",
			},
			{
				ID:   2,
				Nome: "Fornecedor 2",
			},
		},
		Total: 2,
	}

	mockService.On("GetAllFornecedoresService", "1").Return(expectedResponse, nil)

	// Act
	req := httptest.NewRequest("GET", "/fornecedores", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.FornecedorListResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, 2, response.Total)
	assert.Equal(t, 2, len(response.Fornecedores))
	assert.Equal(t, "Fornecedor 1", response.Fornecedores[0].Nome)

	mockService.AssertExpectations(t)
}

func TestController_CreateFornecedor_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Post("/fornecedores", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.CreateFornecedor(c)
	})

	createRequest := dtos.CreateFornecedorRequest{
		Nome:     "Novo Fornecedor",
		Email:    "fornecedor@test.com",
		Telefone: "11999999999",
	}

	mockService.On("CreateFornecedorService", "1", createRequest).Return(true, nil)

	// Act
	jsonBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/fornecedores", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Fornecedor created successfully", response["message"])

	mockService.AssertExpectations(t)
}

// TESTES DE PRODUTOS
func TestController_GetAllProducts_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Get("/produtos", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.GetAllProducts(c)
	})

	expectedResponse := &dtos.ProductListResponse{
		Products: []dtos.ProductResponse{
			{
				ID:          1,
				NomeProduto: "Produto 1",
				PrecoVenda:  10.50,
			},
		},
		Total: 1,
	}

	mockService.On("GetAllProductsService", "1").Return(expectedResponse, nil)

	// Act
	req := httptest.NewRequest("GET", "/produtos", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.ProductListResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, 1, response.Total)
	assert.Equal(t, "Produto 1", response.Products[0].NomeProduto)

	mockService.AssertExpectations(t)
}

// TESTES DE PEDIDOS
func TestController_GetAllPedidos_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Get("/pedidos", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.GetAllPedidos(c)
	})

	expectedResponse := &dtos.PedidoListResponse{
		Pedidos: []dtos.PedidoResponse{
			{
				ID:         1,
				ValorTotal: 100.50,
			},
		},
		Total: 1,
	}

	mockService.On("GetAllPedidosService", "1").Return(expectedResponse, nil)

	// Act
	req := httptest.NewRequest("GET", "/pedidos", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.PedidoListResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, 1, response.Total)
	assert.Equal(t, float64(100.50), response.Pedidos[0].ValorTotal)

	mockService.AssertExpectations(t)
}

// TESTE DE RequestOtherService
func TestController_RequestOtherService(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)
	app.Get("/other", controller.RequestOtherService)

	// Act
	req := httptest.NewRequest("GET", "/other", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Service found successfully", response["message"])
}
