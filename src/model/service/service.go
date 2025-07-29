package service

import (
	"github.com/betine97/back-project.git/cmd/config/exceptions"
	dtos_controllers "github.com/betine97/back-project.git/src/controller/dtos_controllers"
	dtos_models "github.com/betine97/back-project.git/src/model/dtos_models"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/betine97/back-project.git/src/model/persistence"
	"github.com/betine97/back-project.git/src/model/service/crypto"
	"go.uber.org/zap"
)

type ServiceInterface interface {
	CreateUserService(request dtos_controllers.CreateUser) (*entity.User, *exceptions.RestErr)
	LoginUserService(request dtos_controllers.UserLogin) (bool, *exceptions.RestErr)

	GetAllProductsService() (*dtos_models.ProductListResponse, *exceptions.RestErr)

	GetAllPedidosService() (*dtos_models.PedidoListResponse, *exceptions.RestErr)
}

type Service struct {
	crypto crypto.CryptoInterface
	db     persistence.PersistenceInterface
}

func NewServiceInstance(crypto crypto.CryptoInterface, db persistence.PersistenceInterface) ServiceInterface {
	return &Service{
		crypto: crypto,
		db:     db,
	}
}

// FUNÇÕES DE USUÁRIO

//---------------------------------

func (srv *Service) LoginUserService(request dtos_controllers.UserLogin) (bool, *exceptions.RestErr) {
	zap.L().Info("Starting login service")

	user := srv.db.GetUser(request.Email)

	if user.Email == "" {
		zap.L().Warn("User not found", zap.String("email", request.Email))
		return false, exceptions.NewNotFoundError("Account not found")
	}

	_, err := srv.crypto.CheckPassword(request.Password, user.Password)
	if err != nil {
		zap.L().Warn("Incorrect password", zap.String("email", request.Email))
		return false, exceptions.NewUnauthorizedRequestError("The password entered is incorrect")
	}

	zap.L().Info("Successful login", zap.String("email", request.Email))
	return true, nil
}

func (srv *Service) CreateUserService(request dtos_controllers.CreateUser) (*entity.User, *exceptions.RestErr) {
	zap.L().Info("Starting user creation service")

	emailExists := srv.db.GetUser(request.Email)

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

	dbErr := srv.db.CreateUser(*user)
	if dbErr != nil {
		zap.L().Error("Error creating user in database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Internal server error")
	}

	zap.L().Info("User created successfully", zap.String("email", user.Email))
	return user, nil
}

func buildUserEntity(request dtos_controllers.CreateUser, hashedPassword string) *entity.User {
	return &entity.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		City:      request.City,
		Password:  hashedPassword,
	}
}

// FUNÇÕES DE PRODUTOS

//---------------------------------

func (srv *Service) GetAllProductsService() (*dtos_models.ProductListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all products service")

	products, err := srv.db.GetAllProducts()
	if err != nil {
		zap.L().Error("Error getting products from database", zap.Error(err))
		return nil, exceptions.NewInternalServerError("Error retrieving products")
	}

	productResponses := make([]dtos_models.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = dtos_models.ProductResponse{
			ID:            product.ID,
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

	response := &dtos_models.ProductListResponse{
		Products: productResponses,
		Total:    len(productResponses),
		Page:     1,
		Limit:    len(productResponses),
	}

	zap.L().Info("Successfully retrieved all products", zap.Int("count", len(products)))
	return response, nil
}

// FUNÇÕES DE PEDIDOS

//---------------------------------

func (srv *Service) GetAllPedidosService() (*dtos_models.PedidoListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all pedidos service")

	pedidos, err := srv.db.GetAllPedidos()
	if err != nil {
		zap.L().Error("Error getting pedidos from database", zap.Error(err))
		return nil, exceptions.NewInternalServerError("Error retrieving pedidos")
	}

	pedidoResponses := make([]dtos_models.PedidoResponse, len(pedidos))
	for i, pedido := range pedidos {
		pedidoResponses[i] = dtos_models.PedidoResponse{
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

	response := &dtos_models.PedidoListResponse{
		Pedidos: pedidoResponses,
		Total:   len(pedidos),
	}

	zap.L().Info("Successfully retrieved all pedidos", zap.Int("count", len(pedidos)))
	return response, nil
}
