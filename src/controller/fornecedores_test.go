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

// TESTES COMPLETOS PARA FORNECEDORES

func TestController_GetAllFornecedores_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Get("/fornecedores", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.GetAllFornecedores(c)
	})

	mockService.On("GetAllFornecedoresService", "1").Return(nil, exceptions.NewInternalServerError("Database connection failed"))

	// Act
	req := httptest.NewRequest("GET", "/fornecedores", nil)
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

func TestController_CreateFornecedor_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Post("/fornecedores", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.CreateFornecedor(c)
	})

	// Act - Enviar JSON inválido
	req := httptest.NewRequest("POST", "/fornecedores", bytes.NewBuffer([]byte("invalid json")))
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

func TestController_CreateFornecedor_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Post("/fornecedores", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.CreateFornecedor(c)
	})

	createRequest := dtos.CreateFornecedorRequest{
		Nome:  "Fornecedor Teste",
		Email: "invalid-email",
	}

	mockService.On("CreateFornecedorService", "1", createRequest).Return(false, exceptions.NewBadRequestError("Invalid email format"))

	// Act
	jsonBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/fornecedores", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Invalid email format", response["error"])

	mockService.AssertExpectations(t)
}

func TestController_CreateFornecedor_ServiceReturnsFalse(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Post("/fornecedores", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.CreateFornecedor(c)
	})

	createRequest := dtos.CreateFornecedorRequest{
		Nome:  "Fornecedor Teste",
		Email: "test@example.com",
	}

	mockService.On("CreateFornecedorService", "1", createRequest).Return(false, nil)

	// Act
	jsonBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/fornecedores", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Error creating fornecedor", response["error"])

	mockService.AssertExpectations(t)
}

func TestController_ChangeStatusFornecedor_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Put("/fornecedores/changestatus/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.ChangeStatusFornecedor(c)
	})

	mockService.On("ChangeStatusFornecedorService", "1", "123").Return(true, nil)

	// Act
	req := httptest.NewRequest("PUT", "/fornecedores/changestatus/123", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Status fornecedor changed successfully", response["message"])

	mockService.AssertExpectations(t)
}

func TestController_ChangeStatusFornecedor_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Put("/fornecedores/changestatus/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.ChangeStatusFornecedor(c)
	})

	mockService.On("ChangeStatusFornecedorService", "1", "999").Return(false, exceptions.NewNotFoundError("Fornecedor not found"))

	// Act
	req := httptest.NewRequest("PUT", "/fornecedores/changestatus/999", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Fornecedor not found", response["error"])

	mockService.AssertExpectations(t)
}

func TestController_UpdateFornecedorField_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Put("/fornecedores/changefields/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.UpdateFornecedorField(c)
	})

	updateRequest := dtos.UpdateFornecedorRequest{
		Campo: "nome",
		Valor: "Novo Nome",
	}

	mockService.On("UpdateFornecedorFieldService", "1", "123", "nome", "Novo Nome").Return(true, nil)

	// Act
	jsonBody, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest("PUT", "/fornecedores/changefields/123", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Fornecedor field updated successfully", response["message"])

	mockService.AssertExpectations(t)
}

func TestController_UpdateFornecedorField_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Put("/fornecedores/changefields/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.UpdateFornecedorField(c)
	})

	// Act - Enviar JSON inválido
	req := httptest.NewRequest("PUT", "/fornecedores/changefields/123", bytes.NewBuffer([]byte("invalid json")))
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

func TestController_UpdateFornecedorField_InvalidField(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Put("/fornecedores/changefields/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.UpdateFornecedorField(c)
	})

	updateRequest := dtos.UpdateFornecedorRequest{
		Campo: "campo_invalido",
		Valor: "Valor",
	}

	mockService.On("UpdateFornecedorFieldService", "1", "123", "campo_invalido", "Valor").Return(false, exceptions.NewBadRequestError("Invalid field to update"))

	// Act
	jsonBody, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest("PUT", "/fornecedores/changefields/123", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Invalid field to update", response["error"])

	mockService.AssertExpectations(t)
}

func TestController_DeleteFornecedor_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Delete("/fornecedores/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.DeleteFornecedor(c)
	})

	mockService.On("DeleteFornecedorService", "1", "123").Return(true, nil)

	// Act
	req := httptest.NewRequest("DELETE", "/fornecedores/123", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Fornecedor deleted successfully", response["message"])

	mockService.AssertExpectations(t)
}

func TestController_DeleteFornecedor_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Delete("/fornecedores/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.DeleteFornecedor(c)
	})

	mockService.On("DeleteFornecedorService", "1", "999").Return(false, exceptions.NewNotFoundError("Fornecedor not found"))

	// Act
	req := httptest.NewRequest("DELETE", "/fornecedores/999", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Fornecedor not found", response["error"])

	mockService.AssertExpectations(t)
}

func TestController_DeleteFornecedor_ServiceReturnsFalse(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)

	app.Delete("/fornecedores/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controller.DeleteFornecedor(c)
	})

	mockService.On("DeleteFornecedorService", "1", "123").Return(false, nil)

	// Act
	req := httptest.NewRequest("DELETE", "/fornecedores/123", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, "Error deleting fornecedor", response["error"])

	mockService.AssertExpectations(t)
}
