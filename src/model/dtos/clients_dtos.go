package dtos

type CreateClientRequest struct {
	CPF            string `json:"cpf"`
	Nome           string `json:"nome"`
	Sobrenome      string `json:"sobrenome"`
	NumeroCelular  string `json:"numero_celular"`
	Email          string `json:"email"`
	DataNascimento string `json:"data_nascimento"`
}

type CreateClientResponse struct {
	ID             int    `json:"id"`
	CPF            string `json:"cpf"`
	Nome           string `json:"nome"`
	Sobrenome      string `json:"sobrenome"`
	NumeroCelular  string `json:"numero_celular"`
	Email          string `json:"email"`
	DataNascimento string `json:"data_nascimento"`
}
