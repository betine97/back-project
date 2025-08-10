package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/betine97/back-project.git/cmd/config/exceptions"
	dtos "github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/betine97/back-project.git/src/model/interfaces"
	"github.com/betine97/back-project.git/src/model/persistence"
	"github.com/betine97/back-project.git/src/model/service/crypto"
	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type ServiceInterface interface {
	CreateUserService(request dtos.CreateUser) (*entity.User, *exceptions.RestErr)
	LoginUserService(request dtos.UserLogin) (string, *exceptions.RestErr)

	GetAllFornecedoresService(userID string, page, limit int) (*dtos.FornecedorListResponse, *exceptions.RestErr)
	CreateFornecedorService(userID string, request dtos.CreateFornecedorRequest) (bool, *exceptions.RestErr)
	ChangeStatusFornecedorService(userID string, id string) (bool, *exceptions.RestErr)
	UpdateFornecedorFieldService(userID string, id string, campo string, valor string) (bool, *exceptions.RestErr)
	DeleteFornecedorService(userID string, id string) (bool, *exceptions.RestErr)

	GetAllProductsService(userID string, page, limit int) (*dtos.ProductListResponse, *exceptions.RestErr)
	CreateProductService(userID string, request dtos.CreateProductRequest) (bool, *exceptions.RestErr)
	DeleteProductService(userID string, id string) (bool, *exceptions.RestErr)

	GetAllPedidosService(userID string, page, limit int) (*dtos.PedidoListResponse, *exceptions.RestErr)
	GetPedidoByIdService(userID string, id string) (*dtos.PedidoResponse, *exceptions.RestErr)
	CreatePedidoService(userID string, request dtos.CreatePedidoRequest) (int, *exceptions.RestErr)

	// Itens de Pedido
	GetItensPedidoService(userID string, idPedido string, page, limit int) (*dtos.DetalhesPedidoListResponse, *exceptions.RestErr)
	CreateItemPedidoService(userID string, idPedido string, request dtos.CreateItemPedidoRequest) (bool, *exceptions.RestErr)

	// Estoque
	GetAllEstoqueService(userID string, page, limit int) (*dtos.DetalhesEstoqueListResponse, *exceptions.RestErr)
	CreateEstoqueService(userID string, request dtos.CreateEstoqueRequest) (bool, *exceptions.RestErr)

	// Clientes
	GetAllClientesService(userID string, page, limit int) (*dtos.ClienteListResponse, *exceptions.RestErr)
	BuscarClientesCriteriosService(userID string, idPublico string) (*dtos.ClienteCriterioListResponse, *exceptions.RestErr)
	AdicionarClientesAoPublicoService(userID string, idPublico string) (*dtos.AdicionarClientesPublicoResponse, *exceptions.RestErr)
	GetClienteByIDService(userID string, id string) (*dtos.ClienteResponse, *exceptions.RestErr)
	CreateClienteService(userID string, request dtos.CreateClienteRequest) (int, *exceptions.RestErr)
	DeleteClienteService(userID string, id string) (bool, *exceptions.RestErr)

	// Tags de Clientes
	AtribuirTagsClienteService(userID string, clienteID string, request dtos.AtribuirTagsClienteRequest) (bool, *exceptions.RestErr)
	RemoverTagsClienteService(userID string, clienteID string, request dtos.RemoverTagsClienteRequest) (bool, *exceptions.RestErr)
	GetTagsClienteService(userID string, clienteID string) (*dtos.TagsClienteListResponse, *exceptions.RestErr)

	// Tags
	GetAllTagsService(userID string, page, limit int) (*dtos.TagListResponse, *exceptions.RestErr)
	CreateTagService(userID string, request dtos.CreateTagRequest) (bool, *exceptions.RestErr)

	// Endereços
	GetAllEnderecosService(userID string, page, limit int) (*dtos.EnderecoListResponse, *exceptions.RestErr)
	CreateEnderecoService(userID string, request dtos.CreateEnderecoRequest) (bool, *exceptions.RestErr)
	DeleteEnderecoService(userID string, idEndereco string) (bool, *exceptions.RestErr)

	// Critérios
	GetAllCriteriosService(userID string) (*dtos.CriterioListResponse, *exceptions.RestErr)

	// Públicos
	GetAllPublicosService(userID string, page, limit int) (*dtos.PublicoListResponse, *exceptions.RestErr)
	CreatePublicoService(userID string, request dtos.CreatePublicoRequest) (int, *exceptions.RestErr)
	AssociarCriteriosPublicoService(userID string, idPublico string, request dtos.AssociarCriteriosRequest) (bool, *exceptions.RestErr)
	GetCriteriosPublicoService(userID string, idPublico string) (*dtos.PublicoCriterioListResponse, *exceptions.RestErr)
	GetClientesDoPublicoService(userID string, idPublico string, page, limit int) (*dtos.ClientesPublicoListResponse, *exceptions.RestErr)

	// Pets
	GetAllPetsService(userID string, page, limit int) (*dtos.PetListResponse, *exceptions.RestErr)
	CreatePetService(userID string, request dtos.CreatePetRequest) (int, *exceptions.RestErr)

	// Completude
	GetCompletudeClientesService(userID string, page, limit int) (*dtos.CompletudeListResponse, *exceptions.RestErr)

	// Campanhas
	GetAllCampanhasService(userID string, page, limit int) (*dtos.CampanhaListResponse, *exceptions.RestErr)
	GetCampanhaByIDService(userID string, id string) (*dtos.CampanhaResponse, *exceptions.RestErr)
	CreateCampanhaService(userID string, request dtos.CreateCampanhaRequest) (int, *exceptions.RestErr)
	AssociarPublicosCampanhaService(userID string, idCampanha string, request dtos.AssociarPublicosCampanhaRequest) (bool, *exceptions.RestErr)
	GetPublicosCampanhaService(userID string, idCampanha string) (*dtos.PublicosCampanhaListResponse, *exceptions.RestErr)
}

var ctx = context.Background()

type Service struct {
	crypto   crypto.CryptoInterface
	dbmaster persistence.PersistenceInterfaceDBMaster
	dbClient persistence.PersistenceInterfaceDBClient
	redis    interfaces.RedisInterface
	tokenGen interfaces.TokenGeneratorInterface
}

func NewServiceInstance(crypto crypto.CryptoInterface, dbmaster persistence.PersistenceInterfaceDBMaster, dbClient persistence.PersistenceInterfaceDBClient, redisClient interfaces.RedisInterface, tokenGen interfaces.TokenGeneratorInterface) ServiceInterface {
	return &Service{
		crypto:   crypto,
		dbmaster: dbmaster,
		dbClient: dbClient,
		redis:    redisClient,
		tokenGen: tokenGen,
	}
}

// TenantIDKey is a custom type for context keys to avoid collisions

