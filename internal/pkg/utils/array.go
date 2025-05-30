// Package utils provides a array utils.
package utils

import reflect "github.com/goccy/go-reflect"

// Contains is a function that checks if an element is in an array.
func Contains[T any](arr []T, x T) bool {
	for _, v := range arr {
		if reflect.ValueOf(v) == reflect.ValueOf(x) {
			return true
		}
	}

	return false
}

// ContainsFunc is a function that checks if an element is in an array.
func ContainsFunc[T any](arr []T, predicate func(T) bool) bool {
	for _, v := range arr {
		if predicate(v) {
			return true
		}
	}

	return false
}
