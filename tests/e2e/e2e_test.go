package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/betine97/back-project.git/cmd/config/exceptions"
	"github.com/betine97/back-project.git/src/controller"
	dtos "github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/betine97/back-project.git/src/model/persistence"
	"github.com/betine97/back-project.git/src/model/service"
	"github.com/betine97/back-project.git/src/model/service/crypto"
	"github.com/betine97/back-project.git/src/routes"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =============================================================================
// SETUP E HELPERS PARA TESTES E2E
// =============================================================================

type E2ETestSuite struct {
	app       *fiber.App
	container testcontainers.Container
	db        *gorm.DB
}

func setupE2ETestSuite(t *testing.T) *E2ETestSuite {
	ctx := context.Background()

	// Setup MySQL container
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
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "3306")
	require.NoError(t, err)

	connectionString := fmt.Sprintf("root:testpassword@tcp(%s:%s)/testdb?charset=utf8mb4&parseTime=True&loc=Local", host, port.Port())

	// Setup database
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Auto migrate
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Tenants{},
		&entity.Fornecedores{},
		&entity.Produto{},
		&entity.Pedido{},
	)
	require.NoError(t, err)

	// Setup application dependencies
	cryptoService := &crypto.Crypto{}

	dbMaster := persistence.NewDBConnectionDBMaster(db)
	clientDBs := map[string]*gorm.DB{
		"db_1": db,
	}
	dbClient := persistence.NewDBConnectionDBClient(clientDBs)

	// Mock Redis and Token Generator for E2E tests
	mockRedis := &MockRedis{}
	mockTokenGen := &MockTokenGenerator{}

	serviceInstance := service.NewServiceInstance(cryptoService, dbMaster, dbClient, mockRedis, mockTokenGen)
	controllerInstance := controller.NewControllerInstance(serviceInstance)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Setup CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Setup routes
	routes.SetupRoutes(app, controllerInstance)

	return &E2ETestSuite{
		app:       app,
		container: container,
		db:        db,
	}
}

func (suite *E2ETestSuite) cleanup() {
	if suite.container != nil {
		suite.container.Terminate(context.Background())
	}
}

// Mock implementations for E2E tests
type MockRedis struct{}

func (m *MockRedis) Ping(ctx context.Context) *redis.StatusCmd {
	return redis.NewStatusCmd(ctx)
}

func (m *MockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	// Simulate cache miss for E2E tests
	cmd := redis.NewStringCmd(ctx)
	cmd.SetErr(redis.Nil)
	return cmd
}

func (m *MockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redis.NewStatusCmd(ctx)
}

type MockTokenGenerator struct{}

func (m *MockTokenGenerator) GenerateToken(userID uint) (string, *exceptions.RestErr) {
	return fmt.Sprintf("mock.jwt.token.%d", userID), nil
}

// =============================================================================
// TESTES E2E - FLUXO COMPLETO DE USUÁRIO
// =============================================================================

