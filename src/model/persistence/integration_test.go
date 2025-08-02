package persistence

import (
	"context"
	"fmt"
	"testing"
	"time"

	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =============================================================================
// TESTES DE INTEGRAÇÃO COM BANCO REAL
// =============================================================================

type MySQLContainer struct {
	testcontainers.Container
	ConnectionString string
}

func setupMySQLContainer(ctx context.Context) (*MySQLContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "testpassword",
			"MYSQL_DATABASE":      "testdb",
		},
		WaitingFor: wait.ForLog("ready for connections").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		return nil, err
	}

	connectionString := fmt.Sprintf("root:testpassword@tcp(%s:%s)/testdb?charset=utf8mb4&parseTime=True&loc=Local", host, port.Port())

	return &MySQLContainer{
		Container:        container,
		ConnectionString: connectionString,
	}, nil
}

func setupTestDB(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Auto migrate tables
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Tenants{},
		&entity.Fornecedores{},
		&entity.Produto{},
		&entity.Pedido{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestDBConnectionDBMaster_Integration_CreateAndGetUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Arrange
	ctx := context.Background()
	container, err := setupMySQLContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	db, err := setupTestDB(container.ConnectionString)
	require.NoError(t, err)

	repo := NewDBConnectionDBMaster(db)

	user := entity.User{
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@integration.test",
		NomeEmpresa: "Empresa Teste",
		Categoria:   "Tecnologia",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "hashedpassword123",
	}

	// Act - Create user
	err = repo.CreateUser(user)
	require.NoError(t, err)

	// Act - Get user
	retrievedUser := repo.GetUser(user.Email)

	// Assert
	assert.NotNil(t, retrievedUser)
	assert.NotEqual(t, uint(0), retrievedUser.ID)
	assert.Equal(t, user.FirstName, retrievedUser.FirstName)
	assert.Equal(t, user.LastName, retrievedUser.LastName)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.NomeEmpresa, retrievedUser.NomeEmpresa)
	assert.Equal(t, user.Password, retrievedUser.Password)
}

func TestDBConnectionDBMaster_Integration_VerifyExist(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Arrange
	ctx := context.Background()
	container, err := setupMySQLContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	db, err := setupTestDB(container.ConnectionString)
	require.NoError(t, err)

	repo := NewDBConnectionDBMaster(db)

	user := entity.User{
		FirstName: "Maria",
		LastName:  "Santos",
		Email:     "maria@integration.test",
		Password:  "hashedpassword456",
	}

	// Act - Create user first
	err = repo.CreateUser(user)
	require.NoError(t, err)

	// Act - Verify existing user
	exists, err := repo.VerifyExist(user.Email)
	require.NoError(t, err)
	assert.True(t, exists)

	// Act - Verify non-existing user
	exists, err = repo.VerifyExist("nonexistent@test.com")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestDBConnectionDBMaster_Integration_CreateTenantAndRetrieve(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Arrange
	ctx := context.Background()
	container, err := setupMySQLContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	db, err := setupTestDB(container.ConnectionString)
	require.NoError(t, err)

	repo := NewDBConnectionDBMaster(db)

	// Create user first
	user := entity.User{
		FirstName: "Carlos",
		LastName:  "Oliveira",
		Email:     "carlos@integration.test",
		Password:  "hashedpassword789",
	}
	err = repo.CreateUser(user)
	require.NoError(t, err)

	// Get the created user to get the ID
	createdUser := repo.GetUser(user.Email)
	require.NotNil(t, createdUser)
	require.NotEqual(t, uint(0), createdUser.ID)

	// Create tenant manually (since we don't have CreateTenant method in the interface)
	tenant := entity.Tenants{
		UserID:      createdUser.ID,
		NomeEmpresa: "Carlos Company",
		DBName:      "carlos_db",
		DBUser:      "carlos_user",
		DBPassword:  "carlos_pass",
		DBHost:      "localhost",
		DBPort:      "3306",
		CreatedAt:   "2024-01-01",
	}

	result := db.Create(&tenant)
	require.NoError(t, result.Error)

	// Act - Get tenant by user ID
	retrievedTenant := repo.GetTenantByUserID(createdUser.ID)

	// Assert
	assert.NotNil(t, retrievedTenant)
	assert.NotEqual(t, uint(0), retrievedTenant.ID)
	assert.Equal(t, createdUser.ID, retrievedTenant.UserID)
	assert.Equal(t, tenant.NomeEmpresa, retrievedTenant.NomeEmpresa)
	assert.Equal(t, tenant.DBName, retrievedTenant.DBName)
}

func TestDBConnectionDBClient_Integration_FornecedorCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Arrange
	ctx := context.Background()
	container, err := setupMySQLContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	db, err := setupTestDB(container.ConnectionString)
	require.NoError(t, err)

	clientDBs := map[string]*gorm.DB{
		"db_1": db,
	}
	repo := NewDBConnectionDBClient(clientDBs)
	userID := "1"

	fornecedor := entity.Fornecedores{
		DataCadastro: "2024-01-01",
		Nome:         "Fornecedor Integração",
		Telefone:     "11999999999",
		Email:        "fornecedor@integration.test",
		Cidade:       "São Paulo",
		Estado:       "SP",
		Status:       "Ativo",
	}

	// Act - Create fornecedor
	err = repo.CreateFornecedor(fornecedor, userID)
	require.NoError(t, err)

	// Act - Get all fornecedores
	fornecedores, err := repo.GetAllFornecedores(userID)
	require.NoError(t, err)
	require.Len(t, fornecedores, 1)

	createdFornecedor := fornecedores[0]
	assert.Equal(t, fornecedor.Nome, createdFornecedor.Nome)
	assert.Equal(t, fornecedor.Email, createdFornecedor.Email)
	assert.Equal(t, fornecedor.Status, createdFornecedor.Status)

	// Act - Get fornecedor by ID
	fornecedorID := fmt.Sprintf("%d", createdFornecedor.ID)
	retrievedFornecedor, err := repo.GetFornecedorById(fornecedorID, userID)
	require.NoError(t, err)
	assert.Equal(t, createdFornecedor.ID, retrievedFornecedor.ID)
	assert.Equal(t, createdFornecedor.Nome, retrievedFornecedor.Nome)

	// Act - Update fornecedor field
	err = repo.UpdateFornecedorField(fornecedorID, "nome", "Fornecedor Atualizado", userID)
	require.NoError(t, err)

	// Verify update
	updatedFornecedor, err := repo.GetFornecedorById(fornecedorID, userID)
	require.NoError(t, err)
	assert.Equal(t, "Fornecedor Atualizado", updatedFornecedor.Nome)

	// Act - Update full fornecedor
	updatedFornecedor.Status = "Inativo"
	err = repo.UpdateFornecedor(*updatedFornecedor, userID)
	require.NoError(t, err)

	// Verify full update
	finalFornecedor, err := repo.GetFornecedorById(fornecedorID, userID)
	require.NoError(t, err)
	assert.Equal(t, "Inativo", finalFornecedor.Status)

	// Act - Delete fornecedor
	err = repo.DeleteFornecedor(fornecedorID, userID)
	require.NoError(t, err)

	// Verify deletion
	fornecedoresAfterDelete, err := repo.GetAllFornecedores(userID)
	require.NoError(t, err)
	assert.Len(t, fornecedoresAfterDelete, 0)
}

