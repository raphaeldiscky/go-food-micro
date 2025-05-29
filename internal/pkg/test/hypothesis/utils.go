// Package hypothesis provides a hypothesis utils.
package hypothesis

// ForT creates a new hypothesis for a given type.
func ForT[T any](condition func(T) bool) Hypothesis[T] {
	return &hypothesis[T]{condition: condition}
}

// For creates a new hypothesis for a given type.
func For(typ interface{}, condition func(interface{}) bool) Hypothesis[interface{}] {
	return &hypothesis[interface{}]{condition: condition}
}
