// Package dtos contains the get order by id request dto.
package dtos

import uuid "github.com/satori/go.uuid"

// GetOrderByIDRequestDto is the request dto for the get order by id endpoint.
type GetOrderByIDRequestDto struct {
	ID uuid.UUID `param:"id" json:"-"`
}
