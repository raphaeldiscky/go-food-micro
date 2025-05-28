package queries

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	uuid "github.com/satori/go.uuid"
)

type GetProductByID struct {
	Id uuid.UUID
}

func NewGetProductByID(id uuid.UUID) (*GetProductByID, error) {
	product := &GetProductByID{Id: id}
	if err := product.Validate(); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *GetProductByID) Validate() error {
	return validation.ValidateStruct(p, validation.Field(&p.Id, validation.Required, is.UUIDv4))
}
