// Package cqrs provides a module for the cqrs.
package cqrs

// notification is a notification.
type notification struct {
	TypeInfo
}

// Notification is a notification.
type Notification interface {
	isNotification()

	TypeInfo
}

// NewNotificationByT creates a new notification by type.
func NewNotificationByT[T any]() Notification {
	n := &notification{
		TypeInfo: NewTypeInfoT[T](),
	}

	return n
}

// isNotification is a notification.
func (c *notification) isNotification() {
}

// IsNotification checks if the object is a notification.
func IsNotification(obj interface{}) bool {
	if _, ok := obj.(Notification); ok {
		return true
	}

	return false
}
