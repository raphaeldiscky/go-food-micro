// Package config provides a grpc options.
package config

import (
	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[GrpcOptions]())

// GrpcOptions is a struct that represents a grpc options.
type GrpcOptions struct {
	Port        string `mapstructure:"port"        env:"TcpPort"`
	Host        string `mapstructure:"host"        env:"Host"`
	Development bool   `mapstructure:"development" env:"Development"`
	Name        string `mapstructure:"name"        env:"ShortTypeName"`
}

// ProvideConfig is a function that provides a grpc options.
func ProvideConfig(environment environment.Environment) (*GrpcOptions, error) {
	return config.BindConfigKey[*GrpcOptions](optionName, environment)
}
