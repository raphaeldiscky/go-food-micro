// Package cqrs provides a module for the cqrs.
package cqrs

// command is a command.
type command struct {
	TypeInfo
	Request
}

// Command is a command.
type Command interface {
	isCommand()

	Request
	TypeInfo
}

// NewCommandByT creates a new command by type.
func NewCommandByT[T any]() Command {
	c := &command{
		TypeInfo: NewTypeInfoT[T](),
		Request:  NewRequest(),
	}

	return c
}

// isCommand is a command.
// https://github.com/EventStore/EventStore-Client-Go/blob/master/esdb/position.go#L29
func (c *command) isCommand() {
}

// IsCommand checks if the object is a command.
func IsCommand(obj interface{}) bool {
	if _, ok := obj.(Command); ok {
		return true
	}

	return false
}
