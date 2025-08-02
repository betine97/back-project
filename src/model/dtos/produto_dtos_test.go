package dtos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TESTES PARA ProductResponse DTO
// =============================================================================

func TestProductResponse_JSONSerialization(t *testing.T) {
	// Arrange
	response := ProductResponse{
		ID:            1,
		CodigoBarra:   "1234567890123",
		NomeProduto:   "Produto Teste",
		SKU:           "SKU001",
		Categoria:     "Eletrônicos",
		DestinadoPara: "Consumidor Final",
		Variacao:      "Azul",
		Marca:         "Marca Teste",
		Descricao:     "Descrição do produto teste",
		Status:        "Ativo",
		PrecoVenda:    99.99,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled ProductResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, response.ID, unmarshaled.ID)
	assert.Equal(t, response.CodigoBarra, unmarshaled.CodigoBarra)
	assert.Equal(t, response.NomeProduto, unmarshaled.NomeProduto)
	assert.Equal(t, response.SKU, unmarshaled.SKU)
	assert.Equal(t, response.Categoria, unmarshaled.Categoria)
	assert.Equal(t, response.DestinadoPara, unmarshaled.DestinadoPara)
	assert.Equal(t, response.Variacao, unmarshaled.Variacao)
	assert.Equal(t, response.Marca, unmarshaled.Marca)
	assert.Equal(t, response.Descricao, unmarshaled.Descricao)
	assert.Equal(t, response.Status, unmarshaled.Status)
	assert.Equal(t, response.PrecoVenda, unmarshaled.PrecoVenda)
}

func TestProductResponse_JSONTags(t *testing.T) {
	// Arrange
	response := ProductResponse{
		ID:          1,
		CodigoBarra: "1234567890123",
		NomeProduto: "Produto Teste",
		SKU:         "SKU001",
		PrecoVenda:  99.99,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	// Assert
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, `"id_produto":1`)
	assert.Contains(t, jsonString, `"codigo_barra":"1234567890123"`)
	assert.Contains(t, jsonString, `"nome_produto":"Produto Teste"`)
	assert.Contains(t, jsonString, `"sku":"SKU001"`)
	assert.Contains(t, jsonString, `"preco_venda":99.99`)
}

func TestProductResponse_PriceFormats(t *testing.T) {
	tests := []struct {
		name       string
		precoVenda float64
		expected   float64
	}{
		{"Integer price", 100.0, 100.0},
		{"Two decimal places", 99.99, 99.99},
		{"One decimal place", 50.5, 50.5},
		{"Zero price", 0.0, 0.0},
		{"High precision", 123.456789, 123.456789},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			response := ProductResponse{
				ID:         1,
				PrecoVenda: tt.precoVenda,
			}

			// Act
			jsonData, err := json.Marshal(response)
			assert.NoError(t, err)

			var unmarshaled ProductResponse
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, unmarshaled.PrecoVenda)
		})
	}
}

// =============================================================================
// TESTES PARA ProductListResponse DTO
// =============================================================================

func TestProductListResponse_JSONSerialization(t *testing.T) {
	// Arrange
	products := []ProductResponse{
		{
			ID:          1,
			NomeProduto: "Produto 1",
			SKU:         "SKU001",
			PrecoVenda:  10.50,
			Status:      "Ativo",
		},
		{
			ID:          2,
			NomeProduto: "Produto 2",
			SKU:         "SKU002",
			PrecoVenda:  25.99,
			Status:      "Ativo",
		},
	}

	response := ProductListResponse{
		Products: products,
		Total:    2,
		Page:     1,
		Limit:    10,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled ProductListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, response.Total, unmarshaled.Total)
	assert.Equal(t, response.Page, unmarshaled.Page)
	assert.Equal(t, response.Limit, unmarshaled.Limit)
	assert.Len(t, unmarshaled.Products, 2)
	assert.Equal(t, response.Products[0].ID, unmarshaled.Products[0].ID)
	assert.Equal(t, response.Products[0].NomeProduto, unmarshaled.Products[0].NomeProduto)
	assert.Equal(t, response.Products[1].ID, unmarshaled.Products[1].ID)
	assert.Equal(t, response.Products[1].NomeProduto, unmarshaled.Products[1].NomeProduto)
}

