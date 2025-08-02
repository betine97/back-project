package dtos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TESTES PARA CreateFornecedorRequest DTO
// =============================================================================

func TestCreateFornecedorRequest_JSONSerialization(t *testing.T) {
	// Arrange
	request := CreateFornecedorRequest{
		Nome:         "Fornecedor Teste",
		Telefone:     "11999999999",
		Email:        "fornecedor@teste.com",
		Cidade:       "São Paulo",
		Estado:       "SP",
		Status:       "Ativo",
		DataCadastro: "2024-01-01",
	}

	// Act
	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)

	var unmarshaled CreateFornecedorRequest
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, request.Nome, unmarshaled.Nome)
	assert.Equal(t, request.Telefone, unmarshaled.Telefone)
	assert.Equal(t, request.Email, unmarshaled.Email)
	assert.Equal(t, request.Cidade, unmarshaled.Cidade)
	assert.Equal(t, request.Estado, unmarshaled.Estado)
	assert.Equal(t, request.Status, unmarshaled.Status)
	assert.Equal(t, request.DataCadastro, unmarshaled.DataCadastro)
}

func TestCreateFornecedorRequest_EmptyValues(t *testing.T) {
	// Arrange
	request := CreateFornecedorRequest{}

	// Act
	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)

	var unmarshaled CreateFornecedorRequest
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "", unmarshaled.Nome)
	assert.Equal(t, "", unmarshaled.Telefone)
	assert.Equal(t, "", unmarshaled.Email)
	assert.Equal(t, "", unmarshaled.Cidade)
	assert.Equal(t, "", unmarshaled.Estado)
	assert.Equal(t, "", unmarshaled.Status)
	assert.Equal(t, "", unmarshaled.DataCadastro)
}

func TestCreateFornecedorRequest_SpecialCharacters(t *testing.T) {
	// Arrange
	request := CreateFornecedorRequest{
		Nome:         "Fornecedor & Cia Ltda.",
		Telefone:     "(11) 99999-9999",
		Email:        "contato@fornecedor-teste.com.br",
		Cidade:       "São José dos Campos",
		Estado:       "SP",
		Status:       "Ativo",
		DataCadastro: "2024-01-01T10:30:00Z",
	}

	// Act
	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)

	var unmarshaled CreateFornecedorRequest
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, request.Nome, unmarshaled.Nome)
	assert.Equal(t, request.Telefone, unmarshaled.Telefone)
	assert.Equal(t, request.Email, unmarshaled.Email)
	assert.Equal(t, request.Cidade, unmarshaled.Cidade)
}

// =============================================================================
// TESTES PARA UpdateFornecedorRequest DTO
// =============================================================================

func TestUpdateFornecedorRequest_JSONSerialization(t *testing.T) {
	// Arrange
	request := UpdateFornecedorRequest{
		Campo: "nome",
		Valor: "Novo Nome do Fornecedor",
	}

	// Act
	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)

	var unmarshaled UpdateFornecedorRequest
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, request.Campo, unmarshaled.Campo)
	assert.Equal(t, request.Valor, unmarshaled.Valor)
}

func TestUpdateFornecedorRequest_ValidFields(t *testing.T) {
	validFields := []struct {
		campo string
		valor string
	}{
		{"nome", "Novo Nome"},
		{"telefone", "11888888888"},
		{"email", "novo@email.com"},
		{"cidade", "Nova Cidade"},
		{"estado", "RJ"},
		{"status", "Inativo"},
	}

	for _, field := range validFields {
		t.Run("Valid field: "+field.campo, func(t *testing.T) {
			// Arrange
			request := UpdateFornecedorRequest{
				Campo: field.campo,
				Valor: field.valor,
			}

			// Act
			jsonData, err := json.Marshal(request)
			assert.NoError(t, err)

			var unmarshaled UpdateFornecedorRequest
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, field.campo, unmarshaled.Campo)
			assert.Equal(t, field.valor, unmarshaled.Valor)
		})
	}
}

// =============================================================================
// TESTES PARA FornecedorResponse DTO
// =============================================================================

func TestFornecedorResponse_JSONSerialization(t *testing.T) {
	// Arrange
	response := FornecedorResponse{
		ID:           1,
		Nome:         "Fornecedor Teste",
		Telefone:     "11999999999",
		Email:        "fornecedor@teste.com",
		Cidade:       "São Paulo",
		Estado:       "SP",
		Status:       "Ativo",
		DataCadastro: "2024-01-01",
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled FornecedorResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, response.ID, unmarshaled.ID)
	assert.Equal(t, response.Nome, unmarshaled.Nome)
	assert.Equal(t, response.Telefone, unmarshaled.Telefone)
	assert.Equal(t, response.Email, unmarshaled.Email)
	assert.Equal(t, response.Cidade, unmarshaled.Cidade)
	assert.Equal(t, response.Estado, unmarshaled.Estado)
	assert.Equal(t, response.Status, unmarshaled.Status)
	assert.Equal(t, response.DataCadastro, unmarshaled.DataCadastro)
}

func TestFornecedorResponse_JSONTags(t *testing.T) {
	// Arrange
	response := FornecedorResponse{
		ID:           1,
		Nome:         "Fornecedor Teste",
		Telefone:     "11999999999",
		Email:        "fornecedor@teste.com",
		Cidade:       "São Paulo",
		Estado:       "SP",
		Status:       "Ativo",
		DataCadastro: "2024-01-01",
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	// Assert
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, `"id_fornecedor":1`)
	assert.Contains(t, jsonString, `"nome":"Fornecedor Teste"`)
	assert.Contains(t, jsonString, `"telefone":"11999999999"`)
	assert.Contains(t, jsonString, `"email":"fornecedor@teste.com"`)
	assert.Contains(t, jsonString, `"data_cadastro":"2024-01-01"`)
}

