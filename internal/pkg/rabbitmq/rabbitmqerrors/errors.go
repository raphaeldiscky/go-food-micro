// Package rabbitmqerrors provides a set of custom error types for RabbitMQ.
package rabbitmqerrors

import (
	"emperror.dev/errors"
)

// ErrDisconnected is a error that represents a disconnected from rabbitmq, trying to reconnect.
var ErrDisconnected = errors.New("disconnected from rabbitmq, trying to reconnect")
