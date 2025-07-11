// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"

	utils "github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	uuid "github.com/satori/go.uuid"
)

// ProductRepository is an autogenerated mock type for the ProductRepository type
type ProductRepository struct {
	mock.Mock
}

type ProductRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *ProductRepository) EXPECT() *ProductRepository_Expecter {
	return &ProductRepository_Expecter{mock: &_m.Mock}
}

// CreateProduct provides a mock function with given fields: ctx, product
func (_m *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	ret := _m.Called(ctx, product)

	if len(ret) == 0 {
		panic("no return value specified for CreateProduct")
	}

	var r0 *models.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Product) (*models.Product, error)); ok {
		return rf(ctx, product)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Product) *models.Product); ok {
		r0 = rf(ctx, product)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Product) error); ok {
		r1 = rf(ctx, product)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductRepository_CreateProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateProduct'
type ProductRepository_CreateProduct_Call struct {
	*mock.Call
}

// CreateProduct is a helper method to define mock.On call
//   - ctx context.Context
//   - product *models.Product
func (_e *ProductRepository_Expecter) CreateProduct(ctx interface{}, product interface{}) *ProductRepository_CreateProduct_Call {
	return &ProductRepository_CreateProduct_Call{Call: _e.mock.On("CreateProduct", ctx, product)}
}

func (_c *ProductRepository_CreateProduct_Call) Run(run func(ctx context.Context, product *models.Product)) *ProductRepository_CreateProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Product))
	})
	return _c
}

func (_c *ProductRepository_CreateProduct_Call) Return(_a0 *models.Product, _a1 error) *ProductRepository_CreateProduct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductRepository_CreateProduct_Call) RunAndReturn(run func(context.Context, *models.Product) (*models.Product, error)) *ProductRepository_CreateProduct_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteProductByID provides a mock function with given fields: ctx, _a1
func (_m *ProductRepository) DeleteProductByID(ctx context.Context, _a1 uuid.UUID) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProductByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProductRepository_DeleteProductByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteProductByID'
type ProductRepository_DeleteProductByID_Call struct {
	*mock.Call
}

// DeleteProductByID is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 uuid.UUID
func (_e *ProductRepository_Expecter) DeleteProductByID(ctx interface{}, _a1 interface{}) *ProductRepository_DeleteProductByID_Call {
	return &ProductRepository_DeleteProductByID_Call{Call: _e.mock.On("DeleteProductByID", ctx, _a1)}
}

func (_c *ProductRepository_DeleteProductByID_Call) Run(run func(ctx context.Context, _a1 uuid.UUID)) *ProductRepository_DeleteProductByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *ProductRepository_DeleteProductByID_Call) Return(_a0 error) *ProductRepository_DeleteProductByID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ProductRepository_DeleteProductByID_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *ProductRepository_DeleteProductByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllProducts provides a mock function with given fields: ctx, listQuery
func (_m *ProductRepository) GetAllProducts(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {
	ret := _m.Called(ctx, listQuery)

	if len(ret) == 0 {
		panic("no return value specified for GetAllProducts")
	}

	var r0 *utils.ListResult[*models.Product]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *utils.ListQuery) (*utils.ListResult[*models.Product], error)); ok {
		return rf(ctx, listQuery)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *utils.ListQuery) *utils.ListResult[*models.Product]); ok {
		r0 = rf(ctx, listQuery)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.ListResult[*models.Product])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *utils.ListQuery) error); ok {
		r1 = rf(ctx, listQuery)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductRepository_GetAllProducts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllProducts'
type ProductRepository_GetAllProducts_Call struct {
	*mock.Call
}

// GetAllProducts is a helper method to define mock.On call
//   - ctx context.Context
//   - listQuery *utils.ListQuery
func (_e *ProductRepository_Expecter) GetAllProducts(ctx interface{}, listQuery interface{}) *ProductRepository_GetAllProducts_Call {
	return &ProductRepository_GetAllProducts_Call{Call: _e.mock.On("GetAllProducts", ctx, listQuery)}
}

func (_c *ProductRepository_GetAllProducts_Call) Run(run func(ctx context.Context, listQuery *utils.ListQuery)) *ProductRepository_GetAllProducts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*utils.ListQuery))
	})
	return _c
}

func (_c *ProductRepository_GetAllProducts_Call) Return(_a0 *utils.ListResult[*models.Product], _a1 error) *ProductRepository_GetAllProducts_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductRepository_GetAllProducts_Call) RunAndReturn(run func(context.Context, *utils.ListQuery) (*utils.ListResult[*models.Product], error)) *ProductRepository_GetAllProducts_Call {
	_c.Call.Return(run)
	return _c
}

