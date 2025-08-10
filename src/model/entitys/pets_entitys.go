package entity

import (
	"time"

	dtos "github.com/betine97/back-project.git/src/model/dtos"
)

// Entidade para a tabela pets
type Pet struct {
	IDPet           int        `gorm:"primaryKey;autoIncrement;column:id_pet" json:"id_pet"`
	ClienteID       int        `gorm:"column:cliente_id;not null" json:"cliente_id"`
	NomePet         string     `gorm:"column:nome_pet;not null" json:"nome_pet"`
	Especie         string     `gorm:"column:especie;not null" json:"especie"`
	Raca            string     `gorm:"column:raca" json:"raca"`
	Porte           string     `gorm:"column:porte;not null;type:enum('Pequeno','Médio','Grande')" json:"porte"`
	DataAniversario *time.Time `gorm:"column:data_aniversario;type:date" json:"data_aniversario"`
	Idade           *int       `gorm:"column:idade" json:"idade"`
	DataRegistro    time.Time  `gorm:"column:data_registro;not null;default:CURRENT_TIMESTAMP" json:"data_registro"`
}

// TableName especifica o nome da tabela para GORM
func (Pet) TableName() string {
	return "pets"
}

// Função para construir entidade Pet a partir do DTO
func BuildPetEntity(request dtos.CreatePetRequest) *Pet {
	pet := &Pet{
		ClienteID:    request.ClienteID,
		NomePet:      request.NomePet,
		Especie:      request.Especie,
		Raca:         request.Raca,
		Porte:        request.Porte,
		DataRegistro: time.Now(),
	}

	// Converter data de aniversário se fornecida
	if request.DataAniversario != "" {
		if dataAniversario, err := time.Parse("2006-01-02", request.DataAniversario); err == nil {
			pet.DataAniversario = &dataAniversario
		}
	}

	// Definir idade se fornecida
	if request.Idade != nil {
		pet.Idade = request.Idade
	}

	return pet
}
