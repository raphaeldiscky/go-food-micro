// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	readmodels "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/readmodels"
	mock "github.com/stretchr/testify/mock"

	utils "github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	uuid "github.com/satori/go.uuid"
)

// OrderElasticRepository is an autogenerated mock type for the OrderElasticRepository type
type OrderElasticRepository struct {
	mock.Mock
}

type OrderElasticRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *OrderElasticRepository) EXPECT() *OrderElasticRepository_Expecter {
	return &OrderElasticRepository_Expecter{mock: &_m.Mock}
}

// CreateOrder provides a mock function with given fields: ctx, order
func (_m *OrderElasticRepository) CreateOrder(ctx context.Context, order *readmodels.OrderReadModel) (*readmodels.OrderReadModel, error) {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrder")
	}

	var r0 *readmodels.OrderReadModel
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *readmodels.OrderReadModel) (*readmodels.OrderReadModel, error)); ok {
		return rf(ctx, order)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *readmodels.OrderReadModel) *readmodels.OrderReadModel); ok {
		r0 = rf(ctx, order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*readmodels.OrderReadModel)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *readmodels.OrderReadModel) error); ok {
		r1 = rf(ctx, order)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderElasticRepository_CreateOrder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateOrder'
type OrderElasticRepository_CreateOrder_Call struct {
	*mock.Call
}

// CreateOrder is a helper method to define mock.On call
//   - ctx context.Context
//   - order *readmodels.OrderReadModel
func (_e *OrderElasticRepository_Expecter) CreateOrder(ctx interface{}, order interface{}) *OrderElasticRepository_CreateOrder_Call {
	return &OrderElasticRepository_CreateOrder_Call{Call: _e.mock.On("CreateOrder", ctx, order)}
}

func (_c *OrderElasticRepository_CreateOrder_Call) Run(run func(ctx context.Context, order *readmodels.OrderReadModel)) *OrderElasticRepository_CreateOrder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*readmodels.OrderReadModel))
	})
	return _c
}

func (_c *OrderElasticRepository_CreateOrder_Call) Return(_a0 *readmodels.OrderReadModel, _a1 error) *OrderElasticRepository_CreateOrder_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderElasticRepository_CreateOrder_Call) RunAndReturn(run func(context.Context, *readmodels.OrderReadModel) (*readmodels.OrderReadModel, error)) *OrderElasticRepository_CreateOrder_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteOrderByID provides a mock function with given fields: ctx, _a1
func (_m *OrderElasticRepository) DeleteOrderByID(ctx context.Context, _a1 uuid.UUID) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOrderByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrderElasticRepository_DeleteOrderByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteOrderByID'
type OrderElasticRepository_DeleteOrderByID_Call struct {
	*mock.Call
}

// DeleteOrderByID is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 uuid.UUID
func (_e *OrderElasticRepository_Expecter) DeleteOrderByID(ctx interface{}, _a1 interface{}) *OrderElasticRepository_DeleteOrderByID_Call {
	return &OrderElasticRepository_DeleteOrderByID_Call{Call: _e.mock.On("DeleteOrderByID", ctx, _a1)}
}

func (_c *OrderElasticRepository_DeleteOrderByID_Call) Run(run func(ctx context.Context, _a1 uuid.UUID)) *OrderElasticRepository_DeleteOrderByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *OrderElasticRepository_DeleteOrderByID_Call) Return(_a0 error) *OrderElasticRepository_DeleteOrderByID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrderElasticRepository_DeleteOrderByID_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *OrderElasticRepository_DeleteOrderByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllOrders provides a mock function with given fields: ctx, listQuery
func (_m *OrderElasticRepository) GetAllOrders(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*readmodels.OrderReadModel], error) {
	ret := _m.Called(ctx, listQuery)

	if len(ret) == 0 {
		panic("no return value specified for GetAllOrders")
	}

	var r0 *utils.ListResult[*readmodels.OrderReadModel]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *utils.ListQuery) (*utils.ListResult[*readmodels.OrderReadModel], error)); ok {
		return rf(ctx, listQuery)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *utils.ListQuery) *utils.ListResult[*readmodels.OrderReadModel]); ok {
		r0 = rf(ctx, listQuery)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.ListResult[*readmodels.OrderReadModel])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *utils.ListQuery) error); ok {
		r1 = rf(ctx, listQuery)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderElasticRepository_GetAllOrders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllOrders'
type OrderElasticRepository_GetAllOrders_Call struct {
	*mock.Call
}

// GetAllOrders is a helper method to define mock.On call
//   - ctx context.Context
//   - listQuery *utils.ListQuery
func (_e *OrderElasticRepository_Expecter) GetAllOrders(ctx interface{}, listQuery interface{}) *OrderElasticRepository_GetAllOrders_Call {
	return &OrderElasticRepository_GetAllOrders_Call{Call: _e.mock.On("GetAllOrders", ctx, listQuery)}
}

func (_c *OrderElasticRepository_GetAllOrders_Call) Run(run func(ctx context.Context, listQuery *utils.ListQuery)) *OrderElasticRepository_GetAllOrders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*utils.ListQuery))
	})
	return _c
}

