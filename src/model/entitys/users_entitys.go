package entity

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string `gorm:"column:first_name;not null" json:"first_name"`
	LastName  string `gorm:"column:last_name;not null" json:"last_name"`
	Email     string `gorm:"not null;unique" json:"email"`
	City      string `gorm:"not null" json:"city"`
	Password  string `gorm:"not null" json:"-"`
}
