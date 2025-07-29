package persistence

import (
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PersistenceInterface interface {
	CreateUser(user entity.User) error
	VerifyExist(email string) (bool, error)
	GetUser(email string) *entity.User

	GetAllFornecedores() ([]entity.Fornecedor, error)

	GetAllProducts() ([]entity.Produto, error)
	GetProductByID(id int) (*entity.Produto, error)

	GetAllPedidos() ([]entity.Pedido, error)
	GetPedidoByID(id int) (*entity.Pedido, error)

	GetAllItemPedidos() ([]entity.ItemPedido, error)
	GetItemPedidoByID(id int) (*entity.ItemPedido, error)
	CreateItemPedido(itemPedido entity.ItemPedido) error

	GetAllHisCmvPrcMarge() ([]entity.HisCmvPrcMarge, error)
}

type DBConnection struct {
	db *gorm.DB
}

func NewDBConnection(db *gorm.DB) PersistenceInterface {
	return &DBConnection{db: db}
}

// Funções de usuário

//---------------------------------

func (repo *DBConnection) CreateUser(user entity.User) error {
	zap.L().Info("Creating user in the database", zap.String("email", user.Email))
	err := repo.db.Create(&user).Error
	if err != nil {
		zap.L().Error("Error creating user in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnection) VerifyExist(email string) (bool, error) {
	zap.L().Info("Checking user existence", zap.String("email", email))
	var count int64
	err := repo.db.Table("users").Where("email = ?", email).Count(&count).Error
	if err != nil {
		zap.L().Error("Error checking user existence", zap.Error(err))
	}
	return count > 0, err
}

func (repo *DBConnection) GetUser(email string) *entity.User {
	zap.L().Info("Getting user from database", zap.String("email", email))
	var user entity.User
	err := repo.db.Table("users").Where("email = ?", email).First(&user).Error
	if err != nil {
		zap.L().Error("User not found in database", zap.Error(err))
	}
	return &user
}

// Funções de fornecedores

//---------------------------------

func (repo *DBConnection) GetAllFornecedores() ([]entity.Fornecedor, error) {
	zap.L().Info("Getting all fornecedores from database")
	var fornecedores []entity.Fornecedor
	err := repo.db.Find(&fornecedores).Error
	if err != nil {
		zap.L().Error("Error getting fornecedores from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Successfully retrieved fornecedores", zap.Int("count", len(fornecedores)))
	return fornecedores, nil
}

// Funções de produtos

//---------------------------------

func (repo *DBConnection) GetAllProducts() ([]entity.Produto, error) {
	zap.L().Info("Getting all products from database")
	var products []entity.Produto
	err := repo.db.Find(&products).Error
	if err != nil {
		zap.L().Error("Error getting products from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Successfully retrieved products", zap.Int("count", len(products)))
	return products, nil
}

func (repo *DBConnection) GetProductByID(id int) (*entity.Produto, error) {
	zap.L().Info("Getting product by ID", zap.Int("id", id))
	var product entity.Produto
	err := repo.db.First(&product, id).Error
	if err != nil {
		zap.L().Error("Product not found in database", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	return &product, nil
}

// Funções de pedidos

//---------------------------------

func (repo *DBConnection) GetAllPedidos() ([]entity.Pedido, error) {
	zap.L().Info("Getting all pedidos from database")
	var pedidos []entity.Pedido
	err := repo.db.Find(&pedidos).Error
	if err != nil {
		zap.L().Error("Error getting pedidos from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Successfully retrieved pedidos", zap.Int("count", len(pedidos)))
	return pedidos, nil
}

func (repo *DBConnection) GetPedidoByID(id int) (*entity.Pedido, error) {
	zap.L().Info("Getting pedido by ID", zap.Int("id", id))
	var pedido entity.Pedido
	err := repo.db.First(&pedido, id).Error
	if err != nil {
		zap.L().Error("Pedido not found in database", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	return &pedido, nil
}

// Funções de itens de pedidos

//---------------------------------

func (repo *DBConnection) GetAllItemPedidos() ([]entity.ItemPedido, error) {
	zap.L().Info("Getting all item pedidos from database")
	var itemPedidos []entity.ItemPedido
	err := repo.db.Find(&itemPedidos).Error
	if err != nil {
		zap.L().Error("Error getting item pedidos from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Successfully retrieved item pedidos", zap.Int("count", len(itemPedidos)))
	return itemPedidos, nil
}

func (repo *DBConnection) GetItemPedidoByID(id int) (*entity.ItemPedido, error) {
	zap.L().Info("Getting item pedido by ID", zap.Int("id", id))
	var itemPedido entity.ItemPedido
	err := repo.db.First(&itemPedido, id).Error
	if err != nil {
		zap.L().Error("Item pedido not found in database", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	return &itemPedido, nil
}

func (repo *DBConnection) CreateItemPedido(itemPedido entity.ItemPedido) error {
	zap.L().Info("Creating item pedido in the database", zap.Int("id_pedido", itemPedido.IDPedido), zap.Int("id_produto", itemPedido.IDProduto))

	err := repo.db.Create(&itemPedido).Error
	if err != nil {
		zap.L().Error("Error creating item pedido in database", zap.Error(err))
		return err
	}

	zap.L().Info("Item pedido created successfully", zap.Int("id_item", itemPedido.IDItem))
	return nil
}

// ... existing code ...

// Funções de histórico de cmv, preço e margem

//---------------------------------

func (repo *DBConnection) GetAllHisCmvPrcMarge() ([]entity.HisCmvPrcMarge, error) {
	zap.L().Info("Getting all his cmv prc marge from database")
	var hisCmvPrcMarge []entity.HisCmvPrcMarge
	err := repo.db.Find(&hisCmvPrcMarge).Error
	if err != nil {
		zap.L().Error("Error getting his cmv prc marge from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Successfully retrieved his cmv prc marge", zap.Int("count", len(hisCmvPrcMarge)))
	return hisCmvPrcMarge, nil
}
