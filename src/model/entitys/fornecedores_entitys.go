package entity

type Fornecedor struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	Nome         string `gorm:"column:nome;not null" json:"nome"`
	Telefone     string `gorm:"column:telefone" json:"telefone"`
	Email        string `gorm:"column:email" json:"email"`
	Endereco     string `gorm:"column:endereco" json:"endereco"`
	Cidade       string `gorm:"column:cidade" json:"cidade"`
	Estado       string `gorm:"column:estado" json:"estado"`
	CEP          string `gorm:"column:cep" json:"cep"`
	DataCadastro string `gorm:"column:data_cadastro" json:"data_cadastro"`
	Status       string `gorm:"column:status" json:"status"`
}

func (Fornecedor) TableName() string {
	return "fornecedores"
}
