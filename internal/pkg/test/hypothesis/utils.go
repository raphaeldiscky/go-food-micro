// Package hypothesis provides a hypothesis utils.
package hypothesis

import "github.com/onsi/ginkgo"

// ForT creates a new hypothesis for a given type.
func ForT[T any](condition func(T) bool) Hypothesis[T] {
	return &hypothesis[T]{condition: condition, t: ginkgo.GinkgoT()}
}

// For creates a new hypothesis for a given type.
func For(_ interface{}, condition func(interface{}) bool) Hypothesis[interface{}] {
	return &hypothesis[interface{}]{condition: condition, t: ginkgo.GinkgoT()}
}
