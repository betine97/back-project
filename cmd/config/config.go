//go:build !test
// +build !test

// cmd/config/config.go
package config

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

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

// getEnvWithDefault retorna o valor da variável de ambiente ou um valor padrão
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvOrFail retorna o valor da variável de ambiente ou falha a aplicação
func getEnvOrFail(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("❌ Required environment variable %s is not set", key)
	}
	return value
}

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("⚠️  Warning: Error loading .env file: %v", err)
		log.Printf("📋 Make sure you have a .env file with required environment variables")
	}

	Cfg = &Config{
		DBDriver:      getEnvWithDefault("DB_DRIVER", "mysql"),
		DBHost:        getEnvOrFail("DB_HOST"),
		DBPort:        getEnvWithDefault("DB_PORT", "3306"),
		DBUser:        getEnvOrFail("DB_USER"),
		DBPassword:    getEnvOrFail("DB_PASSWORD"),
		DBName:        getEnvOrFail("DB_NAME"),
		WebServerPort: getEnvWithDefault("WEB_SERVER_PORT", "8080"),
		JWTSecret:     getEnvOrFail("JWT_SECRET"),
		CORSOrigins:   getEnvWithDefault("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001,http://localhost:8080,http://127.0.0.1:3000,http://127.0.0.1:8080"),
	}

	if expiresIn := os.Getenv("JWT_EXPIRES_IN"); expiresIn != "" {
		Cfg.JWTExpiresIn, _ = strconv.Atoi(expiresIn)
	} else {
		Cfg.JWTExpiresIn = 30 // valor padrão
	}

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
		Addr: getEnvWithDefault("REDIS_ADDR", "localhost:6379"), // Endereço do Redis
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

	return db, nil
}

func ConnectionDBClients() (map[string]*gorm.DB, error) {
	var dbConnections = make(map[string]*gorm.DB)
	file, err := os.Open("cmd/config/dbclients.json")
	if err != nil {
		return nil, fmt.Errorf("error opening dbclients.json: %v", err)
	}
	defer file.Close()

	var clients struct {
		Clients []struct {
			DB_CLIENT       int    `json:"DB_CLIENT"`
			DB_DRIVER       string `json:"DB_DRIVER"`
			DB_HOST         string `json:"DB_HOST"`
			DB_PORT         int    `json:"DB_PORT"`
			DB_USER         string `json:"DB_USER"`
			DB_PASSWORD     string `json:"DB_PASSWORD"`
			DB_NAME         string `json:"DB_NAME"`
			WEB_SERVER_PORT int    `json:"WEB_SERVER_PORT"`
		} `json:"clients"`
	}

	if err := json.NewDecoder(file).Decode(&clients); err != nil {
		return nil, fmt.Errorf("error decoding dbclients.json: %v", err)
	}

	for _, client := range clients.Clients {
		clientID := strconv.Itoa(client.DB_CLIENT) // Convertendo clientID para string
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			client.DB_USER,
			client.DB_PASSWORD,
			client.DB_HOST,
			client.DB_PORT,
			client.DB_NAME)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("error connecting to database %s: %v", clientID, err)
		}

		key := "db_" + clientID
		dbConnections[key] = db
		log.Printf("✅ Conexão criada: %s -> %s", key, client.DB_NAME)
	}

	return dbConnections, nil
}
