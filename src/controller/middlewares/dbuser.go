package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// DatabaseExtractIdUser creates a middleware that extracts user ID from JWT claims
func DatabaseExtractIdUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// O token já foi validado pelo middleware JWTProtected()
		// Agora só precisamos extrair o userID dos claims
		userClaims, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Extrair o userID dos claims (pode ser "sub", "id", ou outro campo dependendo de como o token foi criado)
		var userID string
		if sub, exists := userClaims["sub"]; exists {
			userID = fmt.Sprintf("%v", sub)
		} else if id, exists := userClaims["id"]; exists {
			userID = fmt.Sprintf("%v", id)
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found in token"})
		}

		zap.L().Info("✅ User ID extraído do token", zap.String("user_id", userID))

		// Armazena apenas o userID no contexto
		c.Locals("userID", userID)
		return c.Next()
	}
}
