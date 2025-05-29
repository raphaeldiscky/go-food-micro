// Package cqrs provides a module for the cqrs.
package cqrs

// query is a query.
type query struct {
	TypeInfo
	Request
}

// Query is a query.
type Query interface {
	isQuery()

	Request
	TypeInfo
}

// NewQueryByT creates a new query by type.
func NewQueryByT[T any]() Query {
	return &query{
		TypeInfo: NewTypeInfoT[T](),
		Request:  NewRequest(),
	}
}

// isQuery is a query.
// https://github.com/EventStore/EventStore-Client-Go/blob/master/esdb/position.go#L29
func (q *query) isQuery() {
}

// IsQuery checks if the object is a query.
func IsQuery(obj interface{}) bool {
	if _, ok := obj.(Query); ok {
		return true
	}

	return false
}
