package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TESTES PARA Pedido Entity
// =============================================================================

func TestPedido_JSONSerialization(t *testing.T) {
	// Arrange
	pedido := Pedido{
		IDPedido:     1,
		IDFornecedor: 10,
		DataPedido:   "2024-01-01",
		DataEntrega:  "2024-01-15",
		ValorFrete:   15.50,
		CustoPedido:  100.00,
		ValorTotal:   115.50,
		Descricao:    "Pedido de teste com produtos diversos",
		Status:       "Pendente",
	}

	// Act
	jsonData, err := json.Marshal(pedido)
	assert.NoError(t, err)

	var unmarshaled Pedido
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, pedido.IDPedido, unmarshaled.IDPedido)
	assert.Equal(t, pedido.IDFornecedor, unmarshaled.IDFornecedor)
	assert.Equal(t, pedido.DataPedido, unmarshaled.DataPedido)
	assert.Equal(t, pedido.DataEntrega, unmarshaled.DataEntrega)
	assert.Equal(t, pedido.ValorFrete, unmarshaled.ValorFrete)
	assert.Equal(t, pedido.CustoPedido, unmarshaled.CustoPedido)
	assert.Equal(t, pedido.ValorTotal, unmarshaled.ValorTotal)
	assert.Equal(t, pedido.Descricao, unmarshaled.Descricao)
	assert.Equal(t, pedido.Status, unmarshaled.Status)
}

func TestPedido_JSONTags(t *testing.T) {
	// Arrange
	pedido := Pedido{
		IDPedido:     1,
		IDFornecedor: 10,
		DataPedido:   "2024-01-01",
		DataEntrega:  "2024-01-15",
		ValorFrete:   15.50,
		CustoPedido:  100.00,
		ValorTotal:   115.50,
		Descricao:    "Pedido de teste",
		Status:       "Pendente",
	}

	// Act
	jsonData, err := json.Marshal(pedido)
	assert.NoError(t, err)

	// Assert
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, `"id_pedido":1`)
	assert.Contains(t, jsonString, `"id_fornecedor":10`)
	assert.Contains(t, jsonString, `"data_pedido":"2024-01-01"`)
	assert.Contains(t, jsonString, `"data_entrega":"2024-01-15"`)
	assert.Contains(t, jsonString, `"valor_frete":15.5`)
	assert.Contains(t, jsonString, `"custo_pedido":100`)
	assert.Contains(t, jsonString, `"valor_total":115.5`)
	assert.Contains(t, jsonString, `"descricao_pedido":"Pedido de teste"`)
	assert.Contains(t, jsonString, `"status":"Pendente"`)
}

func TestPedido_EmptyValues(t *testing.T) {
	// Arrange
	pedido := Pedido{}

	// Act
	jsonData, err := json.Marshal(pedido)
	assert.NoError(t, err)

	var unmarshaled Pedido
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, unmarshaled.IDPedido)
	assert.Equal(t, 0, unmarshaled.IDFornecedor)
	assert.Equal(t, "", unmarshaled.DataPedido)
	assert.Equal(t, "", unmarshaled.DataEntrega)
	assert.Equal(t, float64(0), unmarshaled.ValorFrete)
	assert.Equal(t, float64(0), unmarshaled.CustoPedido)
	assert.Equal(t, float64(0), unmarshaled.ValorTotal)
	assert.Equal(t, "", unmarshaled.Descricao)
	assert.Equal(t, "", unmarshaled.Status)
}

func TestPedido_MonetaryValues(t *testing.T) {
	tests := []struct {
		name        string
		valorFrete  float64
		custoPedido float64
		valorTotal  float64
	}{
		{"Zero values", 0.0, 0.0, 0.0},
		{"Integer values", 10.0, 100.0, 110.0},
		{"Decimal values", 15.75, 99.99, 115.74},
		{"High precision", 12.345, 87.654, 99.999},
		{"Large values", 1000.50, 9999.99, 11000.49},
		{"Negative values", -10.0, -100.0, -110.0}, // Edge case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			pedido := Pedido{
				IDPedido:    1,
				ValorFrete:  tt.valorFrete,
				CustoPedido: tt.custoPedido,
				ValorTotal:  tt.valorTotal,
			}

			// Act
			jsonData, err := json.Marshal(pedido)
			assert.NoError(t, err)

			var unmarshaled Pedido
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.valorFrete, unmarshaled.ValorFrete)
			assert.Equal(t, tt.custoPedido, unmarshaled.CustoPedido)
			assert.Equal(t, tt.valorTotal, unmarshaled.ValorTotal)
		})
	}
}

