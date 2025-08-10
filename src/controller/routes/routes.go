package routes

import (
	"github.com/betine97/back-project.git/src/controller"
	"github.com/betine97/back-project.git/src/controller/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userController controller.ControllerInterface) {
	// Health check routes (sem autenticação)
	app.Get("/health", userController.HealthCheck)
	app.Get("/ready", userController.ReadinessCheck)

	// Public routes
	app.Post("/cadastro", middlewares.UserValidationMiddleware, userController.CreateUser)
	app.Post("/login", userController.LoginUser)

	// Rota de teste temporária (SEM autenticação) - REMOVER EM PRODUÇÃO
	app.Get("/test/publicos/:id/clientes", userController.GetClientesDoPublicoTest)

	// Protected routes
	api := app.Group("/api", middlewares.JWTProtected(), middlewares.DatabaseExtractIdUser())

	fornecedores := api.Group("/fornecedores")
	fornecedores.Get("/", userController.GetAllFornecedores)
	fornecedores.Post("/", userController.CreateFornecedor)
	fornecedores.Put("changestatus/:id", userController.ChangeStatusFornecedor)
	fornecedores.Put("changefields/:id", userController.UpdateFornecedorField)
	fornecedores.Delete("/:id", userController.DeleteFornecedor)

	// Protected product routes (com autenticação)
	produtos := api.Group("/produtos")
	produtos.Get("/", userController.GetAllProducts)
	produtos.Post("/", middlewares.ProductValidationMiddleware, userController.CreateProduct)
	produtos.Delete("/:id", userController.DeleteProduct)

	// Protected pedidos routes (com autenticação)
	pedidos := api.Group("/pedidos")
	pedidos.Get("/", userController.GetAllPedidos)
	pedidos.Get("/:id", userController.GetPedidoById)
	pedidos.Post("/", middlewares.PedidoValidationMiddleware, userController.CreatePedido)

	// Protected itens pedido routes (com autenticação)
	pedidos.Get("/:id/itens", userController.GetItensPedido)
	pedidos.Post("/:id/itens", middlewares.ItemPedidoValidationMiddleware, userController.CreateItemPedido)

	// Protected estoque routes (com autenticação)
	estoque := api.Group("/estoque")
	estoque.Get("/", userController.GetAllEstoque)
	estoque.Post("/", middlewares.EstoqueValidationMiddleware, userController.CreateEstoque)

	// Protected clientes routes (com autenticação)
	clientes := api.Group("/clientes")
	clientes.Get("/", userController.GetAllClientes)
	clientes.Get("/buscar-criterios/:id_publico", userController.BuscarClientesCriterios)
	clientes.Post("/adicionar-ao-publico/:id_publico", userController.AdicionarClientesAoPublico)
	clientes.Get("/id/:id", userController.GetClienteByID)
	clientes.Post("/", middlewares.ClienteValidationMiddleware, userController.CreateCliente)
	clientes.Delete("/:id", userController.DeleteCliente)

	// Tags de clientes
	clientes.Post("/:id/tags", middlewares.AtribuirTagsClienteValidationMiddleware, userController.AtribuirTagsCliente)
	clientes.Delete("/:id/tags", middlewares.RemoverTagsClienteValidationMiddleware, userController.RemoverTagsCliente)
	clientes.Get("/:id/tags", userController.GetTagsCliente)

	// Protected enderecos routes (com autenticação)
	enderecos := api.Group("/enderecos")
	enderecos.Get("/", userController.GetAllEnderecos)
	enderecos.Post("/", middlewares.EnderecoValidationMiddleware, userController.CreateEndereco)
	enderecos.Delete("/:id_endereco", userController.DeleteEndereco)

	// Protected criterios routes (com autenticação)
	criterios := api.Group("/criterios")
	criterios.Get("/", userController.GetAllCriterios)

	// Protected publicos routes (com autenticação)
	publicos := api.Group("/publicos")
	publicos.Get("/", userController.GetAllPublicos)
	publicos.Post("/", middlewares.PublicoValidationMiddleware, userController.CreatePublico)
	publicos.Post("/:id/criterios", userController.AssociarCriteriosPublico)
	publicos.Get("/:id/criterios", userController.GetCriteriosPublico)
	publicos.Get("/:id/clientes", userController.GetClientesDoPublico)

	// Protected pets routes (com autenticação)
	pets := api.Group("/pets")
	pets.Get("/", userController.GetAllPets)
	pets.Post("/", middlewares.PetValidationMiddleware, userController.CreatePet)

	// Protected completude routes (com autenticação)
	completude := api.Group("/completude")
	completude.Get("/clientes", userController.GetCompletudeClientes)

	// Protected tags routes (com autenticação)
	tags := api.Group("/tags")
	tags.Get("/", userController.GetAllTags)
	tags.Post("/", middlewares.TagValidationMiddleware, userController.CreateTag)

	// Protected campanhas routes (com autenticação)
	campanhas := api.Group("/campanhas")
	campanhas.Get("/", userController.GetAllCampanhas)
	campanhas.Get("/:id", userController.GetCampanhaByID)
	campanhas.Post("/", middlewares.CampanhaValidationMiddleware, userController.CreateCampanha)
	campanhas.Post("/:id/publicos", userController.AssociarPublicosCampanha)
	campanhas.Get("/:id/publicos", userController.GetPublicosCampanha)

}