func TestDBConnectionDBClient_Integration_ProductCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Arrange
	ctx := context.Background()
	container, err := setupMySQLContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	db, err := setupTestDB(container.ConnectionString)
	require.NoError(t, err)

	clientDBs := map[string]*gorm.DB{
		"db_1": db,
	}
	repo := NewDBConnectionDBClient(clientDBs)
	userID := "1"

	produto := entity.Produto{
		DataCadastro:  "2024-01-01",
		CodigoBarra:   "1234567890123",
		NomeProduto:   "Produto Integração",
		SKU:           "SKU-INT-001",
		Categoria:     "Eletrônicos",
		DestinadoPara: "Consumidor Final",
		Variacao:      "Azul",
		Marca:         "Marca Teste",
		Descricao:     "Produto para teste de integração",
		Status:        "Ativo",
		PrecoVenda:    99.99,
	}

	// Act - Create product
	err = repo.CreateProduct(produto, userID)
	require.NoError(t, err)

	// Act - Get all products
	produtos, err := repo.GetAllProducts(userID)
	require.NoError(t, err)
	require.Len(t, produtos, 1)

	createdProduto := produtos[0]
	assert.Equal(t, produto.NomeProduto, createdProduto.NomeProduto)
	assert.Equal(t, produto.SKU, createdProduto.SKU)
	assert.Equal(t, produto.PrecoVenda, createdProduto.PrecoVenda)

	// Act - Delete product
	produtoID := fmt.Sprintf("%d", createdProduto.IDProduto)
	err = repo.DeleteProduct(produtoID, userID)
	require.NoError(t, err)

	// Verify deletion
	produtosAfterDelete, err := repo.GetAllProducts(userID)
	require.NoError(t, err)
	assert.Len(t, produtosAfterDelete, 0)
}

