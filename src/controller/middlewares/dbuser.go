package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/betine97/back-project.git/src/model/persistence"
	redis "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

var ctx = context.Background()

func DatabaseConnectionMiddleware(db persistence.PersistenceInterface, redisClient *redis.Client) fiber.Handler {
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

		cacheKey := fmt.Sprintf("user:%s:db_info", userID)
		tenantJSON, err := redisClient.Get(ctx, cacheKey).Result()
		if err == redis.Nil {
			zap.L().Error("Tenant not found in cache")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Tenant not found in cache"})
		} else if err != nil {
			zap.L().Error("Error retrieving from Redis", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving from cache"})
		}

		var cachedTenant entity.Tenants
		if err := json.Unmarshal([]byte(tenantJSON), &cachedTenant); err != nil {
			zap.L().Error("Error unmarshaling tenant from JSON", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving tenant data"})
		}

		// Estabeleça a conexão com o banco de dados do cliente
		db, err := NewDatabaseConnectionForTenant(cachedTenant)
		if err != nil {
			zap.L().Error("Error connecting to tenant database", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error connecting to tenant database"})
		}

		// Aqui você pode armazenar a conexão no contexto, se necessário
		c.Locals("db", db)

		return c.Next()
	}
}
