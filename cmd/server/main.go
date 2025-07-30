package main

import (
	"github.com/betine97/back-project.git/cmd/config"
	"github.com/betine97/back-project.git/src/controller"
	"github.com/betine97/back-project.git/src/controller/middlewares"
	"github.com/betine97/back-project.git/src/controller/routes"
	"github.com/betine97/back-project.git/src/model/persistence"
	"github.com/betine97/back-project.git/src/model/service"
	"github.com/betine97/back-project.git/src/model/service/crypto"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {

	_ = config.NewConfig()

	db, err := config.NewDatabaseConnection()
	if err != nil {
		zap.L().Fatal("Failed to connect to database", zap.Error(err))
	}

	userController := initDependencies(db)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			zap.L().Error("Fiber error", zap.Error(err), zap.Int("status", code))
			return ctx.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Configure CORS middleware
	app.Use(middlewares.CORSMiddleware())

	// Configure logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))

	routes.SetupRoutes(app, userController, db)

	if err := app.Listen(":8080"); err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}

}

func initDependencies(database *gorm.DB) controller.ControllerInterface {
	cryptoService := &crypto.Crypto{}
	persistence := persistence.NewDBConnection(database)
	service := service.NewServiceInstance(cryptoService, persistence, config.RedisClient)
	return controller.NewControllerInstance(service)
}
