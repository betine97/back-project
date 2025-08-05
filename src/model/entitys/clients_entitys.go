package entity

import dtos "github.com/betine97/back-project.git/src/model/dtos"

type Client struct {
	CPF            string `json:"cpf"`
	Nome           string `json:"nome"`
	Sobrenome      string `json:"sobrenome"`
	NumeroCelular  string `json:"numero_celular"`
	Email          string `json:"email"`
	DataNascimento string `json:"data_nascimento"`
}

func BuildClientEntity(request dtos.CreateClientRequest) *Client {
	return &Client{
		CPF:            request.CPF,
		Nome:           request.Nome,
		Sobrenome:      request.Sobrenome,
		NumeroCelular:  request.NumeroCelular,
		Email:          request.Email,
		DataNascimento: request.DataNascimento,
	}
}
