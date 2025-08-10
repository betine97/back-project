package entity

import dtos "github.com/betine97/back-project.git/src/model/dtos"

// Entidade para a tabela campanhas
type Campanha struct {
	ID             int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome           string `gorm:"column:nome;not null" json:"nome"`
	Desc           string `gorm:"column:desc" json:"desc"`
	DataCriacao    string `gorm:"column:data_criacao;not null" json:"data_criacao"`
	DataLancamento string `gorm:"column:data_lancamento;not null" json:"data_lancamento"`
	DataFim        string `gorm:"column:data_fim;not null" json:"data_fim"`
	Status         string `gorm:"column:status;default:ativa" json:"status"`
}

// TableName especifica o nome da tabela para GORM
func (Campanha) TableName() string {
	return "campanhas"
}

// Entidade para a tabela campanhas_publicos (relacionamento)
type CampanhaPublico struct {
	ID         int `gorm:"primaryKey;autoIncrement" json:"id"`
	IDCampanha int `gorm:"column:id_campanha;not null" json:"id_campanha"`
	IDPublico  int `gorm:"column:id_publico;not null" json:"id_publico"`
}

// TableName especifica o nome da tabela para GORM
func (CampanhaPublico) TableName() string {
	return "campanhas_publicos"
}

// Função para construir entidade Campanha a partir do DTO
func BuildCampanhaEntity(request dtos.CreateCampanhaRequest) *Campanha {
	return &Campanha{
		Nome:           request.Nome,
		Desc:           request.Desc,
		DataCriacao:    request.DataCriacao,
		DataLancamento: request.DataLancamento,
		DataFim:        request.DataFim,
		Status:         request.Status,
	}
}

// Estrutura para consulta SQL com JOIN entre campanhas_publicos e publicos_clientes
type CampanhaPublicoJoin struct {
	IDCampanha  int    `json:"id_campanha"`
	IDPublico   int    `json:"id_publico"`
	Nome        string `json:"nome"`
	Descricao   string `json:"descricao"`
	DataCriacao string `json:"data_criacao"`
	Status      string `json:"status"`
}

// Função para construir entidade CampanhaPublico
func BuildCampanhaPublicoEntity(idCampanha, idPublico int) *CampanhaPublico {
	return &CampanhaPublico{
		IDCampanha: idCampanha,
		IDPublico:  idPublico,
	}
}
