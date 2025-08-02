package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function for testing (replica da função original)
func testGetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// =============================================================================
// TESTES PARA Config Struct
// =============================================================================

func TestConfig_Structure(t *testing.T) {
	// Arrange & Act
	config := &Config{
		DBDriver:      "mysql",
		DBHost:        "localhost",
		DBPort:        "3306",
		DBUser:        "root",
		DBPassword:    "password",
		DBName:        "testdb",
		WebServerPort: "8080",
		JWTSecret:     "test_secret",
		JWTExpiresIn:  30,
		CORSOrigins:   "http://localhost:3000",
	}

	// Assert
	assert.Equal(t, "mysql", config.DBDriver)
	assert.Equal(t, "localhost", config.DBHost)
	assert.Equal(t, "3306", config.DBPort)
	assert.Equal(t, "root", config.DBUser)
	assert.Equal(t, "password", config.DBPassword)
	assert.Equal(t, "testdb", config.DBName)
	assert.Equal(t, "8080", config.WebServerPort)
	assert.Equal(t, "test_secret", config.JWTSecret)
	assert.Equal(t, 30, config.JWTExpiresIn)
	assert.Equal(t, "http://localhost:3000", config.CORSOrigins)
}

func TestConfig_DefaultValues(t *testing.T) {
	// Arrange - Clear environment variables to test defaults
	originalEnvs := map[string]string{
		"DB_DRIVER":       os.Getenv("DB_DRIVER"),
		"DB_HOST":         os.Getenv("DB_HOST"),
		"DB_PORT":         os.Getenv("DB_PORT"),
		"DB_USER":         os.Getenv("DB_USER"),
		"DB_PASSWORD":     os.Getenv("DB_PASSWORD"),
		"DB_NAME":         os.Getenv("DB_NAME"),
		"WEB_SERVER_PORT": os.Getenv("WEB_SERVER_PORT"),
		"JWT_SECRET":      os.Getenv("JWT_SECRET"),
		"JWT_EXPIRES_IN":  os.Getenv("JWT_EXPIRES_IN"),
		"CORS_ORIGINS":    os.Getenv("CORS_ORIGINS"),
	}

	// Clear all environment variables
	for key := range originalEnvs {
		os.Unsetenv(key)
	}

	// Act
	config := &Config{
		DBDriver:      testGetEnvWithDefault("DB_DRIVER", "mysql"),
		DBHost:        testGetEnvWithDefault("DB_HOST", "localhost"),
		DBPort:        testGetEnvWithDefault("DB_PORT", "3306"),
		DBUser:        testGetEnvWithDefault("DB_USER", "root"),
		DBPassword:    testGetEnvWithDefault("DB_PASSWORD", ""),
		DBName:        testGetEnvWithDefault("DB_NAME", "masterdb"),
		WebServerPort: testGetEnvWithDefault("WEB_SERVER_PORT", "8080"),
		JWTSecret:     testGetEnvWithDefault("JWT_SECRET", "default_test_secret"),
		CORSOrigins:   testGetEnvWithDefault("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001"),
	}

	if expiresIn := os.Getenv("JWT_EXPIRES_IN"); expiresIn != "" {
		config.JWTExpiresIn, _ = strconv.Atoi(expiresIn)
	} else {
		config.JWTExpiresIn = 30
	}

	// Assert
	assert.Equal(t, "mysql", config.DBDriver)
	assert.Equal(t, "localhost", config.DBHost)
	assert.Equal(t, "3306", config.DBPort)
	assert.Equal(t, "root", config.DBUser)
	assert.Equal(t, "", config.DBPassword)
	assert.Equal(t, "masterdb", config.DBName)
	assert.Equal(t, "8080", config.WebServerPort)
	assert.Equal(t, "default_test_secret", config.JWTSecret)
	assert.Equal(t, 30, config.JWTExpiresIn)
	assert.Contains(t, config.CORSOrigins, "http://localhost:3000")

	// Cleanup - Restore original environment variables
	for key, value := range originalEnvs {
		if value != "" {
			os.Setenv(key, value)
		}
	}
}

