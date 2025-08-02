package entity

import (
	"encoding/json"
	"testing"

	dtos "github.com/betine97/back-project.git/src/model/dtos"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TESTES PARA Produto Entity
// =============================================================================

func TestProduto_JSONSerialization(t *testing.T) {
	// Arrange
	produto := Produto{
		IDProduto:     1,
		DataCadastro:  "2024-01-01",
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
	jsonData, err := json.Marshal(produto)
	assert.NoError(t, err)

	var unmarshaled Produto
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, produto.IDProduto, unmarshaled.IDProduto)
	assert.Equal(t, produto.DataCadastro, unmarshaled.DataCadastro)
	assert.Equal(t, produto.CodigoBarra, unmarshaled.CodigoBarra)
	assert.Equal(t, produto.NomeProduto, unmarshaled.NomeProduto)
	assert.Equal(t, produto.SKU, unmarshaled.SKU)
	assert.Equal(t, produto.Categoria, unmarshaled.Categoria)
	assert.Equal(t, produto.DestinadoPara, unmarshaled.DestinadoPara)
	assert.Equal(t, produto.Variacao, unmarshaled.Variacao)
	assert.Equal(t, produto.Marca, unmarshaled.Marca)
	assert.Equal(t, produto.Descricao, unmarshaled.Descricao)
	assert.Equal(t, produto.Status, unmarshaled.Status)
	assert.Equal(t, produto.PrecoVenda, unmarshaled.PrecoVenda)
}

func TestProduto_JSONTags(t *testing.T) {
	// Arrange
	produto := Produto{
		IDProduto:   1,
		CodigoBarra: "1234567890123",
		NomeProduto: "Produto Teste",
		SKU:         "SKU001",
		PrecoVenda:  99.99,
	}

	// Act
	jsonData, err := json.Marshal(produto)
	assert.NoError(t, err)

	// Assert
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, `"id_produto":1`)
	assert.Contains(t, jsonString, `"codigo_barra":"1234567890123"`)
	assert.Contains(t, jsonString, `"nome_produto":"Produto Teste"`)
	assert.Contains(t, jsonString, `"sku":"SKU001"`)
	assert.Contains(t, jsonString, `"preco_venda":99.99`)
}

func TestProduto_EmptyValues(t *testing.T) {
	// Arrange
	produto := Produto{}

	// Act
	jsonData, err := json.Marshal(produto)
	assert.NoError(t, err)

	var unmarshaled Produto
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, unmarshaled.IDProduto)
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

func TestProduto_PriceFormats(t *testing.T) {
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
		{"Negative price", -10.50, -10.50},
		{"Large price", 999999.99, 999999.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			produto := Produto{
				IDProduto:  1,
				PrecoVenda: tt.precoVenda,
			}

			// Act
			jsonData, err := json.Marshal(produto)
			assert.NoError(t, err)

			var unmarshaled Produto
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, unmarshaled.PrecoVenda)
		})
	}
}

func TestProduto_StatusValues(t *testing.T) {
	validStatuses := []string{"Ativo", "Inativo", "Descontinuado", "Em Falta", "Promoção"}

	for _, status := range validStatuses {
		t.Run("Status: "+status, func(t *testing.T) {
			// Arrange
			produto := Produto{
				IDProduto:   1,
				NomeProduto: "Produto Teste",
				Status:      status,
			}

			// Act
			jsonData, err := json.Marshal(produto)
			assert.NoError(t, err)

			var unmarshaled Produto
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, status, unmarshaled.Status)
		})
	}
}

func TestProduto_SpecialCharacters(t *testing.T) {
	// Arrange
	produto := Produto{
		IDProduto:     1,
		DataCadastro:  "2024-01-01T10:30:00Z",
		CodigoBarra:   "789-0123456789-0",
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
	jsonData, err := json.Marshal(produto)
	assert.NoError(t, err)

	var unmarshaled Produto
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, produto.NomeProduto, unmarshaled.NomeProduto)
	assert.Equal(t, produto.SKU, unmarshaled.SKU)
	assert.Equal(t, produto.Categoria, unmarshaled.Categoria)
	assert.Equal(t, produto.Descricao, unmarshaled.Descricao)
	assert.Equal(t, produto.Marca, unmarshaled.Marca)
}

