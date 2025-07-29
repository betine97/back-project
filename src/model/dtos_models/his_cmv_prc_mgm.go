package dtos_models

type HisCmvPrcMargeResponse struct {
	ID           int     `json:"id"`
	IDProduto    int     `json:"id_produto"`
	PrecoVenda   float64 `json:"preco_venda"`
	Cmv          float64 `json:"cmv"`
	Margem       float64 `json:"margem"`
	DataRegistro string  `json:"data_registro"`
}

type HisCmvPrcMargeListResponse struct {
	HisCmvPrcMarge []HisCmvPrcMargeResponse `json:"his_cmv_prc_marge"`
	Total          int                      `json:"total"`
}

type HisCmvPrcMargeQueryParams struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}
