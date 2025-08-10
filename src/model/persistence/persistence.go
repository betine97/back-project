package persistence

import (
	"fmt"
	"strconv"

	entity "github.com/betine97/back-project.git/src/model/entitys"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PersistenceInterfaceDBMaster interface {
	CreateUser(user entity.User) error
	VerifyExist(email string) (bool, error)
	GetUser(email string) *entity.User
	GetTenantByUserID(userID uint) *entity.Tenants
}

type PersistenceInterfaceDBClient interface {
	GetAllFornecedores(userID string) ([]entity.Fornecedores, error)
	GetAllFornecedoresPaginated(userID string, limit, offset int) ([]entity.Fornecedores, int, error)
	CreateFornecedor(fornecedor entity.Fornecedores, userID string) error
	GetFornecedorById(id string, userID string) (*entity.Fornecedores, error)
	UpdateFornecedor(fornecedor entity.Fornecedores, userID string) error
	UpdateFornecedorField(id string, campo string, valor string, userID string) error
	DeleteFornecedor(id string, userID string) error

	GetAllProducts(userID string) ([]entity.Produto, error)
	GetAllProductsPaginated(userID string, limit, offset int) ([]entity.Produto, int, error)
	GetProductByBarcode(barcode string, userID string) *entity.Produto
	CreateProduct(product entity.Produto, userID string) error
	DeleteProduct(id string, userID string) error

	GetAllPedidos(userID string) ([]entity.Pedido, error)
	GetAllPedidosPaginated(userID string, limit, offset int) ([]entity.Pedido, int, error)
	GetPedidoById(id string, userID string) (*entity.Pedido, error)
	CreatePedido(pedido *entity.Pedido, userID string) error

	// Itens de Pedido
	CreateItemPedido(item entity.ItemPedido, userID string) error

	// View Detalhes Pedido
	GetDetalhesPedido(idPedido string, userID string) ([]entity.ViewDetalhesPedido, error)
	GetDetalhesPedidoPaginated(idPedido string, userID string, limit, offset int) ([]entity.ViewDetalhesPedido, int, error)

	// Estoque
	GetAllEstoque(userID string) ([]entity.Estoque, error)
	GetAllEstoquePaginated(userID string, limit, offset int) ([]entity.Estoque, int, error)
	GetAllDetalhesEstoque(userID string) ([]entity.ViewDetalhesEstoque, error)
	GetAllDetalhesEstoquePaginated(userID string, limit, offset int) ([]entity.ViewDetalhesEstoque, int, error)
	CreateEstoque(estoque entity.Estoque, userID string) error

	// Clientes
	GetAllClientes(userID string) ([]entity.Cliente, error)
	GetAllClientesPaginated(userID string, limit, offset int) ([]entity.Cliente, int, error)
	BuscarClientesCriterios(userID string) ([]entity.Cliente, error)
	BuscarClientesPorCriterios(userID string, criterios []entity.PublicoCriterioJoin) ([]entity.Cliente, error)
	GetClienteByID(id string, userID string) *entity.Cliente
	GetClienteByEmail(email string, userID string) *entity.Cliente
	GetClienteByTelefone(telefone string, userID string) *entity.Cliente
	CreateCliente(cliente entity.Cliente, userID string) error
	DeleteCliente(id string, userID string) error

	// Tags de Clientes
	AtribuirTagsCliente(clienteID int, tagIDs []int, userID string) error
	RemoverTagsCliente(clienteID int, tagIDs []int, userID string) error
	GetTagsCliente(clienteID int, userID string) ([]entity.TagClienteJoin, error)

	// Tags
	GetAllTags(userID string) ([]entity.Tag, error)
	GetAllTagsPaginated(userID string, limit, offset int) ([]entity.Tag, int, error)
	CreateTag(tag entity.Tag, userID string) error

	// Endereços
	GetAllEnderecos(userID string) ([]entity.Endereco, error)
	GetAllEnderecosPaginated(userID string, limit, offset int) ([]entity.Endereco, int, error)
	CreateEndereco(endereco entity.Endereco, userID string) error
	DeleteEndereco(idEndereco string, userID string) error

	// Critérios
	GetAllCriterios(userID string) ([]entity.Criterio, error)

	// Públicos
	GetAllPublicos(userID string) ([]entity.PublicoCliente, error)
	GetAllPublicosPaginated(userID string, limit, offset int) ([]entity.PublicoCliente, int, error)
	CreatePublico(publico *entity.PublicoCliente, userID string) error
	AssociarCriteriosPublico(idPublico int, criterios []int, userID string) error
	GetCriteriosPublico(idPublico string, userID string) ([]entity.PublicoCriterioJoin, error)
	AdicionarClientesAoPublico(userID string, idPublico int, clientes []entity.Cliente) (int, int, error)
	GetClientesDoPublico(userID string, idPublico int, limit, offset int) ([]entity.Cliente, int, error)

	// Pets
	GetAllPets(userID string) ([]entity.Pet, error)
	GetAllPetsPaginated(userID string, limit, offset int) ([]entity.Pet, int, error)
	CreatePet(pet *entity.Pet, userID string) error

	// Completude
	GetClientesComPetsEnderecosParaCompletude(userID string, limit, offset int) ([]entity.Cliente, []entity.Pet, []entity.Endereco, int, error)

	// Campanhas
	GetAllCampanhas(userID string) ([]entity.Campanha, error)
	GetAllCampanhasPaginated(userID string, limit, offset int) ([]entity.Campanha, int, error)
	GetCampanhaByID(id string, userID string) *entity.Campanha
	CreateCampanha(campanha *entity.Campanha, userID string) error
	AssociarPublicosCampanha(idCampanha int, publicos []int, userID string) error
	GetPublicosCampanha(idCampanha string, userID string) ([]entity.CampanhaPublicoJoin, error)
}

type DBConnectionDBMaster struct {
	dbmaster *gorm.DB
}

type DBConnectionDBClient struct {
	dbclient map[string]*gorm.DB
}

func NewDBConnectionDBMaster(db *gorm.DB) PersistenceInterfaceDBMaster {
	return &DBConnectionDBMaster{dbmaster: db}
}

func NewDBConnectionDBClient(db map[string]*gorm.DB) PersistenceInterfaceDBClient {
	return &DBConnectionDBClient{dbclient: db}
}

// getClientDB é uma função auxiliar para obter a conexão do banco baseada no userID
func (repo *DBConnectionDBClient) getClientDB(userID string) *gorm.DB {
	clientKey := "db_" + userID

	zap.L().Info("Getting client DB", zap.String("userID", userID), zap.String("clientKey", clientKey))

	db := repo.dbclient[clientKey]
	if db == nil {
		zap.L().Error("Client DB not found", zap.String("clientKey", clientKey))
		// Log todas as chaves disponíveis
		for key := range repo.dbclient {
			zap.L().Info("Available DB key", zap.String("key", key))
		}
	}

	return db
}

// FUNÇÕES DE USUÁRIO DBMASTER ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBMaster) CreateUser(user entity.User) error {
	zap.L().Info("Creating user in the database", zap.String("email", user.Email))
	err := repo.dbmaster.Create(&user).Error
	if err != nil {
		zap.L().Error("Error creating user in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnectionDBMaster) VerifyExist(email string) (bool, error) {
	zap.L().Info("Checking user existence", zap.String("email", email))
	var count int64
	err := repo.dbmaster.Table("users").Where("email = ?", email).Count(&count).Error
	if err != nil {
		zap.L().Error("Error checking user existence", zap.Error(err))
	}
	return count > 0, err
}

func (repo *DBConnectionDBMaster) GetUser(email string) *entity.User {
	zap.L().Info("Getting user from database", zap.String("email", email))
	var user entity.User
	err := repo.dbmaster.Table("users").Where("email = ?", email).First(&user).Error
	if err != nil {
		zap.L().Error("User not found in database", zap.Error(err))
	}
	return &user
}

func (repo *DBConnectionDBMaster) GetTenantByUserID(userID uint) *entity.Tenants {
	zap.L().Info("Getting tenant by user id from database", zap.Uint("user_id", userID))
	var tenant entity.Tenants
	err := repo.dbmaster.Table("tenants").Where("user_id = ?", userID).First(&tenant).Error
	if err != nil {
		zap.L().Error("Tenant not found in database", zap.Uint("user_id", userID), zap.Error(err))
	} else {
		zap.L().Info("Tenant found in database", zap.Uint("tenant_id", tenant.ID), zap.Uint("user_id", tenant.UserID), zap.String("nome_empresa", tenant.NomeEmpresa))
	}
	return &tenant
}

// FUNÇÕES DE FORNECEDORES ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllFornecedores(userID string) ([]entity.Fornecedores, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all fornecedores from database", zap.String("userID", userID))
	var fornecedores []entity.Fornecedores
	err := db.Find(&fornecedores).Error
	if err != nil {
		zap.L().Error("Error getting fornecedores from database", zap.Error(err))
		return nil, err
	}

	// Adicionando log para visualizar os dados retornados do banco de dados
	zap.L().Info("Fornecedores retrieved from database", zap.Any("fornecedores", fornecedores))
	zap.L().Info("Successfully retrieved fornecedores", zap.Int("count", len(fornecedores)))
	return fornecedores, nil
}

func (repo *DBConnectionDBClient) GetAllFornecedoresPaginated(userID string, limit, offset int) ([]entity.Fornecedores, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated fornecedores from database", zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var fornecedores []entity.Fornecedores
	var total int64

	// Contar total de registros
	err := db.Model(&entity.Fornecedores{}).Count(&total).Error
	if err != nil {
		zap.L().Error("Error counting fornecedores", zap.Error(err))
		return nil, 0, err
	}

	// Buscar registros paginados
	err = db.Limit(limit).Offset(offset).Find(&fornecedores).Error
	if err != nil {
		zap.L().Error("Error getting paginated fornecedores from database", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated fornecedores", zap.Int("count", len(fornecedores)), zap.Int64("total", total))
	return fornecedores, int(total), nil
}

func (repo *DBConnectionDBClient) CreateFornecedor(fornecedor entity.Fornecedores, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Creating fornecedor in the database", zap.String("fornecedor", fornecedor.Nome), zap.String("userID", userID))
	err := db.Create(&fornecedor).Error
	if err != nil {
		zap.L().Error("Error creating fornecedor in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnectionDBClient) GetFornecedorById(id string, userID string) (*entity.Fornecedores, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting fornecedor by id from database", zap.String("id", id), zap.String("userID", userID))
	var fornecedor entity.Fornecedores
	err := db.Table("fornecedores").Where("id_fornecedor = ?", id).First(&fornecedor).Error
	if err != nil {
		zap.L().Error("Error getting fornecedor by id from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Fornecedor retrieved from database", zap.Any("fornecedor", fornecedor))
	return &fornecedor, nil
}

func (repo *DBConnectionDBClient) UpdateFornecedor(fornecedor entity.Fornecedores, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Updating fornecedor in the database", zap.String("fornecedor", fornecedor.Nome), zap.String("userID", userID))
	err := db.Save(&fornecedor).Error
	if err != nil {
		zap.L().Error("Error updating fornecedor in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnectionDBClient) UpdateFornecedorField(id string, campo string, valor string, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Updating fornecedor field in the database", zap.String("id", id), zap.String("campo", campo), zap.String("valor", valor), zap.String("userID", userID))

	// Usando GORM para atualizar o campo específico
	err := db.Model(&entity.Fornecedores{}).Where("id_fornecedor = ?", id).Update(campo, valor).Error
	if err != nil {
		zap.L().Error("Error updating fornecedor field in database", zap.Error(err))
		return err
	}
	return nil
}

func (repo *DBConnectionDBClient) DeleteFornecedor(id string, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Deleting fornecedor from database", zap.String("id", id), zap.String("userID", userID))
	err := db.Delete(&entity.Fornecedores{}, id).Error
	if err != nil {
		zap.L().Error("Error deleting fornecedor from database", zap.Error(err))
	}
	return err
}

// FUNÇÕES DE PRODUTOS ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllProducts(userID string) ([]entity.Produto, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all products from database", zap.String("userID", userID))
	var products []entity.Produto
	err := db.Find(&products).Error
	if err != nil {
		zap.L().Error("Error getting products from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Successfully retrieved products", zap.Int("count", len(products)))
	return products, nil
}

func (repo *DBConnectionDBClient) GetAllProductsPaginated(userID string, limit, offset int) ([]entity.Produto, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated products from database", zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var products []entity.Produto
	var total int64

	// Contar total de registros
	err := db.Model(&entity.Produto{}).Count(&total).Error
	if err != nil {
		zap.L().Error("Error counting products", zap.Error(err))
		return nil, 0, err
	}

	// Buscar registros paginados
	err = db.Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		zap.L().Error("Error getting paginated products from database", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated products", zap.Int("count", len(products)), zap.Int64("total", total))
	return products, int(total), nil
}

func (repo *DBConnectionDBClient) GetProductByBarcode(barcode string, userID string) *entity.Produto {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting product by barcode from database", zap.String("barcode", barcode), zap.String("userID", userID))
	var product entity.Produto
	err := db.Where("codigo_barra = ?", barcode).First(&product).Error
	if err != nil {
		zap.L().Error("Product not found by barcode", zap.Error(err))
	}
	return &product
}

func (repo *DBConnectionDBClient) CreateProduct(product entity.Produto, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Creating product in the database", zap.String("product", product.NomeProduto), zap.String("userID", userID))
	err := db.Create(&product).Error
	if err != nil {
		zap.L().Error("Error creating product in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnectionDBClient) DeleteProduct(id string, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Deleting product from database", zap.String("id", id), zap.String("userID", userID))
	err := db.Delete(&entity.Produto{}, id).Error
	if err != nil {
		zap.L().Error("Error deleting product from database", zap.Error(err))
	}
	return err
}

// FUNÇÕES DE PEDIDOS ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllPedidos(userID string) ([]entity.Pedido, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all pedidos from database", zap.String("userID", userID))
	var pedidos []entity.Pedido
	err := db.Table("pedidos").Find(&pedidos).Error
	if err != nil {
		zap.L().Error("Error getting pedidos from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Successfully retrieved pedidos", zap.Int("count", len(pedidos)))
	return pedidos, nil
}

func (repo *DBConnectionDBClient) GetAllPedidosPaginated(userID string, limit, offset int) ([]entity.Pedido, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated pedidos from database", zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var pedidos []entity.Pedido
	var total int64

	// Contar total de registros
	if err := db.Table("pedidos").Count(&total).Error; err != nil {
		zap.L().Error("Error counting pedidos", zap.Error(err))
		return nil, 0, err
	}

	// Buscar pedidos com paginação
	err := db.Table("pedidos").Limit(limit).Offset(offset).Find(&pedidos).Error
	if err != nil {
		zap.L().Error("Error getting paginated pedidos from database", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated pedidos", zap.Int("count", len(pedidos)), zap.Int64("total", total))
	return pedidos, int(total), nil
}

func (repo *DBConnectionDBClient) GetPedidoById(id string, userID string) (*entity.Pedido, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting pedido by ID from database", zap.String("id", id), zap.String("userID", userID))

	var pedido entity.Pedido
	err := db.Table("pedidos").Where("id_pedido = ?", id).First(&pedido).Error
	if err != nil {
		zap.L().Error("Error getting pedido by ID from database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved pedido by ID", zap.String("id", id))
	return &pedido, nil
}

func (repo *DBConnectionDBClient) CreatePedido(pedido *entity.Pedido, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Creating pedido in database", zap.String("userID", userID))

	err := db.Table("pedidos").Create(pedido).Error
	if err != nil {
		zap.L().Error("Error creating pedido in database", zap.Error(err))
		return err
	}

	zap.L().Info("Successfully created pedido", zap.Int("id", pedido.IDPedido))
	return nil
}

// FUNÇÕES DE ITENS DE PEDIDO ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) CreateItemPedido(item entity.ItemPedido, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Creating item pedido in database", zap.String("userID", userID), zap.Int("idPedido", item.IDPedido))

	err := db.Table("itens_pedido").Create(&item).Error
	if err != nil {
		zap.L().Error("Error creating item pedido in database", zap.Error(err))
		return err
	}

	zap.L().Info("Successfully created item pedido", zap.Int("idItem", item.IDItem), zap.Int("idPedido", item.IDPedido))
	return nil
}

// FUNÇÕES DA VIEW DETALHES PEDIDO ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetDetalhesPedido(idPedido string, userID string) ([]entity.ViewDetalhesPedido, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting detalhes pedido from view", zap.String("idPedido", idPedido), zap.String("userID", userID))

	var detalhes []entity.ViewDetalhesPedido
	err := db.Table("view_detalhes_pedido").Where("id_pedido = ?", idPedido).Find(&detalhes).Error
	if err != nil {
		zap.L().Error("Error getting detalhes pedido from view", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved detalhes pedido from view", zap.String("idPedido", idPedido), zap.Int("count", len(detalhes)))
	return detalhes, nil
}

func (repo *DBConnectionDBClient) GetDetalhesPedidoPaginated(idPedido string, userID string, limit, offset int) ([]entity.ViewDetalhesPedido, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated detalhes pedido from view", zap.String("idPedido", idPedido), zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var detalhes []entity.ViewDetalhesPedido
	var total int64

	// Contar total de registros para o pedido específico na view
	if err := db.Table("view_detalhes_pedido").Where("id_pedido = ?", idPedido).Count(&total).Error; err != nil {
		zap.L().Error("Error counting detalhes pedido from view", zap.Error(err))
		return nil, 0, err
	}

	// Buscar detalhes com paginação
	err := db.Table("view_detalhes_pedido").Where("id_pedido = ?", idPedido).Limit(limit).Offset(offset).Find(&detalhes).Error
	if err != nil {
		zap.L().Error("Error getting paginated detalhes pedido from view", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated detalhes pedido from view", zap.String("idPedido", idPedido), zap.Int("count", len(detalhes)), zap.Int64("total", total))
	return detalhes, int(total), nil
}

// Estoque methods
func (repo *DBConnectionDBClient) GetAllEstoque(userID string) ([]entity.Estoque, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all estoque from database", zap.String("userID", userID))
	var estoque []entity.Estoque
	err := db.Find(&estoque).Error
	if err != nil {
		zap.L().Error("Error getting estoque from database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved estoque", zap.Int("count", len(estoque)))
	return estoque, nil
}

func (repo *DBConnectionDBClient) GetAllEstoquePaginated(userID string, limit, offset int) ([]entity.Estoque, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated estoque from database", zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var estoque []entity.Estoque
	var total int64

	// Contar total de registros
	err := db.Model(&entity.Estoque{}).Count(&total).Error
	if err != nil {
		zap.L().Error("Error counting estoque", zap.Error(err))
		return nil, 0, err
	}

	// Buscar registros paginados
	err = db.Limit(limit).Offset(offset).Find(&estoque).Error
	if err != nil {
		zap.L().Error("Error getting paginated estoque from database", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated estoque", zap.Int("count", len(estoque)), zap.Int64("total", total))
	return estoque, int(total), nil
}

func (repo *DBConnectionDBClient) GetAllDetalhesEstoque(userID string) ([]entity.ViewDetalhesEstoque, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all detalhes estoque from view", zap.String("userID", userID))
	var detalhesEstoque []entity.ViewDetalhesEstoque
	err := db.Table("view_detalhes_estoque").Find(&detalhesEstoque).Error
	if err != nil {
		zap.L().Error("Error getting detalhes estoque from view", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved detalhes estoque from view", zap.Int("count", len(detalhesEstoque)))
	return detalhesEstoque, nil
}

func (repo *DBConnectionDBClient) GetAllDetalhesEstoquePaginated(userID string, limit, offset int) ([]entity.ViewDetalhesEstoque, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated detalhes estoque from view", zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var detalhesEstoque []entity.ViewDetalhesEstoque
	var total int64

	// Contar total de registros na view
	err := db.Table("view_detalhes_estoque").Count(&total).Error
	if err != nil {
		zap.L().Error("Error counting detalhes estoque from view", zap.Error(err))
		return nil, 0, err
	}

	// Buscar registros paginados da view
	err = db.Table("view_detalhes_estoque").Limit(limit).Offset(offset).Find(&detalhesEstoque).Error
	if err != nil {
		zap.L().Error("Error getting paginated detalhes estoque from view", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated detalhes estoque from view", zap.Int("count", len(detalhesEstoque)), zap.Int64("total", total))
	return detalhesEstoque, int(total), nil
}

func (repo *DBConnectionDBClient) CreateEstoque(estoque entity.Estoque, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Creating estoque in the database", zap.Int("id_produto", estoque.IDProduto), zap.String("userID", userID))
	err := db.Create(&estoque).Error
	if err != nil {
		zap.L().Error("Error creating estoque in database", zap.Error(err))
	}
	return err
}

// FUNÇÕES DE CLIENTES ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllClientes(userID string) ([]entity.Cliente, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all clientes from database", zap.String("userID", userID))
	var clientes []entity.Cliente
	err := db.Find(&clientes).Error
	if err != nil {
		zap.L().Error("Error getting clientes from database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved clientes", zap.Int("count", len(clientes)))
	return clientes, nil
}

func (repo *DBConnectionDBClient) GetAllClientesPaginated(userID string, limit, offset int) ([]entity.Cliente, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated clientes from database", zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var clientes []entity.Cliente
	var total int64

	// Contar total de registros
	err := db.Model(&entity.Cliente{}).Count(&total).Error
	if err != nil {
		zap.L().Error("Error counting clientes", zap.Error(err))
		return nil, 0, err
	}

	// Buscar registros paginados
	err = db.Limit(limit).Offset(offset).Find(&clientes).Error
	if err != nil {
		zap.L().Error("Error getting paginated clientes from database", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated clientes", zap.Int("count", len(clientes)), zap.Int64("total", total))
	return clientes, int(total), nil
}

func (repo *DBConnectionDBClient) BuscarClientesCriterios(userID string) ([]entity.Cliente, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting clientes for criterios from database", zap.String("userID", userID))
	var clientes []entity.Cliente
	err := db.Select("id, tipo_cliente, sexo").Find(&clientes).Error
	if err != nil {
		zap.L().Error("Error getting clientes for criterios from database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved clientes for criterios", zap.Int("count", len(clientes)))
	return clientes, nil
}

func (repo *DBConnectionDBClient) BuscarClientesPorCriterios(userID string, criterios []entity.PublicoCriterioJoin) ([]entity.Cliente, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting clientes by criterios from database", zap.String("userID", userID), zap.Int("criterios_count", len(criterios)))

	var clientes []entity.Cliente

	// Construir condições OR para os critérios
	var conditions []string
	var args []interface{}

	for _, criterio := range criterios {
		switch criterio.NomeCondicao {
		case "Pessoa Física":
			conditions = append(conditions, "clientes.tipo_cliente = ?")
			args = append(args, "PF")
		case "Pessoa Jurídica":
			conditions = append(conditions, "clientes.tipo_cliente = ?")
			args = append(args, "PJ")
		case "Sexo Masculino":
			conditions = append(conditions, "clientes.sexo = ?")
			args = append(args, "M")
		case "Sexo Feminino":
			conditions = append(conditions, "clientes.sexo = ?")
			args = append(args, "F")
		case "Possui Gato":
			conditions = append(conditions, "EXISTS (SELECT 1 FROM pets WHERE pets.cliente_id = clientes.id AND pets.especie = ?)")
			args = append(args, "Gato")
		case "Possui Cachorro":
			conditions = append(conditions, "EXISTS (SELECT 1 FROM pets WHERE pets.cliente_id = clientes.id AND pets.especie = ?)")
			args = append(args, "Cachorro")
		default:
			zap.L().Warn("Critério não reconhecido", zap.String("criterio", criterio.NomeCondicao))
		}
	}

	if len(conditions) == 0 {
		zap.L().Warn("No valid criterios found")
		return []entity.Cliente{}, nil
	}

	// Construir query com OR entre as condições
	whereClause := "(" + conditions[0]
	for i := 1; i < len(conditions); i++ {
		whereClause += " OR " + conditions[i]
	}
	whereClause += ")"

	query := db.Select("DISTINCT clientes.id, clientes.tipo_cliente, clientes.sexo").
		Table("clientes").
		Where(whereClause, args...)

	err := query.Find(&clientes).Error
	if err != nil {
		zap.L().Error("Error getting clientes by criterios from database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved clientes by criterios", zap.Int("count", len(clientes)))
	return clientes, nil
}

func (repo *DBConnectionDBClient) GetClienteByID(id string, userID string) *entity.Cliente {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting cliente by ID from database", zap.String("id", id), zap.String("userID", userID))
	var cliente entity.Cliente
	err := db.Where("id = ?", id).First(&cliente).Error
	if err != nil {
		zap.L().Error("Cliente not found by ID", zap.Error(err))
	}
	return &cliente
}

func (repo *DBConnectionDBClient) GetClienteByEmail(email string, userID string) *entity.Cliente {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting cliente by email from database", zap.String("email", email), zap.String("userID", userID))
	var cliente entity.Cliente
	err := db.Where("email = ?", email).First(&cliente).Error
	if err != nil {
		zap.L().Error("Cliente not found by email", zap.Error(err))
	}
	return &cliente
}

func (repo *DBConnectionDBClient) GetClienteByTelefone(telefone string, userID string) *entity.Cliente {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting cliente by telefone from database", zap.String("telefone", telefone), zap.String("userID", userID))
	var cliente entity.Cliente
	err := db.Where("numero_celular = ?", telefone).First(&cliente).Error
	if err != nil {
		zap.L().Error("Cliente not found by telefone", zap.Error(err))
	}
	return &cliente
}

func (repo *DBConnectionDBClient) CreateCliente(cliente entity.Cliente, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Creating cliente in the database", zap.String("nome_cliente", cliente.NomeCliente), zap.String("tipo_cliente", cliente.TipoCliente))
	err := db.Create(&cliente).Error
	if err != nil {
		zap.L().Error("Error creating cliente in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnectionDBClient) DeleteCliente(id string, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Deleting cliente from database", zap.String("id", id), zap.String("userID", userID))
	err := db.Delete(&entity.Cliente{}, id).Error
	if err != nil {
		zap.L().Error("Error deleting cliente from database", zap.Error(err))
	}
	return err
}

// FUNÇÕES DE TAGS DE CLIENTES ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) AtribuirTagsCliente(clienteID int, tagIDs []int, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Atribuindo tags ao cliente", zap.Int("clienteID", clienteID), zap.Ints("tagIDs", tagIDs), zap.String("userID", userID))

	// Inserir as associações em lote
	for _, tagID := range tagIDs {
		tagCliente := entity.TagCliente{
			IDTag:     tagID,
			ClienteID: clienteID,
		}

		// Usar ON DUPLICATE KEY IGNORE para evitar duplicatas
		err := db.Where("id_tag = ? AND cliente_id = ?", tagID, clienteID).FirstOrCreate(&tagCliente).Error
		if err != nil {
			zap.L().Error("Error atribuindo tag ao cliente", zap.Error(err), zap.Int("tagID", tagID), zap.Int("clienteID", clienteID))
			return err
		}
	}

	zap.L().Info("Tags atribuídas ao cliente com sucesso", zap.Int("clienteID", clienteID), zap.Ints("tagIDs", tagIDs))
	return nil
}

func (repo *DBConnectionDBClient) RemoverTagsCliente(clienteID int, tagIDs []int, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Removendo tags do cliente", zap.Int("clienteID", clienteID), zap.Ints("tagIDs", tagIDs), zap.String("userID", userID))

	// Remover as associações específicas
	err := db.Where("cliente_id = ? AND id_tag IN ?", clienteID, tagIDs).Delete(&entity.TagCliente{}).Error
	if err != nil {
		zap.L().Error("Error removendo tags do cliente", zap.Error(err))
		return err
	}

	zap.L().Info("Tags removidas do cliente com sucesso", zap.Int("clienteID", clienteID), zap.Ints("tagIDs", tagIDs))
	return nil
}

func (repo *DBConnectionDBClient) GetTagsCliente(clienteID int, userID string) ([]entity.TagClienteJoin, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting tags do cliente", zap.Int("clienteID", clienteID), zap.String("userID", userID))

	var tags []entity.TagClienteJoin
	err := db.Table("tags_clientes tc").
		Select("tc.id, tc.id_tag, tc.cliente_id, t.nome_tag as nome").
		Joins("INNER JOIN tags t ON tc.id_tag = t.id_tag").
		Where("tc.cliente_id = ?", clienteID).
		Find(&tags).Error

	if err != nil {
		zap.L().Error("Error getting tags do cliente", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved tags do cliente", zap.Int("clienteID", clienteID), zap.Int("count", len(tags)))
	return tags, nil
}

// FUNÇÕES DE TAGS ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllTags(userID string) ([]entity.Tag, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all tags from database", zap.String("userID", userID))
	var tags []entity.Tag
	err := db.Find(&tags).Error
	if err != nil {
		zap.L().Error("Error getting tags from database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved all tags", zap.Int("count", len(tags)))
	return tags, nil
}

func (repo *DBConnectionDBClient) GetAllTagsPaginated(userID string, limit, offset int) ([]entity.Tag, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated tags from database", zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var tags []entity.Tag
	var total int64

	// Contar total de registros
	if err := db.Model(&entity.Tag{}).Count(&total).Error; err != nil {
		zap.L().Error("Error counting tags", zap.Error(err))
		return nil, 0, err
	}

	// Buscar registros paginados
	err := db.Limit(limit).Offset(offset).Order("id_tag DESC").Find(&tags).Error
	if err != nil {
		zap.L().Error("Error getting paginated tags from database", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated tags", zap.Int("count", len(tags)), zap.Int64("total", total))
	return tags, int(total), nil
}

func (repo *DBConnectionDBClient) CreateTag(tag entity.Tag, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Creating tag in database", zap.String("nome", tag.NomeTag), zap.String("categoria", tag.CategoriaTag), zap.String("userID", userID))
	err := db.Create(&tag).Error
	if err != nil {
		zap.L().Error("Error creating tag in database", zap.Error(err))
	}
	return err
}

// FUNÇÕES DE ENDEREÇOS ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllEnderecos(userID string) ([]entity.Endereco, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all enderecos from database", zap.String("userID", userID))
	var enderecos []entity.Endereco
	err := db.Find(&enderecos).Error
	if err != nil {
		zap.L().Error("Error getting enderecos from database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved enderecos", zap.Int("count", len(enderecos)))
	return enderecos, nil
}

func (repo *DBConnectionDBClient) GetAllEnderecosPaginated(userID string, limit, offset int) ([]entity.Endereco, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated enderecos from database", zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var enderecos []entity.Endereco
	var total int64

	// Contar total de registros
	err := db.Model(&entity.Endereco{}).Count(&total).Error
	if err != nil {
		zap.L().Error("Error counting enderecos", zap.Error(err))
		return nil, 0, err
	}

	// Buscar registros paginados
	err = db.Limit(limit).Offset(offset).Find(&enderecos).Error
	if err != nil {
		zap.L().Error("Error getting paginated enderecos from database", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated enderecos", zap.Int("count", len(enderecos)), zap.Int64("total", total))
	return enderecos, int(total), nil
}

func (repo *DBConnectionDBClient) CreateEndereco(endereco entity.Endereco, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Creating endereco in the database", zap.Int("id_cliente", endereco.IDCliente), zap.String("cidade", endereco.Cidade), zap.String("userID", userID))
	err := db.Create(&endereco).Error
	if err != nil {
		zap.L().Error("Error creating endereco in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnectionDBClient) DeleteEndereco(idEndereco string, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Deleting endereco from database", zap.String("id_endereco", idEndereco), zap.String("userID", userID))
	err := db.Where("id_endereco = ?", idEndereco).Delete(&entity.Endereco{}).Error
	if err != nil {
		zap.L().Error("Error deleting endereco from database", zap.Error(err))
	}
	return err
}

// FUNÇÕES DE CRITÉRIOS ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllCriterios(userID string) ([]entity.Criterio, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all criterios from database", zap.String("userID", userID))
	var criterios []entity.Criterio
	err := db.Find(&criterios).Error
	if err != nil {
		zap.L().Error("Error getting criterios from database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved criterios", zap.Int("count", len(criterios)))
	return criterios, nil
}

// FUNÇÕES DE PÚBLICOS ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllPublicos(userID string) ([]entity.PublicoCliente, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all publicos from database", zap.String("userID", userID))
	var publicos []entity.PublicoCliente
	err := db.Find(&publicos).Error
	if err != nil {
		zap.L().Error("Error getting publicos from database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved publicos", zap.Int("count", len(publicos)))
	return publicos, nil
}

func (repo *DBConnectionDBClient) GetAllPublicosPaginated(userID string, limit, offset int) ([]entity.PublicoCliente, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated publicos from database", zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var publicos []entity.PublicoCliente
	var total int64

	// Contar total de registros
	err := db.Model(&entity.PublicoCliente{}).Count(&total).Error
	if err != nil {
		zap.L().Error("Error counting publicos", zap.Error(err))
		return nil, 0, err
	}

	// Buscar registros paginados
	err = db.Limit(limit).Offset(offset).Find(&publicos).Error
	if err != nil {
		zap.L().Error("Error getting paginated publicos from database", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated publicos", zap.Int("count", len(publicos)), zap.Int64("total", total))
	return publicos, int(total), nil
}

func (repo *DBConnectionDBClient) CreatePublico(publico *entity.PublicoCliente, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Creating publico in the database", zap.String("nome", publico.Nome), zap.String("userID", userID))
	err := db.Create(publico).Error
	if err != nil {
		zap.L().Error("Error creating publico in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnectionDBClient) AssociarCriteriosPublico(idPublico int, criterios []int, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Associating criterios to publico", zap.Int("idPublico", idPublico), zap.Ints("criterios", criterios), zap.String("userID", userID))

	// Criar as associações em lote
	var publicoCriterios []entity.PublicoCriterio
	for _, idCriterio := range criterios {
		publicoCriterios = append(publicoCriterios, entity.PublicoCriterio{
			IDPublico:  idPublico,
			IDCriterio: idCriterio,
		})
	}

	err := db.Create(&publicoCriterios).Error
	if err != nil {
		zap.L().Error("Error associating criterios to publico", zap.Error(err))
		return err
	}

	zap.L().Info("Successfully associated criterios to publico", zap.Int("idPublico", idPublico), zap.Int("count", len(criterios)))
	return nil
}

func (repo *DBConnectionDBClient) GetCriteriosPublico(idPublico string, userID string) ([]entity.PublicoCriterioJoin, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting criterios for publico", zap.String("idPublico", idPublico), zap.String("userID", userID))

	var criterios []entity.PublicoCriterioJoin
	err := db.Table("publicos_criterios pc").
		Select("pc.id_publico, pc.id_criterio, c.nome_condicao").
		Joins("JOIN criterios c ON pc.id_criterio = c.id").
		Where("pc.id_publico = ?", idPublico).
		Find(&criterios).Error

	if err != nil {
		zap.L().Error("Error getting criterios for publico", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved criterios for publico", zap.String("idPublico", idPublico), zap.Int("count", len(criterios)))
	return criterios, nil
}

func (repo *DBConnectionDBClient) AdicionarClientesAoPublico(userID string, idPublico int, clientes []entity.Cliente) (int, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Adding clientes to publico", zap.String("userID", userID), zap.Int("idPublico", idPublico), zap.Int("clientes_count", len(clientes)))

	clientesAdicionados := 0
	clientesJaExistiam := 0

	for _, cliente := range clientes {
		// Verificar se a associação já existe
		var count int64
		err := db.Table("addclientes_publicos").
			Where("id_publico = ? AND id_cliente = ?", idPublico, cliente.ID).
			Count(&count).Error

		if err != nil {
			zap.L().Error("Error checking existing association", zap.Error(err))
			return clientesAdicionados, clientesJaExistiam, err
		}

		if count > 0 {
			clientesJaExistiam++
			zap.L().Debug("Cliente already associated with publico", zap.Int("clienteID", cliente.ID), zap.Int("publicoID", idPublico))
			continue
		}

		// Criar nova associação
		addClientePublico := entity.AddClientePublico{
			IDPublico: idPublico,
			IDCliente: cliente.ID,
		}

		err = db.Create(&addClientePublico).Error
		if err != nil {
			zap.L().Error("Error creating cliente-publico association", zap.Error(err))
			return clientesAdicionados, clientesJaExistiam, err
		}

		clientesAdicionados++
		zap.L().Debug("Cliente added to publico", zap.Int("clienteID", cliente.ID), zap.Int("publicoID", idPublico))
	}

	zap.L().Info("Successfully processed clientes for publico",
		zap.Int("idPublico", idPublico),
		zap.Int("adicionados", clientesAdicionados),
		zap.Int("ja_existiam", clientesJaExistiam))

	return clientesAdicionados, clientesJaExistiam, nil
}

func (repo *DBConnectionDBClient) GetClientesDoPublico(userID string, idPublico int, limit, offset int) ([]entity.Cliente, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("🔍 DEBUG - Getting clientes do publico from database",
		zap.String("userID", userID),
		zap.Int("idPublico", idPublico),
		zap.Int("limit", limit),
		zap.Int("offset", offset))

	// Verificar se a conexão do banco existe
	if db == nil {
		zap.L().Error("❌ DEBUG - Database connection is nil for userID", zap.String("userID", userID))
		return nil, 0, fmt.Errorf("database connection not found for userID: %s", userID)
	}

	var clientes []entity.Cliente
	var total int64

	// Contar total de clientes no público
	zap.L().Info("🔍 DEBUG - Executando query de contagem...")
	err := db.Table("addclientes_publicos acp").
		Joins("INNER JOIN clientes c ON acp.id_cliente = c.id").
		Where("acp.id_publico = ?", idPublico).
		Count(&total).Error

	if err != nil {
		zap.L().Error("❌ DEBUG - Error counting clientes do publico",
			zap.Error(err),
			zap.String("userID", userID),
			zap.Int("idPublico", idPublico))
		return nil, 0, err
	}

	zap.L().Info("✅ DEBUG - Count query successful", zap.Int64("total", total))

	// Buscar clientes paginados
	zap.L().Info("🔍 DEBUG - Executando query de busca...")
	err = db.Table("clientes c").
		Select("c.*").
		Joins("INNER JOIN addclientes_publicos acp ON c.id = acp.id_cliente").
		Where("acp.id_publico = ?", idPublico).
		Limit(limit).
		Offset(offset).
		Order("c.nome_cliente ASC").
		Find(&clientes).Error

	if err != nil {
		zap.L().Error("❌ DEBUG - Error getting clientes do publico from database",
			zap.Error(err),
			zap.String("userID", userID),
			zap.Int("idPublico", idPublico))
		return nil, 0, err
	}

	zap.L().Info("✅ DEBUG - Successfully retrieved clientes do publico",
		zap.Int("count", len(clientes)),
		zap.Int64("total", total),
		zap.String("userID", userID))
	return clientes, int(total), nil
}

// FUNÇÕES DE PETS ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllPets(userID string) ([]entity.Pet, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all pets from database", zap.String("userID", userID))
	var pets []entity.Pet
	err := db.Find(&pets).Error
	if err != nil {
		zap.L().Error("Error getting pets from database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved pets", zap.Int("count", len(pets)))
	return pets, nil
}

func (repo *DBConnectionDBClient) GetAllPetsPaginated(userID string, limit, offset int) ([]entity.Pet, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated pets from database", zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var pets []entity.Pet
	var total int64

	// Contar total de registros
	err := db.Model(&entity.Pet{}).Count(&total).Error
	if err != nil {
		zap.L().Error("Error counting pets", zap.Error(err))
		return nil, 0, err
	}

	// Buscar registros paginados
	err = db.Limit(limit).Offset(offset).Find(&pets).Error
	if err != nil {
		zap.L().Error("Error getting paginated pets from database", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated pets", zap.Int("count", len(pets)), zap.Int64("total", total))
	return pets, int(total), nil
}

func (repo *DBConnectionDBClient) CreatePet(pet *entity.Pet, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Creating pet in the database", zap.String("nome_pet", pet.NomePet), zap.String("especie", pet.Especie), zap.String("userID", userID))
	err := db.Create(pet).Error
	if err != nil {
		zap.L().Error("Error creating pet in database", zap.Error(err))
	}
	return err
}

// GetClientesComPetsEnderecosParaCompletude busca clientes com seus pets e endereços para análise de completude
func (repo *DBConnectionDBClient) GetClientesComPetsEnderecosParaCompletude(userID string, limit, offset int) ([]entity.Cliente, []entity.Pet, []entity.Endereco, int, error) {
	db := repo.getClientDB(userID)

	var clientes []entity.Cliente
	var pets []entity.Pet
	var enderecos []entity.Endereco
	var total int64

	// Buscar clientes paginados
	err := db.Model(&entity.Cliente{}).Count(&total).Error
	if err != nil {
		zap.L().Error("Error counting clientes for completude", zap.Error(err))
		return nil, nil, nil, 0, err
	}

	err = db.Limit(limit).Offset(offset).Find(&clientes).Error
	if err != nil {
		zap.L().Error("Error fetching clientes for completude", zap.Error(err))
		return nil, nil, nil, 0, err
	}

	// Buscar todos os pets dos clientes encontrados usando uma abordagem mais segura
	if len(clientes) > 0 {
		var clienteIDs []int
		for _, cliente := range clientes {
			clienteIDs = append(clienteIDs, cliente.ID)
		}

		// Usar uma estrutura temporária para mapear os resultados de forma mais segura
		type PetTemp struct {
			IDPet           interface{} `gorm:"column:id_pet"`
			ClienteID       int         `gorm:"column:cliente_id"`
			NomePet         string      `gorm:"column:nome_pet"`
			Especie         string      `gorm:"column:especie"`
			Raca            string      `gorm:"column:raca"`
			Porte           string      `gorm:"column:porte"`
			DataAniversario interface{} `gorm:"column:data_aniversario"`
			Idade           interface{} `gorm:"column:idade"`
			DataRegistro    interface{} `gorm:"column:data_registro"`
		}

		var petsTemp []PetTemp
		err = db.Table("pets").Where("cliente_id IN ?", clienteIDs).Find(&petsTemp).Error
		if err != nil {
			zap.L().Error("Error fetching pets for completude", zap.Error(err))
			// Em caso de erro, continua sem os pets para não quebrar a funcionalidade
			zap.L().Info("Continuing without pets due to data inconsistency")
			return clientes, []entity.Pet{}, []entity.Endereco{}, int(total), nil
		}

		// Converter os resultados temporários para a estrutura Pet
		for i, petTemp := range petsTemp {
			pet := entity.Pet{
				ClienteID: petTemp.ClienteID,
				NomePet:   petTemp.NomePet,
				Especie:   petTemp.Especie,
				Raca:      petTemp.Raca,
				Porte:     petTemp.Porte,
			}

			// Converter IDPet de forma segura - usar índice se ID for inválido
			if idPet, ok := petTemp.IDPet.(int64); ok && idPet > 0 {
				pet.IDPet = int(idPet)
			} else if idPetStr, ok := petTemp.IDPet.(string); ok && idPetStr != "" {
				// Tentar converter string para int
				if idPetInt, convErr := strconv.Atoi(idPetStr); convErr == nil && idPetInt > 0 {
					pet.IDPet = idPetInt
				} else {
					// Usar um ID temporário baseado no índice para pets sem ID válido
					pet.IDPet = -(i + 1) // IDs negativos para pets sem ID válido
				}
			} else {
				// Usar um ID temporário baseado no índice para pets sem ID válido
				pet.IDPet = -(i + 1) // IDs negativos para pets sem ID válido
			}

			pets = append(pets, pet)
		}

		// Buscar endereços dos clientes
		err = db.Where("id_cliente IN ?", clienteIDs).Find(&enderecos).Error
		if err != nil {
			zap.L().Error("Error fetching enderecos for completude", zap.Error(err))
			// Em caso de erro, continua sem os endereços
			zap.L().Info("Continuing without enderecos due to error")
			enderecos = []entity.Endereco{}
		}
	}

	return clientes, pets, enderecos, int(total), nil
}

// FUNÇÕES DE CAMPANHAS ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllCampanhas(userID string) ([]entity.Campanha, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting all campanhas from database", zap.String("userID", userID))
	var campanhas []entity.Campanha
	err := db.Find(&campanhas).Error
	if err != nil {
		zap.L().Error("Error getting campanhas from database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved campanhas", zap.Int("count", len(campanhas)))
	return campanhas, nil
}

func (repo *DBConnectionDBClient) GetAllCampanhasPaginated(userID string, limit, offset int) ([]entity.Campanha, int, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting paginated campanhas from database", zap.String("userID", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	var campanhas []entity.Campanha
	var total int64

	// Contar total de registros
	err := db.Model(&entity.Campanha{}).Count(&total).Error
	if err != nil {
		zap.L().Error("Error counting campanhas", zap.Error(err))
		return nil, 0, err
	}

	// Buscar registros paginados
	err = db.Limit(limit).Offset(offset).Order("id DESC").Find(&campanhas).Error
	if err != nil {
		zap.L().Error("Error getting paginated campanhas from database", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved paginated campanhas", zap.Int("count", len(campanhas)), zap.Int64("total", total))
	return campanhas, int(total), nil
}

func (repo *DBConnectionDBClient) GetCampanhaByID(id string, userID string) *entity.Campanha {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting campanha by ID from database", zap.String("id", id), zap.String("userID", userID))
	var campanha entity.Campanha
	err := db.Where("id = ?", id).First(&campanha).Error
	if err != nil {
		zap.L().Error("Campanha not found by ID", zap.Error(err))
	}
	return &campanha
}

func (repo *DBConnectionDBClient) CreateCampanha(campanha *entity.Campanha, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Creating campanha in the database", zap.String("nome", campanha.Nome), zap.String("userID", userID))
	err := db.Create(campanha).Error
	if err != nil {
		zap.L().Error("Error creating campanha in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnectionDBClient) AssociarPublicosCampanha(idCampanha int, publicos []int, userID string) error {
	db := repo.getClientDB(userID)

	zap.L().Info("Associating publicos to campanha", zap.Int("idCampanha", idCampanha), zap.Ints("publicos", publicos), zap.String("userID", userID))

	// Criar as associações em lote
	var campanhaPublicos []entity.CampanhaPublico
	for _, idPublico := range publicos {
		campanhaPublicos = append(campanhaPublicos, entity.CampanhaPublico{
			IDCampanha: idCampanha,
			IDPublico:  idPublico,
		})
	}

	err := db.Create(&campanhaPublicos).Error
	if err != nil {
		zap.L().Error("Error associating publicos to campanha", zap.Error(err))
		return err
	}

	zap.L().Info("Successfully associated publicos to campanha", zap.Int("idCampanha", idCampanha), zap.Int("count", len(publicos)))
	return nil
}

func (repo *DBConnectionDBClient) GetPublicosCampanha(idCampanha string, userID string) ([]entity.CampanhaPublicoJoin, error) {
	db := repo.getClientDB(userID)

	zap.L().Info("Getting publicos for campanha", zap.String("idCampanha", idCampanha), zap.String("userID", userID))

	var publicos []entity.CampanhaPublicoJoin
	err := db.Table("campanhas_publicos cp").
		Select("cp.id_campanha, cp.id_publico, pc.nome, pc.descricao, pc.data_criacao, pc.status").
		Joins("INNER JOIN publicos_clientes pc ON cp.id_publico = pc.id").
		Where("cp.id_campanha = ?", idCampanha).
		Find(&publicos).Error

	if err != nil {
		zap.L().Error("Error getting publicos for campanha", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved publicos for campanha", zap.String("idCampanha", idCampanha), zap.Int("count", len(publicos)))
	return publicos, nil
}
