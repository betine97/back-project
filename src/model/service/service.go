package service

import (
	"github.com/betine97/back-project.git/cmd/config/exceptions"
	"github.com/betine97/back-project.git/src/controller/dtos"
	modelDtos "github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/betine97/back-project.git/src/model/persistence"
	"github.com/betine97/back-project.git/src/model/service/crypto"
	"go.uber.org/zap"
)

type ServiceInterface interface {
	CreateUserService(request dtos.CreateUser) (*entity.User, *exceptions.RestErr)
	LoginUserService(request dtos.UserLogin) (bool, *exceptions.RestErr)

	GetAllFornecedoresService() (*modelDtos.FornecedorListResponse, *exceptions.RestErr)

	GetAllProductsService() (*modelDtos.ProductListResponse, *exceptions.RestErr)
	GetProductByIDService(id int) (*modelDtos.ProductResponse, *exceptions.RestErr)

	GetAllPedidosService() (*modelDtos.PedidoListResponse, *exceptions.RestErr)
	GetPedidoByIDService(id int) (*modelDtos.PedidoResponse, *exceptions.RestErr)

	GetAllItemPedidosService() (*modelDtos.ItemPedidoListResponse, *exceptions.RestErr)
	GetItemPedidoByIDService(id int) (*modelDtos.ItemPedidoResponse, *exceptions.RestErr)
	CreateItemPedidoService(request modelDtos.CreateItemPedidoRequest) (*modelDtos.ItemPedidoResponse, *exceptions.RestErr)

	GetAllHisCmvPrcMargeService() (*modelDtos.HisCmvPrcMargeListResponse, *exceptions.RestErr)
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

// Funções de usuário

//---------------------------------

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

func buildUserEntity(request dtos.CreateUser, hashedPassword string) *entity.User {
	return &entity.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		City:      request.City,
		Password:  hashedPassword,
	}
}

// Funções de fornecedores

//---------------------------------

func (srv *Service) GetAllFornecedoresService() (*modelDtos.FornecedorListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all fornecedores service")

	fornecedores, err := srv.db.GetAllFornecedores()
	if err != nil {
		zap.L().Error("Error getting fornecedores from database", zap.Error(err))
		return nil, exceptions.NewInternalServerError("Error retrieving fornecedores")
	}

	fornecedorResponses := make([]modelDtos.FornecedorResponse, len(fornecedores))
	for i, fornecedor := range fornecedores {
		fornecedorResponses[i] = modelDtos.FornecedorResponse{
			ID:           fornecedor.ID,
			Nome:         fornecedor.Nome,
			Telefone:     fornecedor.Telefone,
			Email:        fornecedor.Email,
			Endereco:     fornecedor.Endereco,
			Cidade:       fornecedor.Cidade,
			Estado:       fornecedor.Estado,
			CEP:          fornecedor.CEP,
			DataCadastro: fornecedor.DataCadastro,
			Status:       fornecedor.Status,
		}
	}

	response := &modelDtos.FornecedorListResponse{
		Fornecedores: fornecedorResponses,
		Total:        len(fornecedores),
	}

	zap.L().Info("Successfully retrieved all fornecedores", zap.Int("count", len(fornecedores)))
	return response, nil
}

// Funções de produtos

//---------------------------------

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

// Funções de pedidos

//---------------------------------

