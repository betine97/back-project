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
	// Health checks
	HealthCheck(ctx *fiber.Ctx) error
	ReadinessCheck(ctx *fiber.Ctx) error

	// User management
	CreateUser(ctx *fiber.Ctx) error
	LoginUser(ctx *fiber.Ctx) error
	RequestOtherService(ctx *fiber.Ctx) error

	// Fornecedores
	GetAllFornecedores(ctx *fiber.Ctx) error
	CreateFornecedor(ctx *fiber.Ctx) error
	ChangeStatusFornecedor(ctx *fiber.Ctx) error
	UpdateFornecedorField(ctx *fiber.Ctx) error
	DeleteFornecedor(ctx *fiber.Ctx) error

	// Products
	GetAllProducts(ctx *fiber.Ctx) error
	CreateProduct(ctx *fiber.Ctx) error
	DeleteProduct(ctx *fiber.Ctx) error

	GetAllPedidos(ctx *fiber.Ctx) error
	GetPedidoById(ctx *fiber.Ctx) error
	CreatePedido(ctx *fiber.Ctx) error

	// Itens de Pedido
	GetItensPedido(ctx *fiber.Ctx) error
	CreateItemPedido(ctx *fiber.Ctx) error

	// Estoque
	GetAllEstoque(ctx *fiber.Ctx) error
	CreateEstoque(ctx *fiber.Ctx) error
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

	zap.L().Info("🔑 Iniciando processo de login")

	var user dtos.UserLogin

	if err := ctx.BodyParser(&user); err != nil {
		zap.L().Error("❌ Erro ao ler dados de login", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados de login inválidos",
		})
	}

	token, err := ctl.service.LoginUserService(user)
	if err != nil {
		zap.L().Error("❌ Erro durante login", zap.Error(err))
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
	zap.L().Info("📋 Buscando todos os fornecedores")

	// Obter parâmetros de paginação da query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	userID := ctx.Locals("userID").(string)
	fornecedores, err := ctl.service.GetAllFornecedoresService(userID, page, limit)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar fornecedores", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Fornecedores recuperados com sucesso", zap.Int("total", fornecedores.Total), zap.Int("page", page), zap.Int("limit", limit))
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

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.CreateFornecedorService(userID, fornecedor)
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

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.ChangeStatusFornecedorService(userID, id)
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

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.UpdateFornecedorFieldService(userID, id, request.Campo, request.Valor)
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

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.DeleteFornecedorService(userID, id)
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

	// Obter parâmetros de paginação da query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	userID := ctx.Locals("userID").(string)
	products, err := ctl.service.GetAllProductsService(userID, page, limit)
	if err != nil {
		zap.L().Error("Error getting all products", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved all products", zap.Int("count", len(products.Products)), zap.Int("total", products.Total), zap.Int("page", page), zap.Int("limit", limit))
	return ctx.Status(fiber.StatusOK).JSON(products)
}

func (ctl *Controller) CreateProduct(ctx *fiber.Ctx) error {

	createProduct := ctx.Locals("createProduct").(dtos.CreateProductRequest)

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.CreateProductService(userID, createProduct)
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

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.DeleteProductService(userID, id)
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
	zap.L().Info("📋 Buscando todos os pedidos")

	// Obter parâmetros de paginação da query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	userID := ctx.Locals("userID").(string)
	pedidos, err := ctl.service.GetAllPedidosService(userID, page, limit)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar pedidos", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Pedidos recuperados com sucesso", zap.Int("total", pedidos.Total), zap.Int("page", page), zap.Int("limit", limit))
	return ctx.Status(fiber.StatusOK).JSON(pedidos)
}

func (ctl *Controller) GetPedidoById(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get pedido by ID controller")

	id := ctx.Params("id")
	userID := ctx.Locals("userID").(string)

	pedido, err := ctl.service.GetPedidoByIdService(userID, id)
	if err != nil {
		zap.L().Error("Error getting pedido by ID", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved pedido by ID", zap.String("id", id))
	return ctx.Status(fiber.StatusOK).JSON(pedido)
}

func (ctl *Controller) CreatePedido(ctx *fiber.Ctx) error {
	zap.L().Info("Starting create pedido controller")

	createPedido := ctx.Locals("createPedido").(dtos.CreatePedidoRequest)

	userID := ctx.Locals("userID").(string)
	pedidoID, err := ctl.service.CreatePedidoService(userID, createPedido)
	if err != nil {
		zap.L().Error("Error creating pedido", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if pedidoID == 0 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating pedido",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Pedido created successfully",
		"id_pedido": pedidoID,
	})
}

// FUNÇÕES DE ITENS DE PEDIDO ------------------------------------------------------------------------------------------------------------------------------------

func (ctl *Controller) GetItensPedido(ctx *fiber.Ctx) error {
	zap.L().Info("📋 Buscando itens do pedido")

	idPedido := ctx.Params("id")

	// Obter parâmetros de paginação da query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	userID := ctx.Locals("userID").(string)
	itens, err := ctl.service.GetItensPedidoService(userID, idPedido, page, limit)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar itens do pedido", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Itens do pedido recuperados com sucesso", zap.String("idPedido", idPedido), zap.Int("total", itens.Total))
	return ctx.Status(fiber.StatusOK).JSON(itens)
}

func (ctl *Controller) CreateItemPedido(ctx *fiber.Ctx) error {
	zap.L().Info("Starting create item pedido controller")

	idPedido := ctx.Params("id")
	createItemPedido := ctx.Locals("createItemPedido").(dtos.CreateItemPedidoRequest)

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.CreateItemPedidoService(userID, idPedido, createItemPedido)
	if err != nil {
		zap.L().Error("Error creating item pedido", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating item pedido",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Item pedido created successfully",
	})
}

// FUNÇÕES DE ESTOQUE ------------------------------------------------------------------------------------------------------------------------------------

func (ctl *Controller) GetAllEstoque(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get all estoque controller")

	userID := ctx.Locals("userID").(string)
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	estoque, err := ctl.service.GetAllEstoqueService(userID, page, limit)
	if err != nil {
		zap.L().Error("Error getting estoque", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(estoque)
}

func (ctl *Controller) CreateEstoque(ctx *fiber.Ctx) error {
	zap.L().Info("Starting create estoque controller")

	createEstoque := ctx.Locals("createEstoque").(dtos.CreateEstoqueRequest)
	userID := ctx.Locals("userID").(string)

	success, err := ctl.service.CreateEstoqueService(userID, createEstoque)
	if err != nil {
		zap.L().Error("Error creating estoque", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating estoque",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Estoque created successfully",
	})
}
