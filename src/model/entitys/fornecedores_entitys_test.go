package entity

import (
	"encoding/json"
	"testing"

	dtos "github.com/betine97/back-project.git/src/model/dtos"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TESTES PARA Fornecedores Entity
// =============================================================================

func TestFornecedores_JSONSerialization(t *testing.T) {
	// Arrange
	fornecedor := Fornecedores{
		ID:           1,
		DataCadastro: "2024-01-01",
		Nome:         "Fornecedor Teste",
		Telefone:     "11999999999",
		Email:        "fornecedor@teste.com",
		Cidade:       "São Paulo",
		Estado:       "SP",
		Status:       "Ativo",
	}

	// Act
	jsonData, err := json.Marshal(fornecedor)
	assert.NoError(t, err)

	var unmarshaled Fornecedores
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fornecedor.ID, unmarshaled.ID)
	assert.Equal(t, fornecedor.DataCadastro, unmarshaled.DataCadastro)
	assert.Equal(t, fornecedor.Nome, unmarshaled.Nome)
	assert.Equal(t, fornecedor.Telefone, unmarshaled.Telefone)
	assert.Equal(t, fornecedor.Email, unmarshaled.Email)
	assert.Equal(t, fornecedor.Cidade, unmarshaled.Cidade)
	assert.Equal(t, fornecedor.Estado, unmarshaled.Estado)
	assert.Equal(t, fornecedor.Status, unmarshaled.Status)
}

func TestFornecedores_JSONTags(t *testing.T) {
	// Arrange
	fornecedor := Fornecedores{
		ID:           1,
		DataCadastro: "2024-01-01",
		Nome:         "Fornecedor Teste",
		Telefone:     "11999999999",
		Email:        "fornecedor@teste.com",
		Status:       "Ativo",
	}

	// Act
	jsonData, err := json.Marshal(fornecedor)
	assert.NoError(t, err)

	// Assert
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, `"id_fornecedor":1`)
	assert.Contains(t, jsonString, `"data_cadastro":"2024-01-01"`)
	assert.Contains(t, jsonString, `"nome":"Fornecedor Teste"`)
	assert.Contains(t, jsonString, `"telefone":"11999999999"`)
	assert.Contains(t, jsonString, `"email":"fornecedor@teste.com"`)
	assert.Contains(t, jsonString, `"status":"Ativo"`)
}

func TestFornecedores_EmptyValues(t *testing.T) {
	// Arrange
	fornecedor := Fornecedores{}

	// Act
	jsonData, err := json.Marshal(fornecedor)
	assert.NoError(t, err)

	var unmarshaled Fornecedores
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, unmarshaled.ID)
	assert.Equal(t, "", unmarshaled.DataCadastro)
	assert.Equal(t, "", unmarshaled.Nome)
	assert.Equal(t, "", unmarshaled.Telefone)
	assert.Equal(t, "", unmarshaled.Email)
	assert.Equal(t, "", unmarshaled.Cidade)
	assert.Equal(t, "", unmarshaled.Estado)
	assert.Equal(t, "", unmarshaled.Status)
}

func TestFornecedores_SpecialCharacters(t *testing.T) {
	// Arrange
	fornecedor := Fornecedores{
		ID:           1,
		DataCadastro: "2024-01-01T10:30:00Z",
		Nome:         "Fornecedor & Cia Ltda.",
		Telefone:     "(11) 99999-9999",
		Email:        "contato@fornecedor-teste.com.br",
		Cidade:       "São José dos Campos",
		Estado:       "SP",
		Status:       "Ativo",
	}

	// Act
	jsonData, err := json.Marshal(fornecedor)
	assert.NoError(t, err)

	var unmarshaled Fornecedores
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fornecedor.Nome, unmarshaled.Nome)
	assert.Equal(t, fornecedor.Telefone, unmarshaled.Telefone)
	assert.Equal(t, fornecedor.Email, unmarshaled.Email)
	assert.Equal(t, fornecedor.Cidade, unmarshaled.Cidade)
}

func TestFornecedores_StatusValues(t *testing.T) {
	validStatuses := []string{"Ativo", "Inativo", "Pendente", "Bloqueado"}

	for _, status := range validStatuses {
		t.Run("Status: "+status, func(t *testing.T) {
			// Arrange
			fornecedor := Fornecedores{
				ID:     1,
				Nome:   "Fornecedor Teste",
				Status: status,
			}

			// Act
			jsonData, err := json.Marshal(fornecedor)
			assert.NoError(t, err)

			var unmarshaled Fornecedores
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, status, unmarshaled.Status)
		})
	}
}

// =============================================================================
// TESTES PARA BuildFornecedorEntity Function
// =============================================================================

func TestBuildFornecedorEntity_Success(t *testing.T) {
	// Arrange
	request := dtos.CreateFornecedorRequest{
		Nome:         "Novo Fornecedor",
		Telefone:     "11888888888",
		Email:        "novo@fornecedor.com",
		Cidade:       "Rio de Janeiro",
		Estado:       "RJ",
		Status:       "Ativo",
		DataCadastro: "2024-01-15",
	}

	// Act
	fornecedor := BuildFornecedorEntity(request)

	// Assert
	assert.NotNil(t, fornecedor)
	assert.Equal(t, request.Nome, fornecedor.Nome)
	assert.Equal(t, request.Telefone, fornecedor.Telefone)
	assert.Equal(t, request.Email, fornecedor.Email)
	assert.Equal(t, request.Cidade, fornecedor.Cidade)
	assert.Equal(t, request.Estado, fornecedor.Estado)
	assert.Equal(t, request.Status, fornecedor.Status)
	assert.Equal(t, request.DataCadastro, fornecedor.DataCadastro)
	// ID should be 0 (will be set by database)
	assert.Equal(t, 0, fornecedor.ID)
}