func TestProductListResponse_EmptyList(t *testing.T) {
	// Arrange
	response := ProductListResponse{
		Products: []ProductResponse{},
		Total:    0,
		Page:     1,
		Limit:    10,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled ProductListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, unmarshaled.Total)
	assert.Equal(t, 1, unmarshaled.Page)
	assert.Equal(t, 10, unmarshaled.Limit)
	assert.Len(t, unmarshaled.Products, 0)
	assert.NotNil(t, unmarshaled.Products) // Deve ser slice vazio, não nil
}

func TestProductListResponse_Pagination(t *testing.T) {
	tests := []struct {
		name  string
		page  int
		limit int
		total int
	}{
		{"First page", 1, 10, 25},
		{"Middle page", 3, 5, 25},
		{"Last page", 5, 5, 25},
		{"Single item per page", 1, 1, 1},
		{"Large page size", 1, 100, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			response := ProductListResponse{
				Products: []ProductResponse{},
				Total:    tt.total,
				Page:     tt.page,
				Limit:    tt.limit,
			}

			// Act
			jsonData, err := json.Marshal(response)
			assert.NoError(t, err)

			var unmarshaled ProductListResponse
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.total, unmarshaled.Total)
			assert.Equal(t, tt.page, unmarshaled.Page)
			assert.Equal(t, tt.limit, unmarshaled.Limit)
		})
	}
}

// =============================================================================
// TESTES PARA CreateProductRequest DTO
// =============================================================================

func TestCreateProductRequest_JSONSerialization(t *testing.T) {
	// Arrange
	request := CreateProductRequest{
		DataCadastro:  "2024-01-01",
		CodigoBarra:   "1234567890123",
		NomeProduto:   "Novo Produto",
		SKU:           "SKU003",
		Categoria:     "Eletrônicos",
		DestinadoPara: "Consumidor Final",
		Variacao:      "Verde",
		Marca:         "Marca Nova",
		Descricao:     "Descrição do novo produto",
		Status:        "Ativo",
		PrecoVenda:    149.99,
	}

	// Act
	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)

	var unmarshaled CreateProductRequest
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, request.DataCadastro, unmarshaled.DataCadastro)
	assert.Equal(t, request.CodigoBarra, unmarshaled.CodigoBarra)
	assert.Equal(t, request.NomeProduto, unmarshaled.NomeProduto)
	assert.Equal(t, request.SKU, unmarshaled.SKU)
	assert.Equal(t, request.Categoria, unmarshaled.Categoria)
	assert.Equal(t, request.DestinadoPara, unmarshaled.DestinadoPara)
	assert.Equal(t, request.Variacao, unmarshaled.Variacao)
	assert.Equal(t, request.Marca, unmarshaled.Marca)
	assert.Equal(t, request.Descricao, unmarshaled.Descricao)
	assert.Equal(t, request.Status, unmarshaled.Status)
	assert.Equal(t, request.PrecoVenda, unmarshaled.PrecoVenda)
}

func TestCreateProductRequest_EmptyValues(t *testing.T) {
	// Arrange
	request := CreateProductRequest{}

	// Act
	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)

	var unmarshaled CreateProductRequest
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "", unmarshaled.DataCadastro)
	assert.Equal(t, "", unmarshaled.CodigoBarra)
	assert.Equal(t, "", unmarshaled.NomeProduto)
	assert.Equal(t, "", unmarshaled.SKU)
	assert.Equal(t, "", unmarshaled.Categoria)
	assert.Equal(t, "", unmarshaled.DestinadoPara)
	assert.Equal(t, "", unmarshaled.Variacao)
	assert.Equal(t, "", unmarshaled.Marca)
	assert.Equal(t, "", unmarshaled.Descricao)
	assert.Equal(t, "", unmarshaled.Status)
	assert.Equal(t, float64(0), unmarshaled.PrecoVenda)
}

func TestCreateProductRequest_SpecialCharacters(t *testing.T) {
	// Arrange
	request := CreateProductRequest{
		NomeProduto:   "Produto com Acentos & Símbolos™",
		SKU:           "SKU-001/A",
		Categoria:     "Eletrônicos & Informática",
		DestinadoPara: "B2B & B2C",
		Variacao:      "Azul/Verde",
		Marca:         "Marca® Ltda.",
		Descricao:     "Descrição com 'aspas' e \"aspas duplas\" e símbolos: @#$%",
		Status:        "Ativo",
		PrecoVenda:    1299.99,
	}

	// Act
	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)

	var unmarshaled CreateProductRequest
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, request.NomeProduto, unmarshaled.NomeProduto)
	assert.Equal(t, request.SKU, unmarshaled.SKU)
	assert.Equal(t, request.Categoria, unmarshaled.Categoria)
	assert.Equal(t, request.Descricao, unmarshaled.Descricao)
}

