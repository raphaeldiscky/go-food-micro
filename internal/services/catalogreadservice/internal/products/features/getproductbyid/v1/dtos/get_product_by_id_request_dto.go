package dtos

import uuid "github.com/satori/go.uuid"

// GetProductByIDRequestDto is a struct that contains the get product by id request dto.
type GetProductByIDRequestDto struct {
	ID uuid.UUID `param:"id" json:"-"`
}
