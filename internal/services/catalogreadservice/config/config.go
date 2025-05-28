// Package config contains the config for the catalog read service.
package config

import (
	"strings"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
)

// Config is a struct that contains the config for the catalog read service.
type Config struct {
	AppOptions AppOptions `mapstructure:"appOptions" env:"AppOptions"`
}

// NewConfig creates a new Config.
func NewConfig(env environment.Environment) (*Config, error) {
	cfg, err := config.BindConfig[*Config](env)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// AppOptions is a struct that contains the app options for the catalog read service.
type AppOptions struct {
	DeliveryType string `mapstructure:"deliveryType" env:"DeliveryType"`
	ServiceName  string `mapstructure:"serviceName"  env:"serviceName"`
}

// GetMicroserviceNameUpper returns the microservice name in uppercase.
func (cfg *AppOptions) GetMicroserviceNameUpper() string {
	return strings.ToUpper(cfg.ServiceName)
}

// GetMicroserviceName returns the microservice name.
func (cfg *AppOptions) GetMicroserviceName() string {
	return cfg.ServiceName
}
