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

	// Clientes
	GetAllClientes(ctx *fiber.Ctx) error
	BuscarClientesCriterios(ctx *fiber.Ctx) error
	AdicionarClientesAoPublico(ctx *fiber.Ctx) error
	GetClienteByID(ctx *fiber.Ctx) error
	CreateCliente(ctx *fiber.Ctx) error
	DeleteCliente(ctx *fiber.Ctx) error

	// Tags de Clientes
	AtribuirTagsCliente(ctx *fiber.Ctx) error
	RemoverTagsCliente(ctx *fiber.Ctx) error
	GetTagsCliente(ctx *fiber.Ctx) error

	// Tags
	GetAllTags(ctx *fiber.Ctx) error
	CreateTag(ctx *fiber.Ctx) error

	// Campanhas
	GetAllCampanhas(ctx *fiber.Ctx) error
	GetCampanhaByID(ctx *fiber.Ctx) error
	CreateCampanha(ctx *fiber.Ctx) error
	AssociarPublicosCampanha(ctx *fiber.Ctx) error
	GetPublicosCampanha(ctx *fiber.Ctx) error

	// Endereços
	GetAllEnderecos(ctx *fiber.Ctx) error
	CreateEndereco(ctx *fiber.Ctx) error
	DeleteEndereco(ctx *fiber.Ctx) error

	// Critérios
	GetAllCriterios(ctx *fiber.Ctx) error

	// Públicos
	GetAllPublicos(ctx *fiber.Ctx) error
	CreatePublico(ctx *fiber.Ctx) error
	AssociarCriteriosPublico(ctx *fiber.Ctx) error
	GetCriteriosPublico(ctx *fiber.Ctx) error
	GetClientesDoPublico(ctx *fiber.Ctx) error
	GetClientesDoPublicoTest(ctx *fiber.Ctx) error

	// Pets
	GetAllPets(ctx *fiber.Ctx) error
	CreatePet(ctx *fiber.Ctx) error

	// Completude
	GetCompletudeClientes(ctx *fiber.Ctx) error
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

// FUNÇÕES DE CLIENTES ------------------------------------------------------------------------------------------------------------------------------------

func (ctl *Controller) GetAllClientes(ctx *fiber.Ctx) error {
	zap.L().Info("📋 Buscando todos os clientes")

	// Obter parâmetros de paginação da query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	userID := ctx.Locals("userID").(string)
	clientes, err := ctl.service.GetAllClientesService(userID, page, limit)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar clientes", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Clientes recuperados com sucesso", zap.Int("total", clientes.Total), zap.Int("page", page), zap.Int("limit", limit))
	return ctx.Status(fiber.StatusOK).JSON(clientes)
}

func (ctl *Controller) BuscarClientesCriterios(ctx *fiber.Ctx) error {
	zap.L().Info("🔍 Buscando clientes para validação de critérios")

	idPublico := ctx.Params("id_publico")
	userID := ctx.Locals("userID").(string)

	clientes, err := ctl.service.BuscarClientesCriteriosService(userID, idPublico)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar clientes para critérios", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Clientes para critérios recuperados com sucesso", zap.Int("total", clientes.Total))
	return ctx.Status(fiber.StatusOK).JSON(clientes)
}