// FUNÇÕES DE USUÁRIO ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) LoginUserService(request dtos.UserLogin) (string, *exceptions.RestErr) {
	zap.L().Info("Starting login service")

	// Verifica se o usuário existe
	user := srv.dbmaster.GetUser(request.Email)
	if user.Email == "" {
		zap.L().Warn("User not found", zap.String("email", request.Email))
		return "", exceptions.NewNotFoundError("Account not found")
	}

	zap.L().Info("User found for login", zap.Uint("user_id", user.ID), zap.String("email", user.Email))

	// Verifica a senha
	if _, err := srv.crypto.CheckPassword(request.Password, user.Password); err != nil {
		zap.L().Warn("Incorrect password", zap.String("email", request.Email))
		return "", exceptions.NewUnauthorizedRequestError("The password entered is incorrect")
	}

	// Verifica se o Redis está funcionando
	if err := srv.redis.Ping(ctx).Err(); err != nil {
		zap.L().Error("Redis is not reachable", zap.Error(err))
		return "", exceptions.NewInternalServerError("Redis is not reachable")
	}

	// Tenta obter os dados do usuário do Redis
	cacheKey := fmt.Sprintf("user:%d:db_info", user.ID)
	zap.L().Info("Trying to get tenant from Redis", zap.String("cache_key", cacheKey))
	tenantJSON, err := srv.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		zap.L().Info("Tenant not found in Redis cache, searching in database")
		// Dados não estão em cache, consulte a tabela tenants
		zap.L().Info("Searching for tenant in database", zap.Uint("user_id", user.ID))
		tenant := srv.dbmaster.GetTenantByUserID(user.ID)
		zap.L().Info("Tenant search result", zap.Uint("tenant_id", tenant.ID), zap.Uint("tenant_user_id", tenant.UserID))

		if tenant.ID == 0 {
			zap.L().Error("No tenant found for user", zap.Uint("user_id", user.ID))
			return "", exceptions.NewInternalServerError("Tenant not found")
		}

		// Armazene os dados no Redis
		tenantJSONBytes, err := json.Marshal(tenant)
		if err != nil {
			zap.L().Error("Error marshaling tenant to JSON", zap.Error(err))
			return "", exceptions.NewInternalServerError("Error storing tenant in cache")
		}
		tenantJSON = string(tenantJSONBytes)
		zap.L().Info("Tenant marshaled to JSON", zap.String("tenant_json", tenantJSON))

		if err := srv.redis.Set(ctx, cacheKey, tenantJSONBytes, 0*time.Second).Err(); err != nil {
			zap.L().Error("Error storing tenant in Redis", zap.Error(err))
		} else {
			zap.L().Info("Tenant stored in Redis successfully", zap.String("cache_key", cacheKey))
		}
	} else if err != nil {
		zap.L().Error("Error retrieving from Redis", zap.Error(err))
		return "", exceptions.NewInternalServerError("Error retrieving from cache")
	} else {
		zap.L().Info("Tenant found in Redis cache", zap.String("tenant_json", tenantJSON))
	}

	// Deserializa o JSON de volta para a estrutura Tenant
	zap.L().Info("About to unmarshal tenant JSON", zap.String("tenant_json", tenantJSON))
	var cachedTenant entity.Tenants
	if err := json.Unmarshal([]byte(tenantJSON), &cachedTenant); err != nil {
		zap.L().Error("Error unmarshaling tenant from JSON", zap.String("tenant_json", tenantJSON), zap.Error(err))
		return "", exceptions.NewInternalServerError("Error retrieving tenant data")
	}
	zap.L().Info("Tenant unmarshaled successfully", zap.Uint("tenant_id", cachedTenant.ID), zap.Uint("user_id", cachedTenant.UserID))

	// Gera o token com o ID do usuário
	token, erro := srv.tokenGen.GenerateToken(cachedTenant.UserID)
	if erro != nil {
		zap.L().Error("Error signing token", zap.Error(erro))
		return "", erro
	}

	zap.L().Info("Login successful", zap.String("email", request.Email), zap.String("db_info", tenantJSON))
	return token, nil
}

