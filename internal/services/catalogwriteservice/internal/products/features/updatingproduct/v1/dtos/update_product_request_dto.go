// Package dtos contains the update product request dto.
package dtos

import uuid "github.com/satori/go.uuid"

// https://echo.labstack.com/guide/binding/

// UpdateProductRequestDto is a struct that contains the update product request dto.
type UpdateProductRequestDto struct {
	ProductID   uuid.UUID `json:"-"           param:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
}