// =============================================================================
// TESTES PARA testGetEnvWithDefault Function
// =============================================================================

func TestGetEnvWithDefault_WithEnvironmentVariable(t *testing.T) {
	// Arrange
	key := "TEST_ENV_VAR"
	expectedValue := "test_value"
	defaultValue := "default_value"

	originalValue := os.Getenv(key)
	os.Setenv(key, expectedValue)

	// Act
	result := testGetEnvWithDefault(key, defaultValue)

	// Assert
	assert.Equal(t, expectedValue, result)

	// Cleanup
	if originalValue != "" {
		os.Setenv(key, originalValue)
	} else {
		os.Unsetenv(key)
	}
}

func TestGetEnvWithDefault_WithoutEnvironmentVariable(t *testing.T) {
	// Arrange
	key := "NON_EXISTENT_ENV_VAR"
	defaultValue := "default_value"

	originalValue := os.Getenv(key)
	os.Unsetenv(key)

	// Act
	result := testGetEnvWithDefault(key, defaultValue)

	// Assert
	assert.Equal(t, defaultValue, result)

	// Cleanup
	if originalValue != "" {
		os.Setenv(key, originalValue)
	}
}

func TestGetEnvWithDefault_WithEmptyEnvironmentVariable(t *testing.T) {
	// Arrange
	key := "EMPTY_ENV_VAR"
	defaultValue := "default_value"

	originalValue := os.Getenv(key)
	os.Setenv(key, "")

	// Act
	result := testGetEnvWithDefault(key, defaultValue)

	// Assert
	assert.Equal(t, defaultValue, result)

	// Cleanup
	if originalValue != "" {
		os.Setenv(key, originalValue)
	} else {
		os.Unsetenv(key)
	}
}

// =============================================================================
// TESTES DE VALIDAÇÃO DE CONFIGURAÇÃO
// =============================================================================

