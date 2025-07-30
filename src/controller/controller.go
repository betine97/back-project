package controller

import (
	dtos "github.com/betine97/back-project.git/src/model/dtos"
	"github.com/betine97/back-project.git/src/model/service"
	"github.com/betine97/back-project.git/src/view"
	"github.com/gofiber/fiber/v2"
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
	CreateFornecedor(ctx *fiber.Ctx) error
	ChangeStatusFornecedor(ctx *fiber.Ctx) error
	UpdateFornecedorField(ctx *fiber.Ctx) error
	DeleteFornecedor(ctx *fiber.Ctx) error

	GetAllProducts(ctx *fiber.Ctx) error
	CreateProduct(ctx *fiber.Ctx) error
	DeleteProduct(ctx *fiber.Ctx) error

	GetAllPedidos(ctx *fiber.Ctx) error
}

type Controller struct {
	service service.ServiceInterface
}

// FUNÇÕES DE USUÁRIO ------------------------------------------------------------------------------------------------------------------------------------

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

	token, err := ctl.service.LoginUserService(user)
	if err != nil {
		zap.L().Error("Error when logging in", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}

func (ctl *Controller) RequestOtherService(ctx *fiber.Ctx) error {

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Service found successfully",
	})
}

// FUNÇÕES DE FORNECEDORES ------------------------------------------------------------------------------------------------------------------------------------

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

func (ctl *Controller) CreateFornecedor(ctx *fiber.Ctx) error {
	var fornecedor dtos.CreateFornecedorRequest

	if err := ctx.BodyParser(&fornecedor); err != nil {
		zap.L().Error("Error reading request data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to read request data",
		})
	}

	success, err := ctl.service.CreateFornecedorService(fornecedor)
	if err != nil {
		zap.L().Error("Error creating fornecedor", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating fornecedor",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Fornecedor created successfully",
	})
}

func (ctl *Controller) ChangeStatusFornecedor(ctx *fiber.Ctx) error {
	zap.L().Info("Starting change status fornecedor controller")

	id := ctx.Params("id")

	success, err := ctl.service.ChangeStatusFornecedorService(id)
	if err != nil {
		zap.L().Error("Error changing status fornecedor", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error changing status fornecedor",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Status fornecedor changed successfully",
	})
}

func (ctl *Controller) UpdateFornecedorField(ctx *fiber.Ctx) error {
	zap.L().Info("Starting update fornecedor field controller")

	id := ctx.Params("id")
	var request dtos.UpdateFornecedorRequest

	if err := ctx.BodyParser(&request); err != nil {
		zap.L().Error("Error reading request data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to read request data",
		})
	}

	success, err := ctl.service.UpdateFornecedorFieldService(id, request.Campo, request.Valor)
	if err != nil {
		zap.L().Error("Error updating fornecedor field", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating fornecedor field",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Fornecedor field updated successfully",
	})
}

func (ctl *Controller) DeleteFornecedor(ctx *fiber.Ctx) error {
	zap.L().Info("Starting delete fornecedor controller")

	id := ctx.Params("id")

	success, err := ctl.service.DeleteFornecedorService(id)
	if err != nil {
		zap.L().Error("Error deleting fornecedor", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error deleting fornecedor",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Fornecedor deleted successfully",
	})
}

// FUNÇÕES DE PRODUTOS ------------------------------------------------------------------------------------------------------------------------------------

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

func (ctl *Controller) CreateProduct(ctx *fiber.Ctx) error {

	var product dtos.CreateProductRequest

	if err := ctx.BodyParser(&product); err != nil {
		zap.L().Error("Error reading request data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to read request data",
		})
	}

	success, err := ctl.service.CreateProductService(product)
	if err != nil {
		zap.L().Error("Error creating product", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating product",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
	})
}

func (ctl *Controller) DeleteProduct(ctx *fiber.Ctx) error {
	zap.L().Info("Starting delete product controller")

	id := ctx.Params("id")

	success, err := ctl.service.DeleteProductService(id)
	if err != nil {
		zap.L().Error("Error deleting product", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error deleting product",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}

// FUNÇÕES DE PEDIDOS ------------------------------------------------------------------------------------------------------------------------------------

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
