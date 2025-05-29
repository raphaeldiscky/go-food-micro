// Package validation provides a validation.
package validation

// Validator is a interface that represents a validator.
type Validator interface {
	Validate() error
}
