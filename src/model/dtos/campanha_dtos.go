package dtos

// Para GET api/campanhas
type CampanhaResponse struct {
	ID             int    `json:"id"`
	Nome           string `json:"nome"`
	Desc           string `json:"desc"`
	DataCriacao    string `json:"data_criacao"`
	DataLancamento string `json:"data_lancamento"`
	DataFim        string `json:"data_fim"`
	Status         string `json:"status"`
}

type CampanhaListResponse struct {
	Campanhas  []CampanhaResponse `json:"campanhas"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// Para POST api/campanhas
type CreateCampanhaRequest struct {
	Nome           string `json:"nome" validate:"required,min=2,max=255"`
	Desc           string `json:"desc" validate:"max=1000"`
	DataCriacao    string `json:"data_criacao" validate:"required,max=20"`
	DataLancamento string `json:"data_lancamento" validate:"required,max=20"`
	DataFim        string `json:"data_fim" validate:"required,max=20"`
	Status         string `json:"status" validate:"required,oneof=ativa inativa pausada finalizada"`
}

// Para POST api/campanhas/:id/publicos - Associar públicos a uma campanha
type AssociarPublicosCampanhaRequest struct {
	Publicos []int `json:"publicos" validate:"required,min=1"`
}

// Para GET api/campanhas/:id/publicos - Listar públicos de uma campanha
type PublicoCampanhaResponse struct {
	ID          int    `json:"id"`
	Nome        string `json:"nome"`
	Descricao   string `json:"descricao"`
	DataCriacao string `json:"data_criacao"`
	Status      string `json:"status"`
}

type PublicosCampanhaListResponse struct {
	Publicos   []PublicoCampanhaResponse `json:"publicos"`
	Total      int                       `json:"total"`
	IDCampanha int                       `json:"id_campanha"`
}
