// Package dtos contains the delete product request dto.
package dtos

import uuid "github.com/satori/go.uuid"

// DeleteProductRequestDto is a struct that contains the delete product request dto.
type DeleteProductRequestDto struct {
	ProductID uuid.UUID `param:"id" json:"-"`
}
