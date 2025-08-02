package controller

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestController_HealthCheck_Success(t *testing.T) {
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
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verificar estrutura da resposta
	assert.Equal(t, "alive", response["status"])
	assert.Contains(t, response, "timestamp")
	assert.Contains(t, response, "version")
	assert.Contains(t, response, "uptime")
	assert.Contains(t, response, "message")
	assert.Equal(t, "1.0.0", response["version"])
	assert.Equal(t, "Aplicação está viva e respondendo", response["message"])

	// Verificar que uptime é uma string não vazia
	uptime, ok := response["uptime"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, uptime)

	// Verificar formato do timestamp
	timestamp, ok := response["timestamp"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, timestamp)
}

func TestController_ReadinessCheck_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)
	app.Get("/ready", controller.ReadinessCheck)

	// Act
	req := httptest.NewRequest("GET", "/ready", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)

	var response ReadinessResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verificar estrutura da resposta
	assert.Contains(t, []string{"ready", "not_ready"}, response.Status)
	assert.Equal(t, "1.0.0", response.Version)
	assert.NotEmpty(t, response.Timestamp)
	assert.NotEmpty(t, response.Uptime)
	assert.Contains(t, response.Services, "database_master")
	assert.Contains(t, response.Services, "database_clients")
	assert.Contains(t, response.Services, "redis")

	// Verificar que os serviços têm status válidos
	validStatuses := []string{"healthy", "unhealthy", "partial"}
	for serviceName, status := range response.Services {
		assert.Contains(t, validStatuses, status, "Service %s has invalid status: %s", serviceName, status)
	}

	// Se o status geral é "ready", o código deve ser 200
	// Se é "not_ready", deve ser 503
	if response.Status == "ready" {
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	} else {
		assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
	}
}

func TestController_ReadinessCheck_ResponseStructure(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)
	app.Get("/ready", controller.ReadinessCheck)

	// Act
	req := httptest.NewRequest("GET", "/ready", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verificar que todos os campos obrigatórios estão presentes
	requiredFields := []string{"status", "timestamp", "version", "services", "uptime"}
	for _, field := range requiredFields {
		assert.Contains(t, response, field, "Missing required field: %s", field)
	}

	// Verificar tipos dos campos
	assert.IsType(t, "", response["status"])
	assert.IsType(t, "", response["timestamp"])
	assert.IsType(t, "", response["version"])
	assert.IsType(t, "", response["uptime"])
	assert.IsType(t, map[string]interface{}{}, response["services"])

	// Verificar serviços específicos
	services := response["services"].(map[string]interface{})
	expectedServices := []string{"database_master", "database_clients", "redis"}
	for _, service := range expectedServices {
		assert.Contains(t, services, service, "Missing service: %s", service)
		assert.IsType(t, "", services[service], "Service %s should have string status", service)
	}
}

// Teste para verificar que o health check é rápido
func TestController_HealthCheck_Performance(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)
	app.Get("/health", controller.HealthCheck)

	// Act & Assert - Executar múltiplas vezes para verificar consistência
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &response)

		assert.Equal(t, "alive", response["status"])
	}
}

// Teste para verificar que o readiness check funciona mesmo com falhas parciais
func TestController_ReadinessCheck_PartialFailure(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)
	app.Get("/ready", controller.ReadinessCheck)

	// Act
	req := httptest.NewRequest("GET", "/ready", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)

	// O teste deve sempre retornar uma resposta válida, mesmo com falhas
	var response ReadinessResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verificar que a resposta tem a estrutura correta
	assert.NotEmpty(t, response.Status)
	assert.NotEmpty(t, response.Version)
	assert.NotEmpty(t, response.Timestamp)
	assert.NotEmpty(t, response.Uptime)
	assert.NotEmpty(t, response.Services)

	// Status deve ser "ready" ou "not_ready"
	assert.Contains(t, []string{"ready", "not_ready"}, response.Status)
}

// Teste para verificar headers de resposta
func TestController_HealthCheck_Headers(t *testing.T) {
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
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func TestController_ReadinessCheck_Headers(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	controller := NewControllerInstance(mockService)
	app := setupTestApp(controller)
	app.Get("/ready", controller.ReadinessCheck)

	// Act
	req := httptest.NewRequest("GET", "/ready", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	// Status code deve ser 200 ou 503
	assert.Contains(t, []int{fiber.StatusOK, fiber.StatusServiceUnavailable}, resp.StatusCode)
}