func (ctl *Controller) AdicionarClientesAoPublico(ctx *fiber.Ctx) error {
	zap.L().Info("🔍 Adicionando clientes ao público baseado em critérios")

	idPublico := ctx.Params("id_publico")
	userID := ctx.Locals("userID").(string)

	result, err := ctl.service.AdicionarClientesAoPublicoService(userID, idPublico)
	if err != nil {
		zap.L().Error("❌ Erro ao adicionar clientes ao público", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Clientes adicionados ao público com sucesso", zap.Int("total_adicionados", result.ClientesAdicionados))
	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (ctl *Controller) GetClienteByID(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get cliente by ID controller")

	id := ctx.Params("id")
	userID := ctx.Locals("userID").(string)

	cliente, err := ctl.service.GetClienteByIDService(userID, id)
	if err != nil {
		zap.L().Error("Error getting cliente by ID", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved cliente by ID", zap.String("id", id))
	return ctx.Status(fiber.StatusOK).JSON(cliente)
}

func (ctl *Controller) CreateCliente(ctx *fiber.Ctx) error {
	zap.L().Info("Starting create cliente controller")

	createCliente := ctx.Locals("createCliente").(dtos.CreateClienteRequest)
	userID := ctx.Locals("userID").(string)

	clienteID, err := ctl.service.CreateClienteService(userID, createCliente)
	if err != nil {
		zap.L().Error("Error creating cliente", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if clienteID == 0 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating cliente",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Cliente created successfully",
		"id_cliente": clienteID,
	})
}

// FUNÇÕES DE ENDEREÇOS ------------------------------------------------------------------------------------------------------------------------------------

func (ctl *Controller) GetAllEnderecos(ctx *fiber.Ctx) error {
	zap.L().Info("📋 Buscando todos os endereços")

	// Obter parâmetros de paginação da query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	userID := ctx.Locals("userID").(string)
	enderecos, err := ctl.service.GetAllEnderecosService(userID, page, limit)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar endereços", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Endereços recuperados com sucesso", zap.Int("total", enderecos.Total), zap.Int("page", page), zap.Int("limit", limit))
	return ctx.Status(fiber.StatusOK).JSON(enderecos)
}

func (ctl *Controller) CreateEndereco(ctx *fiber.Ctx) error {
	zap.L().Info("Starting create endereco controller")

	createEndereco := ctx.Locals("createEndereco").(dtos.CreateEnderecoRequest)
	userID := ctx.Locals("userID").(string)

	success, err := ctl.service.CreateEnderecoService(userID, createEndereco)
	if err != nil {
		zap.L().Error("Error creating endereco", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating endereco",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Endereco created successfully",
	})
}
func (ctl *Controller) DeleteCliente(ctx *fiber.Ctx) error {
	zap.L().Info("Starting delete cliente controller")

	id := ctx.Params("id")

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.DeleteClienteService(userID, id)
	if err != nil {
		zap.L().Error("Error deleting cliente", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error deleting cliente",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cliente deleted successfully",
	})
}

func (ctl *Controller) DeleteEndereco(ctx *fiber.Ctx) error {
	zap.L().Info("Starting delete endereco controller")

	idEndereco := ctx.Params("id_endereco")

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.DeleteEnderecoService(userID, idEndereco)
	if err != nil {
		zap.L().Error("Error deleting endereco", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error deleting endereco",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Endereco deleted successfully",
	})
}

// FUNÇÕES DE CRITÉRIOS ------------------------------------------------------------------------------------------------------------------------------------

func (ctl *Controller) GetAllCriterios(ctx *fiber.Ctx) error {
	zap.L().Info("📋 Buscando todos os critérios")

	userID := ctx.Locals("userID").(string)
	criterios, err := ctl.service.GetAllCriteriosService(userID)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar critérios", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Critérios recuperados com sucesso", zap.Int("total", criterios.Total))
	return ctx.Status(fiber.StatusOK).JSON(criterios)
}

// FUNÇÕES DE PÚBLICOS ------------------------------------------------------------------------------------------------------------------------------------

func (ctl *Controller) GetAllPublicos(ctx *fiber.Ctx) error {
	zap.L().Info("📋 Buscando todos os públicos")

	// Obter parâmetros de paginação da query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	userID := ctx.Locals("userID").(string)
	publicos, err := ctl.service.GetAllPublicosService(userID, page, limit)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar públicos", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Públicos recuperados com sucesso", zap.Int("total", publicos.Total), zap.Int("page", page), zap.Int("limit", limit))
	return ctx.Status(fiber.StatusOK).JSON(publicos)
}

func (ctl *Controller) CreatePublico(ctx *fiber.Ctx) error {
	zap.L().Info("Starting create publico controller")

	createPublico := ctx.Locals("createPublico").(dtos.CreatePublicoRequest)

	userID := ctx.Locals("userID").(string)
	publicoID, err := ctl.service.CreatePublicoService(userID, createPublico)
	if err != nil {
		zap.L().Error("Error creating publico", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if publicoID == 0 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating publico",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Publico created successfully",
		"id_publico": publicoID,
	})
}

func (ctl *Controller) AssociarCriteriosPublico(ctx *fiber.Ctx) error {
	zap.L().Info("Starting associar criterios publico controller")

	idPublico := ctx.Params("id")
	var request dtos.AssociarCriteriosRequest

	if err := ctx.BodyParser(&request); err != nil {
		zap.L().Error("Error reading request data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to read request data",
		})
	}

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.AssociarCriteriosPublicoService(userID, idPublico, request)
	if err != nil {
		zap.L().Error("Error associating criterios to publico", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error associating criterios to publico",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Criterios associated successfully",
	})
}

func (ctl *Controller) GetCriteriosPublico(ctx *fiber.Ctx) error {
	zap.L().Info("📋 Buscando critérios do público")

	idPublico := ctx.Params("id")

	userID := ctx.Locals("userID").(string)
	criterios, err := ctl.service.GetCriteriosPublicoService(userID, idPublico)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar critérios do público", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Critérios do público recuperados com sucesso", zap.String("idPublico", idPublico), zap.Int("total", criterios.Total))
	return ctx.Status(fiber.StatusOK).JSON(criterios)
}

