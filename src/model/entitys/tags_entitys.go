package entity

import dtos "github.com/betine97/back-project.git/src/model/dtos"

type Tag struct {
	IDTag        int    `gorm:"primaryKey;autoIncrement;column:id_tag" json:"id_tag"`
	CategoriaTag string `gorm:"column:categoria_tag;not null" json:"categoria_tag"`
	NomeTag      string `gorm:"column:nome_tag;not null" json:"nome_tag"`
}

// TableName especifica o nome da tabela
func (Tag) TableName() string {
	return "tags"
}

func BuildTagEntity(request dtos.CreateTagRequest) *Tag {
	return &Tag{
		CategoriaTag: request.CategoriaTag,
		NomeTag:      request.NomeTag,
	}
}
