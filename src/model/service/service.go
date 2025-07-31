package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/betine97/back-project.git/cmd/config"
	"github.com/betine97/back-project.git/cmd/config/exceptions"
	dtos "github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/betine97/back-project.git/src/model/persistence"
	"github.com/betine97/back-project.git/src/model/service/crypto"
	redis "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type ServiceInterface interface {
	CreateUserService(request dtos.CreateUser) (*entity.User, *exceptions.RestErr)
	LoginUserService(request dtos.UserLogin) (string, *exceptions.RestErr)

	GetAllFornecedoresService() (*dtos.FornecedorListResponse, *exceptions.RestErr)
	CreateFornecedorService(request dtos.CreateFornecedorRequest) (bool, *exceptions.RestErr)
	ChangeStatusFornecedorService(id string) (bool, *exceptions.RestErr)
	UpdateFornecedorFieldService(id string, campo string, valor string) (bool, *exceptions.RestErr)
	DeleteFornecedorService(id string) (bool, *exceptions.RestErr)

	GetAllProductsService() (*dtos.ProductListResponse, *exceptions.RestErr)
	CreateProductService(request dtos.CreateProductRequest) (bool, *exceptions.RestErr)
	DeleteProductService(id string) (bool, *exceptions.RestErr)

	GetAllPedidosService() (*dtos.PedidoListResponse, *exceptions.RestErr)
}

var ctx = context.Background()

type Service struct {
	crypto   crypto.CryptoInterface
	dbmaster persistence.PersistenceInterfaceDBMaster
	dbClient persistence.PersistenceInterfaceDBClient
	redis    *redis.Client
}

func NewServiceInstance(crypto crypto.CryptoInterface, dbmaster persistence.PersistenceInterfaceDBMaster, dbClient persistence.PersistenceInterfaceDBClient, redisClient *redis.Client) ServiceInterface {
	return &Service{
		crypto:   crypto,
		dbmaster: dbmaster,
		dbClient: dbClient,
		redis:    redisClient,
	}
}