// =============================================================================
// TESTES PARA BuildProductEntity Function
// =============================================================================

func TestBuildProductEntity_Success(t *testing.T) {
	// Arrange
	request := dtos.CreateProductRequest{
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
	produto := BuildProductEntity(request)

	// Assert
	assert.NotNil(t, produto)
	assert.Equal(t, request.DataCadastro, produto.DataCadastro)
	assert.Equal(t, request.CodigoBarra, produto.CodigoBarra)
	assert.Equal(t, request.NomeProduto, produto.NomeProduto)
	assert.Equal(t, request.SKU, produto.SKU)
	assert.Equal(t, request.Categoria, produto.Categoria)
	assert.Equal(t, request.DestinadoPara, produto.DestinadoPara)
	assert.Equal(t, request.Variacao, produto.Variacao)
	assert.Equal(t, request.Marca, produto.Marca)
	assert.Equal(t, request.Descricao, produto.Descricao)
	assert.Equal(t, request.Status, produto.Status)
	assert.Equal(t, request.PrecoVenda, produto.PrecoVenda)
	// IDProduto should be 0 (will be set by database)
	assert.Equal(t, 0, produto.IDProduto)
}

func TestBuildProductEntity_EmptyRequest(t *testing.T) {
	// Arrange
	request := dtos.CreateProductRequest{}

	// Act
	produto := BuildProductEntity(request)

	// Assert
	assert.NotNil(t, produto)
	assert.Equal(t, "", produto.DataCadastro)
	assert.Equal(t, "", produto.CodigoBarra)
	assert.Equal(t, "", produto.NomeProduto)
	assert.Equal(t, "", produto.SKU)
	assert.Equal(t, "", produto.Categoria)
	assert.Equal(t, "", produto.DestinadoPara)
	assert.Equal(t, "", produto.Variacao)
	assert.Equal(t, "", produto.Marca)
	assert.Equal(t, "", produto.Descricao)
	assert.Equal(t, "", produto.Status)
	assert.Equal(t, float64(0), produto.PrecoVenda)
	assert.Equal(t, 0, produto.IDProduto)
}

func TestBuildProductEntity_SpecialCharacters(t *testing.T) {
	// Arrange
	request := dtos.CreateProductRequest{
		DataCadastro:  "2024-01-01T10:30:00Z",
		CodigoBarra:   "789-0123456789-0",
		NomeProduto:   "Produto com Acentos & Símbolos™",
		SKU:           "SKU-001/A",
		Categoria:     "Eletrônicos & Informática",
		DestinadoPara: "B2B & B2C",
		Variacao:      "Azul/Verde",
		Marca:         "Marca® Ltda.",
		Descricao:     "Descrição com 'aspas' e \"aspas duplas\"",
		Status:        "Ativo",
		PrecoVenda:    1299.99,
	}

	// Act
	produto := BuildProductEntity(request)

	// Assert
	assert.NotNil(t, produto)
	assert.Equal(t, request.NomeProduto, produto.NomeProduto)
	assert.Equal(t, request.SKU, produto.SKU)
	assert.Equal(t, request.Categoria, produto.Categoria)
	assert.Equal(t, request.Descricao, produto.Descricao)
	assert.Equal(t, request.Marca, produto.Marca)
}

func TestBuildProductEntity_PriceValues(t *testing.T) {
	tests := []struct {
		name       string
		precoVenda float64
	}{
		{"Zero price", 0.0},
		{"Normal price", 99.99},
		{"High precision", 123.456789},
		{"Large price", 999999.99},
		{"Negative price", -10.50}, // Edge case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			request := dtos.CreateProductRequest{
				NomeProduto: "Produto Teste",
				PrecoVenda:  tt.precoVenda,
			}

			// Act
			produto := BuildProductEntity(request)

			// Assert
			assert.NotNil(t, produto)
			assert.Equal(t, tt.precoVenda, produto.PrecoVenda)
		})
	}
}

