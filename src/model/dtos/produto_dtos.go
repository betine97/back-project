package dtos

// Para GET api/produtos

type ProductResponse struct {
	ID            int     `json:"id_produto"`
	CodigoBarra   string  `json:"codigo_barra"`
	NomeProduto   string  `json:"nome_produto"`
	SKU           string  `json:"sku"`
	Categoria     string  `json:"categoria"`
	DestinadoPara string  `json:"destinado_para"`
	Variacao      string  `json:"variacao"`
	Marca         string  `json:"marca"`
	Descricao     string  `json:"descricao"`
	Status        string  `json:"status"`
	PrecoVenda    float64 `json:"preco_venda"`
	IDFornecedor  int     `json:"id_fornecedor"`
}

type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// Para POST api/produtos

type CreateProductRequest struct {
	DataCadastro  string  `json:"data_cadastro" validate:"required"`
	CodigoBarra   string  `json:"codigo_barra"`
	NomeProduto   string  `json:"nome_produto" validate:"required"`
	SKU           string  `json:"sku"`
	Categoria     string  `json:"categoria"`
	DestinadoPara string  `json:"destinado_para"`
	Variacao      string  `json:"variacao"`
	Marca         string  `json:"marca"`
	Descricao     string  `json:"descricao"`
	Status        string  `json:"status" validate:"required"`
	PrecoVenda    float64 `json:"preco_venda" validate:"required,gt=0"`
	IDFornecedor  int     `json:"id_fornecedor" validate:"required,gt=0"`
}
