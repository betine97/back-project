package controller

import (
	"time"

	"github.com/betine97/back-project.git/cmd/config"
	"github.com/betine97/back-project.git/cmd/config/exceptions"
	dtos_controllers "github.com/betine97/back-project.git/src/controller/dtos_controllers"
	"github.com/betine97/back-project.git/src/model/service"
	"github.com/betine97/back-project.git/src/view"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func NewControllerInstance(serviceInterface service.ServiceInterface) ControllerInterface {
	return &Controller{
		service: serviceInterface,
	}
}

type ControllerInterface interface {
	CreateUser(ctx *fiber.Ctx) error
	LoginUser(ctx *fiber.Ctx) error
	RequestOtherService(ctx *fiber.Ctx) error

	GetAllProducts(ctx *fiber.Ctx) error

	GetAllPedidos(ctx *fiber.Ctx) error
}

type Controller struct {
	service service.ServiceInterface
}

// FUNÇÕES DE USUÁRIO

//---------------------------------

func (ctl *Controller) CreateUser(ctx *fiber.Ctx) error {

	createUser := ctx.Locals("createUser").(dtos_controllers.CreateUser)

	resp, err := ctl.service.CreateUserService(createUser)
	if err != nil {
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": "Error creating user: " + err.Error(),
		})
	}

	ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"usuário": view.ConvertDomainToResponse(resp),
	})

	return nil
}

func (ctl *Controller) LoginUser(ctx *fiber.Ctx) error {

	zap.L().Info("Starting user login")

	var user dtos_controllers.UserLogin

	if err := ctx.BodyParser(&user); err != nil {
		zap.L().Error("Error reading request data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to read request data",
		})
	}

	_, err := ctl.service.LoginUserService(user)
	if err != nil {
		zap.L().Error("Error when logging in", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := GenerateToken(user)
	if err != nil {
		zap.L().Error("Error generating token", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Login successfully")
	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successfully",
		"token":   token,
	})

	return nil
}

func (ctl *Controller) RequestOtherService(ctx *fiber.Ctx) error {

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Service found successfully",
	})
}

func GenerateToken(user dtos_controllers.UserLogin) (string, *exceptions.RestErr) {

	claims := jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"role":  "user",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(config.PrivateKey)
	if err != nil {
		zap.L().Error("Error signing token", zap.Error(err))
		return "", exceptions.NewInternalServerError("Internal server error")
	}

	return tokenString, nil
}

// FUNÇÕES DE PRODUTOS

//---------------------------------

func (ctl *Controller) GetAllProducts(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get all products controller")

	products, err := ctl.service.GetAllProductsService()
	if err != nil {
		zap.L().Error("Error getting all products", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved all products", zap.Int("count", len(products.Products)))
	return ctx.Status(fiber.StatusOK).JSON(products)
}

// FUNÇÕES DE PEDIDOS

//---------------------------------

func (ctl *Controller) GetAllPedidos(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get all pedidos controller")

	pedidos, err := ctl.service.GetAllPedidosService()
	if err != nil {
		zap.L().Error("Error getting all pedidos", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved all pedidos", zap.Int("count", len(pedidos.Pedidos)))
	return ctx.Status(fiber.StatusOK).JSON(pedidos)
}
