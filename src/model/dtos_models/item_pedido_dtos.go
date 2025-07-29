package dtos_models

type CreateItemPedidoRequest struct {
	IDPedido      int     `json:"id_pedido" validate:"required,min=1"`
	IDProduto     int     `json:"id_produto" validate:"required,min=1"`
	Quantidade    int     `json:"quantidade" validate:"required,min=1"`
	PrecoUnitario float64 `json:"preco_unitario" validate:"required,min=0.01"`
}

type ItemPedidoQueryParams struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type ItemPedidoResponse struct {
	IDItem        int     `json:"id_item"`
	IDPedido      int     `json:"id_pedido"`
	IDProduto     int     `json:"id_produto"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
	Subtotal      float64 `json:"subtotal"`
}

type ItemPedidoListResponse struct {
	ItemPedidos []ItemPedidoResponse `json:"item_pedidos"`
	Total       int                  `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
}
