package controller

import (
	"strconv"
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

	GetAllFornecedores(ctx *fiber.Ctx) error

	GetAllProducts(ctx *fiber.Ctx) error
	GetProductByID(ctx *fiber.Ctx) error

	GetAllPedidos(ctx *fiber.Ctx) error
	GetPedidoByID(ctx *fiber.Ctx) error

	GetAllItemPedidos(ctx *fiber.Ctx) error
	GetItemPedidoByID(ctx *fiber.Ctx) error
	CreateItemPedido(ctx *fiber.Ctx) error

	GetAllHisCmvPrcMarge(ctx *fiber.Ctx) error
}

type Controller struct {
	service service.ServiceInterface
}

// Funções de usuário

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

// Funções de fornecedores

//---------------------------------

func (ctl *Controller) GetAllFornecedores(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get all fornecedores controller")

	fornecedores, err := ctl.service.GetAllFornecedoresService()
	if err != nil {
		zap.L().Error("Error getting all fornecedores", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved all fornecedores", zap.Int("count", len(fornecedores.Fornecedores)))
	return ctx.Status(fiber.StatusOK).JSON(fornecedores)
}

// Funções de produtos

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

// Funções de pedidos

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

func (ctl *Controller) GetPedidoByID(ctx *fiber.Ctx) error {

	idParam := ctx.Params("id")
	id, parseErr := strconv.Atoi(idParam)
	if parseErr != nil {
		zap.L().Error("Invalid pedido ID", zap.String("id", idParam), zap.Error(parseErr))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid pedido ID",
		})
	}

	pedido, err := ctl.service.GetPedidoByIDService(id)
	if err != nil {
		zap.L().Error("Error getting pedido by ID", zap.Int("id", id), zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved pedido by ID", zap.Int("id", id))
	return ctx.Status(fiber.StatusOK).JSON(pedido)
}

// Funções de itens de pedidos

//---------------------------------

func (ctl *Controller) GetAllItemPedidos(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get all item pedidos controller")

	itemPedidos, err := ctl.service.GetAllItemPedidosService()
	if err != nil {
		zap.L().Error("Error getting all item pedidos", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved all item pedidos", zap.Int("count", len(itemPedidos.ItemPedidos)))
	return ctx.Status(fiber.StatusOK).JSON(itemPedidos)
}

func (ctl *Controller) GetItemPedidoByID(ctx *fiber.Ctx) error {

	idParam := ctx.Params("id")
	id, parseErr := strconv.Atoi(idParam)
	if parseErr != nil {
		zap.L().Error("Invalid item pedido ID", zap.String("id", idParam), zap.Error(parseErr))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid item pedido ID",
		})
	}

	itemPedido, err := ctl.service.GetItemPedidoByIDService(id)
	if err != nil {
		zap.L().Error("Error getting item pedido by ID", zap.Int("id", id), zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved item pedido by ID", zap.Int("id", id))
	return ctx.Status(fiber.StatusOK).JSON(itemPedido)
}

// ... existing code no final do arquivo ...

func (ctl *Controller) CreateItemPedido(ctx *fiber.Ctx) error {
	zap.L().Info("Starting create item pedido controller")

	var request modelDtos.CreateItemPedidoRequest

	// Parse JSON body
	if err := ctx.BodyParser(&request); err != nil {
		zap.L().Error("Error parsing request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if request.IDPedido <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id_pedido is required and must be greater than 0",
		})
	}

	if request.IDProduto <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id_produto is required and must be greater than 0",
		})
	}

	if request.Quantidade <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "quantidade is required and must be greater than 0",
		})
	}

	if request.PrecoUnitario <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "preco_unitario is required and must be greater than 0",
		})
	}

	// Call service
	itemPedido, err := ctl.service.CreateItemPedidoService(request)
	if err != nil {
		zap.L().Error("Error creating item pedido", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully created item pedido", zap.Int("id_item", itemPedido.IDItem))
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Item pedido created successfully",
		"item_pedido": itemPedido,
	})
}

// Funções de histórico de cmv, preço e margem

//---------------------------------

func (ctl *Controller) GetAllHisCmvPrcMarge(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get all his cmv prc marge controller")

	hisCmvPrcMarge, err := ctl.service.GetAllHisCmvPrcMargeService()
	if err != nil {
		zap.L().Error("Error getting all his cmv prc marge", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved all his cmv prc marge", zap.Int("count", len(hisCmvPrcMarge.HisCmvPrcMarge)))
	return ctx.Status(fiber.StatusOK).JSON(hisCmvPrcMarge)
}
