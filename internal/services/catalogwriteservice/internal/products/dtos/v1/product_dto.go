package v1

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// ProductDto is a struct that contains the product dto.
type ProductDto struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
