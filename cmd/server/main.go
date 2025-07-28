package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

	fmt.Println("✅ Banco de dados criado em memória com sucesso!")

	createTableSQL := loadSQLFile("sql/create_table.sql")
	insertProductsSQL := loadSQLFile("sql/insert_products.sql")
	createUserSQL := loadSQLFile("sql/create_user.sql")
	createFornecedoresSQL := loadSQLFile("sql/create_fornecedores.sql")
	insertFornecedoresSQL := loadSQLFile("sql/insert_fornecedores.sql")

	fmt.Println("Scripts carregados com sucesso!")

	executeRawSQL(db, createTableSQL, "Tabela produtos criada com sucesso!")
	executeRawSQL(db, insertProductsSQL, "Produtos inseridos com sucesso!")
	executeRawSQL(db, createUserSQL, "Tabela users criada com sucesso!")
	executeRawSQL(db, createFornecedoresSQL, "Tabela fornecedores criada com sucesso!")
	executeRawSQL(db, insertFornecedoresSQL, "Fornecedores inseridos com sucesso!")

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
	service := service.NewServiceInstance(cryptoService, persistence)
	return controller.NewControllerInstance(service)
}

func loadSQLFile(filepath string) string {
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo SQL %s: %v", filepath, err)
	}
	return string(content)
}

func executeRawSQL(db *gorm.DB, sqlContent string, successMessage string) {
	if err := db.Exec(sqlContent).Error; err != nil {
		// Se a tabela já existe, apenas avisa mas não para a execução
		if strings.Contains(err.Error(), "already exists") {
			fmt.Println("⚠️  Tabela já existe, continuando...")
			return
		}
		log.Fatalf("Erro ao executar SQL: %v", err)
	}
	fmt.Println("✅ " + successMessage)
}
