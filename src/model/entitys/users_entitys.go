package entity

type User struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName   string `gorm:"column:first_name;not null" json:"first_name"`
	LastName    string `gorm:"column:last_name;not null" json:"last_name"`
	Email       string `gorm:"column:email;not null;unique" json:"email"`
	NomeEmpresa string `gorm:"column:nome_empresa;not null" json:"nome_empresa"`
	Categoria   string `gorm:"column:categoria;not null" json:"categoria"`
	Segmento    string `gorm:"column:segmento;not null" json:"segmento"`
	City        string `gorm:"column:city;not null" json:"city"`
	State       string `gorm:"column:state;not null" json:"state"`
	Password    string `gorm:"column:password;not null" json:"-"`
}