func (srv *Service) CreateUserService(request dtos.CreateUser) (*entity.User, *exceptions.RestErr) {
	zap.L().Info("Starting user creation service")

	emailExists := srv.dbmaster.GetUser(request.Email)

	if emailExists.Email != "" {
		zap.L().Warn("Email already associated with an existing account", zap.String("email", request.Email))
		return nil, exceptions.NewBadRequestError("Email is already associated with an existing account")
	}

	hashedPassword, err := srv.crypto.HashPassword(request.Password)
	if err != nil {
		zap.L().Error("Error when hashing password", zap.Error(err))
		return nil, exceptions.NewInternalServerError("Internal server error")
	}

	user := buildUserEntity(request, hashedPassword)

	dbErr := srv.dbmaster.CreateUser(*user)
	if dbErr != nil {
		zap.L().Error("Error creating user in database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("User created successfully", zap.String("email", user.Email))
	return user, nil
}

func buildUserEntity(request dtos.CreateUser, hashedPassword string) *entity.User {
	return &entity.User{
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		Email:       request.Email,
		NomeEmpresa: request.NomeEmpresa,
		Categoria:   request.Categoria,
		Segmento:    request.Segmento,
		City:        request.City,
		State:       request.State,
		Password:    hashedPassword,
	}
}

// FUNÇÕES DE FORNECEDORES ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllFornecedoresService(userID string, page, limit int) (*dtos.FornecedorListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all fornecedores service", zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Calcular offset
	offset := (page - 1) * limit

	fornecedores, total, dbErr := srv.dbClient.GetAllFornecedoresPaginated(userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting fornecedores from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving fornecedores")
	}

	fornecedorResponses := make([]dtos.FornecedorResponse, len(fornecedores))
	for i, fornecedor := range fornecedores {
		fornecedorResponses[i] = dtos.FornecedorResponse{
			ID:           fornecedor.ID,
			Nome:         fornecedor.Nome,
			Telefone:     fornecedor.Telefone,
			Email:        fornecedor.Email,
			Cidade:       fornecedor.Cidade,
			Estado:       fornecedor.Estado,
			Status:       fornecedor.Status,
			DataCadastro: fornecedor.DataCadastro,
		}
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit

	response := &dtos.FornecedorListResponse{
		Fornecedores: fornecedorResponses,
		Total:        total,
		Page:         page,
		Limit:        limit,
		TotalPages:   totalPages,
	}

	zap.L().Info("Successfully retrieved fornecedores", zap.Int("count", len(fornecedores)), zap.Int("total", total), zap.Int("page", page))
	return response, nil
}

func (srv *Service) CreateFornecedorService(userID string, request dtos.CreateFornecedorRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting fornecedor creation service")

	fornecedor := entity.BuildFornecedorEntity(request)

	dbErr := srv.dbClient.CreateFornecedor(*fornecedor, userID)
	if dbErr != nil {
		zap.L().Error("Error creating fornecedor in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Fornecedor created successfully", zap.String("fornecedor", fornecedor.Nome))
	return true, nil
}

func (srv *Service) ChangeStatusFornecedorService(userID string, id string) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting change status fornecedor service")

	// Recupera o fornecedor pelo ID
	fornecedor, dbErr := srv.dbClient.GetFornecedorById(id, userID)
	if dbErr != nil {
		zap.L().Error("Error getting fornecedor from database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Error retrieving fornecedor")
	}

	// Alterna o status
	if fornecedor.Status == "Ativo" {
		fornecedor.Status = "Inativo"
	} else {
		fornecedor.Status = "Ativo"
	}

	// Atualiza o fornecedor no banco de dados
	dbErr = srv.dbClient.UpdateFornecedor(*fornecedor, userID)
	if dbErr != nil {
		zap.L().Error("Error updating fornecedor in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Status fornecedor changed successfully", zap.String("fornecedor", fornecedor.Nome))
	return true, nil
}

func (srv *Service) UpdateFornecedorFieldService(userID string, id string, campo string, valor string) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting update fornecedor field service")

	// Verifica se o campo é válido
	validFields := map[string]bool{
		"nome":     true,
		"telefone": true,
		"email":    true,
		"cidade":   true,
		"estado":   true,
	}

	if !validFields[campo] {
		return false, exceptions.NewBadRequestError("Invalid field to update")
	}

	// Atualiza o campo no banco de dados
	dbErr := srv.dbClient.UpdateFornecedorField(id, campo, valor, userID)
	if dbErr != nil {
		zap.L().Error("Error updating fornecedor field in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Fornecedor field updated successfully", zap.String("id", id), zap.String("campo", campo), zap.String("valor", valor))
	return true, nil
}

func (srv *Service) DeleteFornecedorService(userID string, id string) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting delete fornecedor service")

	dbErr := srv.dbClient.DeleteFornecedor(id, userID)
	if dbErr != nil {
		zap.L().Error("Error deleting fornecedor in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Fornecedor deleted successfully", zap.String("fornecedor", id))
	return true, nil
}

// FUNÇÕES DE PRODUTOS ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllProductsService(userID string, page, limit int) (*dtos.ProductListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all products service", zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Calcular offset
	offset := (page - 1) * limit

	products, total, dbErr := srv.dbClient.GetAllProductsPaginated(userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting products from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving products")
	}

	productResponses := make([]dtos.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = dtos.ProductResponse{
			ID:            product.IDProduto,
			CodigoBarra:   product.CodigoBarra,
			NomeProduto:   product.NomeProduto,
			SKU:           product.SKU,
			Categoria:     product.Categoria,
			DestinadoPara: product.DestinadoPara,
			Variacao:      product.Variacao,
			Marca:         product.Marca,
			Descricao:     product.Descricao,
			Status:        product.Status,
			PrecoVenda:    product.PrecoVenda,
			IDFornecedor:  product.IDFornecedor,
		}
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit

	response := &dtos.ProductListResponse{
		Products:   productResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	zap.L().Info("Successfully retrieved products", zap.Int("count", len(products)), zap.Int("total", total), zap.Int("page", page))
	return response, nil
}

func (srv *Service) CreateProductService(userID string, request dtos.CreateProductRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting product creation service")

	// Validar se o código de barras já existe (apenas se não estiver vazio)
	if request.CodigoBarra != "" {
		existingProduct := srv.dbClient.GetProductByBarcode(request.CodigoBarra, userID)
		if existingProduct.CodigoBarra != "" {
			zap.L().Warn("Barcode already associated with an existing product", zap.String("barcode", request.CodigoBarra))
			return false, exceptions.NewBadRequestError("Este código de barras já está sendo usado por outro produto. Por favor, verifique e tente novamente.")
		}
	}

	product := entity.BuildProductEntity(request)

	dbErr := srv.dbClient.CreateProduct(*product, userID)
	if dbErr != nil {
		zap.L().Error("Error creating product in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Product created successfully", zap.String("product", product.NomeProduto))
	return true, nil
}

func (srv *Service) DeleteProductService(userID string, id string) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting delete product service")

	dbErr := srv.dbClient.DeleteProduct(id, userID)
	if dbErr != nil {
		zap.L().Error("Error deleting product in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Product deleted successfully", zap.String("product", id))
	return true, nil
}

// FUNÇÕES DE PEDIDOS ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllPedidosService(userID string, page, limit int) (*dtos.PedidoListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all pedidos service", zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Calcular offset
	offset := (page - 1) * limit

	pedidos, total, dbErr := srv.dbClient.GetAllPedidosPaginated(userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting pedidos from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving pedidos")
	}

	pedidoResponses := make([]dtos.PedidoResponse, len(pedidos))
	for i, pedido := range pedidos {
		pedidoResponses[i] = dtos.PedidoResponse{
			ID:           pedido.IDPedido,
			IDFornecedor: pedido.IDFornecedor,
			DataPedido:   pedido.DataPedido,
			DataEntrega:  pedido.DataEntrega,
			ValorFrete:   pedido.ValorFrete,
			CustoPedido:  pedido.CustoPedido,
			ValorTotal:   pedido.ValorTotal,
			Descricao:    pedido.Descricao,
			Status:       pedido.Status,
		}
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit

	response := &dtos.PedidoListResponse{
		Pedidos:    pedidoResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	zap.L().Info("Successfully retrieved pedidos", zap.Int("count", len(pedidos)), zap.Int("total", total), zap.Int("page", page))
	return response, nil
}

func (srv *Service) GetPedidoByIdService(userID string, id string) (*dtos.PedidoResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get pedido by ID service", zap.String("id", id))

	pedido, dbErr := srv.dbClient.GetPedidoById(id, userID)
	if dbErr != nil {
		zap.L().Error("Error getting pedido from database", zap.Error(dbErr))
		return nil, exceptions.NewNotFoundError("Pedido not found")
	}

	response := &dtos.PedidoResponse{
		ID:           pedido.IDPedido,
		IDFornecedor: pedido.IDFornecedor,
		DataPedido:   pedido.DataPedido,
		DataEntrega:  pedido.DataEntrega,
		ValorFrete:   pedido.ValorFrete,
		CustoPedido:  pedido.CustoPedido,
		ValorTotal:   pedido.ValorTotal,
		Descricao:    pedido.Descricao,
		Status:       pedido.Status,
	}

	zap.L().Info("Successfully retrieved pedido by ID", zap.String("id", id))
	return response, nil
}

func (srv *Service) CreatePedidoService(userID string, request dtos.CreatePedidoRequest) (int, *exceptions.RestErr) {
	zap.L().Info("Starting pedido creation service")

	pedido := entity.BuildPedidoEntity(request)

	dbErr := srv.dbClient.CreatePedido(pedido, userID)
	if dbErr != nil {
		zap.L().Error("Error creating pedido in database", zap.Error(dbErr))
		return 0, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Pedido created successfully", zap.String("descricao", pedido.Descricao), zap.Int("id", pedido.IDPedido))
	return pedido.IDPedido, nil
}

// FUNÇÕES DE ITENS DE PEDIDO ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetItensPedidoService(userID string, idPedido string, page, limit int) (*dtos.DetalhesPedidoListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get detalhes pedido service using view", zap.String("idPedido", idPedido), zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Calcular offset
	offset := (page - 1) * limit

	detalhes, total, dbErr := srv.dbClient.GetDetalhesPedidoPaginated(idPedido, userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting detalhes pedido from view", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving detalhes pedido")
	}

	detalheResponses := make([]dtos.DetalhesPedidoResponse, len(detalhes))
	for i, detalhe := range detalhes {
		detalheResponses[i] = dtos.DetalhesPedidoResponse{
			IDPedido:      detalhe.IDPedido,
			NomeProduto:   detalhe.NomeProduto,
			Quantidade:    detalhe.Quantidade,
			PrecoUnitario: detalhe.PrecoUnitario,
			TotalItem:     detalhe.TotalItem,
		}
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit

	// Converter idPedido para int
	idPedidoInt := 0
	if len(detalhes) > 0 {
		idPedidoInt = detalhes[0].IDPedido
	}

	response := &dtos.DetalhesPedidoListResponse{
		Detalhes:   detalheResponses,
		Total:      total,
		IDPedido:   idPedidoInt,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	zap.L().Info("Successfully retrieved detalhes pedido from view", zap.String("idPedido", idPedido), zap.Int("count", len(detalhes)), zap.Int("total", total))
	return response, nil
}

func (srv *Service) CreateItemPedidoService(userID string, idPedido string, request dtos.CreateItemPedidoRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting item pedido creation service", zap.String("idPedido", idPedido))

	// Converter idPedido string para int
	idPedidoInt := 0
	if _, err := fmt.Sscanf(idPedido, "%d", &idPedidoInt); err != nil {
		zap.L().Error("Error converting idPedido to int", zap.Error(err))
		return false, exceptions.NewBadRequestError("Invalid pedido ID")
	}

	item := entity.BuildItemPedidoEntity(request, idPedidoInt)

	dbErr := srv.dbClient.CreateItemPedido(*item, userID)
	if dbErr != nil {
		zap.L().Error("Error creating item pedido in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Item pedido created successfully", zap.String("idPedido", idPedido), zap.Int("idProduto", item.IDProduto))
	return true, nil
}

// FUNÇÕES DE ESTOQUE ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllEstoqueService(userID string, page, limit int) (*dtos.DetalhesEstoqueListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all detalhes estoque service using view", zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	offset := (page - 1) * limit

	detalhesEstoque, total, dbErr := srv.dbClient.GetAllDetalhesEstoquePaginated(userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting detalhes estoque from view", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Internal server error")
	}

	// Converter entidades da view para DTOs
	var estoqueResponse []dtos.DetalhesEstoqueResponse
	for _, item := range detalhesEstoque {
		estoqueResponse = append(estoqueResponse, dtos.DetalhesEstoqueResponse{
			NomeProduto:         item.NomeProduto,
			Lote:                item.Lote,
			Quantidade:          item.Quantidade,
			DataEntrada:         item.DataEntrada,
			DataSaida:           item.DataSaida,
			Vencimento:          item.Vencimento,
			DocumentoReferencia: item.DocumentoReferencia,
			Status:              item.Status,
		})
	}

	totalPages := (total + limit - 1) / limit

	response := &dtos.DetalhesEstoqueListResponse{
		Estoque:    estoqueResponse,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	zap.L().Info("Detalhes estoque service completed successfully using view", zap.Int("total", total))
	return response, nil
}

func (srv *Service) CreateEstoqueService(userID string, request dtos.CreateEstoqueRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting estoque creation service", zap.Int("id_produto", request.IDProduto))

	estoque := entity.BuildEstoqueEntity(request)

	dbErr := srv.dbClient.CreateEstoque(*estoque, userID)
	if dbErr != nil {
		zap.L().Error("Error creating estoque in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Estoque created successfully", zap.Int("id_produto", estoque.IDProduto))
	return true, nil
}

// FUNÇÕES DE CLIENTES ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllClientesService(userID string, page, limit int) (*dtos.ClienteListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all clientes service", zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Calcular offset
	offset := (page - 1) * limit

	clientes, total, dbErr := srv.dbClient.GetAllClientesPaginated(userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting clientes from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving clientes")
	}

	clienteResponses := make([]dtos.ClienteResponse, len(clientes))
	for i, cliente := range clientes {
		clienteResponses[i] = dtos.ClienteResponse{
			ID:             cliente.ID,
			TipoCliente:    cliente.TipoCliente,
			NomeCliente:    cliente.NomeCliente,
			NumeroCelular:  cliente.NumeroCelular,
			Sexo:           cliente.Sexo,
			Email:          cliente.Email,
			DataNascimento: cliente.DataNascimento,
			DataCadastro:   cliente.DataCadastro,
		}
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit

	response := &dtos.ClienteListResponse{
		Clientes:   clienteResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	zap.L().Info("Successfully retrieved clientes", zap.Int("count", len(clientes)), zap.Int("total", total), zap.Int("page", page))
	return response, nil
}

func (srv *Service) BuscarClientesCriteriosService(userID string, idPublico string) (*dtos.ClienteCriterioListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting buscar clientes criterios service", zap.String("idPublico", idPublico))

	// Primeiro, buscar os critérios do público
	criterios, dbErr := srv.dbClient.GetCriteriosPublico(idPublico, userID)
	if dbErr != nil {
		zap.L().Error("Error getting criterios for publico", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving criterios for publico")
	}

	if len(criterios) == 0 {
		zap.L().Warn("No criterios found for publico", zap.String("idPublico", idPublico))
		return &dtos.ClienteCriterioListResponse{
			Clientes: []dtos.ClienteCriterioResponse{},
			Total:    0,
		}, nil
	}

	// Buscar clientes que atendem aos critérios
	clientes, dbErr := srv.dbClient.BuscarClientesPorCriterios(userID, criterios)
	if dbErr != nil {
		zap.L().Error("Error getting clientes by criterios from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving clientes for criterios")
	}

	clienteResponses := make([]dtos.ClienteCriterioResponse, len(clientes))
	for i, cliente := range clientes {
		clienteResponses[i] = dtos.ClienteCriterioResponse{
			ID:          cliente.ID,
			TipoCliente: cliente.TipoCliente,
			Sexo:        cliente.Sexo,
		}
	}

	response := &dtos.ClienteCriterioListResponse{
		Clientes: clienteResponses,
		Total:    len(clientes),
	}

	zap.L().Info("Successfully retrieved clientes for criterios", zap.Int("total", len(clientes)), zap.Int("criterios_count", len(criterios)))
	return response, nil
}

func (srv *Service) AdicionarClientesAoPublicoService(userID string, idPublico string) (*dtos.AdicionarClientesPublicoResponse, *exceptions.RestErr) {
	zap.L().Info("Starting adicionar clientes ao publico service", zap.String("idPublico", idPublico))

	// Converter idPublico string para int
	idPublicoInt := 0
	if _, err := fmt.Sscanf(idPublico, "%d", &idPublicoInt); err != nil {
		zap.L().Error("Error converting idPublico to int", zap.Error(err))
		return nil, exceptions.NewBadRequestError("Invalid publico ID")
	}

	// Primeiro, buscar os critérios do público
	criterios, dbErr := srv.dbClient.GetCriteriosPublico(idPublico, userID)
	if dbErr != nil {
		zap.L().Error("Error getting criterios for publico", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving criterios for publico")
	}

	if len(criterios) == 0 {
		zap.L().Warn("No criterios found for publico", zap.String("idPublico", idPublico))
		return &dtos.AdicionarClientesPublicoResponse{
			ClientesAdicionados: 0,
			ClientesJaExistiam:  0,
			ClientesEncontrados: []dtos.ClienteCriterioResponse{},
			Total:               0,
		}, nil
	}

	// Buscar clientes que atendem aos critérios
	clientes, dbErr := srv.dbClient.BuscarClientesPorCriterios(userID, criterios)
	if dbErr != nil {
		zap.L().Error("Error getting clientes by criterios from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving clientes for criterios")
	}

	if len(clientes) == 0 {
		zap.L().Info("No clientes found matching criterios", zap.String("idPublico", idPublico))
		return &dtos.AdicionarClientesPublicoResponse{
			ClientesAdicionados: 0,
			ClientesJaExistiam:  0,
			ClientesEncontrados: []dtos.ClienteCriterioResponse{},
			Total:               0,
		}, nil
	}

	// Adicionar clientes ao público
	clientesAdicionados, clientesJaExistiam, dbErr := srv.dbClient.AdicionarClientesAoPublico(userID, idPublicoInt, clientes)
	if dbErr != nil {
		zap.L().Error("Error adding clientes to publico", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error adding clientes to publico")
	}

	// Preparar resposta
	clienteResponses := make([]dtos.ClienteCriterioResponse, len(clientes))
	for i, cliente := range clientes {
		clienteResponses[i] = dtos.ClienteCriterioResponse{
			ID:          cliente.ID,
			TipoCliente: cliente.TipoCliente,
			Sexo:        cliente.Sexo,
		}
	}

	response := &dtos.AdicionarClientesPublicoResponse{
		ClientesAdicionados: clientesAdicionados,
		ClientesJaExistiam:  clientesJaExistiam,
		ClientesEncontrados: clienteResponses,
		Total:               len(clientes),
	}

	zap.L().Info("Successfully added clientes to publico",
		zap.Int("total_encontrados", len(clientes)),
		zap.Int("adicionados", clientesAdicionados),
		zap.Int("ja_existiam", clientesJaExistiam),
		zap.Int("criterios_count", len(criterios)))
	return response, nil
}

func (srv *Service) GetClienteByIDService(userID string, id string) (*dtos.ClienteResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get cliente by ID service", zap.String("id", id))

	cliente := srv.dbClient.GetClienteByID(id, userID)
	if cliente.ID == 0 {
		zap.L().Warn("Cliente not found by ID", zap.String("id", id))
		return nil, exceptions.NewNotFoundError("Cliente not found")
	}

	response := &dtos.ClienteResponse{
		ID:             cliente.ID,
		TipoCliente:    cliente.TipoCliente,
		NomeCliente:    cliente.NomeCliente,
		NumeroCelular:  cliente.NumeroCelular,
		Sexo:           cliente.Sexo,
		Email:          cliente.Email,
		DataNascimento: cliente.DataNascimento,
		DataCadastro:   cliente.DataCadastro,
	}

	zap.L().Info("Successfully retrieved cliente by ID", zap.String("id", id))
	return response, nil
}

func (srv *Service) CreateClienteService(userID string, request dtos.CreateClienteRequest) (int, *exceptions.RestErr) {
	zap.L().Info("Starting cliente creation service")

	// Validar se o email já existe
	if request.Email != "" {
		existingClienteByEmail := srv.dbClient.GetClienteByEmail(request.Email, userID)
		if existingClienteByEmail.Email != "" {
			zap.L().Warn("Email already associated with an existing cliente", zap.String("email", request.Email))
			return 0, exceptions.NewBadRequestError("Este e-mail já está associado a um cliente cadastrado anteriormente. Por favor, verifique e tente novamente.")
		}
	}

	// Validar se o telefone já existe
	if request.NumeroCelular != "" {
		existingClienteByTelefone := srv.dbClient.GetClienteByTelefone(request.NumeroCelular, userID)
		if existingClienteByTelefone.NumeroCelular != "" {
			zap.L().Warn("Telefone already associated with an existing cliente", zap.String("telefone", request.NumeroCelular))
			return 0, exceptions.NewBadRequestError("Este número de telefone já está associado a um cliente cadastrado anteriormente. Por favor, verifique e tente novamente.")
		}
	}

	cliente := entity.BuildClienteEntity(request)

	dbErr := srv.dbClient.CreateCliente(*cliente, userID)
	if dbErr != nil {
		zap.L().Error("Error creating cliente in database", zap.Error(dbErr))
		return 0, exceptions.NewInternalServerError("Internal server error")
	}

	searchIDuser := srv.dbClient.GetClienteByTelefone(cliente.NumeroCelular, userID)
	if searchIDuser != nil {
		zap.L().Info("Cliente created successfully", zap.String("nome_cliente", cliente.NomeCliente), zap.String("tipo_cliente", cliente.TipoCliente), zap.Int("id", searchIDuser.ID))
		return searchIDuser.ID, nil
	} else {
		zap.L().Error("Error retrieving cliente by phone")
		return 0, exceptions.NewInternalServerError("Internal server error")
	}

}

// FUNÇÕES DE ENDEREÇOS ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllEnderecosService(userID string, page, limit int) (*dtos.EnderecoListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all enderecos service", zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Calcular offset
	offset := (page - 1) * limit

	enderecos, total, dbErr := srv.dbClient.GetAllEnderecosPaginated(userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting enderecos from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving enderecos")
	}

	enderecoResponses := make([]dtos.EnderecoResponse, len(enderecos))
	for i, endereco := range enderecos {
		enderecoResponses[i] = dtos.EnderecoResponse{
			IDEndereco:  endereco.IDEndereco,
			IDCliente:   endereco.IDCliente,
			CEP:         endereco.CEP,
			Cidade:      endereco.Cidade,
			Estado:      endereco.Estado,
			Bairro:      endereco.Bairro,
			Logradouro:  endereco.Logradouro,
			Numero:      endereco.Numero,
			Complemento: endereco.Complemento,
			Obs:         endereco.Obs,
		}
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit

	response := &dtos.EnderecoListResponse{
		Enderecos:  enderecoResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	zap.L().Info("Successfully retrieved enderecos", zap.Int("count", len(enderecos)), zap.Int("total", total), zap.Int("page", page))
	return response, nil
}

func (srv *Service) CreateEnderecoService(userID string, request dtos.CreateEnderecoRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting endereco creation service")

	endereco := entity.BuildEnderecoEntity(request)

	dbErr := srv.dbClient.CreateEndereco(*endereco, userID)
	if dbErr != nil {
		zap.L().Error("Error creating endereco in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Endereco created successfully", zap.Int("id_cliente", endereco.IDCliente), zap.String("cidade", endereco.Cidade))
	return true, nil
}

func (srv *Service) DeleteClienteService(userID string, id string) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting delete cliente service")

	dbErr := srv.dbClient.DeleteCliente(id, userID)
	if dbErr != nil {
		zap.L().Error("Error deleting cliente in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Cliente deleted successfully", zap.String("cliente_id", id))
	return true, nil
}

func (srv *Service) DeleteEnderecoService(userID string, idEndereco string) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting delete endereco service")

	dbErr := srv.dbClient.DeleteEndereco(idEndereco, userID)
	if dbErr != nil {
		zap.L().Error("Error deleting endereco in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Endereco deleted successfully", zap.String("endereco_id", idEndereco))
	return true, nil
}

// FUNÇÕES DE CRITÉRIOS ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllCriteriosService(userID string) (*dtos.CriterioListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all criterios service")

	criterios, dbErr := srv.dbClient.GetAllCriterios(userID)
	if dbErr != nil {
		zap.L().Error("Error getting criterios from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving criterios")
	}

	criterioResponses := make([]dtos.CriterioResponse, len(criterios))
	for i, criterio := range criterios {
		criterioResponses[i] = dtos.CriterioResponse{
			ID:           criterio.ID,
			NomeCondicao: criterio.NomeCondicao,
		}
	}

	response := &dtos.CriterioListResponse{
		Criterios: criterioResponses,
		Total:     len(criterios),
	}

	zap.L().Info("Successfully retrieved criterios", zap.Int("count", len(criterios)))
	return response, nil
}

// FUNÇÕES DE PÚBLICOS ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllPublicosService(userID string, page, limit int) (*dtos.PublicoListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all publicos service", zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Calcular offset
	offset := (page - 1) * limit

	publicos, total, dbErr := srv.dbClient.GetAllPublicosPaginated(userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting publicos from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving publicos")
	}

	publicoResponses := make([]dtos.PublicoResponse, len(publicos))
	for i, publico := range publicos {
		publicoResponses[i] = dtos.PublicoResponse{
			ID:          publico.ID,
			Nome:        publico.Nome,
			Descricao:   publico.Descricao,
			DataCriacao: publico.DataCriacao,
			Status:      publico.Status,
		}
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit

	response := &dtos.PublicoListResponse{
		Publicos:   publicoResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	zap.L().Info("Successfully retrieved publicos", zap.Int("count", len(publicos)), zap.Int("total", total), zap.Int("page", page))
	return response, nil
}

func (srv *Service) CreatePublicoService(userID string, request dtos.CreatePublicoRequest) (int, *exceptions.RestErr) {
	zap.L().Info("Starting publico creation service")

	publico := entity.BuildPublicoClienteEntity(request)

	dbErr := srv.dbClient.CreatePublico(publico, userID)
	if dbErr != nil {
		zap.L().Error("Error creating publico in database", zap.Error(dbErr))
		return 0, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Publico created successfully", zap.String("nome", publico.Nome), zap.Int("id", publico.ID))
	return publico.ID, nil
}

func (srv *Service) AssociarCriteriosPublicoService(userID string, idPublico string, request dtos.AssociarCriteriosRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting associar criterios publico service", zap.String("idPublico", idPublico))

	// Converter idPublico string para int
	idPublicoInt := 0
	if _, err := fmt.Sscanf(idPublico, "%d", &idPublicoInt); err != nil {
		zap.L().Error("Error converting idPublico to int", zap.Error(err))
		return false, exceptions.NewBadRequestError("Invalid publico ID")
	}

	dbErr := srv.dbClient.AssociarCriteriosPublico(idPublicoInt, request.Criterios, userID)
	if dbErr != nil {
		zap.L().Error("Error associating criterios to publico", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Criterios associated successfully", zap.String("idPublico", idPublico), zap.Ints("criterios", request.Criterios))
	return true, nil
}

func (srv *Service) GetCriteriosPublicoService(userID string, idPublico string) (*dtos.PublicoCriterioListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get criterios publico service", zap.String("idPublico", idPublico))

	criterios, dbErr := srv.dbClient.GetCriteriosPublico(idPublico, userID)
	if dbErr != nil {
		zap.L().Error("Error getting criterios for publico", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving criterios for publico")
	}

	criterioResponses := make([]dtos.PublicoCriterioResponse, len(criterios))
	for i, criterio := range criterios {
		criterioResponses[i] = dtos.PublicoCriterioResponse{
			IDPublico:    criterio.IDPublico,
			IDCriterio:   criterio.IDCriterio,
			NomeCondicao: criterio.NomeCondicao,
		}
	}

	response := &dtos.PublicoCriterioListResponse{
		Criterios: criterioResponses,
		Total:     len(criterios),
	}

	zap.L().Info("Successfully retrieved criterios for publico", zap.String("idPublico", idPublico), zap.Int("count", len(criterios)))
	return response, nil
}

func (srv *Service) GetClientesDoPublicoService(userID string, idPublico string, page, limit int) (*dtos.ClientesPublicoListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get clientes do publico service", zap.String("idPublico", idPublico), zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Converter idPublico string para int
	idPublicoInt := 0
	if _, err := fmt.Sscanf(idPublico, "%d", &idPublicoInt); err != nil {
		zap.L().Error("Error converting idPublico to int", zap.Error(err))
		return nil, exceptions.NewBadRequestError("Invalid publico ID")
	}

	// Calcular offset
	offset := (page - 1) * limit

	clientes, total, dbErr := srv.dbClient.GetClientesDoPublico(userID, idPublicoInt, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting clientes do publico from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving clientes do publico")
	}

	clienteResponses := make([]dtos.ClientePublicoResponse, len(clientes))
	for i, cliente := range clientes {
		clienteResponses[i] = dtos.ClientePublicoResponse{
			ID:             cliente.ID,
			TipoCliente:    cliente.TipoCliente,
			NomeCliente:    cliente.NomeCliente,
			NumeroCelular:  cliente.NumeroCelular,
			Sexo:           cliente.Sexo,
			Email:          cliente.Email,
			DataNascimento: cliente.DataNascimento,
			DataCadastro:   cliente.DataCadastro,
			DataAdicao:     "N/A", // Será implementado quando tivermos timestamp na tabela
		}
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit

	response := &dtos.ClientesPublicoListResponse{
		Clientes:   clienteResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		IDPublico:  idPublicoInt,
	}

	zap.L().Info("Successfully retrieved clientes do publico", zap.Int("count", len(clientes)), zap.Int("total", total), zap.Int("page", page))
	return response, nil
}

// FUNÇÕES DE PETS ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllPetsService(userID string, page, limit int) (*dtos.PetListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all pets service", zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Calcular offset
	offset := (page - 1) * limit

	pets, total, dbErr := srv.dbClient.GetAllPetsPaginated(userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting pets from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving pets")
	}

	petResponses := make([]dtos.PetResponse, len(pets))
	for i, pet := range pets {
		petResponse := dtos.PetResponse{
			IDPet:        pet.IDPet,
			ClienteID:    pet.ClienteID,
			NomePet:      pet.NomePet,
			Especie:      pet.Especie,
			Raca:         pet.Raca,
			Porte:        pet.Porte,
			Idade:        pet.Idade,
			DataRegistro: pet.DataRegistro.Format("2006-01-02 15:04:05"),
		}

		// Formatar data de aniversário se existir
		if pet.DataAniversario != nil {
			petResponse.DataAniversario = pet.DataAniversario.Format("2006-01-02")
		}

		petResponses[i] = petResponse
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit

	response := &dtos.PetListResponse{
		Pets:       petResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	zap.L().Info("Successfully retrieved pets", zap.Int("count", len(pets)), zap.Int("total", total), zap.Int("page", page))
	return response, nil
}

func (srv *Service) CreatePetService(userID string, request dtos.CreatePetRequest) (int, *exceptions.RestErr) {
	zap.L().Info("Starting pet creation service")

	// Validar se o cliente existe
	cliente := srv.dbClient.GetClienteByID(fmt.Sprintf("%d", request.ClienteID), userID)
	if cliente.ID == 0 {
		zap.L().Warn("Cliente not found for pet creation", zap.Int("cliente_id", request.ClienteID))
		return 0, exceptions.NewBadRequestError("Cliente não encontrado. Verifique o ID do cliente e tente novamente.")
	}

	pet := entity.BuildPetEntity(request)

	dbErr := srv.dbClient.CreatePet(pet, userID)
	if dbErr != nil {
		zap.L().Error("Error creating pet in database", zap.Error(dbErr))
		return 0, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Pet created successfully", zap.String("nome_pet", pet.NomePet), zap.String("especie", pet.Especie), zap.Int("id_pet", pet.IDPet))
	return pet.IDPet, nil
}

// GetCompletudeClientesService calcula a completude do cadastro de clientes e pets
func (srv *Service) GetCompletudeClientesService(userID string, page, limit int) (*dtos.CompletudeListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting completude clientes service")

	offset := (page - 1) * limit

	clientes, pets, enderecos, total, dbErr := srv.dbClient.GetClientesComPetsEnderecosParaCompletude(userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error fetching clientes and pets for completude", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Internal server error")
	}

	// Mapear pets e endereços por cliente_id para facilitar o acesso
	petsPorCliente := make(map[int][]entity.Pet)
	for _, pet := range pets {
		petsPorCliente[pet.ClienteID] = append(petsPorCliente[pet.ClienteID], pet)
	}

	enderecosPorCliente := make(map[int][]entity.Endereco)
	for _, endereco := range enderecos {
		enderecosPorCliente[endereco.IDCliente] = append(enderecosPorCliente[endereco.IDCliente], endereco)
	}

	var completudeClientes []dtos.CompletudeClienteResponse

	for _, cliente := range clientes {
		// Calcular campos preenchidos do cliente
		camposClientePreenchidos := calcularCamposPreenchidosCliente(cliente)
		totalCamposCliente := 7

		// Calcular campos preenchidos dos pets
		petsDoCliente := petsPorCliente[cliente.ID]
		totalCamposPets := len(petsDoCliente) * 6
		camposPetsPreenchidos := 0
		for _, pet := range petsDoCliente {
			camposPetsPreenchidos += calcularCamposPreenchidosPet(pet)
		}

		// Calcular campos preenchidos dos endereços
		enderecosDoCliente := enderecosPorCliente[cliente.ID]
		totalCamposEnderecos := len(enderecosDoCliente) * 6
		camposEnderecosPreenchidos := 0
		for _, endereco := range enderecosDoCliente {
			camposEnderecosPreenchidos += calcularCamposPreenchidosEndereco(endereco)
		}

		// Calcular percentual geral único
		totalCamposGeral := totalCamposCliente + totalCamposPets + totalCamposEnderecos
		camposPreenchidosGeral := camposClientePreenchidos + camposPetsPreenchidos + camposEnderecosPreenchidos

		var percentualCompleto float64
		if totalCamposGeral > 0 {
			// Calcular percentual e arredondar para cima com 2 casas decimais
			percentualBruto := float64(camposPreenchidosGeral) / float64(totalCamposGeral) * 100
			// Multiplicar por 100, arredondar para cima, depois dividir por 100 para ter 2 casas decimais
			percentualCompleto = math.Ceil(percentualBruto*100) / 100
		}

		completudeClientes = append(completudeClientes, dtos.CompletudeClienteResponse{
			ClienteID:          cliente.ID,
			NomeCliente:        cliente.NomeCliente,
			PercentualCompleto: percentualCompleto,
		})
	}

	totalPages := (total + limit - 1) / limit

	response := &dtos.CompletudeListResponse{
		Clientes:   completudeClientes,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	zap.L().Info("Completude clientes service completed successfully", zap.Int("total_clientes", len(completudeClientes)))
	return response, nil
}

// calcularCamposPreenchidosCliente conta quantos campos do cliente estão preenchidos
func calcularCamposPreenchidosCliente(cliente entity.Cliente) int {
	count := 0
	if cliente.TipoCliente != "" {
		count++
	}
	if cliente.NomeCliente != "" {
		count++
	}
	if cliente.NumeroCelular != "" {
		count++
	}
	if cliente.Sexo != "" {
		count++
	}
	if cliente.Email != "" {
		count++
	}
	if cliente.DataNascimento != "" {
		count++
	}
	if cliente.DataCadastro != "" {
		count++
	}
	return count
}

// calcularCamposPreenchidosPet conta quantos campos do pet estão preenchidos
func calcularCamposPreenchidosPet(pet entity.Pet) int {
	count := 0
	if pet.NomePet != "" {
		count++
	}
	if pet.Especie != "" {
		count++
	}
	if pet.Raca != "" {
		count++
	}
	if pet.Porte != "" {
		count++
	}
	if pet.DataAniversario != nil {
		count++
	}
	if pet.Idade != nil {
		count++
	}
	return count
}

// calcularCamposPreenchidosEndereco conta quantos campos do endereço estão preenchidos
func calcularCamposPreenchidosEndereco(endereco entity.Endereco) int {
	count := 0
	if endereco.CEP != "" {
		count++
	}
	if endereco.Logradouro != "" {
		count++
	}
	if endereco.Numero != "" {
		count++
	}
	if endereco.Bairro != "" {
		count++
	}
	if endereco.Cidade != "" {
		count++
	}
	if endereco.Estado != "" {
		count++
	}
	return count
}

// FUNÇÕES DE TAGS DE CLIENTES ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) AtribuirTagsClienteService(userID string, clienteID string, request dtos.AtribuirTagsClienteRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting atribuir tags cliente service", zap.String("clienteID", clienteID))

	// Converter clienteID string para int
	clienteIDInt, err := strconv.Atoi(clienteID)
	if err != nil {
		zap.L().Error("Error converting clienteID to int", zap.Error(err))
		return false, exceptions.NewBadRequestError("Invalid cliente ID")
	}

	// Verificar se o cliente existe
	cliente := srv.dbClient.GetClienteByID(clienteID, userID)
	if cliente.ID == 0 {
		zap.L().Warn("Cliente not found", zap.String("clienteID", clienteID))
		return false, exceptions.NewNotFoundError("Cliente not found")
	}

	dbErr := srv.dbClient.AtribuirTagsCliente(clienteIDInt, request.IDsTags, userID)
	if dbErr != nil {
		zap.L().Error("Error atribuindo tags ao cliente", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Tags atribuídas ao cliente com sucesso", zap.String("clienteID", clienteID), zap.Ints("tagIDs", request.IDsTags))
	return true, nil
}

func (srv *Service) RemoverTagsClienteService(userID string, clienteID string, request dtos.RemoverTagsClienteRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting remover tags cliente service", zap.String("clienteID", clienteID))

	// Converter clienteID string para int
	clienteIDInt, err := strconv.Atoi(clienteID)
	if err != nil {
		zap.L().Error("Error converting clienteID to int", zap.Error(err))
		return false, exceptions.NewBadRequestError("Invalid cliente ID")
	}

	// Verificar se o cliente existe
	cliente := srv.dbClient.GetClienteByID(clienteID, userID)
	if cliente.ID == 0 {
		zap.L().Warn("Cliente not found", zap.String("clienteID", clienteID))
		return false, exceptions.NewNotFoundError("Cliente not found")
	}

	dbErr := srv.dbClient.RemoverTagsCliente(clienteIDInt, request.IDsTags, userID)
	if dbErr != nil {
		zap.L().Error("Error removendo tags do cliente", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Tags removidas do cliente com sucesso", zap.String("clienteID", clienteID), zap.Ints("tagIDs", request.IDsTags))
	return true, nil
}

func (srv *Service) GetTagsClienteService(userID string, clienteID string) (*dtos.TagsClienteListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get tags cliente service", zap.String("clienteID", clienteID))

	// Converter clienteID string para int
	clienteIDInt, err := strconv.Atoi(clienteID)
	if err != nil {
		zap.L().Error("Error converting clienteID to int", zap.Error(err))
		return nil, exceptions.NewBadRequestError("Invalid cliente ID")
	}

	// Verificar se o cliente existe
	cliente := srv.dbClient.GetClienteByID(clienteID, userID)
	if cliente.ID == 0 {
		zap.L().Warn("Cliente not found", zap.String("clienteID", clienteID))
		return nil, exceptions.NewNotFoundError("Cliente not found")
	}

	tags, dbErr := srv.dbClient.GetTagsCliente(clienteIDInt, userID)
	if dbErr != nil {
		zap.L().Error("Error getting tags do cliente", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Internal server error")
	}

	tagResponses := make([]dtos.TagClienteResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = dtos.TagClienteResponse{
			ID:    tag.ID,
			IDTag: tag.IDTag,
			Nome:  tag.Nome,
		}
	}

	response := &dtos.TagsClienteListResponse{
		Tags:      tagResponses,
		ClienteID: clienteIDInt,
		Total:     len(tags),
	}

	zap.L().Info("Successfully retrieved tags do cliente", zap.String("clienteID", clienteID), zap.Int("total", len(tags)))
	return response, nil
}

// FUNÇÕES DE TAGS ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllTagsService(userID string, page, limit int) (*dtos.TagListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all tags service", zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Calcular offset
	offset := (page - 1) * limit

	tags, total, dbErr := srv.dbClient.GetAllTagsPaginated(userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting tags from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving tags")
	}

	tagResponses := make([]dtos.TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = dtos.TagResponse{
			IDTag:        tag.IDTag,
			CategoriaTag: tag.CategoriaTag,
			NomeTag:      tag.NomeTag,
		}
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := &dtos.TagListResponse{
		Tags:       tagResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	zap.L().Info("Successfully retrieved tags", zap.Int("count", len(tags)), zap.Int("total", total), zap.Int("page", page))
	return response, nil
}

func (srv *Service) CreateTagService(userID string, request dtos.CreateTagRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting tag creation service")

	tag := entity.BuildTagEntity(request)

	dbErr := srv.dbClient.CreateTag(*tag, userID)
	if dbErr != nil {
		zap.L().Error("Error creating tag in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Tag created successfully", zap.String("nome", tag.NomeTag), zap.String("categoria", tag.CategoriaTag))
	return true, nil
}

// FUNÇÕES DE CAMPANHAS ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllCampanhasService(userID string, page, limit int) (*dtos.CampanhaListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all campanhas service", zap.Int("page", page), zap.Int("limit", limit))

	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Calcular offset
	offset := (page - 1) * limit

	campanhas, total, dbErr := srv.dbClient.GetAllCampanhasPaginated(userID, limit, offset)
	if dbErr != nil {
		zap.L().Error("Error getting campanhas from database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving campanhas")
	}

	campanhaResponses := make([]dtos.CampanhaResponse, len(campanhas))
	for i, campanha := range campanhas {
		campanhaResponses[i] = dtos.CampanhaResponse{
			ID:             campanha.ID,
			Nome:           campanha.Nome,
			Desc:           campanha.Desc,
			DataCriacao:    campanha.DataCriacao,
			DataLancamento: campanha.DataLancamento,
			DataFim:        campanha.DataFim,
			Status:         campanha.Status,
		}
	}

	// Calcular total de páginas
	totalPages := (total + limit - 1) / limit

	response := &dtos.CampanhaListResponse{
		Campanhas:  campanhaResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	zap.L().Info("Successfully retrieved campanhas", zap.Int("count", len(campanhas)), zap.Int("total", total), zap.Int("page", page))
	return response, nil
}

func (srv *Service) GetCampanhaByIDService(userID string, id string) (*dtos.CampanhaResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get campanha by ID service", zap.String("id", id))

	campanha := srv.dbClient.GetCampanhaByID(id, userID)
	if campanha.ID == 0 {
		zap.L().Warn("Campanha not found by ID", zap.String("id", id))
		return nil, exceptions.NewNotFoundError("Campanha not found")
	}

	response := &dtos.CampanhaResponse{
		ID:             campanha.ID,
		Nome:           campanha.Nome,
		Desc:           campanha.Desc,
		DataCriacao:    campanha.DataCriacao,
		DataLancamento: campanha.DataLancamento,
		DataFim:        campanha.DataFim,
		Status:         campanha.Status,
	}

	zap.L().Info("Successfully retrieved campanha by ID", zap.String("id", id))
	return response, nil
}

func (srv *Service) CreateCampanhaService(userID string, request dtos.CreateCampanhaRequest) (int, *exceptions.RestErr) {
	zap.L().Info("Starting campanha creation service")

	campanha := entity.BuildCampanhaEntity(request)

	dbErr := srv.dbClient.CreateCampanha(campanha, userID)
	if dbErr != nil {
		zap.L().Error("Error creating campanha in database", zap.Error(dbErr))
		return 0, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Campanha created successfully", zap.String("nome", campanha.Nome), zap.Int("id", campanha.ID))
	return campanha.ID, nil
}

func (srv *Service) AssociarPublicosCampanhaService(userID string, idCampanha string, request dtos.AssociarPublicosCampanhaRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting associar publicos campanha service", zap.String("idCampanha", idCampanha))

	// Converter idCampanha string para int
	idCampanhaInt := 0
	if _, err := fmt.Sscanf(idCampanha, "%d", &idCampanhaInt); err != nil {
		zap.L().Error("Error converting idCampanha to int", zap.Error(err))
		return false, exceptions.NewBadRequestError("Invalid campanha ID")
	}

	dbErr := srv.dbClient.AssociarPublicosCampanha(idCampanhaInt, request.Publicos, userID)
	if dbErr != nil {
		zap.L().Error("Error associating publicos to campanha", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Publicos associated successfully", zap.String("idCampanha", idCampanha), zap.Ints("publicos", request.Publicos))
	return true, nil
}

func (srv *Service) GetPublicosCampanhaService(userID string, idCampanha string) (*dtos.PublicosCampanhaListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get publicos campanha service", zap.String("idCampanha", idCampanha))

	publicos, dbErr := srv.dbClient.GetPublicosCampanha(idCampanha, userID)
	if dbErr != nil {
		zap.L().Error("Error getting publicos for campanha", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error retrieving publicos for campanha")
	}

	// Converter idCampanha string para int para a resposta
	idCampanhaInt := 0
	if _, err := fmt.Sscanf(idCampanha, "%d", &idCampanhaInt); err != nil {
		zap.L().Error("Error converting idCampanha to int", zap.Error(err))
		return nil, exceptions.NewBadRequestError("Invalid campanha ID")
	}

	publicoResponses := make([]dtos.PublicoCampanhaResponse, len(publicos))
	for i, publico := range publicos {
		publicoResponses[i] = dtos.PublicoCampanhaResponse{
			ID:          publico.IDPublico,
			Nome:        publico.Nome,
			Descricao:   publico.Descricao,
			DataCriacao: publico.DataCriacao,
			Status:      publico.Status,
		}
	}

	response := &dtos.PublicosCampanhaListResponse{
		Publicos:   publicoResponses,
		Total:      len(publicos),
		IDCampanha: idCampanhaInt,
	}

	zap.L().Info("Successfully retrieved publicos for campanha", zap.String("idCampanha", idCampanha), zap.Int("count", len(publicos)))
	return response, nil
}
