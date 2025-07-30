package entity

import dtos "github.com/betine97/back-project.git/src/model/dtos"

type Fornecedores struct {
	ID           int    `gorm:"primaryKey;column:id_fornecedor" json:"id_fornecedor"`
	DataCadastro string `gorm:"column:data_cadastro" json:"data_cadastro"`
	Nome         string `gorm:"column:nome" json:"nome"`
	Telefone     string `gorm:"column:telefone" json:"telefone"`
	Email        string `gorm:"column:email" json:"email"`
	Cidade       string `gorm:"column:cidade" json:"cidade"`
	Estado       string `gorm:"column:estado" json:"estado"`
	Status       string `gorm:"column:status" json:"status"`
}

func BuildFornecedorEntity(request dtos.CreateFornecedorRequest) *Fornecedores {
	return &Fornecedores{
		Nome:         request.Nome,
		Telefone:     request.Telefone,
		Email:        request.Email,
		Cidade:       request.Cidade,
		Estado:       request.Estado,
		Status:       request.Status,
		DataCadastro: request.DataCadastro,
	}
}
