package controller

import (
	"time"

	"github.com/betine97/back-project.git/cmd/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ReadinessResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
	Uptime    string            `json:"uptime"`
}

var startTime = time.Now()

// HealthCheck (Liveness Probe) - Verifica se a aplicação está viva
func (ctl *Controller) HealthCheck(ctx *fiber.Ctx) error {
	zap.L().Info("💓 Verificação de vida da aplicação (liveness)")

	// Health check básico - apenas verifica se a aplicação está viva
	uptime := time.Since(startTime).Round(time.Second).String()

	response := map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now().Format("2006-01-02T15:04:05Z07:00"),
		"version":   "1.0.0",
		"uptime":    uptime,
		"message":   "Aplicação está viva e respondendo",
	}

	zap.L().Info("✅ Aplicação está viva", zap.String("uptime", uptime))

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// ReadinessCheck (Readiness Probe) - Verifica se a aplicação está pronta para receber tráfego
func (ctl *Controller) ReadinessCheck(ctx *fiber.Ctx) error {
	zap.L().Info("🚀 Verificação de prontidão da aplicação (readiness)")

	services := make(map[string]string)
	overallStatus := "ready"

	// Verificar conexão com banco master
	if db, err := config.NewDatabaseConnection(); err != nil {
		zap.L().Error("❌ Falha na conexão com banco master", zap.Error(err))
		services["database_master"] = "unhealthy"
		overallStatus = "not_ready"
	} else {
		if sqlDB, err := db.DB(); err != nil {
			services["database_master"] = "unhealthy"
			overallStatus = "not_ready"
		} else if err := sqlDB.Ping(); err != nil {
			services["database_master"] = "unhealthy"
			overallStatus = "not_ready"
		} else {
			services["database_master"] = "healthy"
			zap.L().Info("✅ Banco master conectado com sucesso")
		}
	}

	// Verificar conexões com bancos de clientes
	if clientDBs, err := config.ConnectionDBClients(); err != nil {
		zap.L().Error("❌ Falha na conexão com bancos de clientes", zap.Error(err))
		services["database_clients"] = "unhealthy"
		overallStatus = "not_ready"
	} else {
		healthyClients := 0
		totalClients := len(clientDBs)

		for key, db := range clientDBs {
			if sqlDB, err := db.DB(); err != nil {
				zap.L().Warn("⚠️ Falha ao obter conexão SQL", zap.String("cliente", key))
			} else if err := sqlDB.Ping(); err != nil {
				zap.L().Warn("⚠️ Falha no ping do banco", zap.String("cliente", key))
			} else {
				healthyClients++
			}
		}

		if healthyClients == totalClients {
			services["database_clients"] = "healthy"
			zap.L().Info("✅ Todos os bancos de clientes conectados", zap.Int("total", totalClients))
		} else {
			services["database_clients"] = "partial"
			if healthyClients == 0 {
				overallStatus = "not_ready"
			}
			zap.L().Warn("⚠️ Alguns bancos de clientes com problemas",
				zap.Int("saudaveis", healthyClients),
				zap.Int("total", totalClients))
		}
	}

	// Verificar Redis
	if err := config.RedisClient.Ping(ctx.Context()).Err(); err != nil {
		zap.L().Error("❌ Falha na conexão com Redis", zap.Error(err))
		services["redis"] = "unhealthy"
		overallStatus = "not_ready"
	} else {
		services["redis"] = "healthy"
		zap.L().Info("✅ Redis conectado com sucesso")
	}

	uptime := time.Since(startTime).Round(time.Second).String()

	response := ReadinessResponse{
		Status:    overallStatus,
		Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		Version:   "1.0.0",
		Services:  services,
		Uptime:    uptime,
	}

	statusCode := fiber.StatusOK
	if overallStatus == "not_ready" {
		statusCode = fiber.StatusServiceUnavailable
	}

	zap.L().Info("🚀 Verificação de prontidão concluída",
		zap.String("status", overallStatus),
		zap.String("uptime", uptime))

	return ctx.Status(statusCode).JSON(response)
}
