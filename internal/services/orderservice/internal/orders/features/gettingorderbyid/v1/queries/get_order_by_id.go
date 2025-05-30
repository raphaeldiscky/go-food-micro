// Package queries contains the queries for the get order by id.
package queries

import (
	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

// GetOrderByID is the query for the get order by id.
type GetOrderByID struct {
	ID uuid.UUID
}

// NewGetOrderByID creates a new get order by id query.
func NewGetOrderByID(id uuid.UUID) (*GetOrderByID, error) {
	query := &GetOrderByID{ID: id}

	err := query.Validate()
	if err != nil {
		return nil, err
	}

	return query, nil
}

// Validate validates the get order by id query.
func (g GetOrderByID) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.ID, validation.Required),
	)
}
