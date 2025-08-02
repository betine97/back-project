package main

import (
	"github.com/betine97/back-project.git/cmd/config"
	"github.com/betine97/back-project.git/src/controller"
	"github.com/betine97/back-project.git/src/controller/middlewares"
	"github.com/betine97/back-project.git/src/controller/routes"
	"github.com/betine97/back-project.git/src/model/interfaces"
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

	zap.L().Info("🚀 Iniciando aplicação Back Project")

	zap.L().Info("📊 Conectando ao banco de dados master...")
	dbmaster, err := config.NewDatabaseConnection()
	if err != nil {
		zap.L().Fatal("❌ Falha ao conectar com banco master", zap.Error(err))
	}
	zap.L().Info("✅ Banco master conectado com sucesso")

	zap.L().Info("🏢 Conectando aos bancos de clientes...")
	clientDB, err := config.ConnectionDBClients()
	if err != nil {
		zap.L().Fatal("❌ Falha ao conectar com bancos de clientes", zap.Error(err))
	}
	zap.L().Info("✅ Bancos de clientes conectados com sucesso", zap.Int("total", len(clientDB)))

	zap.L().Info("🔧 Inicializando dependências...")
	userController := initDependencies(dbmaster, clientDB)
	zap.L().Info("✅ Dependências inicializadas com sucesso")

	zap.L().Info("🌐 Configurando servidor Fiber...")
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			zap.L().Error("💥 Erro no servidor Fiber",
				zap.Error(err),
				zap.Int("status", code),
				zap.String("path", ctx.Path()),
				zap.String("method", ctx.Method()))
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

	zap.L().Info("🛣️ Configurando rotas da aplicação...")
	routes.SetupRoutes(app, userController)
	zap.L().Info("✅ Rotas configuradas com sucesso")

	cfg := config.NewConfig()
	port := ":" + cfg.WebServerPort
	zap.L().Info("🚀 Iniciando servidor HTTP", zap.String("porta", port))
	if err := app.Listen(port); err != nil {
		zap.L().Fatal("❌ Falha ao iniciar servidor", zap.Error(err))
	}

}

func initDependencies(masterDB *gorm.DB, clientDB map[string]*gorm.DB) controller.ControllerInterface {
	cryptoService := &crypto.Crypto{}
	persistenceDBMASTER := persistence.NewDBConnectionDBMaster(masterDB)
	persistenceDBCLIENT := persistence.NewDBConnectionDBClient(clientDB)

	// Wrap Redis client with interface
	redisWrapper := interfaces.NewRedisWrapper(config.RedisClient)

	// Create token generator
	tokenGenerator := interfaces.NewJWTTokenGenerator()

	service := service.NewServiceInstance(cryptoService, persistenceDBMASTER, persistenceDBCLIENT, redisWrapper, tokenGenerator)
	return controller.NewControllerInstance(service)
}
