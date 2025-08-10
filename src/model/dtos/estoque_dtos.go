package dtos

// Para GET api/estoque
type EstoqueResponse struct {
	IDEstoque           int     `json:"id_estoque"`
	IDProduto           int     `json:"id_produto"`
	IDLote              int     `json:"id_lote"`
	Quantidade          int     `json:"quantidade"`
	Vencimento          string  `json:"vencimento"`
	CustoUnitario       float64 `json:"custo_unitario"`
	DataEntrada         string  `json:"data_entrada"`
	DataSaida           string  `json:"data_saida"`
	DocumentoReferencia string  `json:"documento_referencia"`
	Status              string  `json:"status"`
}

type EstoqueListResponse struct {
	Estoque    []EstoqueResponse `json:"estoque"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// Para GET api/estoque usando view_detalhes_estoque
type DetalhesEstoqueResponse struct {
	NomeProduto         string `json:"nome_produto"`
	Lote                int    `json:"lote"`
	Quantidade          int    `json:"quantidade"`
	DataEntrada         string `json:"data_entrada"`
	DataSaida           string `json:"data_saida"`
	Vencimento          string `json:"vencimento"`
	DocumentoReferencia string `json:"documento_referencia"`
	Status              string `json:"status"`
}

type DetalhesEstoqueListResponse struct {
	Estoque    []DetalhesEstoqueResponse `json:"estoque"`
	Total      int                       `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

// Para POST api/estoque
type CreateEstoqueRequest struct {
	IDProduto           int     `json:"id_produto" validate:"required,gt=0"`
	IDLote              int     `json:"id_lote" validate:"required,gt=0"`
	Quantidade          int     `json:"quantidade" validate:"required,gte=0"`
	Vencimento          string  `json:"vencimento"`
	CustoUnitario       float64 `json:"custo_unitario" validate:"required,gt=0"`
	DataEntrada         string  `json:"data_entrada"`
	DataSaida           string  `json:"data_saida"`
	DocumentoReferencia string  `json:"documento_referencia"`
	Status              string  `json:"status" validate:"required"`
}