func (_c *OrderElasticRepository_GetAllOrders_Call) Return(_a0 *utils.ListResult[*readmodels.OrderReadModel], _a1 error) *OrderElasticRepository_GetAllOrders_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderElasticRepository_GetAllOrders_Call) RunAndReturn(run func(context.Context, *utils.ListQuery) (*utils.ListResult[*readmodels.OrderReadModel], error)) *OrderElasticRepository_GetAllOrders_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrderByID provides a mock function with given fields: ctx, _a1
func (_m *OrderElasticRepository) GetOrderByID(ctx context.Context, _a1 uuid.UUID) (*readmodels.OrderReadModel, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetOrderByID")
	}

	var r0 *readmodels.OrderReadModel
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*readmodels.OrderReadModel, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *readmodels.OrderReadModel); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*readmodels.OrderReadModel)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderElasticRepository_GetOrderByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrderByID'
type OrderElasticRepository_GetOrderByID_Call struct {
	*mock.Call
}

// GetOrderByID is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 uuid.UUID
func (_e *OrderElasticRepository_Expecter) GetOrderByID(ctx interface{}, _a1 interface{}) *OrderElasticRepository_GetOrderByID_Call {
	return &OrderElasticRepository_GetOrderByID_Call{Call: _e.mock.On("GetOrderByID", ctx, _a1)}
}

func (_c *OrderElasticRepository_GetOrderByID_Call) Run(run func(ctx context.Context, _a1 uuid.UUID)) *OrderElasticRepository_GetOrderByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *OrderElasticRepository_GetOrderByID_Call) Return(_a0 *readmodels.OrderReadModel, _a1 error) *OrderElasticRepository_GetOrderByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderElasticRepository_GetOrderByID_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*readmodels.OrderReadModel, error)) *OrderElasticRepository_GetOrderByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrderByOrderID provides a mock function with given fields: ctx, orderID
func (_m *OrderElasticRepository) GetOrderByOrderID(ctx context.Context, orderID uuid.UUID) (*readmodels.OrderReadModel, error) {
	ret := _m.Called(ctx, orderID)

	if len(ret) == 0 {
		panic("no return value specified for GetOrderByOrderID")
	}

	var r0 *readmodels.OrderReadModel
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*readmodels.OrderReadModel, error)); ok {
		return rf(ctx, orderID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *readmodels.OrderReadModel); ok {
		r0 = rf(ctx, orderID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*readmodels.OrderReadModel)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, orderID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderElasticRepository_GetOrderByOrderID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrderByOrderID'
type OrderElasticRepository_GetOrderByOrderID_Call struct {
	*mock.Call
}

// GetOrderByOrderID is a helper method to define mock.On call
//   - ctx context.Context
//   - orderID uuid.UUID
func (_e *OrderElasticRepository_Expecter) GetOrderByOrderID(ctx interface{}, orderID interface{}) *OrderElasticRepository_GetOrderByOrderID_Call {
	return &OrderElasticRepository_GetOrderByOrderID_Call{Call: _e.mock.On("GetOrderByOrderID", ctx, orderID)}
}

func (_c *OrderElasticRepository_GetOrderByOrderID_Call) Run(run func(ctx context.Context, orderID uuid.UUID)) *OrderElasticRepository_GetOrderByOrderID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *OrderElasticRepository_GetOrderByOrderID_Call) Return(_a0 *readmodels.OrderReadModel, _a1 error) *OrderElasticRepository_GetOrderByOrderID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderElasticRepository_GetOrderByOrderID_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*readmodels.OrderReadModel, error)) *OrderElasticRepository_GetOrderByOrderID_Call {
	_c.Call.Return(run)
	return _c
}

