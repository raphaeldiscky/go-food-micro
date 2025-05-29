// Package config provides a set of functions for the rabbitmq options.
package config

import (
	"fmt"
	"time"

	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// RabbitmqOptions is a struct that contains the rabbitmq options.
type RabbitmqOptions struct {
	RabbitmqHostOptions *RabbitmqHostOptions `mapstructure:"rabbitmqHostOptions"`
	DeliveryMode        uint8
	Persisted           bool
	AppId               string
	AutoStart           bool `mapstructure:"autoStart"           default:"true"`
	Reconnecting        bool `mapstructure:"reconnecting"        default:"true"`
}

// RabbitmqHostOptions is a struct that contains the rabbitmq host options.
type RabbitmqHostOptions struct {
	HostName    string    `mapstructure:"hostName"`
	VirtualHost string    `mapstructure:"virtualHost"`
	Port        int       `mapstructure:"port"`
	HttpPort    int       `mapstructure:"httpPort"`
	UserName    string    `mapstructure:"userName"`
	Password    string    `mapstructure:"password"`
	RetryDelay  time.Time `mapstructure:"retryDelay"`
}

// AmqpEndPoint returns the amqp endpoint.
func (h *RabbitmqHostOptions) AmqpEndPoint() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d", h.UserName, h.Password, h.HostName, h.Port)
}

// HttpEndPoint returns the http endpoint.
func (h *RabbitmqHostOptions) HttpEndPoint() string {
	return fmt.Sprintf("http://%s:%d", h.HostName, h.HttpPort)
}

// ProvideConfig provides the rabbitmq options.
func ProvideConfig(environment environment.Environment) (*RabbitmqOptions, error) {
	optionName := strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[RabbitmqOptions]())
	cfg, err := config.BindConfigKey[*RabbitmqOptions](optionName, environment)

	return cfg, err
}
