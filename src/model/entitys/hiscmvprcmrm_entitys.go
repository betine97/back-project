package entity

type HisCmvPrcMarge struct {
	ID           int     `gorm:"primaryKey" json:"id"`
	IDProduto    int     `gorm:"column:id_produto;not null" json:"id_produto"`
	PrecoVenda   float64 `gorm:"column:preco_venda;not null" json:"preco_venda"`
	Cmv          float64 `gorm:"column:cmv;not null" json:"cmv"`
	Margem       float64 `gorm:"column:margem;not null" json:"margem"`
	DataRegistro string  `gorm:"column:data_registro;not null" json:"data_registro"`
}

func (HisCmvPrcMarge) TableName() string {
	return "produto_precos"
}