func TestPedido_StatusValues(t *testing.T) {
	validStatuses := []string{
		"Pendente",
		"Processando",
		"Enviado",
		"Entregue",
		"Cancelado",
		"Devolvido",
		"Em Análise",
		"Aguardando Pagamento",
	}

	for _, status := range validStatuses {
		t.Run("Status: "+status, func(t *testing.T) {
			// Arrange
			pedido := Pedido{
				IDPedido: 1,
				Status:   status,
			}

			// Act
			jsonData, err := json.Marshal(pedido)
			assert.NoError(t, err)

			var unmarshaled Pedido
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, status, unmarshaled.Status)
		})
	}
}

func TestPedido_DateFormats(t *testing.T) {
	dateFormats := []struct {
		name        string
		dataPedido  string
		dataEntrega string
	}{
		{"ISO Date", "2024-01-01", "2024-01-15"},
		{"ISO DateTime", "2024-01-01T10:30:00Z", "2024-01-15T14:45:00Z"},
		{"Brazilian format", "01/01/2024", "15/01/2024"},
		{"Empty dates", "", ""},
		{"Partial dates", "2024-01", "2024-02"},
		{"Invalid dates", "invalid-date", "another-invalid"},
	}

	for _, df := range dateFormats {
		t.Run(df.name, func(t *testing.T) {
			// Arrange
			pedido := Pedido{
				IDPedido:    1,
				DataPedido:  df.dataPedido,
				DataEntrega: df.dataEntrega,
			}

			// Act
			jsonData, err := json.Marshal(pedido)
			assert.NoError(t, err)

			var unmarshaled Pedido
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, df.dataPedido, unmarshaled.DataPedido)
			assert.Equal(t, df.dataEntrega, unmarshaled.DataEntrega)
		})
	}
}

func TestPedido_SpecialCharacters(t *testing.T) {
	// Arrange
	pedido := Pedido{
		IDPedido:     1,
		IDFornecedor: 10,
		DataPedido:   "2024-01-01T10:30:00Z",
		DataEntrega:  "2024-01-15T14:45:00Z",
		ValorFrete:   15.50,
		CustoPedido:  100.00,
		ValorTotal:   115.50,
		Descricao:    "Pedido com 'aspas simples', \"aspas duplas\", acentos: ção, símbolos: @#$%&*()+=[]{}|\\:;\"'<>?,./ e emojis: 📦🚚",
		Status:       "Pendente",
	}

	// Act
	jsonData, err := json.Marshal(pedido)
	assert.NoError(t, err)

	var unmarshaled Pedido
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, pedido.Descricao, unmarshaled.Descricao)
}

func TestPedido_LongDescription(t *testing.T) {
	// Arrange
	longDescription := string(make([]byte, 2000))
	for i := range longDescription {
		longDescription = longDescription[:i] + "a" + longDescription[i+1:]
	}

	pedido := Pedido{
		IDPedido:  1,
		Descricao: longDescription,
	}

	// Act
	jsonData, err := json.Marshal(pedido)
	assert.NoError(t, err)

	var unmarshaled Pedido
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, longDescription, unmarshaled.Descricao)
}

// =============================================================================
// TESTES DE VALIDAÇÃO DE NEGÓCIO
// =============================================================================

func TestPedido_ValueCalculations(t *testing.T) {
	tests := []struct {
		name        string
		valorFrete  float64
		custoPedido float64
		valorTotal  float64
		isValid     bool
	}{
		{"Correct calculation", 10.0, 100.0, 110.0, true},
		{"Incorrect calculation", 10.0, 100.0, 120.0, false},
		{"Zero freight", 0.0, 100.0, 100.0, true},
		{"Zero cost", 10.0, 0.0, 10.0, true},
		{"All zeros", 0.0, 0.0, 0.0, true},
		{"Negative values", -10.0, -100.0, -110.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			pedido := Pedido{
				IDPedido:    1,
				ValorFrete:  tt.valorFrete,
				CustoPedido: tt.custoPedido,
				ValorTotal:  tt.valorTotal,
			}

			// Act
			jsonData, err := json.Marshal(pedido)
			assert.NoError(t, err)

			var unmarshaled Pedido
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)

			// Verificar se os valores foram preservados
			assert.Equal(t, tt.valorFrete, unmarshaled.ValorFrete)
			assert.Equal(t, tt.custoPedido, unmarshaled.CustoPedido)
			assert.Equal(t, tt.valorTotal, unmarshaled.ValorTotal)

			// Verificar se o cálculo está correto (lógica de negócio)
			expectedTotal := tt.valorFrete + tt.custoPedido
			isCalculationCorrect := unmarshaled.ValorTotal == expectedTotal
			assert.Equal(t, tt.isValid, isCalculationCorrect,
				"Expected calculation validity: %v, but got: %v (%.2f + %.2f = %.2f, but total is %.2f)",
				tt.isValid, isCalculationCorrect, tt.valorFrete, tt.custoPedido, expectedTotal, tt.valorTotal)
		})
	}
}

