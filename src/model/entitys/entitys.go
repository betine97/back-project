package entity

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string `gorm:"column:first_name;not null" json:"first_name"`
	LastName  string `gorm:"column:last_name;not null" json:"last_name"`
	Email     string `gorm:"not null;unique" json:"email"`
	City      string `gorm:"not null" json:"city"`
	Password  string `gorm:"not null" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type Produto struct {
	ID            int     `gorm:"primaryKey" json:"id"`
	CodigoBarra   string  `gorm:"column:codigo_barra;not null" json:"codigo_barra"`
	NomeProduto   string  `gorm:"column:nome_produto;not null" json:"nome_produto"`
	SKU           string  `gorm:"not null" json:"sku"`
	Categoria     string  `gorm:"not null" json:"categoria"`
	DestinadoPara string  `gorm:"column:destinado_para;not null" json:"destinado_para"`
	Variacao      string  `json:"variacao"`
	Marca         string  `json:"marca"`
	Descricao     string  `json:"descricao"`
	Status        string  `gorm:"not null" json:"status"`
	PrecoVenda    float64 `gorm:"column:preco_venda;not null" json:"preco_venda"`
}

func (Produto) TableName() string {
	return "produtos"
}
