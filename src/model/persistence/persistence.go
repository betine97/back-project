package persistence

import (
	"context"

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
	GetAllFornecedores(ctx context.Context) ([]entity.Fornecedores, error)
	CreateFornecedor(fornecedor entity.Fornecedores, ctx context.Context) error
	GetFornecedorById(id string, ctx context.Context) (*entity.Fornecedores, error)
	UpdateFornecedor(fornecedor entity.Fornecedores, ctx context.Context) error
	UpdateFornecedorField(id string, campo string, valor string, ctx context.Context) error
	DeleteFornecedor(id string, ctx context.Context) error

	GetAllProducts(ctx context.Context) ([]entity.Produto, error)
	CreateProduct(product entity.Produto, ctx context.Context) error
	DeleteProduct(id string, ctx context.Context) error

	GetAllPedidos(ctx context.Context) ([]entity.Pedido, error)
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

// getClientDB é uma função auxiliar para obter a conexão do banco baseada no tenantID
func (repo *DBConnectionDBClient) getClientDB(ctx context.Context) *gorm.DB {
	tenantID := ctx.Value("tenantID").(string)
	clientKey := "db_" + tenantID

	zap.L().Info("Getting client DB", zap.String("tenantID", tenantID), zap.String("clientKey", clientKey))

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

func (repo *DBConnectionDBClient) GetAllFornecedores(ctx context.Context) ([]entity.Fornecedores, error) {

	db := repo.getClientDB(ctx)

	zap.L().Info("Getting all fornecedores from database")
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

func (repo *DBConnectionDBClient) CreateFornecedor(fornecedor entity.Fornecedores, ctx context.Context) error {
	db := repo.getClientDB(ctx)

	zap.L().Info("Creating fornecedor in the database", zap.String("fornecedor", fornecedor.Nome))
	err := db.Create(&fornecedor).Error
	if err != nil {
		zap.L().Error("Error creating fornecedor in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnectionDBClient) GetFornecedorById(id string, ctx context.Context) (*entity.Fornecedores, error) {

	db := repo.getClientDB(ctx)

	zap.L().Info("Getting fornecedor by id from database", zap.String("id", id))
	var fornecedor entity.Fornecedores
	err := db.Table("fornecedores").Where("id_fornecedor = ?", id).First(&fornecedor).Error
	if err != nil {
		zap.L().Error("Error getting fornecedor by id from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Fornecedor retrieved from database", zap.Any("fornecedor", fornecedor))
	return &fornecedor, nil
}

func (repo *DBConnectionDBClient) UpdateFornecedor(fornecedor entity.Fornecedores, ctx context.Context) error {

	db := repo.getClientDB(ctx)

	zap.L().Info("Updating fornecedor in the database", zap.String("fornecedor", fornecedor.Nome))
	err := db.Save(&fornecedor).Error
	if err != nil {
		zap.L().Error("Error updating fornecedor in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnectionDBClient) UpdateFornecedorField(id string, campo string, valor string, ctx context.Context) error {

	db := repo.getClientDB(ctx)

	zap.L().Info("Updating fornecedor field in the database", zap.String("id", id), zap.String("campo", campo), zap.String("valor", valor))

	// Usando GORM para atualizar o campo específico
	err := db.Model(&entity.Fornecedores{}).Where("id_fornecedor = ?", id).Update(campo, valor).Error
	if err != nil {
		zap.L().Error("Error updating fornecedor field in database", zap.Error(err))
		return err
	}
	return nil
}

func (repo *DBConnectionDBClient) DeleteFornecedor(id string, ctx context.Context) error {

	db := repo.getClientDB(ctx)

	zap.L().Info("Deleting fornecedor from database", zap.String("id", id))
	err := db.Delete(&entity.Fornecedores{}, id).Error
	if err != nil {
		zap.L().Error("Error deleting fornecedor from database", zap.Error(err))
	}
	return err
}

// FUNÇÕES DE PRODUTOS ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllProducts(ctx context.Context) ([]entity.Produto, error) {

	db := repo.getClientDB(ctx)

	zap.L().Info("Getting all products from database")
	var products []entity.Produto
	err := db.Find(&products).Error
	if err != nil {
		zap.L().Error("Error getting products from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Successfully retrieved products", zap.Int("count", len(products)))
	return products, nil
}

func (repo *DBConnectionDBClient) CreateProduct(product entity.Produto, ctx context.Context) error {

	db := repo.getClientDB(ctx)

	zap.L().Info("Creating product in the database", zap.String("product", product.NomeProduto))
	err := db.Create(&product).Error
	if err != nil {
		zap.L().Error("Error creating product in database", zap.Error(err))
	}
	return err
}

func (repo *DBConnectionDBClient) DeleteProduct(id string, ctx context.Context) error {

	db := repo.getClientDB(ctx)

	zap.L().Info("Deleting product from database", zap.String("id", id))
	err := db.Delete(&entity.Produto{}, id).Error
	if err != nil {
		zap.L().Error("Error deleting product from database", zap.Error(err))
	}
	return err
}

// FUNÇÕES DE PEDIDOS ------------------------------------------------------------------------------------------------------------------------------------

func (repo *DBConnectionDBClient) GetAllPedidos(ctx context.Context) ([]entity.Pedido, error) {

	db := repo.getClientDB(ctx)
	zap.L().Info("Getting all pedidos from database")
	var pedidos []entity.Pedido
	err := db.Table("pedidos").Find(&pedidos).Error // Especificar a tabela correta
	if err != nil {
		zap.L().Error("Error getting pedidos from database", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Successfully retrieved pedidos", zap.Int("count", len(pedidos)))
	return pedidos, nil
}
