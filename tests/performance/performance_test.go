package performance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/betine97/back-project.git/cmd/config/exceptions"
	"github.com/betine97/back-project.git/src/controller"
	dtos "github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/betine97/back-project.git/src/model/persistence"
	"github.com/betine97/back-project.git/src/model/service"
	"github.com/betine97/back-project.git/src/model/service/crypto"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =============================================================================
// SETUP PARA TESTES DE PERFORMANCE
// =============================================================================

type PerformanceTestSuite struct {
	app       *fiber.App
	container testcontainers.Container
	db        *gorm.DB
	service   service.ServiceInterface
}

// TestingInterface define métodos comuns entre *testing.T e *testing.B
type TestingInterface interface {
	Helper()
	Errorf(format string, args ...interface{})
	FailNow()
}

// requireNoError é um helper para verificar erros de forma compatível
func requireNoError(t TestingInterface, err error) {
	if err != nil {
		t.Helper()
		t.Errorf("Unexpected error: %v", err)
		t.FailNow()
	}
}

func setupPerformanceTestSuite(t TestingInterface) *PerformanceTestSuite {
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
	requireNoError(t, err)

	host, err := container.Host(ctx)
	requireNoError(t, err)

	port, err := container.MappedPort(ctx, "3306")
	requireNoError(t, err)

	connectionString := fmt.Sprintf("root:testpassword@tcp(%s:%s)/testdb?charset=utf8mb4&parseTime=True&loc=Local", host, port.Port())

	// Setup database
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	requireNoError(t, err)

	// Auto migrate
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Tenants{},
		&entity.Fornecedores{},
		&entity.Produto{},
		&entity.Pedido{},
	)
	requireNoError(t, err)

	// Setup dependencies
	cryptoService := &crypto.Crypto{}
	dbMaster := persistence.NewDBConnectionDBMaster(db)
	clientDBs := map[string]*gorm.DB{"db_1": db}
	dbClient := persistence.NewDBConnectionDBClient(clientDBs)

	mockRedis := &MockRedis{}
	mockTokenGen := &MockTokenGenerator{}

	serviceInstance := service.NewServiceInstance(cryptoService, dbMaster, dbClient, mockRedis, mockTokenGen)
	controllerInstance := controller.NewControllerInstance(serviceInstance)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Post("/cadastro", controllerInstance.CreateUser)
	app.Post("/login", controllerInstance.LoginUser)
	app.Get("/api/fornecedores", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controllerInstance.GetAllFornecedores(c)
	})
	app.Post("/api/fornecedores", func(c *fiber.Ctx) error {
		c.Locals("userID", "1")
		return controllerInstance.CreateFornecedor(c)
	})

	return &PerformanceTestSuite{
		app:       app,
		container: container,
		db:        db,
		service:   serviceInstance,
	}
}

func (suite *PerformanceTestSuite) cleanup() {
	if suite.container != nil {
		suite.container.Terminate(context.Background())
	}
}

// Mock implementations
type MockRedis struct{}

func (m *MockRedis) Ping(ctx context.Context) *redis.StatusCmd {
	return redis.NewStatusCmd(ctx)
}

