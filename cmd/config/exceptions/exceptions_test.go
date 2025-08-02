package exceptions

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TESTES PARA RestErr
// =============================================================================

func TestRestErr_Error(t *testing.T) {
	// Arrange
	restErr := &RestErr{
		Message: "Test error message",
		Err:     "test_error",
		Code:    500,
	}

	// Act
	errorString := restErr.Error()

	// Assert
	assert.Equal(t, "Test error message", errorString)
}

func TestRestErr_JSONSerialization(t *testing.T) {
	// Arrange
	causes := []Causes{
		{
			Field:        "email",
			FieldMessage: "email is required",
		},
		{
			Field:        "password",
			FieldMessage: "password must be at least 6 characters",
		},
	}

	restErr := &RestErr{
		Message: "Validation failed",
		Err:     "bad_request",
		Code:    400,
		Causes:  causes,
	}

	// Act
	jsonData, err := json.Marshal(restErr)
	assert.NoError(t, err)

	var unmarshaled RestErr
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, restErr.Message, unmarshaled.Message)
	assert.Equal(t, restErr.Err, unmarshaled.Err)
	assert.Equal(t, restErr.Code, unmarshaled.Code)
	assert.Len(t, unmarshaled.Causes, 2)
	assert.Equal(t, restErr.Causes[0].Field, unmarshaled.Causes[0].Field)
	assert.Equal(t, restErr.Causes[0].FieldMessage, unmarshaled.Causes[0].FieldMessage)
}

func TestRestErr_JSONTags(t *testing.T) {
	// Arrange
	restErr := &RestErr{
		Message: "Test message",
		Err:     "test_error",
		Code:    500,
		Causes:  []Causes{},
	}

	// Act
	jsonData, err := json.Marshal(restErr)
	assert.NoError(t, err)

	// Assert
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, `"message":"Test message"`)
	assert.Contains(t, jsonString, `"error":"test_error"`)
	assert.Contains(t, jsonString, `"code":500`)
	assert.Contains(t, jsonString, `"causes":[]`)
}

// =============================================================================
// TESTES PARA Causes
// =============================================================================

func TestCauses_JSONSerialization(t *testing.T) {
	// Arrange
	cause := Causes{
		Field:        "username",
		FieldMessage: "username is required and must be unique",
	}

	// Act
	jsonData, err := json.Marshal(cause)
	assert.NoError(t, err)

	var unmarshaled Causes
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, cause.Field, unmarshaled.Field)
	assert.Equal(t, cause.FieldMessage, unmarshaled.FieldMessage)
}

func TestCauses_JSONTags(t *testing.T) {
	// Arrange
	cause := Causes{
		Field:        "email",
		FieldMessage: "email format is invalid",
	}

	// Act
	jsonData, err := json.Marshal(cause)
	assert.NoError(t, err)

	// Assert
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, `"field":"email"`)
	assert.Contains(t, jsonString, `"message":"email format is invalid"`)
}

// =============================================================================
// TESTES PARA FUNÇÕES DE CRIAÇÃO DE ERROS
// =============================================================================

func TestNewBadRequestError(t *testing.T) {
	// Arrange
	message := "Invalid request data"

	// Act
	err := NewBadRequestError(message)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, "bad_request", err.Err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Nil(t, err.Causes)
}

func TestNewUnauthorizedRequestError(t *testing.T) {
	// Arrange
	message := "Access denied"

	// Act
	err := NewUnauthorizedRequestError(message)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, "unauthorized", err.Err)
	assert.Equal(t, http.StatusUnauthorized, err.Code)
	assert.Nil(t, err.Causes)
}

func TestNewBadRequestValidationError(t *testing.T) {
	// Arrange
	message := "Validation failed"
	causes := []Causes{
		{
			Field:        "name",
			FieldMessage: "name is required",
		},
		{
			Field:        "age",
			FieldMessage: "age must be a positive number",
		},
	}

	// Act
	err := NewBadRequestValidationError(message, causes)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, "bad_request", err.Err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Len(t, err.Causes, 2)
	assert.Equal(t, causes[0].Field, err.Causes[0].Field)
	assert.Equal(t, causes[0].FieldMessage, err.Causes[0].FieldMessage)
	assert.Equal(t, causes[1].Field, err.Causes[1].Field)
	assert.Equal(t, causes[1].FieldMessage, err.Causes[1].FieldMessage)
}

func TestNewInternalServerError(t *testing.T) {
	// Arrange
	message := "Internal server error occurred"

	// Act
	err := NewInternalServerError(message)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, "internal_server_error", err.Err)
	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Nil(t, err.Causes)
}

func TestNewNotFoundError(t *testing.T) {
	// Arrange
	message := "Resource not found"

	// Act
	err := NewNotFoundError(message)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, "not_found", err.Err)
	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Nil(t, err.Causes)
}

