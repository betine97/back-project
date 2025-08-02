package dtos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TESTES PARA PedidoResponse DTO
// =============================================================================

func TestPedidoResponse_JSONSerialization(t *testing.T) {
	// Arrange
	response := PedidoResponse{
		ID:           1,
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
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled PedidoResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, response.ID, unmarshaled.ID)
	assert.Equal(t, response.IDFornecedor, unmarshaled.IDFornecedor)
	assert.Equal(t, response.DataPedido, unmarshaled.DataPedido)
	assert.Equal(t, response.DataEntrega, unmarshaled.DataEntrega)
	assert.Equal(t, response.ValorFrete, unmarshaled.ValorFrete)
	assert.Equal(t, response.CustoPedido, unmarshaled.CustoPedido)
	assert.Equal(t, response.ValorTotal, unmarshaled.ValorTotal)
	assert.Equal(t, response.Descricao, unmarshaled.Descricao)
	assert.Equal(t, response.Status, unmarshaled.Status)
}

func TestPedidoResponse_JSONTags(t *testing.T) {
	// Arrange
	response := PedidoResponse{
		ID:           1,
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
	jsonData, err := json.Marshal(response)
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

func TestPedidoResponse_MonetaryValues(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			response := PedidoResponse{
				ID:          1,
				ValorFrete:  tt.valorFrete,
				CustoPedido: tt.custoPedido,
				ValorTotal:  tt.valorTotal,
			}

			// Act
			jsonData, err := json.Marshal(response)
			assert.NoError(t, err)

			var unmarshaled PedidoResponse
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.valorFrete, unmarshaled.ValorFrete)
			assert.Equal(t, tt.custoPedido, unmarshaled.CustoPedido)
			assert.Equal(t, tt.valorTotal, unmarshaled.ValorTotal)
		})
	}
}

func TestPedidoResponse_StatusValues(t *testing.T) {
	validStatuses := []string{
		"Pendente",
		"Processando",
		"Enviado",
		"Entregue",
		"Cancelado",
		"Devolvido",
	}

	for _, status := range validStatuses {
		t.Run("Status: "+status, func(t *testing.T) {
			// Arrange
			response := PedidoResponse{
				ID:     1,
				Status: status,
			}

			// Act
			jsonData, err := json.Marshal(response)
			assert.NoError(t, err)

			var unmarshaled PedidoResponse
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, status, unmarshaled.Status)
		})
	}
}

func TestPedidoResponse_DateFormats(t *testing.T) {
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
	}

	for _, df := range dateFormats {
		t.Run(df.name, func(t *testing.T) {
			// Arrange
			response := PedidoResponse{
				ID:          1,
				DataPedido:  df.dataPedido,
				DataEntrega: df.dataEntrega,
			}

			// Act
			jsonData, err := json.Marshal(response)
			assert.NoError(t, err)

			var unmarshaled PedidoResponse
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, df.dataPedido, unmarshaled.DataPedido)
			assert.Equal(t, df.dataEntrega, unmarshaled.DataEntrega)
		})
	}
}

// =============================================================================
// TESTES PARA PedidoListResponse DTO
// =============================================================================

func TestPedidoListResponse_JSONSerialization(t *testing.T) {
	// Arrange
	pedidos := []PedidoResponse{
		{
			ID:           1,
			IDFornecedor: 10,
			DataPedido:   "2024-01-01",
			DataEntrega:  "2024-01-15",
			ValorFrete:   15.50,
			CustoPedido:  100.00,
			ValorTotal:   115.50,
			Descricao:    "Primeiro pedido",
			Status:       "Pendente",
		},
		{
			ID:           2,
			IDFornecedor: 20,
			DataPedido:   "2024-01-02",
			DataEntrega:  "2024-01-16",
			ValorFrete:   20.00,
			CustoPedido:  200.00,
			ValorTotal:   220.00,
			Descricao:    "Segundo pedido",
			Status:       "Entregue",
		},
	}

	response := PedidoListResponse{
		Pedidos: pedidos,
		Total:   2,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled PedidoListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, response.Total, unmarshaled.Total)
	assert.Len(t, unmarshaled.Pedidos, 2)
	assert.Equal(t, response.Pedidos[0].ID, unmarshaled.Pedidos[0].ID)
	assert.Equal(t, response.Pedidos[0].IDFornecedor, unmarshaled.Pedidos[0].IDFornecedor)
	assert.Equal(t, response.Pedidos[0].ValorTotal, unmarshaled.Pedidos[0].ValorTotal)
	assert.Equal(t, response.Pedidos[1].ID, unmarshaled.Pedidos[1].ID)
	assert.Equal(t, response.Pedidos[1].Status, unmarshaled.Pedidos[1].Status)
}

func TestPedidoListResponse_EmptyList(t *testing.T) {
	// Arrange
	response := PedidoListResponse{
		Pedidos: []PedidoResponse{},
		Total:   0,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled PedidoListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, unmarshaled.Total)
	assert.Len(t, unmarshaled.Pedidos, 0)
	assert.NotNil(t, unmarshaled.Pedidos) // Deve ser slice vazio, não nil
}

