// Package rabbitmqerrors provides a set of custom error types for RabbitMQ.
package rabbitmqerrors

import (
	"emperror.dev/errors"
)

var ErrDisconnected = errors.New("disconnected from rabbitmq, trying to reconnect")
