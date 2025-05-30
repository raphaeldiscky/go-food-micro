// Package cqrs provides a module for the cqrs.
package cqrs

import (
	"reflect"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// TypeInfo is a type info.
type TypeInfo interface {
	ShortTypeName() string
	FullTypeName() string
	Type() reflect.Type
}

// typeInfo is a type info.
type typeInfo struct {
	shortTypeName string
	fullTypeName  string
	typ           reflect.Type
}

// NewTypeInfoT creates a new type info by type.
func NewTypeInfoT[T any]() TypeInfo {
	name := typemapper.GetGenericTypeNameByT[T]()
	fullName := typemapper.GetGenericFullTypeNameByT[T]()
	typ := typemapper.GetGenericTypeByT[T]()

	return &typeInfo{fullTypeName: fullName, typ: typ, shortTypeName: name}
}

// ShortTypeName gets the short type name.
func (t *typeInfo) ShortTypeName() string {
	return t.shortTypeName
}

// FullTypeName gets the full type name.
func (t *typeInfo) FullTypeName() string {
	return t.fullTypeName
}

// Type gets the type.
func (t *typeInfo) Type() reflect.Type {
	return t.typ
}
