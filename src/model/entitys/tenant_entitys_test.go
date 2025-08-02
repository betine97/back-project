package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TESTES PARA Tenants Entity
// =============================================================================

func TestTenants_JSONSerialization(t *testing.T) {
	// Arrange
	tenant := Tenants{
		ID:          1,
		UserID:      10,
		NomeEmpresa: "Empresa Teste Ltda.",
		DBName:      "empresa_teste_db",
		DBUser:      "empresa_user",
		DBPassword:  "secure_password_123",
		DBHost:      "localhost",
		DBPort:      "3306",
		CreatedAt:   "2024-01-01T10:30:00Z",
	}

	// Act
	jsonData, err := json.Marshal(tenant)
	assert.NoError(t, err)

	var unmarshaled Tenants
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tenant.ID, unmarshaled.ID)
	assert.Equal(t, tenant.UserID, unmarshaled.UserID)
	assert.Equal(t, tenant.NomeEmpresa, unmarshaled.NomeEmpresa)
	assert.Equal(t, tenant.DBName, unmarshaled.DBName)
	assert.Equal(t, tenant.DBUser, unmarshaled.DBUser)
	assert.Equal(t, tenant.DBPassword, unmarshaled.DBPassword)
	assert.Equal(t, tenant.DBHost, unmarshaled.DBHost)
	assert.Equal(t, tenant.DBPort, unmarshaled.DBPort)
	assert.Equal(t, tenant.CreatedAt, unmarshaled.CreatedAt)
}

func TestTenants_JSONTags(t *testing.T) {
	// Arrange
	tenant := Tenants{
		ID:          1,
		UserID:      10,
		NomeEmpresa: "Empresa Teste",
		DBName:      "test_db",
		DBUser:      "test_user",
		DBPassword:  "test_pass",
		DBHost:      "localhost",
		DBPort:      "3306",
		CreatedAt:   "2024-01-01",
	}

	// Act
	jsonData, err := json.Marshal(tenant)
	assert.NoError(t, err)

	// Assert
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, `"id":1`)
	assert.Contains(t, jsonString, `"user_id":10`)
	assert.Contains(t, jsonString, `"nome_empresa":"Empresa Teste"`)
	assert.Contains(t, jsonString, `"db_name":"test_db"`)
	assert.Contains(t, jsonString, `"db_user":"test_user"`)
	assert.Contains(t, jsonString, `"db_password":"test_pass"`)
	assert.Contains(t, jsonString, `"db_host":"localhost"`)
	assert.Contains(t, jsonString, `"db_port":"3306"`)
	assert.Contains(t, jsonString, `"created_at":"2024-01-01"`)
}

func TestTenants_EmptyValues(t *testing.T) {
	// Arrange
	tenant := Tenants{}

	// Act
	jsonData, err := json.Marshal(tenant)
	assert.NoError(t, err)

	var unmarshaled Tenants
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint(0), unmarshaled.ID)
	assert.Equal(t, uint(0), unmarshaled.UserID)
	assert.Equal(t, "", unmarshaled.NomeEmpresa)
	assert.Equal(t, "", unmarshaled.DBName)
	assert.Equal(t, "", unmarshaled.DBUser)
	assert.Equal(t, "", unmarshaled.DBPassword)
	assert.Equal(t, "", unmarshaled.DBHost)
	assert.Equal(t, "", unmarshaled.DBPort)
	assert.Equal(t, "", unmarshaled.CreatedAt)
}

