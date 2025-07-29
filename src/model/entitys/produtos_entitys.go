package entity

type Produto struct {
	ID            int     `gorm:"primaryKey" json:"id"`
	CodigoBarra   string  `gorm:"column:codigo_barra;not null" json:"codigo_barra"`
	NomeProduto   string  `gorm:"column:nome_produto;not null" json:"nome_produto"`
	SKU           string  `gorm:"not null" json:"sku"`
	Categoria     string  `gorm:"not null" json:"categoria"`
	DestinadoPara string  `gorm:"column:destinado_para;not null" json:"destinado_para"`
	Variacao      string  `json:"variacao"`
	Marca         string  `json:"marca"`
	Descricao     string  `json:"descricao"`
	Status        string  `gorm:"not null" json:"status"`
	PrecoVenda    float64 `gorm:"column:preco_venda;not null" json:"preco_venda"`
}

func (Produto) TableName() string {
	return "produtos"
}
