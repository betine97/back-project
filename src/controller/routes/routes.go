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

	// Public product routes (sem autenticação)
	app.Get("/produtos", userController.GetAllProducts)                // GET /produtos - Public access
	app.Get("/produtos/search", userController.GetProductsWithFilters) // GET /produtos/search - Public access
	app.Get("/produtos/:id", userController.GetProductByID)            // GET /produtos/:id - Public access

	// Protected routes
	api := app.Group("/api", middlewares.JWTProtected())
	api.Get("/otherservice", middlewares.JWTClaimsRequired("role", "user"), userController.RequestOtherService)

	// Protected product routes (com autenticação)
	products := api.Group("/produtos")
	products.Get("/", userController.GetAllProducts)               // GET /api/produtos - Protected
	products.Get("/search", userController.GetProductsWithFilters) // GET /api/produtos/search - Protected
	products.Get("/:id", userController.GetProductByID)            // GET /api/produtos/:id - Protected
}