func (m *MockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
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
// TESTES DE PERFORMANCE - CRIAÇÃO DE USUÁRIOS
// =============================================================================

func TestPerformance_UserCreation_Sequential(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	suite := setupPerformanceTestSuite(t)
	defer suite.cleanup()

	numUsers := 100
	start := time.Now()

	for i := 0; i < numUsers; i++ {
		userRequest := dtos.CreateUser{
			FirstName:   fmt.Sprintf("User%d", i),
			LastName:    "Performance",
			Email:       fmt.Sprintf("user%d@performance.test", i),
			NomeEmpresa: "Performance Company",
			Categoria:   "Tech",
			Segmento:    "Software",
			City:        "São Paulo",
			State:       "SP",
			Password:    "senha123!",
		}

		requestBody, _ := json.Marshal(userRequest)
		req := httptest.NewRequest("POST", "/cadastro", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.app.Test(req, 10000)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
	}

	duration := time.Since(start)
	avgTimePerUser := duration / time.Duration(numUsers)

	t.Logf("Created %d users in %v", numUsers, duration)
	t.Logf("Average time per user: %v", avgTimePerUser)
	t.Logf("Users per second: %.2f", float64(numUsers)/duration.Seconds())

	// Performance assertions
	assert.Less(t, duration, 30*time.Second, "Creating 100 users should take less than 30 seconds")
	assert.Less(t, avgTimePerUser, 300*time.Millisecond, "Average time per user should be less than 300ms")
}

func TestPerformance_UserCreation_Concurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	suite := setupPerformanceTestSuite(t)
	defer suite.cleanup()

	numUsers := 50
	concurrency := 10
	start := time.Now()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)
	errors := make(chan error, numUsers)

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(userIndex int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			userRequest := dtos.CreateUser{
				FirstName:   fmt.Sprintf("ConcurrentUser%d", userIndex),
				LastName:    "Performance",
				Email:       fmt.Sprintf("concurrent%d@performance.test", userIndex),
				NomeEmpresa: "Performance Company",
				Categoria:   "Tech",
				Segmento:    "Software",
				City:        "São Paulo",
				State:       "SP",
				Password:    "senha123!",
			}

			requestBody, _ := json.Marshal(userRequest)
			req := httptest.NewRequest("POST", "/cadastro", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := suite.app.Test(req, 10000)
			if err != nil {
				errors <- err
				return
			}
			if resp.StatusCode != 201 {
				errors <- fmt.Errorf("expected status 201, got %d", resp.StatusCode)
				return
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	var errorCount int
	for err := range errors {
		if err != nil {
			t.Logf("Error: %v", err)
			errorCount++
		}
	}

	duration := time.Since(start)
	successfulUsers := numUsers - errorCount

	t.Logf("Created %d users concurrently (concurrency: %d) in %v", successfulUsers, concurrency, duration)
	t.Logf("Users per second: %.2f", float64(successfulUsers)/duration.Seconds())
	t.Logf("Error rate: %.2f%%", float64(errorCount)/float64(numUsers)*100)

	// Performance assertions
	assert.Less(t, duration, 15*time.Second, "Creating 50 users concurrently should take less than 15 seconds")
	assert.Less(t, errorCount, numUsers/10, "Error rate should be less than 10%")
}

// =============================================================================
// TESTES DE PERFORMANCE - OPERAÇÕES DE FORNECEDORES
// =============================================================================

func TestPerformance_FornecedorOperations_Bulk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	suite := setupPerformanceTestSuite(t)
	defer suite.cleanup()

	// Setup user and tenant
	setupTestUserAndTenant(t, suite)

	numFornecedores := 200

	// Test bulk creation
	t.Run("Bulk Creation", func(t *testing.T) {
		start := time.Now()

		for i := 0; i < numFornecedores; i++ {
			fornecedorRequest := dtos.CreateFornecedorRequest{
				Nome:         fmt.Sprintf("Fornecedor Performance %d", i),
				Telefone:     "11999999999",
				Email:        fmt.Sprintf("fornecedor%d@performance.test", i),
				Cidade:       "São Paulo",
				Estado:       "SP",
				Status:       "Ativo",
				DataCadastro: time.Now().Format("2006-01-02"),
			}

			requestBody, _ := json.Marshal(fornecedorRequest)
			req := httptest.NewRequest("POST", "/api/fornecedores", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := suite.app.Test(req, 5000)
			require.NoError(t, err)
			assert.Equal(t, 201, resp.StatusCode)
		}

		duration := time.Since(start)
		t.Logf("Created %d fornecedores in %v", numFornecedores, duration)
		t.Logf("Average time per fornecedor: %v", duration/time.Duration(numFornecedores))
		t.Logf("Fornecedores per second: %.2f", float64(numFornecedores)/duration.Seconds())

		assert.Less(t, duration, 60*time.Second, "Creating 200 fornecedores should take less than 60 seconds")
	})

	// Test bulk reading
	t.Run("Bulk Reading", func(t *testing.T) {
		numReads := 100
		start := time.Now()

		for i := 0; i < numReads; i++ {
			req := httptest.NewRequest("GET", "/api/fornecedores", nil)
			resp, err := suite.app.Test(req, 5000)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
		}

		duration := time.Since(start)
		t.Logf("Performed %d reads in %v", numReads, duration)
		t.Logf("Average time per read: %v", duration/time.Duration(numReads))
		t.Logf("Reads per second: %.2f", float64(numReads)/duration.Seconds())

		assert.Less(t, duration, 10*time.Second, "100 reads should take less than 10 seconds")
	})
}

func TestPerformance_FornecedorOperations_Concurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	suite := setupPerformanceTestSuite(t)
	defer suite.cleanup()

	// Setup user and tenant
	setupTestUserAndTenant(t, suite)

	// Test concurrent reads
	t.Run("Concurrent Reads", func(t *testing.T) {
		numReads := 100
		concurrency := 20
		start := time.Now()

		var wg sync.WaitGroup
		semaphore := make(chan struct{}, concurrency)
		errors := make(chan error, numReads)

		for i := 0; i < numReads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				req := httptest.NewRequest("GET", "/api/fornecedores", nil)
				resp, err := suite.app.Test(req, 5000)
				if err != nil {
					errors <- err
					return
				}
				if resp.StatusCode != 200 {
					errors <- fmt.Errorf("expected status 200, got %d", resp.StatusCode)
					return
				}
			}()
		}

		wg.Wait()
		close(errors)

		var errorCount int
		for err := range errors {
			if err != nil {
				t.Logf("Error: %v", err)
				errorCount++
			}
		}

		duration := time.Since(start)
		successfulReads := numReads - errorCount

		t.Logf("Performed %d concurrent reads (concurrency: %d) in %v", successfulReads, concurrency, duration)
		t.Logf("Reads per second: %.2f", float64(successfulReads)/duration.Seconds())
		t.Logf("Error rate: %.2f%%", float64(errorCount)/float64(numReads)*100)

		assert.Less(t, duration, 5*time.Second, "100 concurrent reads should take less than 5 seconds")
		assert.Less(t, errorCount, numReads/20, "Error rate should be less than 5%")
	})
}

