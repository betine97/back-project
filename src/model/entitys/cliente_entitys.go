package entity

import dtos "github.com/betine97/back-project.git/src/model/dtos"

type Cliente struct {
	ID             int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TipoCliente    string `gorm:"column:tipo_cliente;not null;unique" json:"tipo_cliente"`
	NomeCliente    string `gorm:"column:nome_cliente;not null" json:"nome_cliente"`
	NumeroCelular  string `gorm:"column:numero_celular;not null" json:"numero_celular"`
	Sexo           string `gorm:"column:sexo;not null" json:"sexo"`
	Email          string `gorm:"column:email;unique" json:"email"`
	DataNascimento string `gorm:"column:data_nascimento" json:"data_nascimento"`
	DataCadastro   string `gorm:"column:data_cadastro;not null" json:"data_cadastro"`
}

func BuildClienteEntity(request dtos.CreateClienteRequest) *Cliente {
	return &Cliente{
		TipoCliente:    request.TipoCliente,
		NomeCliente:    request.NomeCliente,
		NumeroCelular:  request.NumeroCelular,
		Sexo:           request.Sexo,
		Email:          request.Email,
		DataNascimento: request.DataNascimento,
		DataCadastro:   request.DataCadastro,
	}
}

// TagCliente representa a tabela tags_clientes
type TagCliente struct {
	ID        int `gorm:"primaryKey;autoIncrement;column:id;table:tags_clientes" json:"id"`
	IDTag     int `gorm:"column:id_tag;not null" json:"id_tag"`
	ClienteID int `gorm:"column:cliente_id;not null" json:"cliente_id"`
}

// TableName especifica o nome da tabela
func (TagCliente) TableName() string {
	return "tags_clientes"
}

// TagClienteJoin representa o resultado do JOIN entre tags_clientes e tags
type TagClienteJoin struct {
	ID        int    `json:"id"`
	IDTag     int    `json:"id_tag"`
	ClienteID int    `json:"cliente_id"`
	Nome      string `json:"nome"`
}
