package dtos

// DTO para resposta de pets (GET /api/pets)
type PetResponse struct {
	IDPet           int    `json:"id_pet"`
	ClienteID       int    `json:"cliente_id"`
	NomePet         string `json:"nome_pet"`
	Especie         string `json:"especie"`
	Raca            string `json:"raca"`
	Porte           string `json:"porte"`
	DataAniversario string `json:"data_aniversario"`
	Idade           *int   `json:"idade"`
	DataRegistro    string `json:"data_registro"`
}

type PetListResponse struct {
	Pets       []PetResponse `json:"pets"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}

// DTO para criação de pet (POST /api/pets)
type CreatePetRequest struct {
	ClienteID       int    `json:"cliente_id" validate:"required,gt=0"`
	NomePet         string `json:"nome_pet" validate:"required,min=2,max=100"`
	Especie         string `json:"especie" validate:"required,min=2,max=50"`
	Raca            string `json:"raca" validate:"max=50"`
	Porte           string `json:"porte" validate:"required,oneof=Pequeno Médio Grande"`
	DataAniversario string `json:"data_aniversario" validate:"omitempty"`
	Idade           *int   `json:"idade" validate:"omitempty,min=0,max=50"`
}
