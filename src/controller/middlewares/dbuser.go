package middlewares

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func DatabaseExtractIdUser() fiber.Handler {

	err := godotenv.Load("C:/Users/erons/life/back_life/cmd/config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))

	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is required"})
		}

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})

		}

		userID := claims.Id

		c.Locals("userID", userID) // Opcional, se você ainda quiser usar o contexto do Fiber

		return c.Next()
	}
}
