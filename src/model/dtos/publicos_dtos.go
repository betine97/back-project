package dtos

// DTO para resposta de critérios (GET /api/criterios)
type CriterioResponse struct {
	ID           int    `json:"id"`
	NomeCondicao string `json:"nome_condicao"`
}

type CriterioListResponse struct {
	Criterios []CriterioResponse `json:"criterios"`
	Total     int                `json:"total"`
}

// DTO para resposta de públicos (GET /api/publicos)
type PublicoResponse struct {
	ID          int    `json:"id"`
	Nome        string `json:"nome"`
	Descricao   string `json:"descricao"`
	DataCriacao string `json:"data_criacao"`
	Status      string `json:"status"`
}

type PublicoListResponse struct {
	Publicos   []PublicoResponse `json:"publicos"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// DTO para criação de público (POST /api/publicos)
type CreatePublicoRequest struct {
	Nome        string `json:"nome" validate:"required,min=2,max=100"`
	Descricao   string `json:"descricao" validate:"required,min=2,max=100"`
	DataCriacao string `json:"data_criacao" validate:"required,max=100"`
	Status      string `json:"status" validate:"required,max=50"`
}

// DTO para associar critérios ao público (POST /api/publicos/:id/criterios)
type AssociarCriteriosRequest struct {
	Criterios []int `json:"criterios" validate:"required,min=1"`
}

// DTO para resposta dos critérios de um público (GET /api/publicos/:id/criterios)
type PublicoCriterioResponse struct {
	IDPublico    int    `json:"id_publico"`
	IDCriterio   int    `json:"id_criterio"`
	NomeCondicao string `json:"nome_condicao"`
}

type PublicoCriterioListResponse struct {
	Criterios []PublicoCriterioResponse `json:"criterios"`
	Total     int                       `json:"total"`
}
