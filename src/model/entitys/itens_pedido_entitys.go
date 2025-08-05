package entity

import (
	dtos "github.com/betine97/back-project.git/src/model/dtos"
)

type ViewDetalhesPedido struct {
	IDPedido      int     `gorm:"column:id_pedido" json:"id_pedido"`
	NomeProduto   string  `gorm:"column:nome_produto" json:"nome_produto"`
	Quantidade    int     `gorm:"column:quantidade" json:"quantidade"`
	PrecoUnitario float64 `gorm:"column:preco_unitario" json:"preco_unitario"`
	TotalItem     float64 `gorm:"column:total_item" json:"total_item"`
}

func (ViewDetalhesPedido) TableName() string {
	return "view_detalhes_pedido"
}

type ItemPedido struct {
	IDItem        int     `gorm:"primaryKey;column:id_item;autoIncrement" json:"id_item"`
	IDPedido      int     `gorm:"column:id_pedido;not null" json:"id_pedido"`
	IDProduto     int     `gorm:"column:id_produto;not null" json:"id_produto"`
	Quantidade    int     `gorm:"column:quantidade;not null" json:"quantidade"`
	PrecoUnitario float64 `gorm:"column:preco_unitario;type:decimal(10,2);not null" json:"preco_unitario"`
	Subtotal      float64 `gorm:"column:subtotal;type:decimal(10,2);not null" json:"subtotal"`
}

func (ItemPedido) TableName() string {
	return "itens_pedido"
}

func BuildItemPedidoEntity(request dtos.CreateItemPedidoRequest, idPedido int) *ItemPedido {
	return &ItemPedido{
		IDPedido:      idPedido,
		IDProduto:     request.IDProduto,
		Quantidade:    request.Quantidade,
		PrecoUnitario: request.PrecoUnitario,
		Subtotal:      request.Subtotal,
	}
}
