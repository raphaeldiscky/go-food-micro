// Package utils provides utils for the core package.
package utils

import (
	"reflect"

	linq "github.com/ahmetb/go-linq/v3"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/events"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// GetAllDomainEventTypes is a function that returns all domain event types.
func GetAllDomainEventTypes() []reflect.Type {
	var types []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).
		SelectManyT(func(i linq.KeyValue) linq.Query {
			return linq.From(i.Value)
		})
	d.ToSlice(&types)
	res := typeMapper.TypesImplementedInterfaceWithFilterTypes[domain.IDomainEvent](types)
	linq.From(res).Distinct().ToSlice(&types)

	return types
}

// GetAllEventTypes is a function that returns all event types.
func GetAllEventTypes() []reflect.Type {
	var types []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).
		SelectManyT(func(i linq.KeyValue) linq.Query {
			return linq.From(i.Value)
		})
	d.ToSlice(&types)
	res := typeMapper.TypesImplementedInterfaceWithFilterTypes[events.IEvent](types)
	linq.From(res).Distinct().ToSlice(&types)

	return types
}
