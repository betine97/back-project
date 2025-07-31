package middlewares

import (
	"fmt"
	"strconv"

	"github.com/betine97/back-project.git/cmd/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

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

		// Buscar o tenant_id baseado no user_id
		masterDB, err := config.NewDatabaseConnection()
		if err != nil {
			zap.L().Error("❌ Falha ao conectar com banco master", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de conexão com banco"})
		}

		var tenant struct {
			ID uint `json:"id"`
		}

		userIDInt, _ := strconv.Atoi(userID)
		err = masterDB.Table("tenants").Select("id").Where("user_id = ?", userIDInt).First(&tenant).Error
		if err != nil {
			zap.L().Error("❌ Tenant não encontrado", zap.String("user_id", userID), zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Tenant não encontrado"})
		}

		tenantID := fmt.Sprintf("%d", tenant.ID)
		zap.L().Info("✅ Tenant encontrado", zap.String("user_id", userID), zap.String("tenant_id", tenantID))

		c.Locals("userID", userID)
		c.Locals("tenantID", tenantID)
		return c.Next()
	}
}
