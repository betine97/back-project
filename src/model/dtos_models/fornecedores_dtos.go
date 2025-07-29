package dtos_models

type Fornecedor struct {
	ID           int    `json:"id_fornecedor"`
	Nome         string `json:"nome"`
	Telefone     string `json:"telefone"`
	Email        string `json:"email"`
	Endereco     string `json:"endereco"`
	Cidade       string `json:"cidade"`
	Estado       string `json:"estado"`
	CEP          string `json:"cep"`
	DataCadastro string `json:"data_cadastro"`
	Status       string `json:"status"`
}

type FornecedorListResponse struct {
	Fornecedores []FornecedorResponse `json:"fornecedores"`
	Total        int                  `json:"total"`
}

type FornecedorResponse struct {
	ID           int    `json:"id"`
	Nome         string `json:"nome"`
	Telefone     string `json:"telefone"`
	Email        string `json:"email"`
	Endereco     string `json:"endereco"`
	Cidade       string `json:"cidade"`
	Estado       string `json:"estado"`
	CEP          string `json:"cep"`
	DataCadastro string `json:"data_cadastro"`
	Status       string `json:"status"`
}

type FornecedorQueryParams struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}
