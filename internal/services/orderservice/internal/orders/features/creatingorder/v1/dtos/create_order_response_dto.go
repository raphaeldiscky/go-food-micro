// Package dtos contains the create order response dto.
package dtos

import uuid "github.com/satori/go.uuid"

// https://echo.labstack.com/guide/response/

// CreateOrderResponseDto is the response dto for the create order command.
type CreateOrderResponseDto struct {
	OrderID uuid.UUID `json:"ID"`
}
