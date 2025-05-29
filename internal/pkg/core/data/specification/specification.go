// Package specification provides a module for the specification.
package specification

import (
	"fmt"
	"strings"
)

// Specification is a specification.
type Specification interface {
	GetQuery() string
	GetValues() []any
}

// joinSpecification is a join specification.
type joinSpecification struct {
	specifications []Specification
	separator      string
}

// GetQuery gets the query.
func (s joinSpecification) GetQuery() string {
	queries := make([]string, 0, len(s.specifications))

	for _, spec := range s.specifications {
		queries = append(queries, spec.GetQuery())
	}

	return strings.Join(queries, fmt.Sprintf(" %s ", s.separator))
}

// GetValues gets the values.
func (s joinSpecification) GetValues() []any {
	values := make([]any, 0)

	for _, spec := range s.specifications {
		values = append(values, spec.GetValues()...)
	}

	return values
}

// And joins specifications with AND.
func And(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		separator:      "AND",
	}
}

// Or joins specifications with OR.
func Or(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		separator:      "OR",
	}
}

// notSpecification is a not specification.
type notSpecification struct {
	Specification
}

// GetQuery gets the query.
func (s notSpecification) GetQuery() string {
	return fmt.Sprintf(" NOT (%s)", s.Specification.GetQuery())
}

// Not negates a specification.
func Not(specification Specification) Specification {
	return notSpecification{
		specification,
	}
}

// binaryOperatorSpecification is a binary operator specification.
type binaryOperatorSpecification[T any] struct {
	field    string
	operator string
	value    T
}

// GetQuery gets the query.
func (s binaryOperatorSpecification[T]) GetQuery() string {
	return fmt.Sprintf("%s %s ?", s.field, s.operator)
}

// GetValues gets the values.
func (s binaryOperatorSpecification[T]) GetValues() []any {
	return []any{s.value}
}

// Equal creates a new equal specification.
func Equal[T any](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "=",
		value:    value,
	}
}

// GreaterThan creates a new greater than specification.
func GreaterThan[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: ">",
		value:    value,
	}
}

// GreaterOrEqual creates a new greater or equal specification.
func GreaterOrEqual[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: ">=",
		value:    value,
	}
}

// LessThan creates a new less than specification.
func LessThan[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "<",
		value:    value,
	}
}

// LessOrEqual creates a new less or equal specification.
func LessOrEqual[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: ">=",
		value:    value,
	}
}

// stringSpecification is a string specification.
type stringSpecification string

// GetQuery gets the query.
func (s stringSpecification) GetQuery() string {
	return string(s)
}

// GetValues gets the values.
func (s stringSpecification) GetValues() []any {
	return nil
}

// IsNull creates a new is null specification.
func IsNull(field string) Specification {
	return stringSpecification(fmt.Sprintf("%s IS NULL", field))
}
