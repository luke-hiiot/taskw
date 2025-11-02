package api

import (
	"sync"
	"time"

	"{{.Module}}/internal/pkg/config"

	"github.com/gofiber/fiber/v2"
)

var (
	appInstance *fiber.App
	appOnce     sync.Once
)

// ProvideFiberApp provides Fiber application singleton instance
func ProvideFiberApp(cfg *config.Config) *fiber.App {
	appOnce.Do(func() {
		appInstance = newFiberApp(cfg)
	})
	return appInstance
}

// newFiberApp builds Fiber application based on configuration
func newFiberApp(cfg *config.Config) *fiber.App {
	var (
		appName      string
		readTimeout  time.Duration
		writeTimeout time.Duration
		idleTimeout  time.Duration
	)

	if cfg != nil {
		if cfg.App.Name != "" {
			appName = cfg.App.Name
		}
		if cfg.Server.ReadTimeoutSec > 0 {
			readTimeout = time.Duration(cfg.Server.ReadTimeoutSec) * time.Second
		}
		if cfg.Server.WriteTimeoutSec > 0 {
			writeTimeout = time.Duration(cfg.Server.WriteTimeoutSec) * time.Second
		}
		if cfg.Server.IdleTimeoutSec > 0 {
			idleTimeout = time.Duration(cfg.Server.IdleTimeoutSec) * time.Second
		}
	}

	return fiber.New(fiber.Config{
		AppName:      appName,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
}

func (ar *Router) GetApp() *fiber.App {
	return ar.app
}
