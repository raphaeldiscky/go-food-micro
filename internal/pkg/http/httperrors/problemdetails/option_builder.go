// Package problemdetails provides problem details.
package problemdetails

import "reflect"

// OptionBuilder is a struct that represents a option builder.
type OptionBuilder struct {
	internalErrors map[reflect.Type]func(err error) ProblemDetailErr
}

// NewOptionBuilder creates a new option builder.
func NewOptionBuilder() *OptionBuilder {
	return &OptionBuilder{}
}

// Map maps an error type to a problem detail function.
func (p *OptionBuilder) Map(
	srcErrorType reflect.Type,
	problem ProblemDetailFunc[error],
) *OptionBuilder {
	internalErrorMaps[srcErrorType] = problem

	return p
}

// Build builds the option builder.
func (p *OptionBuilder) Build() map[reflect.Type]func(err error) ProblemDetailErr {
	return p.internalErrors
}