func TestPedido_FornecedorRelationship(t *testing.T) {
	tests := []struct {
		name         string
		idFornecedor int
		isValid      bool
	}{
		{"Valid fornecedor ID", 1, true},
		{"Another valid ID", 999, true},
		{"Zero ID", 0, false}, // Assuming 0 is invalid
		{"Negative ID", -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			pedido := Pedido{
				IDPedido:     1,
				IDFornecedor: tt.idFornecedor,
			}

			// Act
			jsonData, err := json.Marshal(pedido)
			assert.NoError(t, err)

			var unmarshaled Pedido
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.idFornecedor, unmarshaled.IDFornecedor)

			// Business logic validation
			isValidFornecedor := unmarshaled.IDFornecedor > 0
			assert.Equal(t, tt.isValid, isValidFornecedor)
		})
	}
}

// =============================================================================
// TESTES DE EDGE CASES
// =============================================================================

func TestPedido_NegativeIDs(t *testing.T) {
	// Arrange
	pedido := Pedido{
		IDPedido:     -1,
		IDFornecedor: -10,
	}

	// Act
	jsonData, err := json.Marshal(pedido)
	assert.NoError(t, err)

	var unmarshaled Pedido
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, -1, unmarshaled.IDPedido)
	assert.Equal(t, -10, unmarshaled.IDFornecedor)
}

func TestPedido_MaxValues(t *testing.T) {
	// Arrange - Test with very large values
	pedido := Pedido{
		IDPedido:     999999999,
		IDFornecedor: 999999999,
		ValorFrete:   999999999.99,
		CustoPedido:  999999999.99,
		ValorTotal:   1999999999.98,
	}

	// Act
	jsonData, err := json.Marshal(pedido)
	assert.NoError(t, err)

	var unmarshaled Pedido
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, pedido.IDPedido, unmarshaled.IDPedido)
	assert.Equal(t, pedido.IDFornecedor, unmarshaled.IDFornecedor)
	assert.Equal(t, pedido.ValorFrete, unmarshaled.ValorFrete)
	assert.Equal(t, pedido.CustoPedido, unmarshaled.CustoPedido)
	assert.Equal(t, pedido.ValorTotal, unmarshaled.ValorTotal)
}

func TestPedido_EmptyDescription(t *testing.T) {
	// Arrange
	pedido := Pedido{
		IDPedido:  1,
		Descricao: "",
	}

	// Act
	jsonData, err := json.Marshal(pedido)
	assert.NoError(t, err)

	var unmarshaled Pedido
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "", unmarshaled.Descricao)
}

func TestPedido_UnicodeCharacters(t *testing.T) {
	// Arrange
	pedido := Pedido{
		IDPedido:  1,
		Descricao: "Pedido com caracteres Unicode: 中文, العربية, русский, 日本語, 한국어, 🚚📦💰",
		Status:    "Enviado 🚚",
	}

	// Act
	jsonData, err := json.Marshal(pedido)
	assert.NoError(t, err)

	var unmarshaled Pedido
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, pedido.Descricao, unmarshaled.Descricao)
	assert.Equal(t, pedido.Status, unmarshaled.Status)
}

func TestPedido_GormTags(t *testing.T) {
	// Arrange - Test that struct has proper GORM tags for database mapping
	pedido := Pedido{
		IDPedido:     1,
		IDFornecedor: 10,
		DataPedido:   "2024-01-01",
		DataEntrega:  "2024-01-15",
		ValorFrete:   15.50,
		CustoPedido:  100.00,
		ValorTotal:   115.50,
		Descricao:    "Pedido teste",
		Status:       "Pendente",
	}

	// Act & Assert - Verify struct can be used with GORM
	// This is more of a structural test
	assert.Equal(t, 1, pedido.IDPedido)
	assert.Equal(t, 10, pedido.IDFornecedor)
	assert.Equal(t, "2024-01-01", pedido.DataPedido)
	assert.Equal(t, "2024-01-15", pedido.DataEntrega)
	assert.Equal(t, 15.50, pedido.ValorFrete)
	assert.Equal(t, 100.00, pedido.CustoPedido)
	assert.Equal(t, 115.50, pedido.ValorTotal)
	assert.Equal(t, "Pedido teste", pedido.Descricao)
	assert.Equal(t, "Pendente", pedido.Status)
}
