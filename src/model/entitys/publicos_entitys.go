package entity

import dtos "github.com/betine97/back-project.git/src/model/dtos"

// Entidade para a tabela criterios
type Criterio struct {
	ID           int    `gorm:"primaryKey;autoIncrement;table:criterios" json:"id"`
	NomeCondicao string `gorm:"column:nome_condicao;not null" json:"nome_condicao"`
}

// TableName especifica o nome da tabela para GORM
func (Criterio) TableName() string {
	return "criterios"
}

// Entidade para a tabela publicos_clientes
type PublicoCliente struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome        string `gorm:"column:nome;not null" json:"nome"`
	Descricao   string `gorm:"column:descricao;not null" json:"descricao"`
	DataCriacao string `gorm:"column:data_criacao;not null" json:"data_criacao"`
	Status      string `gorm:"column:status;not null" json:"status"`
}

// TableName especifica o nome da tabela para GORM
func (PublicoCliente) TableName() string {
	return "publicos_clientes"
}

// Entidade para a tabela publicos_criterios (relacionamento)
type PublicoCriterio struct {
	ID         int `gorm:"primaryKey;autoIncrement" json:"id"`
	IDPublico  int `gorm:"column:id_publico;not null" json:"id_publico"`
	IDCriterio int `gorm:"column:id_criterio;not null" json:"id_criterio"`
}

// TableName especifica o nome da tabela para GORM
func (PublicoCriterio) TableName() string {
	return "publicos_criterios"
}

// Estrutura para consulta SQL com JOIN
type PublicoCriterioJoin struct {
	IDPublico    int    `json:"id_publico"`
	IDCriterio   int    `json:"id_criterio"`
	NomeCondicao string `json:"nome_condicao"`
}

// Função para construir entidade PublicoCliente a partir do DTO
func BuildPublicoClienteEntity(request dtos.CreatePublicoRequest) *PublicoCliente {
	return &PublicoCliente{
		Nome:        request.Nome,
		Descricao:   request.Descricao,
		DataCriacao: request.DataCriacao,
		Status:      request.Status,
	}
}

// Entidade para a tabela addclientes_publicos (relacionamento cliente-público)
type AddClientePublico struct {
	ID        int `gorm:"primaryKey;autoIncrement" json:"id"`
	IDPublico int `gorm:"column:id_publico;not null" json:"id_publico"`
	IDCliente int `gorm:"column:id_cliente;not null" json:"id_cliente"`
}

// TableName especifica o nome da tabela para GORM
func (AddClientePublico) TableName() string {
	return "addclientes_publicos"
}

// Função para construir entidade PublicoCriterio
func BuildPublicoCriterioEntity(idPublico, idCriterio int) *PublicoCriterio {
	return &PublicoCriterio{
		IDPublico:  idPublico,
		IDCriterio: idCriterio,
	}
}

// Função para construir entidade AddClientePublico
func BuildAddClientePublicoEntity(idPublico, idCliente int) *AddClientePublico {
	return &AddClientePublico{
		IDPublico: idPublico,
		IDCliente: idCliente,
	}
}
