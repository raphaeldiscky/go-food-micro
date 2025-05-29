// Package cqrs provides a module for the cqrs.
package cqrs

// TxRequest is a tx request.
// https://www.mohitkhare.com/blog/go-naming-conventions/
// https://github.com/EventStore/EventStore-Client-Go/blob/master/esdb/position.go
type TxRequest interface {
	Request

	isTxRequest()
}
