package dtos

import uuid "github.com/satori/go.uuid"

type GetProductByIDRequestDto struct {
	Id uuid.UUID `param:"id" json:"-"`
}
