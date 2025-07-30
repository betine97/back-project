package entity

import (
	dtos "github.com/betine97/back-project.git/src/model/dtos"
)

type Produto struct {
	IDProduto     int     `gorm:"primaryKey;column:id_produto" json:"id_produto"`
	DataCadastro  string  `gorm:"column:data_cadastro" json:"data_cadastro"`
	CodigoBarra   string  `gorm:"column:codigo_barra" json:"codigo_barra"`
	NomeProduto   string  `gorm:"column:nome_produto" json:"nome_produto"`
	SKU           string  `gorm:"column:sku" json:"sku"`
	Categoria     string  `gorm:"column:categoria" json:"categoria"`
	DestinadoPara string  `gorm:"column:destinado_para" json:"destinado_para"`
	Variacao      string  `gorm:"column:variacao" json:"variacao"`
	Marca         string  `gorm:"column:marca" json:"marca"`
	Descricao     string  `gorm:"column:descricao" json:"descricao"`
	Status        string  `gorm:"column:status" json:"status"`
	PrecoVenda    float64 `gorm:"column:preco_venda" json:"preco_venda"`
}

func BuildProductEntity(request dtos.CreateProductRequest) *Produto {
	return &Produto{
		DataCadastro:  request.DataCadastro,
		CodigoBarra:   request.CodigoBarra,
		NomeProduto:   request.NomeProduto,
		SKU:           request.SKU,
		Categoria:     request.Categoria,
		DestinadoPara: request.DestinadoPara,
		Variacao:      request.Variacao,
		Marca:         request.Marca,
		Descricao:     request.Descricao,
		Status:        request.Status,
		PrecoVenda:    request.PrecoVenda,
	}
}
