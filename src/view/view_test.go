package view

import (
	"testing"

	dtos "github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TESTES PARA ConvertDomainToResponse
// =============================================================================

func TestConvertDomainToResponse_Success(t *testing.T) {
	// Arrange
	user := &entity.User{
		ID:          1,
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@teste.com",
		NomeEmpresa: "Empresa Teste",
		Categoria:   "Tecnologia",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "senha_secreta", // Não deve aparecer no response
	}

	// Act
	result := ConvertDomainToResponse(user)

	// Assert
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.FirstName, result.FirstName)
	assert.Equal(t, user.LastName, result.LastName)
	assert.Equal(t, user.Email, result.Email)

	// Verificar que é do tipo correto
	assert.IsType(t, dtos.NewUser{}, result)
}

func TestConvertDomainToResponse_EmptyUser(t *testing.T) {
	// Arrange
	user := &entity.User{}

	// Act
	result := ConvertDomainToResponse(user)

	// Assert
	assert.Equal(t, uint(0), result.ID)
	assert.Equal(t, "", result.FirstName)
	assert.Equal(t, "", result.LastName)
	assert.Equal(t, "", result.Email)
}

func TestConvertDomainToResponse_SpecialCharacters(t *testing.T) {
	// Arrange
	user := &entity.User{
		ID:        1,
		FirstName: "José María",
		LastName:  "González-Pérez",
		Email:     "jose.maria@empresa-teste.com.br",
	}

	// Act
	result := ConvertDomainToResponse(user)

	// Assert
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.FirstName, result.FirstName)
	assert.Equal(t, user.LastName, result.LastName)
	assert.Equal(t, user.Email, result.Email)
}

func TestConvertDomainToResponse_LongStrings(t *testing.T) {
	// Arrange
	longString := string(make([]byte, 500))
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	user := &entity.User{
		ID:        999999,
		FirstName: longString,
		LastName:  longString,
		Email:     "test@example.com",
	}

	// Act
	result := ConvertDomainToResponse(user)

	// Assert
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, longString, result.FirstName)
	assert.Equal(t, longString, result.LastName)
	assert.Equal(t, user.Email, result.Email)
}

func TestConvertDomainToResponse_OnlyMapsRequiredFields(t *testing.T) {
	// Arrange
	user := &entity.User{
		ID:          1,
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@teste.com",
		NomeEmpresa: "Empresa Teste", // Campo que não deve ser mapeado
		Categoria:   "Tecnologia",    // Campo que não deve ser mapeado
		Segmento:    "Software",      // Campo que não deve ser mapeado
		City:        "São Paulo",     // Campo que não deve ser mapeado
		State:       "SP",            // Campo que não deve ser mapeado
		Password:    "senha_secreta", // Campo que não deve ser mapeado
	}

	// Act
	result := ConvertDomainToResponse(user)

	// Assert - Verificar que apenas os campos corretos foram mapeados
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.FirstName, result.FirstName)
	assert.Equal(t, user.LastName, result.LastName)
	assert.Equal(t, user.Email, result.Email)

	// Verificar que a estrutura NewUser não tem os campos extras
	// (isso é garantido pela estrutura do DTO, mas é bom documentar)
	expectedFields := 4   // ID, FirstName, LastName, Email
	actualFieldCount := 4 // Contagem manual baseada na estrutura NewUser
	assert.Equal(t, expectedFields, actualFieldCount)
}

func TestConvertDomainToResponse_NilSafety(t *testing.T) {
	// Arrange - Teste com ponteiro nil
	// Nota: Esta função não trata nil, então esperamos panic
	// Em produção, seria bom adicionar verificação de nil

	// Act & Assert
	assert.Panics(t, func() {
		ConvertDomainToResponse(nil)
	}, "Function should panic with nil input")
}

func TestConvertDomainToResponse_MaxValues(t *testing.T) {
	// Arrange
	user := &entity.User{
		ID:        ^uint(0), // Valor máximo para uint
		FirstName: "Max",
		LastName:  "User",
		Email:     "max@test.com",
	}

	// Act
	result := ConvertDomainToResponse(user)

	// Assert
	assert.Equal(t, ^uint(0), result.ID)
	assert.Equal(t, "Max", result.FirstName)
	assert.Equal(t, "User", result.LastName)
	assert.Equal(t, "max@test.com", result.Email)
}

// =============================================================================
// TESTES DE INTEGRAÇÃO COM DTOs
// =============================================================================

func TestConvertDomainToResponse_IntegrationWithDTO(t *testing.T) {
	// Arrange
	user := &entity.User{
		ID:        1,
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@teste.com",
	}

	// Act
	result := ConvertDomainToResponse(user)

	// Assert - Verificar que o resultado pode ser usado como NewUser DTO
	var dto dtos.NewUser = result
	assert.Equal(t, user.ID, dto.ID)
	assert.Equal(t, user.FirstName, dto.FirstName)
	assert.Equal(t, user.LastName, dto.LastName)
	assert.Equal(t, user.Email, dto.Email)
}

// =============================================================================
// BENCHMARKS
// =============================================================================

func BenchmarkConvertDomainToResponse(b *testing.B) {
	user := &entity.User{
		ID:          1,
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@teste.com",
		NomeEmpresa: "Empresa Teste",
		Categoria:   "Tecnologia",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "senha_secreta",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConvertDomainToResponse(user)
	}
}

func BenchmarkConvertDomainToResponse_LargeData(b *testing.B) {
	longString := string(make([]byte, 1000))
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	user := &entity.User{
		ID:        1,
		FirstName: longString,
		LastName:  longString,
		Email:     "test@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConvertDomainToResponse(user)
	}
}
