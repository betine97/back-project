package dtos

// Para GET api/tags
type TagResponse struct {
	IDTag        int    `json:"id_tag"`
	CategoriaTag string `json:"categoria_tag"`
	NomeTag      string `json:"nome_tag"`
}

type TagListResponse struct {
	Tags       []TagResponse `json:"tags"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}

// Para POST api/tags
type CreateTagRequest struct {
	CategoriaTag string `json:"categoria_tag" validate:"required,min=2,max=100"`
	NomeTag      string `json:"nome_tag" validate:"required,min=2,max=100"`
}
