package dtos_models

type Pedido struct {
	ID           int     `json:"id"`
	IDFornecedor int     `json:"id_fornecedor"`
	DataPedido   string  `json:"data_pedido"`
	DataEntrega  string  `json:"data_entrega"`
	ValorFrete   float64 `json:"valor_frete"`
}

type PedidoListResponse struct {
	Pedidos []PedidoResponse `json:"pedidos"`
	Total   int              `json:"total"`
}

type PedidoResponse struct {
	ID           int     `json:"id"`
	IDFornecedor int     `json:"id_fornecedor"`
	DataPedido   string  `json:"data_pedido"`
	DataEntrega  string  `json:"data_entrega"`
	ValorFrete   float64 `json:"valor_frete"`
	CustoPedido  float64 `json:"custo_pedido"`
	ValorTotal   float64 `json:"valor_total"`
	Descricao    string  `json:"descricao"`
	Status       string  `json:"status"`
}

type PedidoQueryParams struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}
