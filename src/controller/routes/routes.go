package routes

import (
	"github.com/betine97/back-project.git/src/controller"
	"github.com/betine97/back-project.git/src/controller/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, userController controller.ControllerInterface, db *gorm.DB) {
	// Public routes
	app.Post("/register", middlewares.UserValidationMiddleware, userController.CreateUser)
	app.Post("/login", userController.LoginUser)

	// Protected routes
	api := app.Group("/api", middlewares.JWTProtected())
	api.Get("/otherservice", middlewares.JWTClaimsRequired("role", "user"), userController.RequestOtherService)

	// Protected pedidos routes (com autenticação)
	pedidos := api.Group("/pedidos")
	pedidos.Get("/", userController.GetAllPedidos)
	pedidos.Get("/:id", userController.GetPedidoByID)

	// Protected item pedidos routes (com autenticação)
	itemPedidos := api.Group("/item-pedidos")
	itemPedidos.Get("/", userController.GetAllItemPedidos)
	itemPedidos.Get("/:id", userController.GetItemPedidoByID)
	itemPedidos.Post("/", userController.CreateItemPedido)

	// Protected fornecedores routes (com autenticação)
	fornecedores := api.Group("/fornecedores")
	fornecedores.Get("/", userController.GetAllFornecedores)

	// Protected product routes (com autenticação)
	products := api.Group("/produtos")
	products.Get("/", userController.GetAllProducts)
	products.Get("/:id", userController.GetProductByID)

	// Protected his cmv prc marge routes (com autenticação)
	hisCmvPrcMarge := api.Group("/his-cmv-prc-marge")
	hisCmvPrcMarge.Get("/", userController.GetAllHisCmvPrcMarge)
}
