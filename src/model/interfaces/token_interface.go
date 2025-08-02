package interfaces

import (
	"time"

	"github.com/betine97/back-project.git/cmd/config"
	"github.com/betine97/back-project.git/cmd/config/exceptions"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// TokenGeneratorInterface define a interface para geração de tokens
type TokenGeneratorInterface interface {
	GenerateToken(userID uint) (string, *exceptions.RestErr)
}

// JWTTokenGenerator implementa a geração de tokens JWT
type JWTTokenGenerator struct{}

func NewJWTTokenGenerator() TokenGeneratorInterface {
	return &JWTTokenGenerator{}
}

func (j *JWTTokenGenerator) GenerateToken(userID uint) (string, *exceptions.RestErr) {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(config.PrivateKey)
	if err != nil {
		zap.L().Error("Error signing token", zap.Error(err))
		return "", exceptions.NewInternalServerError("Internal server error")
	}

	return tokenString, nil
}
