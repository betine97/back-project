package entity

import dtos "github.com/betine97/back-project.git/src/model/dtos"

type Endereco struct {
	IDEndereco  int    `gorm:"primaryKey;autoIncrement;column:id_endereco" json:"id_endereco"`
	IDCliente   int    `gorm:"column:id_cliente;not null" json:"id_cliente"`
	CEP         string `gorm:"column:cep;not null" json:"cep"`
	Cidade      string `gorm:"column:cidade;not null" json:"cidade"`
	Estado      string `gorm:"column:estado;not null" json:"estado"`
	Bairro      string `gorm:"column:bairro;not null" json:"bairro"`
	Logradouro  string `gorm:"column:logradouro;not null" json:"logradouro"`
	Numero      string `gorm:"column:numero;not null" json:"numero"`
	Complemento string `gorm:"column:complemento" json:"complemento"`
	Obs         string `gorm:"column:obs" json:"obs"`
}

func BuildEnderecoEntity(request dtos.CreateEnderecoRequest) *Endereco {
	return &Endereco{
		IDCliente:   request.IDCliente,
		CEP:         request.CEP,
		Cidade:      request.Cidade,
		Estado:      request.Estado,
		Bairro:      request.Bairro,
		Logradouro:  request.Logradouro,
		Numero:      request.Numero,
		Complemento: request.Complemento,
		Obs:         request.Obs,
	}
}
