package entity

import dtos "github.com/betine97/back-project.git/src/model/dtos"

type Estoque struct {
	IDEstoque           int     `gorm:"primaryKey;column:id_estoque" json:"id_estoque"`
	IDProduto           int     `gorm:"column:id_produto;not null" json:"id_produto"`
	IDLote              int     `gorm:"column:id_lote;not null" json:"id_lote"`
	Quantidade          int     `gorm:"column:quantidade;not null" json:"quantidade"`
	Vencimento          string  `gorm:"column:vencimento" json:"vencimento"`
	CustoUnitario       float64 `gorm:"column:custo_unitario;not null" json:"custo_unitario"`
	DataEntrada         string  `gorm:"column:data_entrada" json:"data_entrada"`
	DataSaida           string  `gorm:"column:data_saida" json:"data_saida"`
	DocumentoReferencia string  `gorm:"column:documento_referencia" json:"documento_referencia"`
	Status              string  `gorm:"column:status;not null" json:"status"`
}

func BuildEstoqueEntity(request dtos.CreateEstoqueRequest) *Estoque {
	return &Estoque{
		IDProduto:           request.IDProduto,
		IDLote:              request.IDLote,
		Quantidade:          request.Quantidade,
		Vencimento:          request.Vencimento,
		CustoUnitario:       request.CustoUnitario,
		DataEntrada:         request.DataEntrada,
		DataSaida:           request.DataSaida,
		DocumentoReferencia: request.DocumentoReferencia,
		Status:              request.Status,
	}
}
