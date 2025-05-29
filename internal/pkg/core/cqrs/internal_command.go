// Package cqrs provides a module for the cqrs.
package cqrs

// internalCommand is a internal command.
type internalCommand struct {
	Command
}

// InternalCommand is a internal command.
type InternalCommand interface {
	Command
	isInternalCommand()
}

// NewInternalCommandByT creates a new internal command by type.
func NewInternalCommandByT[T any]() InternalCommand {
	return &internalCommand{Command: NewCommandByT[T]()}
}

// isInternalCommand is a internal command.
func (c *internalCommand) isInternalCommand() {
}

// IsInternalCommand checks if the object is a internal command.
func IsInternalCommand(obj interface{}) bool {
	if _, ok := obj.(InternalCommand); ok {
		return true
	}

	return false
}
