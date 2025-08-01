package crypto

import (
	"crypto/rand"
	"math/big"
	"testing"
)

// generateRandomPassword gera uma senha aleatória com o tamanho especificado
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)

	for i := range password {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[num.Int64()]
	}

	return string(password)
}

func TestCrypto_HashAndCheckPassword(t *testing.T) {
	crypto := &Crypto{}

	// Tabela de testes com 5 casos diferentes
	testCases := []struct {
		name     string
		password string
	}{
		{
			name:     "Short password",
			password: generateRandomPassword(8),
		},
		{
			name:     "Medium password",
			password: generateRandomPassword(16),
		},
		{
			name:     "Long password",
			password: generateRandomPassword(32),
		},
		{
			name:     "Password with special chars",
			password: generateRandomPassword(20),
		},
		{
			name:     "Very long password",
			password: generateRandomPassword(64),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			originalPassword := tc.password

			// Act - Hash the password
			hashedPassword, err := crypto.HashPassword(originalPassword)

			// Assert - Hash should succeed
			if err != nil {
				t.Errorf("HashPassword failed: %v", err)
				return
			}

			if hashedPassword == "" {
				t.Error("HashPassword returned empty string")
				return
			}

			// Act - Check the password
			isValid, err := crypto.CheckPassword(originalPassword, hashedPassword)

			// Assert - Check should succeed and return true
			if err != nil {
				t.Errorf("CheckPassword failed: %v", err)
				return
			}

			if !isValid {
				t.Error("CheckPassword returned false for correct password")
				return
			}

			// Additional test - wrong password should fail
			wrongPassword := originalPassword + "wrong"
			isValidWrong, err := crypto.CheckPassword(wrongPassword, hashedPassword)

			// CheckPassword should return false and an error for wrong password
			if err == nil {
				t.Error("CheckPassword should return error for incorrect password")
				return
			}

			if isValidWrong {
				t.Error("CheckPassword returned true for incorrect password")
			}
		})
	}
}

func TestCrypto_HashPassword_EmptyPassword(t *testing.T) {
	crypto := &Crypto{}

	// Test with empty password
	hashedPassword, err := crypto.HashPassword("")

	if err != nil {
		t.Errorf("HashPassword with empty string failed: %v", err)
	}

	if hashedPassword == "" {
		t.Error("HashPassword returned empty string for empty input")
	}
}

func TestCrypto_CheckPassword_InvalidHash(t *testing.T) {
	crypto := &Crypto{}

	// Test with invalid hash
	isValid, err := crypto.CheckPassword("password", "invalid_hash")

	if err == nil {
		t.Error("CheckPassword should fail with invalid hash")
	}

	if isValid {
		t.Error("CheckPassword should return false with invalid hash")
	}
}
