package controller

import (
	"strconv"
	"time"

	"github.com/betine97/back-project.git/cmd/config"
	"github.com/betine97/back-project.git/cmd/config/exceptions"
	"github.com/betine97/back-project.git/src/controller/dtos"
	modelDtos "github.com/betine97/back-project.git/src/model/dtos"
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
	GetProductByID(ctx *fiber.Ctx) error
	GetProductsWithFilters(ctx *fiber.Ctx) error
}

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

func (ctl *Controller) GetProductByID(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get product by ID controller")

	idParam := ctx.Params("id")
	id, parseErr := strconv.Atoi(idParam)
	if parseErr != nil {
		zap.L().Error("Invalid product ID", zap.String("id", idParam), zap.Error(parseErr))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	product, err := ctl.service.GetProductByIDService(id)
	if err != nil {
		zap.L().Error("Error getting product by ID", zap.Int("id", id), zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved product by ID", zap.Int("id", id))
	return ctx.Status(fiber.StatusOK).JSON(product)
}

func (ctl *Controller) GetProductsWithFilters(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get products with filters controller")

	// Parse query parameters
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	// Parse filters
	filters := &modelDtos.ProductFilters{
		Categoria:     ctx.Query("categoria"),
		DestinadoPara: ctx.Query("destinado_para"),
		Marca:         ctx.Query("marca"),
		Variacao:      ctx.Query("variacao"),
		Status:        ctx.Query("status", "ativo"),
		Search:        ctx.Query("search"),
	}

	// Parse price filters
	if minPriceStr := ctx.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = minPrice
		}
	}
	if maxPriceStr := ctx.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = maxPrice
		}
	}

	params := modelDtos.ProductQueryParams{
		Page:   page,
		Limit:  limit,
		Filter: filters,
	}

	products, err := ctl.service.GetProductsWithFiltersService(params)
	if err != nil {
		zap.L().Error("Error getting filtered products", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved filtered products",
		zap.Int("count", len(products.Products)),
		zap.Int("total", products.Total))

	return ctx.Status(fiber.StatusOK).JSON(products)
}

type Controller struct {
	service service.ServiceInterface
}

func (ctl *Controller) CreateUser(ctx *fiber.Ctx) error {

	createUser := ctx.Locals("createUser").(dtos.CreateUser)

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

	var user dtos.UserLogin

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

func GenerateToken(user dtos.UserLogin) (string, *exceptions.RestErr) {

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
