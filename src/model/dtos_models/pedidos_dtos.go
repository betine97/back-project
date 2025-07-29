package dtos_models

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
	Pedidos []PedidoResponse `json:"descricao_pedido"`
	Total   int              `json:"total"`
}
