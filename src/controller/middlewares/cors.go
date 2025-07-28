package middlewares

import (
	"github.com/betine97/back-project.git/cmd/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

// CORSMiddleware configures and returns CORS middleware
func CORSMiddleware() fiber.Handler {
	cfg := config.NewConfig()

	zap.L().Info("Configuring CORS middleware",
		zap.String("allowed_origins", cfg.CORSOrigins))

	return cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-CSRF-Token",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type",
		MaxAge:           86400, // 24 hours in seconds
	})
}

// CORSMiddlewareStrict returns a more restrictive CORS configuration for production
func CORSMiddlewareStrict() fiber.Handler {
	cfg := config.NewConfig()

	zap.L().Info("Configuring strict CORS middleware",
		zap.String("allowed_origins", cfg.CORSOrigins))

	return cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           3600, // 1 hour in seconds
	})
}
