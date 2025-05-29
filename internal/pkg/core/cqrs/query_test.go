package cqrs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// TestQuery tests the query.
func TestQuery(t *testing.T) {
	query := &GetProductByID{
		Query:     NewQueryByT[*GetProductByID](),
		ProductID: uuid.NewV4(),
	}

	isImplementedQuery := typemapper.ImplementedInterfaceT[Query](query)
	assert.True(t, isImplementedQuery)

	var i interface{} = query
	_, isQuery := i.(Query)
	_, isTypeInfo := i.(TypeInfo)
	_, isCommand := i.(Command)
	_, isRequest := i.(Request)

	assert.True(t, isQuery)
	assert.False(t, isCommand)
	assert.True(t, isTypeInfo)
	assert.True(t, isRequest)

	assert.True(t, IsQuery(query))
	assert.False(t, IsCommand(query))
	assert.True(t, IsRequest(query))

	assert.Equal(t, query.ShortTypeName(), "*GetProductByID")
	assert.Equal(t, query.FullTypeName(), "*cqrs.GetProductByID")
}

// GetProductByID is a struct that represents a get product by id.
type GetProductByID struct {
	Query

	ProductID uuid.UUID
}