func (ctl *Controller) GetClientesDoPublico(ctx *fiber.Ctx) error {
	zap.L().Info("📋 Buscando clientes do público")

	idPublico := ctx.Params("id")
	userID := ctx.Locals("userID").(string)

	// Obter parâmetros de paginação da query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	clientes, err := ctl.service.GetClientesDoPublicoService(userID, idPublico, page, limit)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar clientes do público", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Clientes do público recuperados com sucesso", zap.String("idPublico", idPublico), zap.Int("total", clientes.Total))
	return ctx.Status(fiber.StatusOK).JSON(clientes)
}

func (ctl *Controller) GetClientesDoPublicoTest(ctx *fiber.Ctx) error {
	zap.L().Info("🧪 TESTE - Buscando clientes do público (sem autenticação)")

	idPublico := ctx.Params("id")
	userID := "2" // Sabemos que funciona com userID = "2"

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	clientes, err := ctl.service.GetClientesDoPublicoService(userID, idPublico, page, limit)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar clientes do público (teste)", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ TESTE - Clientes do público recuperados com sucesso", zap.String("idPublico", idPublico), zap.Int("total", clientes.Total))
	return ctx.Status(fiber.StatusOK).JSON(clientes)
}

// FUNÇÕES DE PETS ------------------------------------------------------------------------------------------------------------------------------------

func (ctl *Controller) GetAllPets(ctx *fiber.Ctx) error {
	zap.L().Info("📋 Buscando todos os pets")

	// Obter parâmetros de paginação da query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	userID := ctx.Locals("userID").(string)
	pets, err := ctl.service.GetAllPetsService(userID, page, limit)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar pets", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Pets recuperados com sucesso", zap.Int("total", pets.Total), zap.Int("page", page), zap.Int("limit", limit))
	return ctx.Status(fiber.StatusOK).JSON(pets)
}

func (ctl *Controller) CreatePet(ctx *fiber.Ctx) error {
	zap.L().Info("Starting create pet controller")

	createPet := ctx.Locals("createPet").(dtos.CreatePetRequest)

	userID := ctx.Locals("userID").(string)
	petID, err := ctl.service.CreatePetService(userID, createPet)
	if err != nil {
		zap.L().Error("Error creating pet", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if petID == 0 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating pet",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pet created successfully",
		"id_pet":  petID,
	})
}

