package middlewares

import (
	"strings"

	"github.com/betine97/back-project.git/cmd/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		zap.L().Info("🔐 Iniciando proteção JWT", zap.String("ip", ip))
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			zap.L().Warn("⚠️ Token ausente ou malformado", zap.String("ip", ip))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token ausente ou malformado"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				zap.L().Warn("⚠️ Método de assinatura inesperado", zap.String("ip", ip))
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Método de assinatura inesperado")
			}
			return &config.PrivateKey.PublicKey, nil

		})

		if err != nil || !token.Valid {
			zap.L().Warn("❌ Token JWT inválido", zap.Error(err), zap.String("ip", ip))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token inválido"})
		}

		zap.L().Info("✅ Token JWT válido", zap.String("ip", ip))
		c.Locals("user", token.Claims)
		return c.Next()
	}
}

func JWTClaimsRequired(claimKey string, claimValue string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		zap.L().Info("🔍 Verificando claims JWT", zap.String("ip", ip), zap.String("claim", claimKey), zap.String("valor", claimValue))
		userClaims := c.Locals("user").(jwt.MapClaims)
		if userClaims[claimKey] != claimValue {
			zap.L().Warn("🚫 Permissão negada", zap.String("ip", ip), zap.String("claim", claimKey), zap.String("valor_esperado", claimValue))
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permissão negada"})
		}
		zap.L().Info("✅ Claims JWT válidos", zap.String("ip", ip), zap.String("claim", claimKey))
		return c.Next()
	}
}
