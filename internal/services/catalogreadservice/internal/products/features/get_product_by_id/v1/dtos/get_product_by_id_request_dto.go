package dtos

import uuid "github.com/satori/go.uuid"

type GetProductByIDRequestDto struct {
	ID uuid.UUID `param:"id" json:"-"`
}