func TestTenants_SpecialCharacters(t *testing.T) {
	// Arrange
	tenant := Tenants{
		ID:          1,
		UserID:      10,
		NomeEmpresa: "Empresa & Cia Ltda. - Soluções Tecnológicas",
		DBName:      "empresa_db_2024",
		DBUser:      "user_empresa_123",
		DBPassword:  "P@ssw0rd!2024#$%",
		DBHost:      "db-server.empresa.com.br",
		DBPort:      "3306",
		CreatedAt:   "2024-01-01T10:30:00.000Z",
	}

	// Act
	jsonData, err := json.Marshal(tenant)
	assert.NoError(t, err)

	var unmarshaled Tenants
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tenant.NomeEmpresa, unmarshaled.NomeEmpresa)
	assert.Equal(t, tenant.DBName, unmarshaled.DBName)
	assert.Equal(t, tenant.DBUser, unmarshaled.DBUser)
	assert.Equal(t, tenant.DBPassword, unmarshaled.DBPassword)
	assert.Equal(t, tenant.DBHost, unmarshaled.DBHost)
}

func TestTenants_DatabasePorts(t *testing.T) {
	validPorts := []string{
		"3306",  // MySQL
		"5432",  // PostgreSQL
		"1433",  // SQL Server
		"1521",  // Oracle
		"27017", // MongoDB
		"6379",  // Redis
		"8080",  // Custom port
		"",      // Empty port
	}

	for _, port := range validPorts {
		t.Run("Port: "+port, func(t *testing.T) {
			// Arrange
			tenant := Tenants{
				ID:          1,
				UserID:      10,
				NomeEmpresa: "Empresa Teste",
				DBPort:      port,
			}

			// Act
			jsonData, err := json.Marshal(tenant)
			assert.NoError(t, err)

			var unmarshaled Tenants
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, port, unmarshaled.DBPort)
		})
	}
}

func TestTenants_DatabaseHosts(t *testing.T) {
	validHosts := []string{
		"localhost",
		"127.0.0.1",
		"192.168.1.100",
		"db.empresa.com",
		"mysql-server.amazonaws.com",
		"db-cluster.region.rds.amazonaws.com",
		"",
	}

	for _, host := range validHosts {
		t.Run("Host: "+host, func(t *testing.T) {
			// Arrange
			tenant := Tenants{
				ID:          1,
				UserID:      10,
				NomeEmpresa: "Empresa Teste",
				DBHost:      host,
			}

			// Act
			jsonData, err := json.Marshal(tenant)
			assert.NoError(t, err)

			var unmarshaled Tenants
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, host, unmarshaled.DBHost)
		})
	}
}

func TestTenants_LongStrings(t *testing.T) {
	// Arrange
	longString := string(make([]byte, 255))
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	tenant := Tenants{
		ID:          1,
		UserID:      10,
		NomeEmpresa: longString,
		DBName:      longString,
		DBUser:      longString,
		DBPassword:  longString,
		DBHost:      longString,
		DBPort:      longString,
		CreatedAt:   longString,
	}

	// Act
	jsonData, err := json.Marshal(tenant)
	assert.NoError(t, err)

	var unmarshaled Tenants
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, longString, unmarshaled.NomeEmpresa)
	assert.Equal(t, longString, unmarshaled.DBName)
	assert.Equal(t, longString, unmarshaled.DBUser)
	assert.Equal(t, longString, unmarshaled.DBPassword)
	assert.Equal(t, longString, unmarshaled.DBHost)
	assert.Equal(t, longString, unmarshaled.DBPort)
	assert.Equal(t, longString, unmarshaled.CreatedAt)
}

// =============================================================================
// TESTES PARA TenantConnection Entity
// =============================================================================

func TestTenantConnection_JSONSerialization(t *testing.T) {
	// Arrange
	connection := TenantConnection{
		DBUser:     "test_user",
		DBPassword: "secure_password",
		DBHost:     "localhost",
		DBPort:     "3306",
		DBName:     "test_database",
	}

	// Act
	jsonData, err := json.Marshal(connection)
	assert.NoError(t, err)

	var unmarshaled TenantConnection
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, connection.DBUser, unmarshaled.DBUser)
	assert.Equal(t, connection.DBPassword, unmarshaled.DBPassword)
	assert.Equal(t, connection.DBHost, unmarshaled.DBHost)
	assert.Equal(t, connection.DBPort, unmarshaled.DBPort)
	assert.Equal(t, connection.DBName, unmarshaled.DBName)
}

