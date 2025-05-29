// Package types provides a set of types for the rabbitmq package.
package types

import amqp091 "github.com/rabbitmq/amqp091-go"

// ExchangeType is a type that represents the type of exchange.
type ExchangeType string

// ExchangeType constants.
const (
	ExchangeFanout ExchangeType = amqp091.ExchangeFanout
	ExchangeDirect              = amqp091.ExchangeDirect
	ExchangeTopic               = amqp091.ExchangeTopic
)