// =============================================================================
// TESTES DE PERFORMANCE - OPERAÇÕES DE BANCO DE DADOS
// =============================================================================

func TestPerformance_DatabaseOperations_DirectAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	suite := setupPerformanceTestSuite(t)
	defer suite.cleanup()

	// Test direct database operations performance
	t.Run("Direct User Creation", func(t *testing.T) {
		numUsers := 1000
		start := time.Now()

		users := make([]entity.User, numUsers)
		for i := 0; i < numUsers; i++ {
			users[i] = entity.User{
				FirstName: fmt.Sprintf("DirectUser%d", i),
				LastName:  "Performance",
				Email:     fmt.Sprintf("direct%d@performance.test", i),
				Password:  "hashedpassword",
			}
		}

		// Bulk insert
		err := suite.db.CreateInBatches(users, 100).Error
		require.NoError(t, err)

		duration := time.Since(start)
		t.Logf("Direct created %d users in %v", numUsers, duration)
		t.Logf("Users per second: %.2f", float64(numUsers)/duration.Seconds())

		assert.Less(t, duration, 5*time.Second, "Direct creation of 1000 users should take less than 5 seconds")
	})

	t.Run("Direct User Query", func(t *testing.T) {
		numQueries := 1000
		start := time.Now()

		for i := 0; i < numQueries; i++ {
			var user entity.User
			err := suite.db.Where("email = ?", fmt.Sprintf("direct%d@performance.test", i%100)).First(&user).Error
			if err != nil && err.Error() != "record not found" {
				require.NoError(t, err)
			}
		}

		duration := time.Since(start)
		t.Logf("Performed %d direct queries in %v", numQueries, duration)
		t.Logf("Queries per second: %.2f", float64(numQueries)/duration.Seconds())

		assert.Less(t, duration, 2*time.Second, "1000 direct queries should take less than 2 seconds")
	})
}

