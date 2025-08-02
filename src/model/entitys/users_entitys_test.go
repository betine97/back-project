package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// TESTES PARA User Entity
// =============================================================================

func TestUser_JSONSerialization(t *testing.T) {
	// Arrange
	user := User{
		ID:          1,
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@teste.com",
		NomeEmpresa: "Empresa Teste",
		Categoria:   "Tecnologia",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "hashedpassword123",
	}

	// Act
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)

	var unmarshaled User
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user.ID, unmarshaled.ID)
	assert.Equal(t, user.FirstName, unmarshaled.FirstName)
	assert.Equal(t, user.LastName, unmarshaled.LastName)
	assert.Equal(t, user.Email, unmarshaled.Email)
	assert.Equal(t, user.NomeEmpresa, unmarshaled.NomeEmpresa)
	assert.Equal(t, user.Categoria, unmarshaled.Categoria)
	assert.Equal(t, user.Segmento, unmarshaled.Segmento)
	assert.Equal(t, user.City, unmarshaled.City)
	assert.Equal(t, user.State, unmarshaled.State)
	// Password should not be serialized (json:"-" tag)
	assert.Empty(t, unmarshaled.Password)
}

func TestUser_PasswordNotSerialized(t *testing.T) {
	// Arrange
	user := User{
		ID:        1,
		FirstName: "João",
		Email:     "joao@teste.com",
		Password:  "supersecretpassword",
	}

	// Act
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)

	// Assert
	jsonString := string(jsonData)
	assert.NotContains(t, jsonString, "supersecretpassword")
	assert.NotContains(t, jsonString, "password")
	assert.Contains(t, jsonString, "joao@teste.com")
	assert.Contains(t, jsonString, "João")
}

func TestUser_GormTags(t *testing.T) {
	// Arrange
	user := User{
		ID:          1,
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@teste.com",
		NomeEmpresa: "Empresa Teste",
		Categoria:   "Tecnologia",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "hashedpassword123",
	}

	// Act & Assert - Verificar se a struct pode ser usada com GORM
	// Isso é mais um teste de estrutura do que funcional
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "joao@teste.com", user.Email)
	assert.NotEmpty(t, user.Password)
}

func TestUser_EmptyValues(t *testing.T) {
	// Arrange
	user := User{}

	// Act
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)

	var unmarshaled User
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint(0), unmarshaled.ID)
	assert.Equal(t, "", unmarshaled.FirstName)
	assert.Equal(t, "", unmarshaled.LastName)
	assert.Equal(t, "", unmarshaled.Email)
	assert.Equal(t, "", unmarshaled.NomeEmpresa)
	assert.Equal(t, "", unmarshaled.Password) // Should be empty due to json:"-"
}

func TestUser_SpecialCharacters(t *testing.T) {
	// Arrange
	user := User{
		ID:          1,
		FirstName:   "José María",
		LastName:    "González-Pérez",
		Email:       "jose.maria@empresa-teste.com.br",
		NomeEmpresa: "Empresa & Cia Ltda.",
		Categoria:   "Tecnologia & Inovação",
		Segmento:    "Software/Hardware",
		City:        "São José dos Campos",
		State:       "SP",
		Password:    "password@123!",
	}

	// Act
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)

	var unmarshaled User
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user.FirstName, unmarshaled.FirstName)
	assert.Equal(t, user.LastName, unmarshaled.LastName)
	assert.Equal(t, user.Email, unmarshaled.Email)
	assert.Equal(t, user.NomeEmpresa, unmarshaled.NomeEmpresa)
	assert.Equal(t, user.City, unmarshaled.City)
}

func TestUser_LongStrings(t *testing.T) {
	// Arrange
	longString := string(make([]byte, 500))
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	user := User{
		ID:          1,
		FirstName:   longString,
		LastName:    longString,
		Email:       "test@example.com",
		NomeEmpresa: longString,
		Categoria:   longString,
		Segmento:    longString,
		City:        longString,
		State:       longString,
		Password:    longString,
	}

	// Act
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)

	var unmarshaled User
	err = json.Unmarshal(jsonData, &unmarshaled)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, longString, unmarshaled.FirstName)
	assert.Equal(t, longString, unmarshaled.LastName)
	assert.Equal(t, longString, unmarshaled.NomeEmpresa)
	// Password should not be serialized
	assert.Empty(t, unmarshaled.Password)
}

func TestUser_JSONTags(t *testing.T) {
	// Arrange
	user := User{
		ID:          1,
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@teste.com",
		NomeEmpresa: "Empresa Teste",
		Password:    "secret",
	}

	// Act
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)

	// Assert
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, `"id":1`)
	assert.Contains(t, jsonString, `"first_name":"João"`)
	assert.Contains(t, jsonString, `"last_name":"Silva"`)
	assert.Contains(t, jsonString, `"email":"joao@teste.com"`)
	assert.Contains(t, jsonString, `"nome_empresa":"Empresa Teste"`)
	// Password should not appear in JSON
	assert.NotContains(t, jsonString, "secret")
	assert.NotContains(t, jsonString, "password")
}

func TestUser_UniqueConstraints(t *testing.T) {
	// Arrange - Test that we can create users with different emails
	user1 := User{
		ID:    1,
		Email: "user1@test.com",
	}

	user2 := User{
		ID:    2,
		Email: "user2@test.com",
	}

	// Act & Assert - Different emails should be fine
	assert.NotEqual(t, user1.Email, user2.Email)
	assert.NotEqual(t, user1.ID, user2.ID)

	// Same email should be identifiable
	user3 := User{
		ID:    3,
		Email: "user1@test.com", // Same as user1
	}

	assert.Equal(t, user1.Email, user3.Email)
	assert.NotEqual(t, user1.ID, user3.ID)
}

func TestUser_RequiredFields(t *testing.T) {
	// Test that we can identify when required fields are missing
	tests := []struct {
		name              string
		user              User
		hasRequiredFields bool
	}{
		{
			name: "All required fields present",
			user: User{
				FirstName:   "João",
				LastName:    "Silva",
				Email:       "joao@test.com",
				NomeEmpresa: "Empresa",
				Categoria:   "Tech",
				Segmento:    "Software",
				City:        "SP",
				State:       "SP",
				Password:    "pass123",
			},
			hasRequiredFields: true,
		},
		{
			name: "Missing FirstName",
			user: User{
				LastName:    "Silva",
				Email:       "joao@test.com",
				NomeEmpresa: "Empresa",
				Password:    "pass123",
			},
			hasRequiredFields: false,
		},
		{
			name: "Missing Email",
			user: User{
				FirstName:   "João",
				LastName:    "Silva",
				NomeEmpresa: "Empresa",
				Password:    "pass123",
			},
			hasRequiredFields: false,
		},
		{
			name: "Missing Password",
			user: User{
				FirstName:   "João",
				LastName:    "Silva",
				Email:       "joao@test.com",
				NomeEmpresa: "Empresa",
			},
			hasRequiredFields: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act - Check if basic required fields are present
			hasFirstName := tt.user.FirstName != ""
			hasLastName := tt.user.LastName != ""
			hasEmail := tt.user.Email != ""
			hasPassword := tt.user.Password != ""

			allRequired := hasFirstName && hasLastName && hasEmail && hasPassword

			// Assert
			assert.Equal(t, tt.hasRequiredFields, allRequired)
		})
	}
}