// =============================================================================
// TESTES PARA FornecedorListResponse DTO
// =============================================================================

func TestFornecedorListResponse_JSONSerialization(t *testing.T) {
	// Arrange
	fornecedores := []FornecedorResponse{
		{
			ID:           1,
			Nome:         "Fornecedor 1",
			Telefone:     "11999999999",
			Email:        "fornecedor1@teste.com",
			Cidade:       "São Paulo",
			Estado:       "SP",
			Status:       "Ativo",
			DataCadastro: "2024-01-01",
		},
		{
			ID:           2,
			Nome:         "Fornecedor 2",
			Telefone:     "11888888888",
			Email:        "fornecedor2@teste.com",
			Cidade:       "Rio de Janeiro",
			Estado:       "RJ",
			Status:       "Inativo",
			DataCadastro: "2024-01-02",
		},
	}

	response := FornecedorListResponse{
		Fornecedores: fornecedores,
		Total:        2,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled FornecedorListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, response.Total, unmarshaled.Total)
	assert.Len(t, unmarshaled.Fornecedores, 2)
	assert.Equal(t, response.Fornecedores[0].ID, unmarshaled.Fornecedores[0].ID)
	assert.Equal(t, response.Fornecedores[0].Nome, unmarshaled.Fornecedores[0].Nome)
	assert.Equal(t, response.Fornecedores[1].ID, unmarshaled.Fornecedores[1].ID)
	assert.Equal(t, response.Fornecedores[1].Nome, unmarshaled.Fornecedores[1].Nome)
}

func TestFornecedorListResponse_EmptyList(t *testing.T) {
	// Arrange
	response := FornecedorListResponse{
		Fornecedores: []FornecedorResponse{},
		Total:        0,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled FornecedorListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, unmarshaled.Total)
	assert.Len(t, unmarshaled.Fornecedores, 0)
	assert.NotNil(t, unmarshaled.Fornecedores) // Deve ser slice vazio, não nil
}

func TestFornecedorListResponse_LargeList(t *testing.T) {
	// Arrange
	fornecedores := make([]FornecedorResponse, 100)
	for i := 0; i < 100; i++ {
		fornecedores[i] = FornecedorResponse{
			ID:           i + 1,
			Nome:         "Fornecedor " + string(rune(i+1)),
			Telefone:     "11999999999",
			Email:        "fornecedor" + string(rune(i+1)) + "@teste.com",
			Cidade:       "São Paulo",
			Estado:       "SP",
			Status:       "Ativo",
			DataCadastro: "2024-01-01",
		}
	}

	response := FornecedorListResponse{
		Fornecedores: fornecedores,
		Total:        100,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled FornecedorListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 100, unmarshaled.Total)
	assert.Len(t, unmarshaled.Fornecedores, 100)
	assert.Equal(t, 1, unmarshaled.Fornecedores[0].ID)
	assert.Equal(t, 100, unmarshaled.Fornecedores[99].ID)
}

// =============================================================================
// TESTES DE CONVERSÃO ENTRE DTOs
// =============================================================================

func TestCreateFornecedorRequest_ToFornecedorResponse_Conversion(t *testing.T) {
	// Arrange
	request := CreateFornecedorRequest{
		Nome:         "Fornecedor Teste",
		Telefone:     "11999999999",
		Email:        "fornecedor@teste.com",
		Cidade:       "São Paulo",
		Estado:       "SP",
		Status:       "Ativo",
		DataCadastro: "2024-01-01",
	}

	// Act - Simular conversão que seria feita no service
	response := FornecedorResponse{
		ID:           1, // Seria gerado pelo banco
		Nome:         request.Nome,
		Telefone:     request.Telefone,
		Email:        request.Email,
		Cidade:       request.Cidade,
		Estado:       request.Estado,
		Status:       request.Status,
		DataCadastro: request.DataCadastro,
	}

	// Assert
	assert.Equal(t, request.Nome, response.Nome)
	assert.Equal(t, request.Telefone, response.Telefone)
	assert.Equal(t, request.Email, response.Email)
	assert.Equal(t, request.Cidade, response.Cidade)
	assert.Equal(t, request.Estado, response.Estado)
	assert.Equal(t, request.Status, response.Status)
	assert.Equal(t, request.DataCadastro, response.DataCadastro)
	assert.Equal(t, 1, response.ID)
}

// =============================================================================
// TESTES DE EDGE CASES
// =============================================================================

func TestFornecedorResponse_NilValues(t *testing.T) {
	// Arrange
	response := FornecedorResponse{}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled FornecedorResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, unmarshaled.ID)
	assert.Equal(t, "", unmarshaled.Nome)
	assert.Equal(t, "", unmarshaled.Telefone)
	assert.Equal(t, "", unmarshaled.Email)
}

func TestUpdateFornecedorRequest_EmptyValues(t *testing.T) {
	// Arrange
	request := UpdateFornecedorRequest{}

	// Act
	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)

	var unmarshaled UpdateFornecedorRequest
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "", unmarshaled.Campo)
	assert.Equal(t, "", unmarshaled.Valor)
}

func TestFornecedorListResponse_NilSlice(t *testing.T) {
	// Arrange
	response := FornecedorListResponse{
		Fornecedores: nil,
		Total:        0,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled FornecedorListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, unmarshaled.Total)
	assert.Nil(t, unmarshaled.Fornecedores) // nil slice é preservado
}
