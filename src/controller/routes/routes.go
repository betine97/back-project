package routes

import (
	"github.com/betine97/back-project.git/src/controller"
	"github.com/betine97/back-project.git/src/controller/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userController controller.ControllerInterface) {
	app.Post("/register", middlewares.UserValidationMiddleware, userController.CreateUser)
	app.Post("/login", userController.LoginUser)

	api := app.Group("/api", middlewares.JWTProtected())
	api.Get("/otherservice", middlewares.JWTClaimsRequired("role", "user"), userController.RequestOtherService)
}