func TestBuildFornecedorEntity_EmptyRequest(t *testing.T) {
	// Arrange
	request := dtos.CreateFornecedorRequest{}

	// Act
	fornecedor := BuildFornecedorEntity(request)

	// Assert
	assert.NotNil(t, fornecedor)
	assert.Equal(t, "", fornecedor.Nome)
	assert.Equal(t, "", fornecedor.Telefone)
	assert.Equal(t, "", fornecedor.Email)
	assert.Equal(t, "", fornecedor.Cidade)
	assert.Equal(t, "", fornecedor.Estado)
	assert.Equal(t, "", fornecedor.Status)
	assert.Equal(t, "", fornecedor.DataCadastro)
	assert.Equal(t, 0, fornecedor.ID)
}

func TestBuildFornecedorEntity_SpecialCharacters(t *testing.T) {
	// Arrange
	request := dtos.CreateFornecedorRequest{
		Nome:         "Fornecedor & Cia Ltda.",
		Telefone:     "(11) 99999-9999",
		Email:        "contato@fornecedor-especial.com.br",
		Cidade:       "São José dos Campos",
		Estado:       "SP",
		Status:       "Ativo",
		DataCadastro: "2024-01-01T10:30:00Z",
	}

	// Act
	fornecedor := BuildFornecedorEntity(request)

	// Assert
	assert.NotNil(t, fornecedor)
	assert.Equal(t, request.Nome, fornecedor.Nome)
	assert.Equal(t, request.Telefone, fornecedor.Telefone)
	assert.Equal(t, request.Email, fornecedor.Email)
	assert.Equal(t, request.Cidade, fornecedor.Cidade)
	assert.Equal(t, request.Estado, fornecedor.Estado)
	assert.Equal(t, request.DataCadastro, fornecedor.DataCadastro)
}

func TestBuildFornecedorEntity_LongStrings(t *testing.T) {
	// Arrange
	longString := string(make([]byte, 255))
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	request := dtos.CreateFornecedorRequest{
		Nome:         longString,
		Telefone:     longString,
		Email:        "test@example.com",
		Cidade:       longString,
		Estado:       longString,
		Status:       longString,
		DataCadastro: longString,
	}

	// Act
	fornecedor := BuildFornecedorEntity(request)

	// Assert
	assert.NotNil(t, fornecedor)
	assert.Equal(t, longString, fornecedor.Nome)
	assert.Equal(t, longString, fornecedor.Telefone)
	assert.Equal(t, "test@example.com", fornecedor.Email)
	assert.Equal(t, longString, fornecedor.Cidade)
	assert.Equal(t, longString, fornecedor.Estado)
	assert.Equal(t, longString, fornecedor.Status)
	assert.Equal(t, longString, fornecedor.DataCadastro)
}

func TestBuildFornecedorEntity_ReturnType(t *testing.T) {
	// Arrange
	request := dtos.CreateFornecedorRequest{
		Nome:   "Fornecedor Teste",
		Status: "Ativo",
	}

	// Act
	fornecedor := BuildFornecedorEntity(request)

	// Assert
	assert.IsType(t, &Fornecedores{}, fornecedor)
	assert.NotNil(t, fornecedor)
}

// =============================================================================
// TESTES DE EDGE CASES
// =============================================================================

func TestFornecedores_NegativeID(t *testing.T) {
	// Arrange
	fornecedor := Fornecedores{
		ID:   -1,
		Nome: "Fornecedor Teste",
	}

	// Act
	jsonData, err := json.Marshal(fornecedor)
	assert.NoError(t, err)

	var unmarshaled Fornecedores
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, -1, unmarshaled.ID)
}

func TestFornecedores_EmailFormats(t *testing.T) {
	emailFormats := []string{
		"simple@example.com",
		"user.name@domain.co.uk",
		"user+tag@example.org",
		"123@example.com",
		"user@sub.domain.com",
		"invalid-email", // Should still be stored as-is
		"",              // Empty email
	}

	for _, email := range emailFormats {
		t.Run("Email: "+email, func(t *testing.T) {
			// Arrange
			fornecedor := Fornecedores{
				ID:    1,
				Nome:  "Fornecedor Teste",
				Email: email,
			}

			// Act
			jsonData, err := json.Marshal(fornecedor)
			assert.NoError(t, err)

			var unmarshaled Fornecedores
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, email, unmarshaled.Email)
		})
	}
}

func TestFornecedores_PhoneFormats(t *testing.T) {
	phoneFormats := []string{
		"11999999999",
		"(11) 99999-9999",
		"+55 11 99999-9999",
		"11 9 9999-9999",
		"119999999999",
		"",
		"invalid-phone",
	}

	for _, phone := range phoneFormats {
		t.Run("Phone: "+phone, func(t *testing.T) {
			// Arrange
			fornecedor := Fornecedores{
				ID:       1,
				Nome:     "Fornecedor Teste",
				Telefone: phone,
			}

			// Act
			jsonData, err := json.Marshal(fornecedor)
			assert.NoError(t, err)

			var unmarshaled Fornecedores
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, phone, unmarshaled.Telefone)
		})
	}
}

func TestBuildFornecedorEntity_NilSafety(t *testing.T) {
	// Arrange - Test that function doesn't panic with zero values
	var request dtos.CreateFornecedorRequest

	// Act
	fornecedor := BuildFornecedorEntity(request)

	// Assert
	assert.NotNil(t, fornecedor)
	assert.IsType(t, &Fornecedores{}, fornecedor)
}
