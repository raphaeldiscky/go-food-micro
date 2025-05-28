package types

import amqp091 "github.com/rabbitmq/amqp091-go"

type ExchangeType string

const (
	ExchangeFanout ExchangeType = amqp091.ExchangeFanout
	ExchangeDirect              = amqp091.ExchangeDirect
	ExchangeTopic               = amqp091.ExchangeTopic
)
