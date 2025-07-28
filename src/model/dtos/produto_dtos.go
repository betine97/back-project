package dtos

// ProductResponse represents the response structure for product data
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

// ProductListResponse represents the response for product list
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

// ProductFilters represents filters for product search
type ProductFilters struct {
	Categoria     string  `json:"categoria,omitempty"`
	DestinadoPara string  `json:"destinado_para,omitempty"`
	Marca         string  `json:"marca,omitempty"`
	Variacao      string  `json:"variacao,omitempty"`
	Status        string  `json:"status,omitempty"`
	MinPrice      float64 `json:"min_price,omitempty"`
	MaxPrice      float64 `json:"max_price,omitempty"`
	Search        string  `json:"search,omitempty"`
}

// ProductQueryParams represents query parameters for product listing
type ProductQueryParams struct {
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
	Filter *ProductFilters `json:"filter,omitempty"`
}