func TestPedidoListResponse_LargeList(t *testing.T) {
	// Arrange
	pedidos := make([]PedidoResponse, 50)
	for i := 0; i < 50; i++ {
		pedidos[i] = PedidoResponse{
			ID:           i + 1,
			IDFornecedor: (i % 10) + 1, // Varia entre 1 e 10
			DataPedido:   "2024-01-01",
			DataEntrega:  "2024-01-15",
			ValorFrete:   float64(i+1) * 5.0,
			CustoPedido:  float64(i+1) * 50.0,
			ValorTotal:   float64(i+1) * 55.0,
			Descricao:    "Pedido número " + string(rune(i+1)),
			Status:       "Pendente",
		}
	}

	response := PedidoListResponse{
		Pedidos: pedidos,
		Total:   50,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled PedidoListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 50, unmarshaled.Total)
	assert.Len(t, unmarshaled.Pedidos, 50)
	assert.Equal(t, 1, unmarshaled.Pedidos[0].ID)
	assert.Equal(t, 50, unmarshaled.Pedidos[49].ID)
	assert.Equal(t, 1, unmarshaled.Pedidos[0].IDFornecedor)
	assert.Equal(t, 10, unmarshaled.Pedidos[49].IDFornecedor)
}

// =============================================================================
// TESTES DE EDGE CASES
// =============================================================================

func TestPedidoResponse_EmptyValues(t *testing.T) {
	// Arrange
	response := PedidoResponse{}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled PedidoResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, unmarshaled.ID)
	assert.Equal(t, 0, unmarshaled.IDFornecedor)
	assert.Equal(t, "", unmarshaled.DataPedido)
	assert.Equal(t, "", unmarshaled.DataEntrega)
	assert.Equal(t, float64(0), unmarshaled.ValorFrete)
	assert.Equal(t, float64(0), unmarshaled.CustoPedido)
	assert.Equal(t, float64(0), unmarshaled.ValorTotal)
	assert.Equal(t, "", unmarshaled.Descricao)
	assert.Equal(t, "", unmarshaled.Status)
}

func TestPedidoResponse_NegativeValues(t *testing.T) {
	// Arrange
	response := PedidoResponse{
		ID:           -1,
		IDFornecedor: -10,
		ValorFrete:   -15.50,
		CustoPedido:  -100.00,
		ValorTotal:   -115.50,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled PedidoResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, -1, unmarshaled.ID)
	assert.Equal(t, -10, unmarshaled.IDFornecedor)
	assert.Equal(t, -15.50, unmarshaled.ValorFrete)
	assert.Equal(t, -100.00, unmarshaled.CustoPedido)
	assert.Equal(t, -115.50, unmarshaled.ValorTotal)
}

func TestPedidoResponse_LongDescription(t *testing.T) {
	// Arrange
	longDescription := string(make([]byte, 2000))
	for i := range longDescription {
		longDescription = longDescription[:i] + "a" + longDescription[i+1:]
	}

	response := PedidoResponse{
		ID:        1,
		Descricao: longDescription,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled PedidoResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, longDescription, unmarshaled.Descricao)
}

func TestPedidoResponse_SpecialCharactersInDescription(t *testing.T) {
	// Arrange
	specialDescription := "Pedido com 'aspas simples', \"aspas duplas\", acentos: ção, símbolos: @#$%&*()+=[]{}|\\:;\"'<>?,./"

	response := PedidoResponse{
		ID:        1,
		Descricao: specialDescription,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled PedidoResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, specialDescription, unmarshaled.Descricao)
}

func TestPedidoListResponse_NilSlice(t *testing.T) {
	// Arrange
	response := PedidoListResponse{
		Pedidos: nil,
		Total:   0,
	}

	// Act
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled PedidoListResponse
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, unmarshaled.Total)
	assert.Nil(t, unmarshaled.Pedidos) // nil slice é preservado
}

// =============================================================================
// TESTES DE CÁLCULOS E VALIDAÇÕES DE NEGÓCIO
// =============================================================================

func TestPedidoResponse_ValueCalculations(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			response := PedidoResponse{
				ID:          1,
				ValorFrete:  tt.valorFrete,
				CustoPedido: tt.custoPedido,
				ValorTotal:  tt.valorTotal,
			}

			// Act
			jsonData, err := json.Marshal(response)
			assert.NoError(t, err)

			var unmarshaled PedidoResponse
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

func TestPedidoListResponse_TotalConsistency(t *testing.T) {
	// Arrange
	pedidos := []PedidoResponse{
		{ID: 1, Status: "Pendente"},
		{ID: 2, Status: "Entregue"},
		{ID: 3, Status: "Cancelado"},
	}

	tests := []struct {
		name    string
		pedidos []PedidoResponse
		total   int
		isValid bool
	}{
		{"Correct total", pedidos, 3, true},
		{"Incorrect total - higher", pedidos, 5, false},
		{"Incorrect total - lower", pedidos, 1, false},
		{"Empty list with zero total", []PedidoResponse{}, 0, true},
		{"Empty list with non-zero total", []PedidoResponse{}, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			response := PedidoListResponse{
				Pedidos: tt.pedidos,
				Total:   tt.total,
			}

			// Act
			jsonData, err := json.Marshal(response)
			assert.NoError(t, err)

			var unmarshaled PedidoListResponse
			err = json.Unmarshal(jsonData, &unmarshaled)

			// Assert
			assert.NoError(t, err)

			// Verificar se os valores foram preservados
			assert.Equal(t, tt.total, unmarshaled.Total)
			assert.Len(t, unmarshaled.Pedidos, len(tt.pedidos))

			// Verificar consistência entre total e tamanho da lista
			isConsistent := unmarshaled.Total == len(unmarshaled.Pedidos)
			assert.Equal(t, tt.isValid, isConsistent,
				"Expected consistency: %v, but got: %v (total: %d, list length: %d)",
				tt.isValid, isConsistent, unmarshaled.Total, len(unmarshaled.Pedidos))
		})
	}
}
