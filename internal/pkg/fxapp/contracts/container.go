// Package contracts provides a module for the contracts.
package contracts

// Container is a contract for the container.
type Container interface {
	ResolveFunc(function interface{})
	ResolveFuncWithParamTag(function interface{}, paramTagName string)
}
