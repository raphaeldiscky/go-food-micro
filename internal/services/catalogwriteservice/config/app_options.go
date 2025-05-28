// Package config contains the app options.
package config

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"

	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// AppOptions is a struct that contains the app options.
type AppOptions struct {
	DeliveryType string `mapstructure:"deliveryType" env:"DeliveryType"`
	ServiceName  string `mapstructure:"serviceName"  env:"serviceName"`
}

// NewAppOptions is a constructor for the AppOptions.
func NewAppOptions(env environment.Environment) (*AppOptions, error) {
	optionName := strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[AppOptions]())
	cfg, err := config.BindConfigKey[*AppOptions](optionName, env)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetMicroserviceNameUpper is a method that returns the microservice name in uppercase.
func (cfg *AppOptions) GetMicroserviceNameUpper() string {
	return strings.ToUpper(cfg.ServiceName)
}

// GetMicroserviceName is a method that returns the microservice name.
func (cfg *AppOptions) GetMicroserviceName() string {
	return cfg.ServiceName
}
