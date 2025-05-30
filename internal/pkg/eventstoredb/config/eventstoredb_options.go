// Package config provides a event store db options.
package config

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// https://developers.eventstore.com/clients/dotnet/21.2/#connect-to-eventstoredb
// https://developers.eventstore.com/clients/http-api/v5
// https://developers.eventstore.com/clients/grpc/
// https://developers.eventstore.com/server/v20.10/networking.html#http-configuration

// EventStoreDbOptions is a struct that represents the event store db options.
type EventStoreDbOptions struct {
	Host    string `mapstructure:"host"`
	TcpPort int    `mapstructure:"tcpPort"`
	// HTTP is the primary protocol for EventStoreDB. It is used in gRPC communication and HTTP APIs (management, gossip and diagnostics).
	HttpPort     int           `mapstructure:"httpPort"`
	Subscription *Subscription `mapstructure:"subscription"`
}

// https://developers.eventstore.com/server/v20.10/networking.html#http-configuration
// https://developers.eventstore.com/clients/grpc/#connection-string

// GrpcEndPoint returns the grpc end point.
func (e *EventStoreDbOptions) GrpcEndPoint() string {
	return fmt.Sprintf("esdb://%s:%d?tls=false", e.Host, e.HttpPort)
}

// https://developers.eventstore.com/clients/dotnet/21.2/#connect-to-eventstoredb
// https://developers.eventstore.com/server/v20.10/networking.html#external

// TCPEndPoint returns the tcp end point.
func (e *EventStoreDbOptions) TCPEndPoint() string {
	return fmt.Sprintf("tcp://%s:%d?tls=false", e.Host, e.TcpPort)
}

// https://developers.eventstore.com/server/v20.10/networking.html#http-configuration
// https://developers.eventstore.com/clients/http-api/v5

// HTTPEndPoint returns the http end point.
func (e *EventStoreDbOptions) HTTPEndPoint() string {
	return fmt.Sprintf("http://%s:%d", e.Host, e.HttpPort)
}

// Subscription is a struct that represents the subscription.
type Subscription struct {
	Prefix         []string `mapstructure:"prefix"         validate:"required"`
	SubscriptionId string   `mapstructure:"subscriptionId" validate:"required"`
}

// ProvideConfig provides the event store db options.
func ProvideConfig(environment environment.Environment) (*EventStoreDbOptions, error) {
	optionName := strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[EventStoreDbOptions]())

	return config.BindConfigKey[*EventStoreDbOptions](optionName, environment)
}