// =============================================================================
// TESTES DE CONVERSÃO ENTRE DTOs
// =============================================================================

func TestCreateProductRequest_ToProductResponse_Conversion(t *testing.T) {
	// Arrange
	request := CreateProductRequest{
		DataCadastro:  "2024-01-01",
		CodigoBarra:   "1234567890123",
		NomeProduto:   "Produto Teste",
		SKU:           "SKU001",
		Categoria:     "Eletrônicos",
		DestinadoPara: "Consumidor Final",
		Variacao:      "Azul",
		Marca:         "Marca Teste",
		Descricao:     "Descrição do produto",
		Status:        "Ativo",
		PrecoVenda:    99.99,
	}

	// Act - Simular conversão que seria feita no service
	response := ProductResponse{
		ID:            1, // Seria gerado pelo banco
		CodigoBarra:   request.CodigoBarra,
		NomeProduto:   request.NomeProduto,
		SKU:           request.SKU,
		Categoria:     request.Categoria,
		DestinadoPara: request.DestinadoPara,
		Variacao:      request.Variacao,
		Marca:         request.Marca,
		Descricao:     request.Descricao,
		Status:        request.Status,
		PrecoVenda:    request.PrecoVenda,
	}

	// Assert
	assert.Equal(t, request.CodigoBarra, response.CodigoBarra)
	assert.Equal(t, request.NomeProduto, response.NomeProduto)
	assert.Equal(t, request.SKU, response.SKU)
	assert.Equal(t, request.Categoria, response.Categoria)
	assert.Equal(t, request.DestinadoPara, response.DestinadoPara)
	assert.Equal(t, request.Variacao, response.Variacao)
	assert.Equal(t, request.Marca, response.Marca)
	assert.Equal(t, request.Descricao, response.Descricao)
	assert.Equal(t, request.Status, response.Status)
	assert.Equal(t, request.PrecoVenda, response.PrecoVenda)
	assert.Equal(t, 1, response.ID)
}

// =============================================================================
// TESTES DE EDGE CASES
// =============================================================================

func TestProductResponse_NegativePrice(t *testing.T) {
	// Arrange
	response := ProductResponse{
		ID:         1,
		PrecoVenda: -10.50,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled ProductResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, -10.50, unmarshaled.PrecoVenda)
}

func TestProductListResponse_NilSlice(t *testing.T) {
	// Arrange
	response := ProductListResponse{
		Products: nil,
		Total:    0,
		Page:     1,
		Limit:    10,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled ProductListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, unmarshaled.Total)
	assert.Equal(t, 1, unmarshaled.Page)
	assert.Equal(t, 10, unmarshaled.Limit)
	assert.Nil(t, unmarshaled.Products) // nil slice é preservado
}

func TestCreateProductRequest_LongStrings(t *testing.T) {
	// Arrange
	longString := string(make([]byte, 1000))
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	request := CreateProductRequest{
		NomeProduto: longString,
		Descricao:   longString,
		SKU:         longString,
	}

	// Act
	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)

	var unmarshaled CreateProductRequest
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, longString, unmarshaled.NomeProduto)
	assert.Equal(t, longString, unmarshaled.Descricao)
	assert.Equal(t, longString, unmarshaled.SKU)
}

func TestProductListResponse_LargeList(t *testing.T) {
	// Arrange
	products := make([]ProductResponse, 1000)
	for i := 0; i < 1000; i++ {
		products[i] = ProductResponse{
			ID:          i + 1,
			NomeProduto: "Produto " + string(rune(i+1)),
			SKU:         "SKU" + string(rune(i+1)),
			PrecoVenda:  float64(i+1) * 10.50,
		}
	}

	response := ProductListResponse{
		Products: products,
		Total:    1000,
		Page:     1,
		Limit:    1000,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled ProductListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1000, unmarshaled.Total)
	assert.Len(t, unmarshaled.Products, 1000)
	assert.Equal(t, 1, unmarshaled.Products[0].ID)
	assert.Equal(t, 1000, unmarshaled.Products[999].ID)
}