func TestE2E_UserRegistrationAndLogin_CompleteFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	// Test data
	userRequest := dtos.CreateUser{
		FirstName:   "João",
		LastName:    "Silva",
		Email:       "joao@e2e.test",
		NomeEmpresa: "Empresa E2E",
		Categoria:   "Tecnologia",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "senha123!",
	}

	// Step 1: Register user
	t.Run("Register User", func(t *testing.T) {
		requestBody, _ := json.Marshal(userRequest)
		req := httptest.NewRequest("POST", "/cadastro", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.app.Test(req, 30000) // 30 second timeout
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		// Verify user was created in database
		var user entity.User
		err = suite.db.Where("email = ?", userRequest.Email).First(&user).Error
		require.NoError(t, err)
		assert.Equal(t, userRequest.FirstName, user.FirstName)
		assert.Equal(t, userRequest.Email, user.Email)
	})

	// Step 2: Create tenant for the user
	t.Run("Create Tenant", func(t *testing.T) {
		var user entity.User
		err := suite.db.Where("email = ?", userRequest.Email).First(&user).Error
		require.NoError(t, err)

		tenant := entity.Tenants{
			UserID:      user.ID,
			NomeEmpresa: userRequest.NomeEmpresa,
			DBName:      "test_db",
			DBUser:      "test_user",
			DBPassword:  "test_pass",
			DBHost:      "localhost",
			DBPort:      "3306",
			CreatedAt:   time.Now().Format("2006-01-02"),
		}

		err = suite.db.Create(&tenant).Error
		require.NoError(t, err)
	})

	// Step 3: Login user
	var authToken string
	t.Run("Login User", func(t *testing.T) {
		loginRequest := dtos.UserLogin{
			Email:    userRequest.Email,
			Password: userRequest.Password,
		}

		requestBody, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.app.Test(req, 30000)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		token, exists := response["token"]
		require.True(t, exists)
		authToken = token.(string)
		assert.NotEmpty(t, authToken)
	})

	// Step 4: Access protected endpoints with token
	t.Run("Access Protected Endpoints", func(t *testing.T) {
		// Test health check (public)
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := suite.app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Test protected endpoint without token (should fail)
		req = httptest.NewRequest("GET", "/api/fornecedores", nil)
		resp, err = suite.app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode) // Unauthorized

		// Test protected endpoint with token (should work)
		req = httptest.NewRequest("GET", "/api/fornecedores", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		resp, err = suite.app.Test(req, 30000)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestE2E_FornecedorCRUD_CompleteFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	// Setup user and get auth token
	authToken := setupUserAndGetToken(t, suite)

	var fornecedorID int

	// Step 1: Create Fornecedor
	t.Run("Create Fornecedor", func(t *testing.T) {
		fornecedorRequest := dtos.CreateFornecedorRequest{
			Nome:         "Fornecedor E2E",
			Telefone:     "11999999999",
			Email:        "fornecedor@e2e.test",
			Cidade:       "São Paulo",
			Estado:       "SP",
			Status:       "Ativo",
			DataCadastro: time.Now().Format("2006-01-02"),
		}

		requestBody, _ := json.Marshal(fornecedorRequest)
		req := httptest.NewRequest("POST", "/api/fornecedores", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := suite.app.Test(req, 30000)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		// Get the created fornecedor ID from database
		var fornecedor entity.Fornecedores
		err = suite.db.Where("email = ?", fornecedorRequest.Email).First(&fornecedor).Error
		require.NoError(t, err)
		fornecedorID = fornecedor.ID
	})

	// Step 2: Get All Fornecedores
	t.Run("Get All Fornecedores", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/fornecedores", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := suite.app.Test(req, 30000)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var response dtos.FornecedorListResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, 1, response.Total)
		assert.Len(t, response.Fornecedores, 1)
		assert.Equal(t, "Fornecedor E2E", response.Fornecedores[0].Nome)
	})

	// Step 3: Update Fornecedor Field
	t.Run("Update Fornecedor Field", func(t *testing.T) {
		updateRequest := dtos.UpdateFornecedorRequest{
			Campo: "nome",
			Valor: "Fornecedor E2E Atualizado",
		}

		requestBody, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/fornecedores/changefields/%d", fornecedorID), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := suite.app.Test(req, 30000)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Verify update in database
		var fornecedor entity.Fornecedores
		err = suite.db.First(&fornecedor, fornecedorID).Error
		require.NoError(t, err)
		assert.Equal(t, "Fornecedor E2E Atualizado", fornecedor.Nome)
	})

	// Step 4: Change Status
	t.Run("Change Fornecedor Status", func(t *testing.T) {
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/fornecedores/changestatus/%d", fornecedorID), nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := suite.app.Test(req, 30000)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Verify status change in database
		var fornecedor entity.Fornecedores
		err = suite.db.First(&fornecedor, fornecedorID).Error
		require.NoError(t, err)
		assert.Equal(t, "Inativo", fornecedor.Status) // Should toggle from Ativo to Inativo
	})

	// Step 5: Delete Fornecedor
	t.Run("Delete Fornecedor", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/fornecedores/%d", fornecedorID), nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := suite.app.Test(req, 30000)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Verify deletion in database
		var count int64
		suite.db.Model(&entity.Fornecedores{}).Where("id_fornecedor = ?", fornecedorID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestE2E_ProductCRUD_CompleteFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	// Setup user and get auth token
	authToken := setupUserAndGetToken(t, suite)

	var productID int

	// Step 1: Create Product
	t.Run("Create Product", func(t *testing.T) {
		productRequest := dtos.CreateProductRequest{
			DataCadastro:  time.Now().Format("2006-01-02"),
			CodigoBarra:   "1234567890123",
			NomeProduto:   "Produto E2E",
			SKU:           "SKU-E2E-001",
			Categoria:     "Eletrônicos",
			DestinadoPara: "Consumidor Final",
			Variacao:      "Azul",
			Marca:         "Marca E2E",
			Descricao:     "Produto para teste E2E",
			Status:        "Ativo",
			PrecoVenda:    99.99,
		}

		requestBody, _ := json.Marshal(productRequest)
		req := httptest.NewRequest("POST", "/api/produtos", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := suite.app.Test(req, 30000)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		// Get the created product ID from database
		var product entity.Produto
		err = suite.db.Where("sku = ?", productRequest.SKU).First(&product).Error
		require.NoError(t, err)
		productID = product.IDProduto
	})

	// Step 2: Get All Products
	t.Run("Get All Products", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/produtos", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := suite.app.Test(req, 30000)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var response dtos.ProductListResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, 1, response.Total)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, "Produto E2E", response.Products[0].NomeProduto)
	})

	// Step 3: Delete Product
	t.Run("Delete Product", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/produtos/%d", productID), nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := suite.app.Test(req, 30000)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Verify deletion in database
		var count int64
		suite.db.Model(&entity.Produto{}).Where("id_produto = ?", productID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestE2E_PedidosRead_Flow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	// Setup user and get auth token
	authToken := setupUserAndGetToken(t, suite)

	// Create a pedido for testing
	pedido := entity.Pedido{
		IDFornecedor: 1,
		DataPedido:   "2024-01-01",
		DataEntrega:  "2024-01-15",
		ValorFrete:   15.50,
		CustoPedido:  100.00,
		ValorTotal:   115.50,
		Descricao:    "Pedido E2E Test",
		Status:       "Pendente",
	}
	err := suite.db.Create(&pedido).Error
	require.NoError(t, err)

	// Test Get All Pedidos
	t.Run("Get All Pedidos", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/pedidos", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := suite.app.Test(req, 30000)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var response dtos.PedidoListResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, 1, response.Total)
		assert.Len(t, response.Pedidos, 1)
		assert.Equal(t, "Pedido E2E Test", response.Pedidos[0].Descricao)
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func setupUserAndGetToken(t *testing.T, suite *E2ETestSuite) string {
	// Create user
	userRequest := dtos.CreateUser{
		FirstName:   "Test",
		LastName:    "User",
		Email:       "test@e2e.test",
		NomeEmpresa: "Test Company",
		Categoria:   "Tech",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "senha123!",
	}

	requestBody, _ := json.Marshal(userRequest)
	req := httptest.NewRequest("POST", "/cadastro", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req, 30000)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	// Create tenant
	var user entity.User
	err = suite.db.Where("email = ?", userRequest.Email).First(&user).Error
	require.NoError(t, err)

	tenant := entity.Tenants{
		UserID:      user.ID,
		NomeEmpresa: userRequest.NomeEmpresa,
		DBName:      "test_db",
		DBUser:      "test_user",
		DBPassword:  "test_pass",
		DBHost:      "localhost",
		DBPort:      "3306",
		CreatedAt:   time.Now().Format("2006-01-02"),
	}

	err = suite.db.Create(&tenant).Error
	require.NoError(t, err)

	// Login and get token
	loginRequest := dtos.UserLogin{
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}

	requestBody, _ = json.Marshal(loginRequest)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.app.Test(req, 30000)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	token, exists := response["token"]
	require.True(t, exists)
	return token.(string)
}