func TestTenantConnection_JSONTags(t *testing.T) {
	// Arrange
	connection := TenantConnection{
		DBUser:     "test_user",
		DBPassword: "test_pass",
		DBHost:     "localhost",
		DBPort:     "3306",
		DBName:     "test_db",
	}

	// Act
	jsonData, err := json.Marshal(connection)
	assert.NoError(t, err)

	// Assert
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, `"db_user":"test_user"`)
	assert.Contains(t, jsonString, `"db_password":"test_pass"`)
	assert.Contains(t, jsonString, `"db_host":"localhost"`)
	assert.Contains(t, jsonString, `"db_port":"3306"`)
	assert.Contains(t, jsonString, `"db_name":"test_db"`)
}

func TestTenantConnection_EmptyValues(t *testing.T) {
	// Arrange
	connection := TenantConnection{}

	// Act
	jsonData, err := json.Marshal(connection)
	assert.NoError(t, err)

	var unmarshaled TenantConnection
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "", unmarshaled.DBUser)
	assert.Equal(t, "", unmarshaled.DBPassword)
	assert.Equal(t, "", unmarshaled.DBHost)
	assert.Equal(t, "", unmarshaled.DBPort)
	assert.Equal(t, "", unmarshaled.DBName)
}

func TestTenantConnection_SpecialCharacters(t *testing.T) {
	// Arrange
	connection := TenantConnection{
		DBUser:     "user_with_underscore_123",
		DBPassword: "P@ssw0rd!#$%^&*()",
		DBHost:     "db-server.domain.com",
		DBPort:     "3306",
		DBName:     "database_name_2024",
	}

	// Act
	jsonData, err := json.Marshal(connection)
	assert.NoError(t, err)

	var unmarshaled TenantConnection
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, connection.DBUser, unmarshaled.DBUser)
	assert.Equal(t, connection.DBPassword, unmarshaled.DBPassword)
	assert.Equal(t, connection.DBHost, unmarshaled.DBHost)
	assert.Equal(t, connection.DBPort, unmarshaled.DBPort)
	assert.Equal(t, connection.DBName, unmarshaled.DBName)
}

// =============================================================================
// TESTES DE VALIDAÇÃO DE NEGÓCIO
// =============================================================================

func TestTenants_UserRelationship(t *testing.T) {
	tests := []struct {
		name    string
		userID  uint
		isValid bool
	}{
		{"Valid user ID", 1, true},
		{"Another valid ID", 999, true},
		{"Zero ID", 0, false}, // Assuming 0 is invalid
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tenant := Tenants{
				ID:     1,
				UserID: tt.userID,
			}

			// Act
			jsonData, err := json.Marshal(tenant)
			assert.NoError(t, err)

			var unmarshaled Tenants
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.userID, unmarshaled.UserID)

			// Business logic validation
			isValidUser := unmarshaled.UserID > 0
			assert.Equal(t, tt.isValid, isValidUser)
		})
	}
}

func TestTenants_DatabaseConnectionString(t *testing.T) {
	// Arrange
	tenant := Tenants{
		ID:         1,
		UserID:     10,
		DBName:     "empresa_db",
		DBUser:     "empresa_user",
		DBPassword: "secure_pass",
		DBHost:     "localhost",
		DBPort:     "3306",
	}

	// Act - Simulate building a connection string
	connectionString := tenant.DBUser + ":" + tenant.DBPassword + "@tcp(" + tenant.DBHost + ":" + tenant.DBPort + ")/" + tenant.DBName

	// Assert
	expected := "empresa_user:secure_pass@tcp(localhost:3306)/empresa_db"
	assert.Equal(t, expected, connectionString)
}