func TestBuildProductEntity_LongStrings(t *testing.T) {
	// Arrange
	longString := string(make([]byte, 500))
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	request := dtos.CreateProductRequest{
		DataCadastro:  longString,
		CodigoBarra:   longString,
		NomeProduto:   longString,
		SKU:           longString,
		Categoria:     longString,
		DestinadoPara: longString,
		Variacao:      longString,
		Marca:         longString,
		Descricao:     longString,
		Status:        longString,
		PrecoVenda:    99.99,
	}

	// Act
	produto := BuildProductEntity(request)

	// Assert
	assert.NotNil(t, produto)
	assert.Equal(t, longString, produto.DataCadastro)
	assert.Equal(t, longString, produto.CodigoBarra)
	assert.Equal(t, longString, produto.NomeProduto)
	assert.Equal(t, longString, produto.SKU)
	assert.Equal(t, longString, produto.Categoria)
	assert.Equal(t, longString, produto.DestinadoPara)
	assert.Equal(t, longString, produto.Variacao)
	assert.Equal(t, longString, produto.Marca)
	assert.Equal(t, longString, produto.Descricao)
	assert.Equal(t, longString, produto.Status)
	assert.Equal(t, 99.99, produto.PrecoVenda)
}

func TestBuildProductEntity_ReturnType(t *testing.T) {
	// Arrange
	request := dtos.CreateProductRequest{
		NomeProduto: "Produto Teste",
		Status:      "Ativo",
	}

	// Act
	produto := BuildProductEntity(request)

	// Assert
	assert.IsType(t, &Produto{}, produto)
	assert.NotNil(t, produto)
}

// =============================================================================
// TESTES DE EDGE CASES
// =============================================================================

func TestProduto_NegativeID(t *testing.T) {
	// Arrange
	produto := Produto{
		IDProduto:   -1,
		NomeProduto: "Produto Teste",
	}

	// Act
	jsonData, err := json.Marshal(produto)
	assert.NoError(t, err)

	var unmarshaled Produto
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, -1, unmarshaled.IDProduto)
}

func TestProduto_BarcodeFormats(t *testing.T) {
	barcodeFormats := []string{
		"1234567890123",    // EAN-13
		"123456789012",     // UPC-A
		"12345678",         // EAN-8
		"789-0123456789-0", // With separators
		"",                 // Empty
		"invalid-barcode",  // Invalid format
		"ABCD1234EFGH5678", // Alphanumeric
	}

	for _, barcode := range barcodeFormats {
		t.Run("Barcode: "+barcode, func(t *testing.T) {
			// Arrange
			produto := Produto{
				IDProduto:   1,
				NomeProduto: "Produto Teste",
				CodigoBarra: barcode,
			}

			// Act
			jsonData, err := json.Marshal(produto)
			assert.NoError(t, err)

			var unmarshaled Produto
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, barcode, unmarshaled.CodigoBarra)
		})
	}
}

func TestProduto_SKUFormats(t *testing.T) {
	skuFormats := []string{
		"SKU001",
		"SKU-001",
		"SKU_001",
		"SKU.001",
		"SKU/001",
		"PROD-ABC-123",
		"",
		"123456789012345678901234567890", // Very long SKU
	}

	for _, sku := range skuFormats {
		t.Run("SKU: "+sku, func(t *testing.T) {
			// Arrange
			produto := Produto{
				IDProduto:   1,
				NomeProduto: "Produto Teste",
				SKU:         sku,
			}

			// Act
			jsonData, err := json.Marshal(produto)
			assert.NoError(t, err)

			var unmarshaled Produto
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, sku, unmarshaled.SKU)
		})
	}
}

func TestBuildProductEntity_NilSafety(t *testing.T) {
	// Arrange - Test that function doesn't panic with zero values
	var request dtos.CreateProductRequest

	// Act
	produto := BuildProductEntity(request)

	// Assert
	assert.NotNil(t, produto)
	assert.IsType(t, &Produto{}, produto)
}