// SearchOrders provides a mock function with given fields: ctx, searchText, listQuery
func (_m *OrderElasticRepository) SearchOrders(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*readmodels.OrderReadModel], error) {
	ret := _m.Called(ctx, searchText, listQuery)

	if len(ret) == 0 {
		panic("no return value specified for SearchOrders")
	}

	var r0 *utils.ListResult[*readmodels.OrderReadModel]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *utils.ListQuery) (*utils.ListResult[*readmodels.OrderReadModel], error)); ok {
		return rf(ctx, searchText, listQuery)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *utils.ListQuery) *utils.ListResult[*readmodels.OrderReadModel]); ok {
		r0 = rf(ctx, searchText, listQuery)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.ListResult[*readmodels.OrderReadModel])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *utils.ListQuery) error); ok {
		r1 = rf(ctx, searchText, listQuery)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderElasticRepository_SearchOrders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SearchOrders'
type OrderElasticRepository_SearchOrders_Call struct {
	*mock.Call
}

// SearchOrders is a helper method to define mock.On call
//   - ctx context.Context
//   - searchText string
//   - listQuery *utils.ListQuery
func (_e *OrderElasticRepository_Expecter) SearchOrders(ctx interface{}, searchText interface{}, listQuery interface{}) *OrderElasticRepository_SearchOrders_Call {
	return &OrderElasticRepository_SearchOrders_Call{Call: _e.mock.On("SearchOrders", ctx, searchText, listQuery)}
}

func (_c *OrderElasticRepository_SearchOrders_Call) Run(run func(ctx context.Context, searchText string, listQuery *utils.ListQuery)) *OrderElasticRepository_SearchOrders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*utils.ListQuery))
	})
	return _c
}

func (_c *OrderElasticRepository_SearchOrders_Call) Return(_a0 *utils.ListResult[*readmodels.OrderReadModel], _a1 error) *OrderElasticRepository_SearchOrders_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderElasticRepository_SearchOrders_Call) RunAndReturn(run func(context.Context, string, *utils.ListQuery) (*utils.ListResult[*readmodels.OrderReadModel], error)) *OrderElasticRepository_SearchOrders_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateOrder provides a mock function with given fields: ctx, order
func (_m *OrderElasticRepository) UpdateOrder(ctx context.Context, order *readmodels.OrderReadModel) (*readmodels.OrderReadModel, error) {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrder")
	}

	var r0 *readmodels.OrderReadModel
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *readmodels.OrderReadModel) (*readmodels.OrderReadModel, error)); ok {
		return rf(ctx, order)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *readmodels.OrderReadModel) *readmodels.OrderReadModel); ok {
		r0 = rf(ctx, order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*readmodels.OrderReadModel)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *readmodels.OrderReadModel) error); ok {
		r1 = rf(ctx, order)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderElasticRepository_UpdateOrder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateOrder'
type OrderElasticRepository_UpdateOrder_Call struct {
	*mock.Call
}

// UpdateOrder is a helper method to define mock.On call
//   - ctx context.Context
//   - order *readmodels.OrderReadModel
func (_e *OrderElasticRepository_Expecter) UpdateOrder(ctx interface{}, order interface{}) *OrderElasticRepository_UpdateOrder_Call {
	return &OrderElasticRepository_UpdateOrder_Call{Call: _e.mock.On("UpdateOrder", ctx, order)}
}

func (_c *OrderElasticRepository_UpdateOrder_Call) Run(run func(ctx context.Context, order *readmodels.OrderReadModel)) *OrderElasticRepository_UpdateOrder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*readmodels.OrderReadModel))
	})
	return _c
}

func (_c *OrderElasticRepository_UpdateOrder_Call) Return(_a0 *readmodels.OrderReadModel, _a1 error) *OrderElasticRepository_UpdateOrder_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderElasticRepository_UpdateOrder_Call) RunAndReturn(run func(context.Context, *readmodels.OrderReadModel) (*readmodels.OrderReadModel, error)) *OrderElasticRepository_UpdateOrder_Call {
	_c.Call.Return(run)
	return _c
}

// NewOrderElasticRepository creates a new instance of OrderElasticRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderElasticRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrderElasticRepository {
	mock := &OrderElasticRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
