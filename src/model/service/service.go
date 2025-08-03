package service

import (
	"context"
	"encoding/json"
	"fmt"
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

	GetAllPedidosService(userID string) (*dtos.PedidoListResponse, *exceptions.RestErr)
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

func (srv *Service) GetAllPedidosService(userID string) (*dtos.PedidoListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all pedidos service")

	pedidos, dbErr := srv.dbClient.GetAllPedidos(userID)
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

	response := &dtos.PedidoListResponse{
		Pedidos: pedidoResponses,
		Total:   len(pedidos),
	}

	zap.L().Info("Successfully retrieved all pedidos", zap.Int("count", len(pedidos)))
	return response, nil
}
