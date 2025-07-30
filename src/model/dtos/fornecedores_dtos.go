package dtos

type CreateFornecedorRequest struct {
	Nome         string `json:"nome"`
	Telefone     string `json:"telefone"`
	Email        string `json:"email"`
	Cidade       string `json:"cidade"`
	Estado       string `json:"estado"`
	Status       string `json:"status"`
	DataCadastro string `json:"data_cadastro"`
}

type UpdateFornecedorRequest struct {
	Campo string `json:"campo"` // Nome do campo a ser alterado
	Valor string `json:"valor"` // Novo valor para o campo
}

type FornecedorResponse struct {
	ID           int    `json:"id_fornecedor"`
	Nome         string `json:"nome"`
	Telefone     string `json:"telefone"`
	Email        string `json:"email"`
	Cidade       string `json:"cidade"`
	Estado       string `json:"estado"`
	Status       string `json:"status"`
	DataCadastro string `json:"data_cadastro"`
}

type FornecedorListResponse struct {
	Fornecedores []FornecedorResponse `json:"fornecedores"`
	Total        int                  `json:"total"`
}
