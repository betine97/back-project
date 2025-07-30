package routes

import (
	"github.com/betine97/back-project.git/src/controller"
	"github.com/betine97/back-project.git/src/controller/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, userController controller.ControllerInterface, db *gorm.DB) {
	// Public routes
	app.Post("/cadastro", middlewares.UserValidationMiddleware, userController.CreateUser)
	app.Post("/login", userController.LoginUser)

	// Protected routes
	api := app.Group("/api", middlewares.JWTProtected(), middlewares.DatabaseConnectionMiddleware(db, redisClient))

	fornecedores := api.Group("/fornecedores")
	fornecedores.Get("/", userController.GetAllFornecedores)
	fornecedores.Post("/", userController.CreateFornecedor)
	fornecedores.Put("changestatus/:id", userController.ChangeStatusFornecedor)
	fornecedores.Put("changefields/:id", userController.UpdateFornecedorField)
	fornecedores.Delete("/:id", userController.DeleteFornecedor)

	// Protected product routes (com autenticação)
	produtos := api.Group("/produtos")
	produtos.Get("/", userController.GetAllProducts)
	produtos.Post("/", userController.CreateProduct)
	produtos.Delete("/:id", userController.DeleteProduct)

	// Protected pedidos routes (com autenticação)
	pedidos := api.Group("/pedidos")
	pedidos.Get("/", userController.GetAllPedidos)

}
