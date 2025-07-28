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
	GetProductByID(id int) (*entity.Produto, error)
	GetProductsWithFilters(filters map[string]interface{}, limit, offset int) ([]entity.Produto, int, error)
}

type DBConnection struct {
	db *gorm.DB
}

func NewDBConnection(db *gorm.DB) PersistenceInterface {
	return &DBConnection{db: db}
}

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

func (repo *DBConnection) GetProductsWithFilters(filters map[string]interface{}, limit, offset int) ([]entity.Produto, int, error) {
	zap.L().Info("Getting products with filters", zap.Any("filters", filters), zap.Int("limit", limit), zap.Int("offset", offset))

	var products []entity.Produto
	var total int64

	query := repo.db.Model(&entity.Produto{})

	// Apply filters
	for key, value := range filters {
		switch key {
		case "categoria":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("categoria = ?", v)
			}
		case "destinado_para":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("destinado_para = ?", v)
			}
		case "marca":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("marca = ?", v)
			}
		case "variacao":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("variacao = ?", v)
			}
		case "status":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("status = ?", v)
			}
		case "min_price":
			if v, ok := value.(float64); ok && v > 0 {
				query = query.Where("preco_venda >= ?", v)
			}
		case "max_price":
			if v, ok := value.(float64); ok && v > 0 {
				query = query.Where("preco_venda <= ?", v)
			}
		case "search":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("nome_produto LIKE ? OR descricao LIKE ? OR sku LIKE ?",
					"%"+v+"%", "%"+v+"%", "%"+v+"%")
			}
		}
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		zap.L().Error("Error counting products", zap.Error(err))
		return nil, 0, err
	}

	// Get products with pagination
	err = query.Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		zap.L().Error("Error getting filtered products", zap.Error(err))
		return nil, 0, err
	}

	zap.L().Info("Successfully retrieved filtered products",
		zap.Int("count", len(products)),
		zap.Int64("total", total))

	return products, int(total), nil
}