func TestNewForbiddenError(t *testing.T) {
	// Arrange
	message := "Access forbidden"

	// Act
	err := NewForbiddenError(message)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, "forbidden", err.Err)
	assert.Equal(t, http.StatusForbidden, err.Code)
	assert.Nil(t, err.Causes)
}

// =============================================================================
// TESTES DE EDGE CASES
// =============================================================================

func TestRestErr_EmptyMessage(t *testing.T) {
	// Arrange
	restErr := &RestErr{
		Message: "",
		Err:     "empty_error",
		Code:    400,
	}

	// Act
	errorString := restErr.Error()

	// Assert
	assert.Equal(t, "", errorString)
}

func TestRestErr_NilCauses(t *testing.T) {
	// Arrange
	restErr := &RestErr{
		Message: "Test error",
		Err:     "test_error",
		Code:    500,
		Causes:  nil,
	}

	// Act
	jsonData, err := json.Marshal(restErr)
	assert.NoError(t, err)

	var unmarshaled RestErr
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, unmarshaled.Causes)
}

func TestRestErr_EmptyCauses(t *testing.T) {
	// Arrange
	restErr := &RestErr{
		Message: "Test error",
		Err:     "test_error",
		Code:    500,
		Causes:  []Causes{},
	}

	// Act
	jsonData, err := json.Marshal(restErr)
	assert.NoError(t, err)

	var unmarshaled RestErr
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, unmarshaled.Causes)
	assert.Len(t, unmarshaled.Causes, 0)
}

func TestCauses_EmptyValues(t *testing.T) {
	// Arrange
	cause := Causes{
		Field:        "",
		FieldMessage: "",
	}

	// Act
	jsonData, err := json.Marshal(cause)
	assert.NoError(t, err)

	var unmarshaled Causes
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "", unmarshaled.Field)
	assert.Equal(t, "", unmarshaled.FieldMessage)
}

func TestCauses_SpecialCharacters(t *testing.T) {
	// Arrange
	cause := Causes{
		Field:        "field_with_underscore",
		FieldMessage: "Message with special chars: @#$%^&*()+=[]{}|\\:;\"'<>?,./ and unicode: 中文",
	}

	// Act
	jsonData, err := json.Marshal(cause)
	assert.NoError(t, err)

	var unmarshaled Causes
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, cause.Field, unmarshaled.Field)
	assert.Equal(t, cause.FieldMessage, unmarshaled.FieldMessage)
}

// =============================================================================
// TESTES DE INTEGRAÇÃO
// =============================================================================

func TestAllErrorTypes_HTTPStatusCodes(t *testing.T) {
	tests := []struct {
		name         string
		errorFunc    func(string) *RestErr
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "BadRequest",
			errorFunc:    NewBadRequestError,
			expectedCode: http.StatusBadRequest,
			expectedErr:  "bad_request",
		},
		{
			name:         "Unauthorized",
			errorFunc:    NewUnauthorizedRequestError,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  "unauthorized",
		},
		{
			name:         "InternalServer",
			errorFunc:    NewInternalServerError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  "internal_server_error",
		},
		{
			name:         "NotFound",
			errorFunc:    NewNotFoundError,
			expectedCode: http.StatusNotFound,
			expectedErr:  "not_found",
		},
		{
			name:         "Forbidden",
			errorFunc:    NewForbiddenError,
			expectedCode: http.StatusForbidden,
			expectedErr:  "forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			message := "Test message for " + tt.name

			// Act
			err := tt.errorFunc(message)

			// Assert
			assert.Equal(t, tt.expectedCode, err.Code)
			assert.Equal(t, tt.expectedErr, err.Err)
			assert.Equal(t, message, err.Message)
		})
	}
}

func TestRestErr_AsError(t *testing.T) {
	// Arrange
	restErr := NewInternalServerError("Something went wrong")

	// Act - Test that RestErr implements error interface
	var err error = restErr
	errorMessage := err.Error()

	// Assert
	assert.Equal(t, "Something went wrong", errorMessage)
}

func TestBadRequestValidationError_MultipleCauses(t *testing.T) {
	// Arrange
	causes := []Causes{
		{Field: "email", FieldMessage: "email is required"},
		{Field: "password", FieldMessage: "password too short"},
		{Field: "name", FieldMessage: "name contains invalid characters"},
		{Field: "age", FieldMessage: "age must be between 18 and 120"},
	}

	// Act
	err := NewBadRequestValidationError("Multiple validation errors", causes)

	// Assert
	assert.Equal(t, "Multiple validation errors", err.Message)
	assert.Equal(t, "bad_request", err.Err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Len(t, err.Causes, 4)

	// Verify all causes are preserved
	for i, expectedCause := range causes {
		assert.Equal(t, expectedCause.Field, err.Causes[i].Field)
		assert.Equal(t, expectedCause.FieldMessage, err.Causes[i].FieldMessage)
	}
}
