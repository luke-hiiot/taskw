//go:build wireinject
// +build wireinject

package api

import (
	"{{.Module}}/internal/pkg/config"
	"{{.Module}}/internal/pkg/logger"

	"github.com/google/wire"
)

// ProviderSet contains infrastructure providers
var ProviderSet = wire.NewSet(
	// config
	config.ProvideConfig,

	// logger
	logger.NewDevelopment,
	logger.NewDevelopmentLogger,

	// generated providers
	GeneratedProviderSet,
)

// InitializeRouter initializes Router dependencies
func InitializeRouter() (*Router, error) {
	wire.Build(ProviderSet)
	return &Router{}, nil
}
