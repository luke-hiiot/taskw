package health

import (
	"context"
	"runtime"
	"time"

	"{{.Module}}/internal/middleware"
	"{{.Module}}/internal/pkg/config"
	"{{.Module}}/internal/pkg/db"

	"github.com/gofiber/fiber/v2"
)

// Handler handles health check requests
type Handler struct{}

// 记录应用启动时间
var startTime = time.Now()

// ProvideHandler creates a new health handler
func ProvideHandler() *Handler {
	return &Handler{}
}

// @Summary Health check
// @Description Get the health status of the API with system information
// @Tags System Management
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/system/health [get]
func (h *Handler) GetHealth(c *fiber.Ctx) error {
	return middleware.SuccessResponse(c, fiber.Map{
		"status":     "healthy",
		"timestamp":  time.Now().Format(time.RFC3339),
		"uptime":     time.Since(startTime).String(),
		"goroutines": runtime.NumGoroutine(),
		"memory_mb":  h.getMemoryMB(),
		"database":   h.getDatabaseStatus(),
		"app":        h.getAppInfo(),
	})
}

// getMemoryMB 获取内存使用量(MB)
func (h *Handler) getMemoryMB() int {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int(m.Alloc / 1024 / 1024)
}

// getDatabaseStatus 获取数据库状态
func (h *Handler) getDatabaseStatus() fiber.Map {
	client := db.GetDBClient()
	if client == nil {
		return fiber.Map{
			"status": "disconnected",
			"error":  "client not initialized",
		}
	}

	// 测试数据库连接
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	_, err := client.User.Query().Count(ctx)
	duration := time.Since(start)

	if err != nil {
		return fiber.Map{
			"status":   "error",
			"error":    err.Error(),
			"duration": duration.Milliseconds(),
		}
	}

	// 获取连接池信息
	cfg := config.GetConfig()
	poolInfo := fiber.Map{}
	if cfg != nil {
		// 获取当前连接池统计信息
		stats := db.GetDBStats()
		poolInfo = fiber.Map{
			"max_open":      cfg.DB.MaxOpen,
			"max_idle":      cfg.DB.MaxIdle,
			"max_lifetime":  cfg.DB.ConnMaxLifetimeSec,
			"max_idle_time": cfg.DB.ConnMaxIdleTimeSec,
		}

		if stats != nil {
			poolInfo["current"] = fiber.Map{
				"open":          stats.OpenConnections,
				"in_use":        stats.InUse,
				"idle":          stats.Idle,
				"wait":          stats.WaitCount,
				"wait_duration": stats.WaitDuration.Milliseconds(),
			}
		}
	}

	return fiber.Map{
		"status":   "connected",
		"duration": duration.Milliseconds(),
		"pool":     poolInfo,
	}
}

// getAppInfo 获取应用信息
func (h *Handler) getAppInfo() fiber.Map {
	cfg := config.GetConfig()

	info := fiber.Map{
		"name":    "unknown",
		"version": "unknown",
		"env":     "unknown",
	}

	if cfg != nil {
		info["name"] = cfg.App.Name
		info["version"] = cfg.App.Version
		info["env"] = cfg.App.Env
	}

	return info
}
