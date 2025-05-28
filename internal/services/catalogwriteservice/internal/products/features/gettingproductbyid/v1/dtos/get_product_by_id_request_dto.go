package dtos

import uuid "github.com/satori/go.uuid"

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// GetProductByIDRequestDto is a struct that contains the get product by id request dto.
type GetProductByIDRequestDto struct {
	ProductId uuid.UUID `param:"id" json:"-"`
}
