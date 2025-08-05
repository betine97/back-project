package dtos

// Para GET api/pedidos
type PedidoResponse struct {
	ID           int     `json:"id_pedido"`
	IDFornecedor int     `json:"id_fornecedor"`
	DataPedido   string  `json:"data_pedido"`
	DataEntrega  string  `json:"data_entrega"`
	ValorFrete   float64 `json:"valor_frete"`
	CustoPedido  float64 `json:"custo_pedido"`
	ValorTotal   float64 `json:"valor_total"`
	Descricao    string  `json:"descricao_pedido"`
	Status       string  `json:"status"`
}

type PedidoListResponse struct {
	Pedidos    []PedidoResponse `json:"pedidos"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

// Para POST api/pedidos
type CreatePedidoRequest struct {
	IDFornecedor int     `json:"id_fornecedor" validate:"required,gt=0"`
	DataPedido   string  `json:"data_pedido" validate:"required"`
	DataEntrega  string  `json:"data_entrega" validate:"required"`
	ValorFrete   float64 `json:"valor_frete" validate:"required,gte=0"`
	CustoPedido  float64 `json:"custo_pedido" validate:"required,gt=0"`
	ValorTotal   float64 `json:"valor_total" validate:"required,gt=0"`
	Descricao    string  `json:"descricao_pedido" validate:"required"`
	Status       string  `json:"status" validate:"required"`
}
