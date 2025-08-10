package entity

type ViewDetalhesEstoque struct {
	NomeProduto         string `gorm:"column:nome_produto" json:"nome_produto"`
	Lote                int    `gorm:"column:lote" json:"lote"`
	Quantidade          int    `gorm:"column:quantidade" json:"quantidade"`
	DataEntrada         string `gorm:"column:data_entrada" json:"data_entrada"`
	DataSaida           string `gorm:"column:data_saida" json:"data_saida"`
	Vencimento          string `gorm:"column:vencimento" json:"vencimento"`
	DocumentoReferencia string `gorm:"column:documento_referencia" json:"documento_referencia"`
	Status              string `gorm:"column:status" json:"status"`
}

func (ViewDetalhesEstoque) TableName() string {
	return "view_detalhes_estoque"
}
