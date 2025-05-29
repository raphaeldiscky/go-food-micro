// Package contracts provides a rabbitmq container contracts.
package contracts

import (
	"context"
	"fmt"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
)

// RabbitMQContainerOptions represents a rabbitmq container options.
type RabbitMQContainerOptions struct {
	Host        string
	VirtualHost string
	Ports       []string
	HostPort    int
	HttpPort    int
	UserName    string
	Password    string
	ImageName   string
	Name        string
	Tag         string
}

// AmqpEndPoint returns the amqp endpoint.
func (h *RabbitMQContainerOptions) AmqpEndPoint() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d",
		h.UserName,
		h.Password,
		h.Host,
		h.HostPort,
	)
}

// HttpEndPoint returns the http endpoint.
func (h *RabbitMQContainerOptions) HttpEndPoint() string {
	return fmt.Sprintf("http://%s:%d", h.Host, h.HttpPort)
}

// RabbitMQContainer is a interface that represents a rabbitmq container.
type RabbitMQContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*RabbitMQContainerOptions,
	) (*config.RabbitmqHostOptions, error)

	Cleanup(ctx context.Context) error
}