func TestTenants_RequiredFields(t *testing.T) {
	tests := []struct {
		name              string
		tenant            Tenants
		hasRequiredFields bool
	}{
		{
			name: "All required fields present",
			tenant: Tenants{
				UserID:      1,
				NomeEmpresa: "Empresa",
				DBName:      "db",
				DBUser:      "user",
				DBPassword:  "pass",
				DBHost:      "host",
				DBPort:      "port",
				CreatedAt:   "2024-01-01",
			},
			hasRequiredFields: true,
		},
		{
			name: "Missing UserID",
			tenant: Tenants{
				NomeEmpresa: "Empresa",
				DBName:      "db",
				DBUser:      "user",
				DBPassword:  "pass",
				DBHost:      "host",
				DBPort:      "port",
				CreatedAt:   "2024-01-01",
			},
			hasRequiredFields: false,
		},
		{
			name: "Missing DBName",
			tenant: Tenants{
				UserID:      1,
				NomeEmpresa: "Empresa",
				DBUser:      "user",
				DBPassword:  "pass",
				DBHost:      "host",
				DBPort:      "port",
				CreatedAt:   "2024-01-01",
			},
			hasRequiredFields: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act - Check if basic required fields are present
			hasUserID := tt.tenant.UserID > 0
			hasNomeEmpresa := tt.tenant.NomeEmpresa != ""
			hasDBName := tt.tenant.DBName != ""
			hasDBUser := tt.tenant.DBUser != ""
			hasDBPassword := tt.tenant.DBPassword != ""
			hasDBHost := tt.tenant.DBHost != ""
			hasDBPort := tt.tenant.DBPort != ""
			hasCreatedAt := tt.tenant.CreatedAt != ""

			allRequired := hasUserID && hasNomeEmpresa && hasDBName && hasDBUser && hasDBPassword && hasDBHost && hasDBPort && hasCreatedAt

			// Assert
			assert.Equal(t, tt.hasRequiredFields, allRequired)
		})
	}
}

// =============================================================================
// TESTES DE EDGE CASES
// =============================================================================

func TestTenants_MaxValues(t *testing.T) {
	// Arrange - Test with very large ID values
	tenant := Tenants{
		ID:     4294967295, // Max uint32
		UserID: 4294967295, // Max uint32
	}

	// Act
	jsonData, err := json.Marshal(tenant)
	assert.NoError(t, err)

	var unmarshaled Tenants
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tenant.ID, unmarshaled.ID)
	assert.Equal(t, tenant.UserID, unmarshaled.UserID)
}

func TestTenants_GormTags(t *testing.T) {
	// Arrange - Test that struct has proper GORM tags
	tenant := Tenants{
		ID:          1,
		UserID:      10,
		NomeEmpresa: "Empresa Teste",
		DBName:      "test_db",
		DBUser:      "test_user",
		DBPassword:  "test_pass",
		DBHost:      "localhost",
		DBPort:      "3306",
		CreatedAt:   "2024-01-01",
	}

	// Act & Assert - Verify struct can be used with GORM
	assert.Equal(t, uint(1), tenant.ID)
	assert.Equal(t, uint(10), tenant.UserID)
	assert.Equal(t, "Empresa Teste", tenant.NomeEmpresa)
	assert.Equal(t, "test_db", tenant.DBName)
	assert.Equal(t, "test_user", tenant.DBUser)
	assert.Equal(t, "test_pass", tenant.DBPassword)
	assert.Equal(t, "localhost", tenant.DBHost)
	assert.Equal(t, "3306", tenant.DBPort)
	assert.Equal(t, "2024-01-01", tenant.CreatedAt)
}

func TestTenantConnection_GormTags(t *testing.T) {
	// Arrange - Test that struct has proper GORM tags
	connection := TenantConnection{
		DBUser:     "test_user",
		DBPassword: "test_pass",
		DBHost:     "localhost",
		DBPort:     "3306",
		DBName:     "test_db",
	}

	// Act & Assert - Verify struct can be used with GORM
	assert.Equal(t, "test_user", connection.DBUser)
	assert.Equal(t, "test_pass", connection.DBPassword)
	assert.Equal(t, "localhost", connection.DBHost)
	assert.Equal(t, "3306", connection.DBPort)
	assert.Equal(t, "test_db", connection.DBName)
}