// GetProductByID provides a mock function with given fields: ctx, _a1
func (_m *ProductRepository) GetProductByID(ctx context.Context, _a1 uuid.UUID) (*models.Product, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetProductByID")
	}

	var r0 *models.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.Product, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.Product); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductRepository_GetProductByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProductByID'
type ProductRepository_GetProductByID_Call struct {
	*mock.Call
}

// GetProductByID is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 uuid.UUID
func (_e *ProductRepository_Expecter) GetProductByID(ctx interface{}, _a1 interface{}) *ProductRepository_GetProductByID_Call {
	return &ProductRepository_GetProductByID_Call{Call: _e.mock.On("GetProductByID", ctx, _a1)}
}

func (_c *ProductRepository_GetProductByID_Call) Run(run func(ctx context.Context, _a1 uuid.UUID)) *ProductRepository_GetProductByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *ProductRepository_GetProductByID_Call) Return(_a0 *models.Product, _a1 error) *ProductRepository_GetProductByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductRepository_GetProductByID_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*models.Product, error)) *ProductRepository_GetProductByID_Call {
	_c.Call.Return(run)
	return _c
}

// SearchProducts provides a mock function with given fields: ctx, searchText, listQuery
func (_m *ProductRepository) SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {
	ret := _m.Called(ctx, searchText, listQuery)

	if len(ret) == 0 {
		panic("no return value specified for SearchProducts")
	}

	var r0 *utils.ListResult[*models.Product]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *utils.ListQuery) (*utils.ListResult[*models.Product], error)); ok {
		return rf(ctx, searchText, listQuery)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *utils.ListQuery) *utils.ListResult[*models.Product]); ok {
		r0 = rf(ctx, searchText, listQuery)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.ListResult[*models.Product])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *utils.ListQuery) error); ok {
		r1 = rf(ctx, searchText, listQuery)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductRepository_SearchProducts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SearchProducts'
type ProductRepository_SearchProducts_Call struct {
	*mock.Call
}

// SearchProducts is a helper method to define mock.On call
//   - ctx context.Context
//   - searchText string
//   - listQuery *utils.ListQuery
func (_e *ProductRepository_Expecter) SearchProducts(ctx interface{}, searchText interface{}, listQuery interface{}) *ProductRepository_SearchProducts_Call {
	return &ProductRepository_SearchProducts_Call{Call: _e.mock.On("SearchProducts", ctx, searchText, listQuery)}
}

func (_c *ProductRepository_SearchProducts_Call) Run(run func(ctx context.Context, searchText string, listQuery *utils.ListQuery)) *ProductRepository_SearchProducts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*utils.ListQuery))
	})
	return _c
}

func (_c *ProductRepository_SearchProducts_Call) Return(_a0 *utils.ListResult[*models.Product], _a1 error) *ProductRepository_SearchProducts_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductRepository_SearchProducts_Call) RunAndReturn(run func(context.Context, string, *utils.ListQuery) (*utils.ListResult[*models.Product], error)) *ProductRepository_SearchProducts_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateProduct provides a mock function with given fields: ctx, product
func (_m *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	ret := _m.Called(ctx, product)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProduct")
	}

	var r0 *models.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Product) (*models.Product, error)); ok {
		return rf(ctx, product)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Product) *models.Product); ok {
		r0 = rf(ctx, product)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Product) error); ok {
		r1 = rf(ctx, product)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductRepository_UpdateProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateProduct'
type ProductRepository_UpdateProduct_Call struct {
	*mock.Call
}

// UpdateProduct is a helper method to define mock.On call
//   - ctx context.Context
//   - product *models.Product
func (_e *ProductRepository_Expecter) UpdateProduct(ctx interface{}, product interface{}) *ProductRepository_UpdateProduct_Call {
	return &ProductRepository_UpdateProduct_Call{Call: _e.mock.On("UpdateProduct", ctx, product)}
}

func (_c *ProductRepository_UpdateProduct_Call) Run(run func(ctx context.Context, product *models.Product)) *ProductRepository_UpdateProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Product))
	})
	return _c
}

func (_c *ProductRepository_UpdateProduct_Call) Return(_a0 *models.Product, _a1 error) *ProductRepository_UpdateProduct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductRepository_UpdateProduct_Call) RunAndReturn(run func(context.Context, *models.Product) (*models.Product, error)) *ProductRepository_UpdateProduct_Call {
	_c.Call.Return(run)
	return _c
}

// NewProductRepository creates a new instance of ProductRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProductRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProductRepository {
	mock := &ProductRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
