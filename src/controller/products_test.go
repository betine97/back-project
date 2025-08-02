package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/betine97/back-project.git/cmd/config/exceptions"
	"github.com/betine97/back-project.git/src/model/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TESTES COMPLETOS PARA PRODUTOS

func TestController_GetAllProducts_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Get("/produtos", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.GetAllProducts(c)
	})

	mockService.On("GetAllProductsService", "1").Return(nil, exceptions.NewInternalServerError("Database connection failed"))

	// Act
	req := httptest.NewRequest("GET", "/produtos", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Database connection failed", response["error"])

	mockService.AssertExpectations(t)
}

func TestController_CreateProduct_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Post("/produtos", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.CreateProduct(c)
	})

	createRequest := dtos.CreateProductRequest{
		NomeProduto: "Produto Teste",
		SKU:         "SKU001",
		PrecoVenda:  29.99,
		Status:      "Ativo",
	}

	mockService.On("CreateProductService", "1", createRequest).Return(true, nil)

	// Act
	jsonBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/produtos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Product created successfully", response["message"])

	mockService.AssertExpectations(t)
}

func TestController_CreateProduct_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Post("/produtos", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.CreateProduct(c)
	})

	// Act - Enviar JSON inválido
	req := httptest.NewRequest("POST", "/produtos", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Unable to read request data", response["error"])
}

func TestController_CreateProduct_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Post("/produtos", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.CreateProduct(c)
	})

	createRequest := dtos.CreateProductRequest{
		NomeProduto: "",
		SKU:         "SKU001",
		PrecoVenda:  -10.00, // Preço inválido
	}

	mockService.On("CreateProductService", "1", createRequest).Return(false, exceptions.NewBadRequestError("Invalid product data"))

	// Act
	jsonBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/produtos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Invalid product data", response["error"])

	mockService.AssertExpectations(t)
}

func TestController_CreateProduct_ServiceReturnsFalse(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Post("/produtos", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.CreateProduct(c)
	})

	createRequest := dtos.CreateProductRequest{
		NomeProduto: "Produto Teste",
		SKU:         "SKU001",
		PrecoVenda:  29.99,
	}

	mockService.On("CreateProductService", "1", createRequest).Return(false, nil)

	// Act
	jsonBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/produtos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Error creating product", response["error"])

	mockService.AssertExpectations(t)
}

func TestController_DeleteProduct_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Delete("/produtos/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.DeleteProduct(c)
	})

	mockService.On("DeleteProductService", "1", "123").Return(true, nil)

	// Act
	req := httptest.NewRequest("DELETE", "/produtos/123", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Product deleted successfully", response["message"])

	mockService.AssertExpectations(t)
}

func TestController_DeleteProduct_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Delete("/produtos/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.DeleteProduct(c)
	})

	mockService.On("DeleteProductService", "1", "999").Return(false, exceptions.NewNotFoundError("Product not found"))

	// Act
	req := httptest.NewRequest("DELETE", "/produtos/999", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Product not found", response["error"])

	mockService.AssertExpectations(t)
}

func TestController_DeleteProduct_ServiceReturnsFalse(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Delete("/produtos/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.DeleteProduct(c)
	})

	mockService.On("DeleteProductService", "1", "123").Return(false, nil)

	// Act
	req := httptest.NewRequest("DELETE", "/produtos/123", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Error deleting product", response["error"])

	mockService.AssertExpectations(t)
}

// TESTES PARA PEDIDOS

func TestController_GetAllPedidos_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Get("/pedidos", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.GetAllPedidos(c)
	})

	mockService.On("GetAllPedidosService", "1").Return(nil, exceptions.NewInternalServerError("Database connection failed"))

	// Act
	req := httptest.NewRequest("GET", "/pedidos", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Database connection failed", response["error"])

	mockService.AssertExpectations(t)
}

func TestController_GetAllPedidos_EmptyList(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Get("/pedidos", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.GetAllPedidos(c)
	})

	expectedResponse := &dtos.PedidoListResponse{
		Pedidos: []dtos.PedidoResponse{},
		Total:   0,
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

	assert.Equal(t, 0, response.Total)
	assert.Equal(t, 0, len(response.Pedidos))

	mockService.AssertExpectations(t)
}

func TestController_GetAllPedidos_MultipleItems(t *testing.T) {
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
				ID:           1,
				IDFornecedor: 1,
				DataPedido:   "2024-01-01",
				ValorTotal:   100.50,
				Status:       "Pendente",
			},
			{
				ID:           2,
				IDFornecedor: 2,
				DataPedido:   "2024-01-02",
				ValorTotal:   250.99,
				Status:       "Concluído",
			},
		},
		Total: 2,
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

	assert.Equal(t, 2, response.Total)
	assert.Equal(t, 2, len(response.Pedidos))
	assert.Equal(t, "Pendente", response.Pedidos[0].Status)
	assert.Equal(t, "Concluído", response.Pedidos[1].Status)
	assert.Equal(t, float64(100.50), response.Pedidos[0].ValorTotal)
	assert.Equal(t, float64(250.99), response.Pedidos[1].ValorTotal)

	mockService.AssertExpectations(t)
}