// =============================================================================
// BENCHMARKS
// =============================================================================

func BenchmarkUserCreation_Service(b *testing.B) {
	suite := setupPerformanceTestSuite(b)
	defer suite.cleanup()

	userRequest := dtos.CreateUser{
		FirstName:   "Benchmark",
		LastName:    "User",
		Email:       "benchmark@test.com",
		NomeEmpresa: "Benchmark Company",
		Categoria:   "Tech",
		Segmento:    "Software",
		City:        "São Paulo",
		State:       "SP",
		Password:    "senha123!",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userRequest.Email = fmt.Sprintf("benchmark%d@test.com", i)
		_, err := suite.service.CreateUserService(userRequest)
		if err != nil {
			b.Fatalf("Error creating user: %v", err)
		}
	}
}

func BenchmarkFornecedorCreation_Service(b *testing.B) {
	suite := setupPerformanceTestSuite(b)
	defer suite.cleanup()

	// Setup user and tenant
	setupTestUserAndTenant(b, suite)

	fornecedorRequest := dtos.CreateFornecedorRequest{
		Nome:         "Benchmark Fornecedor",
		Telefone:     "11999999999",
		Email:        "benchmark@fornecedor.test",
		Cidade:       "São Paulo",
		Estado:       "SP",
		Status:       "Ativo",
		DataCadastro: time.Now().Format("2006-01-02"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fornecedorRequest.Email = fmt.Sprintf("benchmark%d@fornecedor.test", i)
		_, err := suite.service.CreateFornecedorService("1", fornecedorRequest)
		if err != nil {
			b.Fatalf("Error creating fornecedor: %v", err)
		}
	}
}

func BenchmarkFornecedorRead_Service(b *testing.B) {
	suite := setupPerformanceTestSuite(b)
	defer suite.cleanup()

	// Setup user and tenant
	setupTestUserAndTenant(b, suite)

	// Create some fornecedores for reading
	for i := 0; i < 10; i++ {
		fornecedorRequest := dtos.CreateFornecedorRequest{
			Nome:         fmt.Sprintf("Read Fornecedor %d", i),
			Email:        fmt.Sprintf("read%d@fornecedor.test", i),
			Telefone:     "11999999999",
			Status:       "Ativo",
			DataCadastro: time.Now().Format("2006-01-02"),
		}
		_, err := suite.service.CreateFornecedorService("1", fornecedorRequest)
		if err != nil {
			b.Fatalf("Error creating fornecedor: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := suite.service.GetAllFornecedoresService("1")
		if err != nil {
			b.Fatalf("Error reading fornecedores: %v", err)
		}
	}
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func setupTestUserAndTenant(t testing.TB, suite *PerformanceTestSuite) {
	// Create user
	user := entity.User{
		FirstName: "Performance",
		LastName:  "User",
		Email:     "performance@test.com",
		Password:  "hashedpassword",
	}
	err := suite.db.Create(&user).Error
	require.NoError(t, err)

	// Create tenant
	tenant := entity.Tenants{
		UserID:      user.ID,
		NomeEmpresa: "Performance Company",
		DBName:      "test_db",
		DBUser:      "test_user",
		DBPassword:  "test_pass",
		DBHost:      "localhost",
		DBPort:      "3306",
		CreatedAt:   time.Now().Format("2006-01-02"),
	}
	err = suite.db.Create(&tenant).Error
	require.NoError(t, err)
}
