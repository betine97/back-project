package service

import (
	"testing"

	"github.com/betine97/back-project.git/src/model/dtos"
	"github.com/stretchr/testify/assert"
)

// Teste isolado para a função buildUserEntity
func TestBuildUserEntity_Isolated(t *testing.T) {
	// Arrange
	request := dtos.CreateUser{
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john@example.com",
		NomeEmpresa: "Test Company",
		Categoria:   "Tech",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "password123",
	}
	hashedPassword := "hashedpassword123"

	// Act
	user := buildUserEntity(request, hashedPassword)

	// Assert
	assert.Equal(t, request.FirstName, user.FirstName)
	assert.Equal(t, request.LastName, user.LastName)
	assert.Equal(t, request.Email, user.Email)
	assert.Equal(t, request.NomeEmpresa, user.NomeEmpresa)
	assert.Equal(t, request.Categoria, user.Categoria)
	assert.Equal(t, request.Segmento, user.Segmento)
	assert.Equal(t, request.City, user.City)
	assert.Equal(t, request.State, user.State)
	assert.Equal(t, hashedPassword, user.Password)
	assert.Equal(t, uint(0), user.ID) // ID deve ser 0 para novos usuários
}
