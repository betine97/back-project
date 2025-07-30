// cmd/config/config.go
package config

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strconv"

	entity "github.com/betine97/back-project.git/src/model/entitys"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Cfg         *Config
	PrivateKey  *rsa.PrivateKey
	Logger      *zap.Logger
	RedisClient *redis.Client
)

type Config struct {
	DBDriver      string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	WebServerPort string
	JWTSecret     string
	JWTExpiresIn  int
	CORSOrigins   string
}

func NewConfig() *Config {
	return Cfg
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Warning: Error loading .env file: %v", err)
	}

	Cfg = &Config{
		DBDriver:      os.Getenv("DB_DRIVER"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBName:        os.Getenv("DB_NAME"),
		WebServerPort: os.Getenv("WEB_SERVER_PORT"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		CORSOrigins:   getEnvWithDefault("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001,http://localhost:8080,http://127.0.0.1:3000,http://127.0.0.1:8080"),
	}

	Cfg.JWTExpiresIn, _ = strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))

	rng := rand.Reader
	PrivateKey, err = rsa.GenerateKey(rng, 2048)
	if err != nil {
		log.Fatalf("Error generating RSA key: %v", err)
	}

	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "ts"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""

	RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"), // Endereço do Redis
	})

	Logger, err = config.Build()
	if err != nil {
		log.Fatalf("Error initializing zap logger: %v", err)
	}
	zap.ReplaceGlobals(Logger)
}

func NewDatabaseConnection() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch Cfg.DBDriver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			Cfg.DBUser,
			Cfg.DBPassword,
			Cfg.DBHost,
			Cfg.DBPort,
			Cfg.DBName)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", Cfg.DBDriver)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewDatabaseConnectionForTenant(tenant entity.Tenants) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		tenant.DBUser,
		tenant.DBPassword,
		tenant.DBHost,
		tenant.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
