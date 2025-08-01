package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Certifique-se de que o UserValidationMiddleware e ValidateUnexpectedFields estão importados corretamente
// ou definidos neste arquivo.

func TestUserValidationMiddleware(t *testing.T) {
	// Setup Fiber app for testing
	app := fiber.New()

	// Adiciona o middleware e um endpoint de teste
	app.Post("/test", UserValidationMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		shouldPass     bool
	}{
		{
			name: "Valid user data",
			requestBody: map[string]interface{}{
				"first_name":   "John",
				"last_name":    "Doe",
				"email":        "john.doe@example.com",
				"nome_empresa": "Test Company",
				"categoria":    "Technology",
				"segmento":     "Software",
				"city":         "São Paulo",
				"state":        "SP",
				"password":     "password123!",
			},
			expectedStatus: 200,
			shouldPass:     true,
		},
		// ... outros casos de teste ...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			jsonBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest("POST", "/test", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Act
			resp, err := app.Test(req)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.shouldPass {
				// Se a validação deve passar, verifique se obtemos a resposta de sucesso
				body, _ := io.ReadAll(resp.Body)
				var response map[string]interface{}
				json.Unmarshal(body, &response)
				assert.Equal(t, true, response["success"])
			} else {
				// Se a validação deve falhar, verifique se obtemos a resposta de erro
				body, _ := io.ReadAll(resp.Body)
				var response map[string]interface{}
				json.Unmarshal(body, &response)
				assert.Contains(t, response, "error")
			}
		})
	}
}

func TestValidateUnexpectedFields(t *testing.T) {
	// Setup Fiber context for testing
	app := fiber.New()
	c := app.AcquireCtx(&fiber.Ctx{}) // Corrigido para usar o tipo correto
	defer app.ReleaseCtx(c)

	testCases := []struct {
		name        string
		requestData map[string]interface{}
		shouldError bool
	}{
		{
			name: "No unexpected fields",
			requestData: map[string]interface{}{
				"first_name":   "John",
				"last_name":    "Doe",
				"email":        "john@example.com",
				"nome_empresa": "Company",
				"categoria":    "Tech",
				"segmento":     "Software",
				"city":         "SP",
				"state":        "SP",
				"password":     "pass123!",
			},
			shouldError: false,
		},
		{
			name: "With unexpected fields",
			requestData: map[string]interface{}{
				"first_name":       "John",
				"last_name":        "Doe",
				"email":            "john@example.com",
				"unexpected_field": "should not be here",
			},
			shouldError: true,
		},
		{
			name: "Invalid JSON",
			requestData: map[string]interface{}{
				"invalid": make(chan int), // Isso causará falha na serialização JSON
			},
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			var jsonData []byte
			var err error

			if tc.name == "Invalid JSON" {
				// Cria JSON inválido manualmente
				jsonData = []byte(`{"invalid": invalid}`)
			} else {
				jsonData, err = json.Marshal(tc.requestData)
				assert.NoError(t, err)
			}

			// Act
			validationErr := ValidateUnexpectedFields(c, jsonData)

			// Assert
			if tc.shouldError {
				assert.Error(t, validationErr)
			} else {
				assert.NoError(t, validationErr)
			}
		})
	}
}
