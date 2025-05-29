// Package config provides a echo http options.
package config

import (
	"fmt"
	"net/url"

	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[EchoHTTPOptions]())

// EchoHTTPOptions is a struct that represents a echo http options.
type EchoHTTPOptions struct {
	Port                string   `mapstructure:"port"                validate:"required" env:"TcpPort"`
	Development         bool     `mapstructure:"development"                             env:"Development"`
	BasePath            string   `mapstructure:"basePath"            validate:"required" env:"BasePath"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"                     env:"DebugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout"                                 env:"Timeout"`
	Host                string   `mapstructure:"host"                                    env:"Host"`
	Name                string   `mapstructure:"name"                                    env:"ShortTypeName"`
}

// Address is a function that returns the address.
func (c *EchoHTTPOptions) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}

// BasePathAddress is a function that returns the base path address.
func (c *EchoHTTPOptions) BasePathAddress() string {
	path, err := url.JoinPath(c.Address(), c.BasePath)
	if err != nil {
		return ""
	}

	return path
}

// ProvideConfig is a function that provides the config.
func ProvideConfig(environment environment.Environment) (*EchoHTTPOptions, error) {
	return config.BindConfigKey[*EchoHTTPOptions](optionName, environment)
}
