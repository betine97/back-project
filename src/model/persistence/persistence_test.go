package persistence

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Helper function para setup do mock DB
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create gorm DB: %v", err)
	}

	return gormDB, mock, mock
}

// Teste básico para GetAllFornecedores
func TestDBConnectionDBClient_GetAllFornecedores_Basic(t *testing.T) {
	// Arrange
	gormDB, mock, _ := setupMockDB(t)
	defer func() {
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()
	}()

	clientDBs := map[string]*gorm.DB{
		"db_1": gormDB,
	}
	repo := NewDBConnectionDBClient(clientDBs)
	userID := "1"

	rows := sqlmock.NewRows([]string{"id_fornecedor", "data_cadastro", "nome", "telefone", "email", "cidade", "estado", "status"}).
		AddRow(1, "2024-01-01", "Fornecedor 1", "11999999999", "fornecedor1@test.com", "São Paulo", "SP", "Ativo")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `fornecedores`")).
		WillReturnRows(rows)

	// Act
	fornecedores, err := repo.GetAllFornecedores(userID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, fornecedores, 1)
	assert.Equal(t, "Fornecedor 1", fornecedores[0].Nome)
}

// Teste básico para CreateFornecedor
func TestDBConnectionDBClient_CreateFornecedor_Basic(t *testing.T) {
	// Arrange
	gormDB, mock, _ := setupMockDB(t)
	defer func() {
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()
	}()

	clientDBs := map[string]*gorm.DB{
		"db_1": gormDB,
	}
	repo := NewDBConnectionDBClient(clientDBs)
	userID := "1"

	fornecedor := entity.Fornecedores{
		Nome:     "Novo Fornecedor",
		Email:    "novo@fornecedor.com",
		Telefone: "11999999999",
		Status:   "Ativo",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `fornecedores`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	err := repo.CreateFornecedor(fornecedor, userID)

	// Assert
	assert.NoError(t, err)
}
