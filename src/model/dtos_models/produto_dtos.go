package dtos_models

type ProductResponse struct {
	ID            int     `json:"id"`
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
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}
