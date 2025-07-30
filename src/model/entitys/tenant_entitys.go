package entity

type Tenants struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint   `gorm:"not null" json:"user_id"`
	NomeEmpresa string `gorm:"not null" json:"nome_empresa"`
	DBName      string `gorm:"not null" json:"db_name"`
	DBUser      string `gorm:"not null" json:"db_user"`
	DBPassword  string `gorm:"not null" json:"db_password"`
	DBHost      string `gorm:"not null" json:"db_host"`
	CreatedAt   string `gorm:"not null" json:"created_at"`
}