// FUNÇÕES DE USUÁRIO ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) LoginUserService(request dtos.UserLogin) (string, *exceptions.RestErr) {
	zap.L().Info("Starting login service")

	// Verifica se o usuário existe
	user := srv.dbmaster.GetUser(request.Email)
	if user.Email == "" {
		zap.L().Warn("User not found", zap.String("email", request.Email))
		return "", exceptions.NewNotFoundError("Account not found")
	}

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
	tenantJSON, err := srv.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// Dados não estão em cache, consulte a tabela tenants
		tenant := srv.dbmaster.GetTenantByUserID(user.ID)
		if tenant == nil {
			return "", exceptions.NewInternalServerError("Tenant not found")
		}

		// Armazene os dados no Redis
		tenantJSON, err := json.Marshal(tenant)
		if err != nil {
			zap.L().Error("Error marshaling tenant to JSON", zap.Error(err))
			return "", exceptions.NewInternalServerError("Error storing tenant in cache")
		}

		if err := srv.redis.Set(ctx, cacheKey, tenantJSON, 0).Err(); err != nil {
			zap.L().Error("Error storing tenant in Redis", zap.Error(err))
		}
	} else if err != nil {
		zap.L().Error("Error retrieving from Redis", zap.Error(err))
		return "", exceptions.NewInternalServerError("Error retrieving from cache")
	}

	// Deserializa o JSON de volta para a estrutura Tenant
	var cachedTenant entity.Tenants
	if err := json.Unmarshal([]byte(tenantJSON), &cachedTenant); err != nil {
		zap.L().Error("Error unmarshaling tenant from JSON", zap.Error(err))
		return "", exceptions.NewInternalServerError("Error retrieving tenant data")
	}

	// Gera o token com o ID do usuário
	token, erro := GenerateToken(cachedTenant.UserID)
	if erro != nil {
		zap.L().Error("Error signing token", zap.Error(err))
		return "", exceptions.NewInternalServerError("Internal server error")
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

func GenerateToken(userID uint) (string, *exceptions.RestErr) {

	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(config.PrivateKey)
	if err != nil {
		zap.L().Error("Error signing token", zap.Error(err))
		return "", exceptions.NewInternalServerError("Internal server error")
	}

	return tokenString, nil
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

func (srv *Service) GetAllFornecedoresService() (*dtos.FornecedorListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all fornecedores service")

	fornecedores, err := srv.dbClient.GetAllFornecedores(ctx)
	if err != nil {
		zap.L().Error("Error getting fornecedores from database", zap.Error(err))
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
			DataCadastro: fornecedor.DataCadastro, // Certifique-se de incluir este campo se necessário
		}
	}

	response := &dtos.FornecedorListResponse{
		Fornecedores: fornecedorResponses,
		Total:        len(fornecedores),
	}

	zap.L().Info("Successfully retrieved all fornecedores", zap.Int("count", len(fornecedores)))
	return response, nil
}

func (srv *Service) CreateFornecedorService(request dtos.CreateFornecedorRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting fornecedor creation service")

	fornecedor := entity.BuildFornecedorEntity(request)

	dbErr := srv.dbClient.CreateFornecedor(*fornecedor, ctx)
	if dbErr != nil {
		zap.L().Error("Error creating fornecedor in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Fornecedor created successfully", zap.String("fornecedor", fornecedor.Nome))
	return true, nil
}

func (srv *Service) ChangeStatusFornecedorService(id string) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting change status fornecedor service")

	// Recupera o fornecedor pelo ID
	fornecedor, err := srv.dbClient.GetFornecedorById(id, ctx)
	if err != nil {
		zap.L().Error("Error getting fornecedor from database", zap.Error(err))
		return false, exceptions.NewInternalServerError("Error retrieving fornecedor")
	}

	// Alterna o status
	if fornecedor.Status == "Ativo" {
		fornecedor.Status = "Inativo"
	} else {
		fornecedor.Status = "Ativo"
	}

	// Atualiza o fornecedor no banco de dados
	dbErr := srv.dbClient.UpdateFornecedor(*fornecedor, ctx)
	if dbErr != nil {
		zap.L().Error("Error updating fornecedor in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Status fornecedor changed successfully", zap.String("fornecedor", fornecedor.Nome))
	return true, nil
}

func (srv *Service) UpdateFornecedorFieldService(id string, campo string, valor string) (bool, *exceptions.RestErr) {
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
	dbErr := srv.dbClient.UpdateFornecedorField(id, campo, valor, ctx)
	if dbErr != nil {
		zap.L().Error("Error updating fornecedor field in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Fornecedor field updated successfully", zap.String("id", id), zap.String("campo", campo), zap.String("valor", valor))
	return true, nil
}

func (srv *Service) DeleteFornecedorService(id string) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting delete fornecedor service")

	dbErr := srv.dbClient.DeleteFornecedor(id, ctx)
	if dbErr != nil {
		zap.L().Error("Error deleting fornecedor in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Fornecedor deleted successfully", zap.String("fornecedor", id))
	return true, nil
}

// FUNÇÕES DE PRODUTOS ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllProductsService() (*dtos.ProductListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all products service")

	products, err := srv.dbClient.GetAllProducts(ctx)
	if err != nil {
		zap.L().Error("Error getting products from database", zap.Error(err))
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
		}
	}

	response := &dtos.ProductListResponse{
		Products: productResponses,
		Total:    len(productResponses),
		Page:     1,
		Limit:    len(productResponses),
	}

	zap.L().Info("Successfully retrieved all products", zap.Int("count", len(products)))
	return response, nil
}

func (srv *Service) CreateProductService(request dtos.CreateProductRequest) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting product creation service")

	product := entity.BuildProductEntity(request)

	dbErr := srv.dbClient.CreateProduct(*product, ctx)
	if dbErr != nil {
		zap.L().Error("Error creating product in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Product created successfully", zap.String("product", product.NomeProduto))
	return true, nil
}

func (srv *Service) DeleteProductService(id string) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting delete product service")

	dbErr := srv.dbClient.DeleteProduct(id, ctx)
	if dbErr != nil {
		zap.L().Error("Error deleting product in database", zap.Error(dbErr))
		return false, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("Product deleted successfully", zap.String("product", id))
	return true, nil
}

// FUNÇÕES DE PEDIDOS ------------------------------------------------------------------------------------------------------------------------------------

func (srv *Service) GetAllPedidosService() (*dtos.PedidoListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all pedidos service")

	pedidos, err := srv.dbClient.GetAllPedidos(ctx)
	if err != nil {
		zap.L().Error("Error getting pedidos from database", zap.Error(err))
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
