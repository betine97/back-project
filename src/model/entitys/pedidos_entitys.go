package entity

type Pedido struct {
	IDPedido     int     `gorm:"primaryKey;column:id_pedido" json:"id_pedido"`
	IDFornecedor int     `gorm:"column:id_fornecedor;not null" json:"id_fornecedor"`
	DataPedido   string  `gorm:"column:data_pedido;not null" json:"data_pedido"`
	DataEntrega  string  `gorm:"column:data_entrega;not null" json:"data_entrega"`
	ValorFrete   float64 `gorm:"column:valor_frete;not null" json:"valor_frete"`
	CustoPedido  float64 `gorm:"column:custo_pedido;not null" json:"custo_pedido"`
	ValorTotal   float64 `gorm:"column:valor_total;not null" json:"valor_total"`
	Descricao    string  `gorm:"column:descricao_pedido;not null" json:"descricao_pedido"`
	Status       string  `gorm:"column:status;not null" json:"status"`
}