// GetCompletudeClientes retorna a completude do cadastro de clientes e pets
func (ctl *Controller) GetCompletudeClientes(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get completude clientes controller")

	userID := ctx.Locals("userID").(string)

	// Parâmetros de paginação
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	response, err := ctl.service.GetCompletudeClientesService(userID, page, limit)
	if err != nil {
		zap.L().Error("Error getting completude clientes", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// FUNÇÕES DE TAGS DE CLIENTES ------------------------------------------------------------------------------------------------------------------------------------

func (ctl *Controller) AtribuirTagsCliente(ctx *fiber.Ctx) error {
	zap.L().Info("Starting atribuir tags cliente controller")

	clienteID := ctx.Params("id")
	request := ctx.Locals("atribuirTagsCliente").(dtos.AtribuirTagsClienteRequest)

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.AtribuirTagsClienteService(userID, clienteID, request)
	if err != nil {
		zap.L().Error("Error atribuindo tags ao cliente", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error atribuindo tags ao cliente",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tags atribuídas ao cliente com sucesso",
	})
}

func (ctl *Controller) RemoverTagsCliente(ctx *fiber.Ctx) error {
	zap.L().Info("Starting remover tags cliente controller")

	clienteID := ctx.Params("id")
	request := ctx.Locals("removerTagsCliente").(dtos.RemoverTagsClienteRequest)

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.RemoverTagsClienteService(userID, clienteID, request)
	if err != nil {
		zap.L().Error("Error removendo tags do cliente", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error removendo tags do cliente",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tags removidas do cliente com sucesso",
	})
}

func (ctl *Controller) GetTagsCliente(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get tags cliente controller")

	clienteID := ctx.Params("id")
	userID := ctx.Locals("userID").(string)

	tags, err := ctl.service.GetTagsClienteService(userID, clienteID)
	if err != nil {
		zap.L().Error("Error getting tags do cliente", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved tags do cliente", zap.String("clienteID", clienteID))
	return ctx.Status(fiber.StatusOK).JSON(tags)
}

// FUNÇÕES DE TAGS ------------------------------------------------------------------------------------------------------------------------------------

func (ctl *Controller) GetAllTags(ctx *fiber.Ctx) error {
	zap.L().Info("📋 Buscando todas as tags")

	// Obter parâmetros de paginação da query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	userID := ctx.Locals("userID").(string)
	tags, err := ctl.service.GetAllTagsService(userID, page, limit)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar tags", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Tags recuperadas com sucesso", zap.Int("total", tags.Total), zap.Int("page", page), zap.Int("limit", limit))
	return ctx.Status(fiber.StatusOK).JSON(tags)
}

func (ctl *Controller) CreateTag(ctx *fiber.Ctx) error {
	zap.L().Info("Starting create tag controller")

	createTag := ctx.Locals("createTag").(dtos.CreateTagRequest)
	userID := ctx.Locals("userID").(string)

	success, err := ctl.service.CreateTagService(userID, createTag)
	if err != nil {
		zap.L().Error("Error creating tag", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating tag",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tag created successfully",
	})
}

// FUNÇÕES DE CAMPANHAS ------------------------------------------------------------------------------------------------------------------------------------

func (ctl *Controller) GetAllCampanhas(ctx *fiber.Ctx) error {
	zap.L().Info("📋 Buscando todas as campanhas")

	// Obter parâmetros de paginação da query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 30)

	userID := ctx.Locals("userID").(string)
	campanhas, err := ctl.service.GetAllCampanhasService(userID, page, limit)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar campanhas", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Campanhas recuperadas com sucesso", zap.Int("total", campanhas.Total), zap.Int("page", page), zap.Int("limit", limit))
	return ctx.Status(fiber.StatusOK).JSON(campanhas)
}

func (ctl *Controller) GetCampanhaByID(ctx *fiber.Ctx) error {
	zap.L().Info("Starting get campanha by ID controller")

	id := ctx.Params("id")
	userID := ctx.Locals("userID").(string)

	campanha, err := ctl.service.GetCampanhaByIDService(userID, id)
	if err != nil {
		zap.L().Error("Error getting campanha by ID", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("Successfully retrieved campanha by ID", zap.String("id", id))
	return ctx.Status(fiber.StatusOK).JSON(campanha)
}

func (ctl *Controller) CreateCampanha(ctx *fiber.Ctx) error {
	zap.L().Info("Starting create campanha controller")

	createCampanha := ctx.Locals("createCampanha").(dtos.CreateCampanhaRequest)
	userID := ctx.Locals("userID").(string)

	campanhaID, err := ctl.service.CreateCampanhaService(userID, createCampanha)
	if err != nil {
		zap.L().Error("Error creating campanha", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if campanhaID == 0 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating campanha",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Campanha created successfully",
		"id_campanha": campanhaID,
	})
}

func (ctl *Controller) AssociarPublicosCampanha(ctx *fiber.Ctx) error {
	zap.L().Info("Starting associar publicos campanha controller")

	idCampanha := ctx.Params("id")
	var request dtos.AssociarPublicosCampanhaRequest

	if err := ctx.BodyParser(&request); err != nil {
		zap.L().Error("Error reading request data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to read request data",
		})
	}

	userID := ctx.Locals("userID").(string)
	success, err := ctl.service.AssociarPublicosCampanhaService(userID, idCampanha, request)
	if err != nil {
		zap.L().Error("Error associating publicos to campanha", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error associating publicos to campanha",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Publicos associated successfully",
	})
}

func (ctl *Controller) GetPublicosCampanha(ctx *fiber.Ctx) error {
	zap.L().Info("📋 Buscando públicos da campanha")

	idCampanha := ctx.Params("id")
	userID := ctx.Locals("userID").(string)

	publicos, err := ctl.service.GetPublicosCampanhaService(userID, idCampanha)
	if err != nil {
		zap.L().Error("❌ Erro ao buscar públicos da campanha", zap.Error(err))
		return ctx.Status(err.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zap.L().Info("✅ Públicos da campanha recuperados com sucesso", zap.String("idCampanha", idCampanha), zap.Int("total", publicos.Total))
	return ctx.Status(fiber.StatusOK).JSON(publicos)
}
