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

	// Product routes
	products := api.Group("/produtos")
	products.Get("/", userController.GetAllProducts)               // GET /api/produtos - Get all products
	products.Get("/search", userController.GetProductsWithFilters) // GET /api/produtos/search - Get products with filters
	products.Get("/:id", userController.GetProductByID)            // GET /api/produtos/:id - Get product by ID
}
