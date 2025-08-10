package dtos

// Para GET api/enderecos
type EnderecoResponse struct {
	IDEndereco  int    `json:"id_endereco"`
	IDCliente   int    `json:"id_cliente"`
	CEP         string `json:"cep"`
	Cidade      string `json:"cidade"`
	Estado      string `json:"estado"`
	Bairro      string `json:"bairro"`
	Logradouro  string `json:"logradouro"`
	Numero      string `json:"numero"`
	Complemento string `json:"complemento"`
	Obs         string `json:"obs"`
}

type EnderecoListResponse struct {
	Enderecos  []EnderecoResponse `json:"enderecos"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// Para POST api/enderecos
type CreateEnderecoRequest struct {
	IDCliente   int    `json:"id_cliente" validate:"required,gt=0"`
	CEP         string `json:"cep" validate:"required,max=20"`
	Cidade      string `json:"cidade" validate:"required,max=100"`
	Estado      string `json:"estado" validate:"required,max=50"`
	Bairro      string `json:"bairro" validate:"required,max=100"`
	Logradouro  string `json:"logradouro" validate:"required,max=255"`
	Numero      string `json:"numero" validate:"required,max=20"`
	Complemento string `json:"complemento" validate:"max=255"`
	Obs         string `json:"obs" validate:"max=255"`
}