func TestDBConnectionDBClient_Integration_PedidosRead(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Arrange
	ctx := context.Background()
	container, err := setupMySQLContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	db, err := setupTestDB(container.ConnectionString)
	require.NoError(t, err)

	clientDBs := map[string]*gorm.DB{
		"db_1": db,
	}
	repo := NewDBConnectionDBClient(clientDBs)
	userID := "1"

	// Create a pedido manually for testing
	pedido := entity.Pedido{
		IDFornecedor: 1,
		DataPedido:   "2024-01-01",
		DataEntrega:  "2024-01-15",
		ValorFrete:   15.50,
		CustoPedido:  100.00,
		ValorTotal:   115.50,
		Descricao:    "Pedido de teste de integração",
		Status:       "Pendente",
	}

	result := db.Create(&pedido)
	require.NoError(t, result.Error)

	// Act - Get all pedidos
	pedidos, err := repo.GetAllPedidos(userID)
	require.NoError(t, err)
	require.Len(t, pedidos, 1)

	retrievedPedido := pedidos[0]
	assert.Equal(t, pedido.IDFornecedor, retrievedPedido.IDFornecedor)
	assert.Equal(t, pedido.ValorTotal, retrievedPedido.ValorTotal)
	assert.Equal(t, pedido.Status, retrievedPedido.Status)
}

// =============================================================================
// TESTES DE PERFORMANCE E CARGA
// =============================================================================

func TestDBConnectionDBMaster_Integration_Performance_CreateMultipleUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Arrange
	ctx := context.Background()
	container, err := setupMySQLContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	db, err := setupTestDB(container.ConnectionString)
	require.NoError(t, err)

	repo := NewDBConnectionDBMaster(db)

	// Act - Create 100 users
	start := time.Now()
	for i := 0; i < 100; i++ {
		user := entity.User{
			FirstName: fmt.Sprintf("User%d", i),
			LastName:  "Test",
			Email:     fmt.Sprintf("user%d@performance.test", i),
			Password:  "hashedpassword",
		}
		err := repo.CreateUser(user)
		require.NoError(t, err)
	}
	duration := time.Since(start)

	// Assert - Should complete within reasonable time (adjust as needed)
	assert.Less(t, duration, 10*time.Second, "Creating 100 users took too long")

	// Verify all users were created
	var count int64
	db.Model(&entity.User{}).Count(&count)
	assert.Equal(t, int64(100), count)
}

func TestDBConnectionDBClient_Integration_Performance_BulkFornecedores(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Arrange
	ctx := context.Background()
	container, err := setupMySQLContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	db, err := setupTestDB(container.ConnectionString)
	require.NoError(t, err)

	clientDBs := map[string]*gorm.DB{
		"db_1": db,
	}
	repo := NewDBConnectionDBClient(clientDBs)
	userID := "1"

	// Act - Create 50 fornecedores
	start := time.Now()
	for i := 0; i < 50; i++ {
		fornecedor := entity.Fornecedores{
			Nome:     fmt.Sprintf("Fornecedor %d", i),
			Email:    fmt.Sprintf("fornecedor%d@performance.test", i),
			Telefone: "11999999999",
			Status:   "Ativo",
		}
		err := repo.CreateFornecedor(fornecedor, userID)
		require.NoError(t, err)
	}
	createDuration := time.Since(start)

	// Act - Read all fornecedores
	start = time.Now()
	fornecedores, err := repo.GetAllFornecedores(userID)
	readDuration := time.Since(start)

	// Assert
	require.NoError(t, err)
	assert.Len(t, fornecedores, 50)
	assert.Less(t, createDuration, 5*time.Second, "Creating 50 fornecedores took too long")
	assert.Less(t, readDuration, 1*time.Second, "Reading 50 fornecedores took too long")
}

// =============================================================================
// TESTES DE TRANSAÇÕES E CONSISTÊNCIA
// =============================================================================

func TestDBConnectionDBMaster_Integration_TransactionRollback(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Arrange
	ctx := context.Background()
	container, err := setupMySQLContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	db, err := setupTestDB(container.ConnectionString)
	require.NoError(t, err)

	// Act - Try to create user with duplicate email in transaction
	tx := db.Begin()

	user1 := entity.User{
		FirstName: "User1",
		Email:     "duplicate@test.com",
		Password:  "password1",
	}

	user2 := entity.User{
		FirstName: "User2",
		Email:     "duplicate@test.com", // Same email - should cause constraint violation
		Password:  "password2",
	}

	// Create first user
	err = tx.Create(&user1).Error
	require.NoError(t, err)

	// Try to create second user with same email
	err = tx.Create(&user2).Error
	assert.Error(t, err) // Should fail due to unique constraint

	// Rollback transaction
	tx.Rollback()

	// Assert - No users should exist after rollback
	var count int64
	db.Model(&entity.User{}).Where("email = ?", "duplicate@test.com").Count(&count)
	assert.Equal(t, int64(0), count)
}

// =============================================================================
// HELPER FUNCTIONS FOR INTEGRATION TESTS
// =============================================================================

func cleanupDatabase(db *gorm.DB) error {
	// Clean up in reverse order of dependencies
	if err := db.Exec("DELETE FROM pedidos").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM produtos").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM fornecedores").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM tenants").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		return err
	}
	return nil
}
