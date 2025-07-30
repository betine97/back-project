package dtos

type CreateUser struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=100" example:"John"`
	LastName    string `json:"last_name" validate:"required,min=2,max=100" example:"Doe"`
	Email       string `json:"email" validate:"required,email" example:"test@test.com"`
	NomeEmpresa string `json:"nome_empresa" validate:"required,min=2,max=100" example:"Empresa Teste"`
	Categoria   string `json:"categoria" validate:"required,min=2,max=100" example:"Categoria Teste"`
	Segmento    string `json:"segmento" validate:"required,min=2,max=100" example:"Segmento Teste"`
	City        string `json:"city" validate:"required,min=2,max=100" example:"São Paulo"`
	State       string `json:"state" validate:"required,min=2,max=100" example:"SP"`
	Password    string `json:"password" validate:"required,min=6,containsany=!@#$%*" example:"password#@#@!2121"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email" example:"test@test.com"`
	Password string `json:"password" binding:"required,min=6,containsany=!@#$%*" example:"password#@#@!2121"`
}

type NewUser struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