func (srv *Service) GetAllPedidosService() (*modelDtos.PedidoListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all pedidos service")

	pedidos, err := srv.db.GetAllPedidos()
	if err != nil {
		zap.L().Error("Error getting pedidos from database", zap.Error(err))
		return nil, exceptions.NewInternalServerError("Error retrieving pedidos")
	}

	pedidoResponses := make([]modelDtos.PedidoResponse, len(pedidos))
	for i, pedido := range pedidos {
		pedidoResponses[i] = modelDtos.PedidoResponse{
			ID:           pedido.ID,
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

	response := &modelDtos.PedidoListResponse{
		Pedidos: pedidoResponses,
		Total:   len(pedidos),
	}

	zap.L().Info("Successfully retrieved all pedidos", zap.Int("count", len(pedidos)))
	return response, nil
}

func (srv *Service) GetPedidoByIDService(id int) (*modelDtos.PedidoResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get pedido by ID service", zap.Int("id", id))

	pedido, err := srv.db.GetPedidoByID(id)
	if err != nil {
		zap.L().Error("Pedido not found", zap.Error(err), zap.Int("id", id))
		return nil, exceptions.NewNotFoundError("Pedido not found")
	}

	response := &modelDtos.PedidoResponse{
		ID:           pedido.ID,
		IDFornecedor: pedido.IDFornecedor,
		DataPedido:   pedido.DataPedido,
		DataEntrega:  pedido.DataEntrega,
		ValorFrete:   pedido.ValorFrete,
		CustoPedido:  pedido.CustoPedido,
		ValorTotal:   pedido.ValorTotal,
		Descricao:    pedido.Descricao,
		Status:       pedido.Status,
	}

	zap.L().Info("Successfully retrieved pedido by ID", zap.Int("id", id))
	return response, nil
}

// Funções de itens de pedidos

//---------------------------------

func (srv *Service) GetAllItemPedidosService() (*modelDtos.ItemPedidoListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all item pedidos service")

	itemPedidos, err := srv.db.GetAllItemPedidos()
	if err != nil {
		zap.L().Error("Error getting item pedidos from database", zap.Error(err))
		return nil, exceptions.NewInternalServerError("Error retrieving item pedidos")
	}

	itemPedidoResponses := make([]modelDtos.ItemPedidoResponse, len(itemPedidos))
	for i, itemPedido := range itemPedidos {
		itemPedidoResponses[i] = modelDtos.ItemPedidoResponse{
			IDItem:        itemPedido.IDItem,
			IDPedido:      itemPedido.IDPedido,
			IDProduto:     itemPedido.IDProduto,
			Quantidade:    itemPedido.Quantidade,
			PrecoUnitario: itemPedido.PrecoUnitario,
			Subtotal:      itemPedido.Subtotal,
		}
	}

	response := &modelDtos.ItemPedidoListResponse{
		ItemPedidos: itemPedidoResponses,
		Total:       len(itemPedidos),
	}

	zap.L().Info("Successfully retrieved all item pedidos", zap.Int("count", len(itemPedidos)))
	return response, nil
}

func (srv *Service) GetItemPedidoByIDService(id int) (*modelDtos.ItemPedidoResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get item pedido by ID service", zap.Int("id", id))

	itemPedido, err := srv.db.GetItemPedidoByID(id)
	if err != nil {
		zap.L().Error("Item pedido not found", zap.Error(err), zap.Int("id", id))
		return nil, exceptions.NewNotFoundError("Item pedido not found")
	}

	response := &modelDtos.ItemPedidoResponse{
		IDItem:        itemPedido.IDItem,
		IDPedido:      itemPedido.IDPedido,
		IDProduto:     itemPedido.IDProduto,
		Quantidade:    itemPedido.Quantidade,
		PrecoUnitario: itemPedido.PrecoUnitario,
		Subtotal:      itemPedido.Subtotal,
	}

	zap.L().Info("Successfully retrieved item pedido by ID", zap.Int("id", id))
	return response, nil
}

func (srv *Service) CreateItemPedidoService(request modelDtos.CreateItemPedidoRequest) (*modelDtos.ItemPedidoResponse, *exceptions.RestErr) {
	zap.L().Info("Starting create item pedido service", zap.Int("id_pedido", request.IDPedido), zap.Int("id_produto", request.IDProduto))

	// Validar se o pedido existe
	_, err := srv.db.GetPedidoByID(request.IDPedido)
	if err != nil {
		zap.L().Error("Pedido not found", zap.Error(err), zap.Int("id_pedido", request.IDPedido))
		return nil, exceptions.NewNotFoundError("Pedido not found")
	}

	// Validar se o produto existe
	_, err = srv.db.GetProductByID(request.IDProduto)
	if err != nil {
		zap.L().Error("Product not found", zap.Error(err), zap.Int("id_produto", request.IDProduto))
		return nil, exceptions.NewNotFoundError("Product not found")
	}

	// Calcular subtotal
	subtotal := float64(request.Quantidade) * request.PrecoUnitario

	// Criar entidade
	itemPedido := entity.ItemPedido{
		IDPedido:      request.IDPedido,
		IDProduto:     request.IDProduto,
		Quantidade:    request.Quantidade,
		PrecoUnitario: request.PrecoUnitario,
		Subtotal:      subtotal,
	}

	// Salvar no banco
	dbErr := srv.db.CreateItemPedido(itemPedido)
	if dbErr != nil {
		zap.L().Error("Error creating item pedido in database", zap.Error(dbErr))
		return nil, exceptions.NewInternalServerError("Error creating item pedido")
	}

	// Criar response
	response := &modelDtos.ItemPedidoResponse{
		IDItem:        itemPedido.IDItem,
		IDPedido:      itemPedido.IDPedido,
		IDProduto:     itemPedido.IDProduto,
		Quantidade:    itemPedido.Quantidade,
		PrecoUnitario: itemPedido.PrecoUnitario,
		Subtotal:      itemPedido.Subtotal,
	}

	zap.L().Info("Item pedido created successfully", zap.Int("id_item", itemPedido.IDItem))
	return response, nil
}

// Funções de histórico de cmv, preço e margem

//---------------------------------

func (srv *Service) GetAllHisCmvPrcMargeService() (*modelDtos.HisCmvPrcMargeListResponse, *exceptions.RestErr) {
	zap.L().Info("Starting get all his cmv prc marge service")

	hisCmvPrcMarge, err := srv.db.GetAllHisCmvPrcMarge()
	if err != nil {
		zap.L().Error("Error getting his cmv prc marge from database", zap.Error(err))
		return nil, exceptions.NewInternalServerError("Error retrieving his cmv prc marge")
	}

	hisCmvPrcMargeResponses := make([]modelDtos.HisCmvPrcMargeResponse, len(hisCmvPrcMarge))
	for i, hisCmvPrcMarge := range hisCmvPrcMarge {
		hisCmvPrcMargeResponses[i] = modelDtos.HisCmvPrcMargeResponse{
			ID:           hisCmvPrcMarge.ID,
			IDProduto:    hisCmvPrcMarge.IDProduto,
			PrecoVenda:   hisCmvPrcMarge.PrecoVenda,
			Cmv:          hisCmvPrcMarge.Cmv,
			Margem:       hisCmvPrcMarge.Margem,
			DataRegistro: hisCmvPrcMarge.DataRegistro,
		}
	}

	response := &modelDtos.HisCmvPrcMargeListResponse{
		HisCmvPrcMarge: hisCmvPrcMargeResponses,
		Total:          len(hisCmvPrcMarge),
	}

	zap.L().Info("Successfully retrieved all his cmv prc marge", zap.Int("count", len(hisCmvPrcMargeResponses)))
	return response, nil
}
