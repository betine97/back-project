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
	GetTenantByUserID(userID uint) *entity.Tenants

	GetAllFornecedores() ([]entity.Fornecedores, error)
	CreateFornecedor(fornecedor entity.Fornecedores) error
	GetFornecedorById(id string) (*entity.Fornecedores, error)
	UpdateFornecedor(fornecedor entity.Fornecedores) error
	UpdateFornecedorField(id string, campo string, valor string) error
	DeleteFornecedor(id string) error

	GetAllProducts() ([]entity.Produto, error)
	CreateProduct(product entity.Produto) error
	DeleteProduct(id string) error

	GetAllPedidos() ([]entity.Pedido, error)
}

type DBConnection struct {
	db *gorm.DB
}

func NewDBConnection(db *gorm.DB) PersistenceInterface {
	return &DBConnection{db: db}
}

// FUNÇÕES DE USUÁRIO ------------------------------------------------------------------------------------------------------------------------------------

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

func (repo *DBConnection) GetTenantByUserID(userID uint) *entity.Tenants {
	zap.L().Info("Getting tenant by user id from database", zap.Uint("user_id", userID))
	var tenant entity.Tenants
	err := repo.db.Table("tenants").Where("user_id = ?", userID).First(&tenant).Error
	if err != nil {
		zap.L().Error("Tenant not found in database", zap.Error(err))
	}
	return &tenant
}

// FUNÇÕES DE FORNECEDORES ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnection) GetAllFornecedores() ([]entity.Fornecedores, error) {
	zap.L().Info("Getting all fornecedores from database")
	var fornecedores []entity.Fornecedores
	err := repo.db.Find(&fornecedores).Error
	if err != nil {
		zap.L().Error("Error getting fornecedores from database", zap.Error(err))
		return nil, err
	}

	// Adicionando log para visualizar os dados retornados do banco de dados
	zap.L().Info("Fornecedores retrieved from database", zap.Any("fornecedores", fornecedores))

	zap.L().Info("Successfully retrieved fornecedores", zap.Int("count", len(fornecedores)))
	return fornecedores, nil
}

func (repo *DBConnection) CreateFornecedor(fornecedor entity.Fornecedores) error {
	zap.L().Info("Creating fornecedor in the database", zap.String("fornecedor", fornecedor.Nome))
	err := repo.db.Create(&fornecedor).Error
	if err != nil {
		zap.L().Error("Error creating fornecedor in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnection) GetFornecedorById(id string) (*entity.Fornecedores, error) {
	zap.L().Info("Getting fornecedor by id from database", zap.String("id", id))
	var fornecedor entity.Fornecedores
	err := repo.db.Table("fornecedores").Where("id_fornecedor = ?", id).First(&fornecedor).Error
	if err != nil {
		zap.L().Error("Error getting fornecedor by id from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Fornecedor retrieved from database", zap.Any("fornecedor", fornecedor))
	return &fornecedor, nil
}

func (repo *DBConnection) UpdateFornecedor(fornecedor entity.Fornecedores) error {
	zap.L().Info("Updating fornecedor in the database", zap.String("fornecedor", fornecedor.Nome))
	err := repo.db.Save(&fornecedor).Error
	if err != nil {
		zap.L().Error("Error updating fornecedor in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnection) UpdateFornecedorField(id string, campo string, valor string) error {
	zap.L().Info("Updating fornecedor field in the database", zap.String("id", id), zap.String("campo", campo), zap.String("valor", valor))

	// Usando GORM para atualizar o campo específico
	err := repo.db.Model(&entity.Fornecedores{}).Where("id_fornecedor = ?", id).Update(campo, valor).Error
	if err != nil {
		zap.L().Error("Error updating fornecedor field in database", zap.Error(err))
		return err
	}
	return nil
}

func (repo *DBConnection) DeleteFornecedor(id string) error {
	zap.L().Info("Deleting fornecedor from database", zap.String("id", id))
	err := repo.db.Delete(&entity.Fornecedores{}, id).Error
	if err != nil {
		zap.L().Error("Error deleting fornecedor from database", zap.Error(err))
	}
	return err
}

// FUNÇÕES DE PRODUTOS ------------------------------------------------------------------------------------------------------------------------------------

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

func (repo *DBConnection) CreateProduct(product entity.Produto) error {
	zap.L().Info("Creating product in the database", zap.String("product", product.NomeProduto))
	err := repo.db.Create(&product).Error
	if err != nil {
		zap.L().Error("Error creating product in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnection) DeleteProduct(id string) error {
	zap.L().Info("Deleting product from database", zap.String("id", id))
	err := repo.db.Delete(&entity.Produto{}, id).Error
	if err != nil {
		zap.L().Error("Error deleting product from database", zap.Error(err))
	}
	return err
}

// FUNÇÕES DE PEDIDOS ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnection) GetAllPedidos() ([]entity.Pedido, error) {
	zap.L().Info("Getting all pedidos from database")
	var pedidos []entity.Pedido
	err := repo.db.Table("pedidos").Find(&pedidos).Error // Especificar a tabela correta
	if err != nil {
		zap.L().Error("Error getting pedidos from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Successfully retrieved pedidos", zap.Int("count", len(pedidos)))
	return pedidos, nil
}
