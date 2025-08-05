package dtos

// Para GET api/pedidos/:id/itens (usando view_detalhes_pedido)
type DetalhesPedidoResponse struct {
	IDPedido      int     `json:"id_pedido"`
	NomeProduto   string  `json:"nome_produto"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
	TotalItem     float64 `json:"total_item"`
}

// Para GET api/pedidos/:id/itens
type ItemPedidoResponse struct {
	ID            int     `json:"id_item"`
	IDPedido      int     `json:"id_pedido"`
	IDProduto     int     `json:"id_produto"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
	Subtotal      float64 `json:"subtotal"`
}

type DetalhesPedidoListResponse struct {
	Detalhes   []DetalhesPedidoResponse `json:"detalhes"`
	Total      int                      `json:"total"`
	IDPedido   int                      `json:"id_pedido"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}

type ItemPedidoListResponse struct {
	Itens      []ItemPedidoResponse `json:"itens"`
	Total      int                  `json:"total"`
	IDPedido   int                  `json:"id_pedido"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"total_pages"`
}

// Para POST api/pedidos/:id/itens
type CreateItemPedidoRequest struct {
	IDProduto     int     `json:"id_produto" validate:"required,gt=0"`
	Quantidade    int     `json:"quantidade" validate:"required,gt=0"`
	PrecoUnitario float64 `json:"preco_unitario" validate:"required,gt=0"`
	Subtotal      float64 `json:"subtotal" validate:"required,gt=0"`
}
