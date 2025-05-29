// Package cqrs provides a module for the cqrs.
package cqrs

// request is a request.
type request struct{}

// Request is a request.
type Request interface {
	isRequest()
}

// NewRequest creates a new request.
func NewRequest() Request {
	return &request{}
}

// isRequest is a request.
func (r *request) isRequest() {
}

// IsRequest checks if the object is a request.
func IsRequest(obj interface{}) bool {
	if _, ok := obj.(Request); ok {
		return true
	}

	return false
}
