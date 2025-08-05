package entity

import (
	dtos "github.com/betine97/back-project.git/src/model/dtos"
)

type Pedido struct {
	IDPedido     int     `gorm:"primaryKey;column:id_pedido;autoIncrement" json:"id_pedido"`
	IDFornecedor int     `gorm:"column:id_fornecedor;not null" json:"id_fornecedor"`
	DataPedido   string  `gorm:"column:data_pedido;size:20;not null" json:"data_pedido"`
	DataEntrega  string  `gorm:"column:data_entrega;size:20;not null" json:"data_entrega"`
	ValorFrete   float64 `gorm:"column:valor_frete;type:decimal(10,2);not null" json:"valor_frete"`
	CustoPedido  float64 `gorm:"column:custo_pedido;type:decimal(10,2);not null" json:"custo_pedido"`
	ValorTotal   float64 `gorm:"column:valor_total;type:decimal(10,2);not null" json:"valor_total"`
	Descricao    string  `gorm:"column:descricao_pedido;type:text;not null" json:"descricao_pedido"`
	Status       string  `gorm:"column:status;size:100;not null" json:"status"`
}

func (Pedido) TableName() string {
	return "pedidos"
}

func BuildPedidoEntity(request dtos.CreatePedidoRequest) *Pedido {
	return &Pedido{
		IDFornecedor: request.IDFornecedor,
		DataPedido:   request.DataPedido,
		DataEntrega:  request.DataEntrega,
		ValorFrete:   request.ValorFrete,
		CustoPedido:  request.CustoPedido,
		ValorTotal:   request.ValorTotal,
		Descricao:    request.Descricao,
		Status:       request.Status,
	}
}
