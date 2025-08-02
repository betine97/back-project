package service

import (
	"testing"

	"github.com/betine97/back-project.git/src/model/dtos"
	"github.com/stretchr/testify/assert"
)

// Teste simples que não depende de configurações externas
func TestBuildUserEntity_Simple(t *testing.T) {
	// Arrange
	request := dtos.CreateUser{
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@teste.com",
		NomeEmpresa: "Empresa Teste",
		Categoria:   "Tecnologia",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "senha123",
	}
	hashedPassword := "senha_hash_123"

	// Act
	user := buildUserEntity(request, hashedPassword)

	// Assert
	assert.Equal(t, "João", user.FirstName)
	assert.Equal(t, "Silva", user.LastName)
	assert.Equal(t, "joao@teste.com", user.Email)
	assert.Equal(t, "Empresa Teste", user.NomeEmpresa)
	assert.Equal(t, "Tecnologia", user.Categoria)
	assert.Equal(t, "Software", user.Segmento)
	assert.Equal(t, "São Paulo", user.City)
	assert.Equal(t, "SP", user.State)
	assert.Equal(t, "senha_hash_123", user.Password)
	assert.Equal(t, uint(0), user.ID)
}
