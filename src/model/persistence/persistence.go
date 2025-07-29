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

	GetAllProducts() ([]entity.Produto, error)

	GetAllPedidos() ([]entity.Pedido, error)
}

type DBConnection struct {
	db *gorm.DB
}

func NewDBConnection(db *gorm.DB) PersistenceInterface {
	return &DBConnection{db: db}
}

// FUNÇÕES DE USUÁRIO

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

// FUNÇÕES DE PRODUTOS

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

// FUNÇÕES DE PEDIDOS

//---------------------------------

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
