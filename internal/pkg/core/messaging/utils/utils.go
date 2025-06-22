// Package utils provides utils for messaging.
package utils

import (
	"reflect"

	"github.com/iancoleman/strcase"

	linq "github.com/ahmetb/go-linq/v3"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// GetMessageName is a function that returns the message name.
func GetMessageName(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}

	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

// GetMessageNameFromType is a function that returns the message name from type.
func GetMessageNameFromType(message reflect.Type) string {
	if message.Kind() == reflect.Pointer {
		return strcase.ToSnake(message.Elem().Name())
	}

	return strcase.ToSnake(message.Name())
}

// GetMessageBaseReflectTypeFromType is a function that returns the message base reflect type from type.
func GetMessageBaseReflectTypeFromType(message reflect.Type) reflect.Type {
	return typeMapper.GetBaseReflectType(message)
}

// GetMessageBaseReflectType is a function that returns the message base reflect type.
func GetMessageBaseReflectType(message interface{}) reflect.Type {
	return typeMapper.GetBaseReflectType(message)
}

// GetTopicOrExchangeName is a function that returns the topic or exchange name.
func GetTopicOrExchangeName(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}

	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

// GetTopicOrExchangeNameFromType is a function that returns the topic or exchange name from type.
func GetTopicOrExchangeNameFromType(message reflect.Type) string {
	if message.Kind() == reflect.Pointer {
		return strcase.ToSnake(message.Elem().Name())
	}

	return strcase.ToSnake(message.Name())
}

// GetQueueName is a function that returns the queue name.
func GetQueueName(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}

	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

// GetQueueNameFromType is a function that returns the queue name from type.
func GetQueueNameFromType(message reflect.Type) string {
	if message.Kind() == reflect.Pointer {
		return strcase.ToSnake(message.Elem().Name())
	}

	return strcase.ToSnake(message.Name())
}

// GetRoutingKey is a function that returns the routing key.
func GetRoutingKey(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}

	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

// GetRoutingKeyFromType is a function that returns the routing key from type.
func GetRoutingKeyFromType(message reflect.Type) string {
	if message.Kind() == reflect.Pointer {
		return strcase.ToSnake(message.Elem().Name())
	}

	return strcase.ToSnake(message.Name())
}

// RegisterCustomMessageTypesToRegistry is a function that registers custom message types to registrty.
func RegisterCustomMessageTypesToRegistry(typesMap map[string]types.IMessage) {
	if len(typesMap) == 0 {
		return
	}

	for k, v := range typesMap {
		typeMapper.RegisterTypeWithKey(k, typeMapper.GetReflectType(v))
	}
}

// GetAllMessageTypes is a function that returns all message types.
func GetAllMessageTypes() []reflect.Type {
	var squares []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).
		SelectManyT(func(i linq.KeyValue) linq.Query {
			return linq.From(i.Value)
		})
	d.ToSlice(&squares)
	res := typeMapper.TypesImplementedInterfaceWithFilterTypes[types.IMessage](squares)
	linq.From(res).Distinct().ToSlice(&squares)

	return squares
}
