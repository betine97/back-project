package dtos

// DTO para resposta simples de completude do cadastro
type CompletudeClienteResponse struct {
	ClienteID          int     `json:"cliente_id"`
	NomeCliente        string  `json:"nome_cliente"`
	PercentualCompleto float64 `json:"percentual_completo"`
}

type CompletudeListResponse struct {
	Clientes   []CompletudeClienteResponse `json:"clientes"`
	Total      int                         `json:"total"`
	Page       int                         `json:"page"`
	Limit      int                         `json:"limit"`
	TotalPages int                         `json:"total_pages"`
}
