package dtos

import uuid "github.com/satori/go.uuid"

type GetOrderByIdRequestDto struct {
	ID uuid.UUID `param:"id" json:"-"`
}
