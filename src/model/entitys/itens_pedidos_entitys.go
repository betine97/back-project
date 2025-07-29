package entity

type ItemPedido struct {
	IDItem        int     `gorm:"primaryKey" json:"id_item"`
	IDPedido      int     `gorm:"column:id_pedido;not null" json:"id_pedido"`
	IDProduto     int     `gorm:"column:id_produto;not null" json:"id_produto"`
	Quantidade    int     `gorm:"column:quantidade;not null" json:"quantidade"`
	PrecoUnitario float64 `gorm:"column:preco_unitario;not null" json:"preco_unitario"`
	Subtotal      float64 `gorm:"column:subtotal;not null" json:"subtotal"`
}

func (ItemPedido) TableName() string {
	return "item_pedido"
}