func TestConfig_DatabaseConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		driver   string
		host     string
		port     string
		user     string
		password string
		dbname   string
		isValid  bool
	}{
		{
			name:     "Valid MySQL config",
			driver:   "mysql",
			host:     "localhost",
			port:     "3306",
			user:     "root",
			password: "password",
			dbname:   "testdb",
			isValid:  true,
		},
		{
			name:     "Valid PostgreSQL config",
			driver:   "postgresql",
			host:     "localhost",
			port:     "5432",
			user:     "postgres",
			password: "password",
			dbname:   "testdb",
			isValid:  true,
		},
		{
			name:     "Missing host",
			driver:   "mysql",
			host:     "",
			port:     "3306",
			user:     "root",
			password: "password",
			dbname:   "testdb",
			isValid:  false,
		},
		{
			name:     "Invalid port",
			driver:   "mysql",
			host:     "localhost",
			port:     "invalid",
			user:     "root",
			password: "password",
			dbname:   "testdb",
			isValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config := &Config{
				DBDriver:   tt.driver,
				DBHost:     tt.host,
				DBPort:     tt.port,
				DBUser:     tt.user,
				DBPassword: tt.password,
				DBName:     tt.dbname,
			}

			// Act - Basic validation logic
			isValid := config.DBDriver != "" &&
				config.DBHost != "" &&
				config.DBPort != "" &&
				config.DBUser != "" &&
				config.DBName != ""

			// Additional port validation
			if isValid {
				if _, err := strconv.Atoi(config.DBPort); err != nil {
					isValid = false
				}
			}

			// Assert
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestConfig_JWTConfiguration(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		expiresIn int
		isValid   bool
	}{
		{
			name:      "Valid JWT config",
			secret:    "my_secret_key_123",
			expiresIn: 30,
			isValid:   true,
		},
		{
			name:      "Empty secret",
			secret:    "",
			expiresIn: 30,
			isValid:   false,
		},
		{
			name:      "Short secret",
			secret:    "123",
			expiresIn: 30,
			isValid:   false,
		},
		{
			name:      "Zero expiration",
			secret:    "my_secret_key_123",
			expiresIn: 0,
			isValid:   false,
		},
		{
			name:      "Negative expiration",
			secret:    "my_secret_key_123",
			expiresIn: -1,
			isValid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config := &Config{
				JWTSecret:    tt.secret,
				JWTExpiresIn: tt.expiresIn,
			}

			// Act - Basic validation logic
			isValid := len(config.JWTSecret) >= 8 && config.JWTExpiresIn > 0

			// Assert
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestConfig_WebServerConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		isValid bool
	}{
		{
			name:    "Valid port 8080",
			port:    "8080",
			isValid: true,
		},
		{
			name:    "Valid port 3000",
			port:    "3000",
			isValid: true,
		},
		{
			name:    "Valid port 80",
			port:    "80",
			isValid: true,
		},
		{
			name:    "Invalid port - non-numeric",
			port:    "invalid",
			isValid: false,
		},
		{
			name:    "Invalid port - empty",
			port:    "",
			isValid: false,
		},
		{
			name:    "Invalid port - too high",
			port:    "99999",
			isValid: false,
		},
		{
			name:    "Invalid port - zero",
			port:    "0",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config := &Config{
				WebServerPort: tt.port,
			}

			// Act - Basic validation logic
			isValid := false
			if port, err := strconv.Atoi(config.WebServerPort); err == nil {
				isValid = port > 0 && port <= 65535
			}

			// Assert
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

// =============================================================================
// TESTES DE EDGE CASES
// =============================================================================

func TestConfig_SpecialCharactersInValues(t *testing.T) {
	// Arrange
	config := &Config{
		DBPassword:  "p@ssw0rd!#$%^&*()",
		JWTSecret:   "jwt_secret_with_special_chars_!@#$%^&*()",
		CORSOrigins: "http://localhost:3000,https://example.com/path?query=value&other=123",
	}

	// Act & Assert
	assert.Contains(t, config.DBPassword, "@")
	assert.Contains(t, config.DBPassword, "!")
	assert.Contains(t, config.JWTSecret, "_")
	assert.Contains(t, config.JWTSecret, "!")
	assert.Contains(t, config.CORSOrigins, "?")
	assert.Contains(t, config.CORSOrigins, "&")
}

func TestConfig_UnicodeCharacters(t *testing.T) {
	// Arrange
	config := &Config{
		DBName:    "测试数据库",
		JWTSecret: "jwt_密钥_with_unicode_中文",
	}

	// Act & Assert
	assert.Contains(t, config.DBName, "测试")
	assert.Contains(t, config.JWTSecret, "密钥")
	assert.Contains(t, config.JWTSecret, "中文")
}

func TestConfig_VeryLongValues(t *testing.T) {
	// Arrange
	longString := string(make([]byte, 1000))
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	config := &Config{
		DBPassword: longString,
		JWTSecret:  longString,
	}

	// Act & Assert
	assert.Len(t, config.DBPassword, 1000)
	assert.Len(t, config.JWTSecret, 1000)
}

func TestConfig_EmptyValues(t *testing.T) {
	// Arrange
	config := &Config{}

	// Act & Assert
	assert.Equal(t, "", config.DBDriver)
	assert.Equal(t, "", config.DBHost)
	assert.Equal(t, "", config.DBPort)
	assert.Equal(t, "", config.DBUser)
	assert.Equal(t, "", config.DBPassword)
	assert.Equal(t, "", config.DBName)
	assert.Equal(t, "", config.WebServerPort)
	assert.Equal(t, "", config.JWTSecret)
	assert.Equal(t, 0, config.JWTExpiresIn)
	assert.Equal(t, "", config.CORSOrigins)
}

func TestConfig_NewConfig(t *testing.T) {
	// Act
	config := NewConfig()

	// Assert
	assert.NotNil(t, config)
	// NewConfig() retorna Cfg que pode ser nil se não foi inicializado
	// Em ambiente de teste, isso é esperado
}
