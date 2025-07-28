package service

import (
	"github.com/betine97/back-project.git/cmd/config/exceptions"
	"github.com/betine97/back-project.git/src/controller/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/betine97/back-project.git/src/model/persistence"
	"github.com/betine97/back-project.git/src/model/service/crypto"
	"go.uber.org/zap"
)

type ServiceInterface interface {
	CreateUserService(request dtos.CreateUser) (*entity.User, *exceptions.RestErr)
	LoginUserService(request dtos.UserLogin) (bool, *exceptions.RestErr)
	GetAllProductsService() (*dtos.ProductListResponse, *exceptions.RestErr)
	GetProductByIDService(id int) (*dtos.ProductResponse, *exceptions.RestErr)
	GetProductsWithFiltersService(params dtos.ProductQueryParams) (*dtos.ProductListResponse, *exceptions.RestErr)
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

func (srv *Service) LoginUserService(request dtos.UserLogin) (bool, *exceptions.RestErr) {
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

func buildUserEntity(request dtos.CreateUser, hashedPassword string) *entity.User {
	return &entity.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		City:      request.City,
		Password:  hashedPassword,
	}
}

func (srv *Service) CreateUserService(request dtos.CreateUser) (*entity.User, *exceptions.RestErr) {
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

func (srv *Service) GetAllProductsService() (*modelDtos.ProductListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all products service")

	products, err := srv.db.GetAllProducts()
	if err != nil {
		zap.L().Error("Error getting products from database", zap.Error(err))
		return nil, exceptions.NewInternalServerError("Error retrieving products")
	}

	productResponses := make([]modelDtos.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = modelDtos.ProductResponse{
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

	response := &modelDtos.ProductListResponse{
		Products: productResponses,
		Total:    len(productResponses),
		Page:     1,
		Limit:    len(productResponses),
	}

	zap.L().Info("Successfully retrieved all products", zap.Int("count", len(products)))
	return response, nil
}

func (srv *Service) GetProductByIDService(id int) (*modelDtos.ProductResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get product by ID service", zap.Int("id", id))

	product, err := srv.db.GetProductByID(id)
	if err != nil {
		zap.L().Error("Product not found", zap.Error(err), zap.Int("id", id))
		return nil, exceptions.NewNotFoundError("Product not found")
	}

	response := &modelDtos.ProductResponse{
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

	zap.L().Info("Successfully retrieved product by ID", zap.Int("id", id))
	return response, nil
}

func (srv *Service) GetProductsWithFiltersService(params modelDtos.ProductQueryParams) (*modelDtos.ProductListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get products with filters service", zap.Any("params", params))

	// Set default values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	// Convert filters to map
	filters := make(map[string]interface{})
	if params.Filter != nil {
		if params.Filter.Categoria != "" {
			filters["categoria"] = params.Filter.Categoria
		}
		if params.Filter.DestinadoPara != "" {
			filters["destinado_para"] = params.Filter.DestinadoPara
		}
		if params.Filter.Marca != "" {
			filters["marca"] = params.Filter.Marca
		}
		if params.Filter.Variacao != "" {
			filters["variacao"] = params.Filter.Variacao
		}
		if params.Filter.Status != "" {
			filters["status"] = params.Filter.Status
		}
		if params.Filter.MinPrice > 0 {
			filters["min_price"] = params.Filter.MinPrice
		}
		if params.Filter.MaxPrice > 0 {
			filters["max_price"] = params.Filter.MaxPrice
		}
		if params.Filter.Search != "" {
			filters["search"] = params.Filter.Search
		}
	}

	products, total, err := srv.db.GetProductsWithFilters(filters, params.Limit, offset)
	if err != nil {
		zap.L().Error("Error getting filtered products", zap.Error(err))
		return nil, exceptions.NewInternalServerError("Error retrieving filtered products")
	}

	productResponses := make([]modelDtos.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = modelDtos.ProductResponse{
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

	response := &modelDtos.ProductListResponse{
		Products: productResponses,
		Total:    total,
		Page:     params.Page,
		Limit:    params.Limit,
	}

	zap.L().Info("Successfully retrieved filtered products",
		zap.Int("count", len(products)),
		zap.Int("total", total),
		zap.Int("page", params.Page))

	return response, nil
}
