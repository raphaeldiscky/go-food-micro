// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	types "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	mock "github.com/stretchr/testify/mock"
)

// BusControl is an autogenerated mock type for the BusControl type
type BusControl struct {
	mock.Mock
}

type BusControl_Expecter struct {
	mock *mock.Mock
}

func (_m *BusControl) EXPECT() *BusControl_Expecter {
	return &BusControl_Expecter{mock: &_m.Mock}
}

// IsConsumed provides a mock function with given fields: _a0
func (_m *BusControl) IsConsumed(_a0 func(types.IMessage)) {
	_m.Called(_a0)
}

// BusControl_IsConsumed_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsConsumed'
type BusControl_IsConsumed_Call struct {
	*mock.Call
}

// IsConsumed is a helper method to define mock.On call
//   - _a0 func(types.IMessage)
func (_e *BusControl_Expecter) IsConsumed(_a0 interface{}) *BusControl_IsConsumed_Call {
	return &BusControl_IsConsumed_Call{Call: _e.mock.On("IsConsumed", _a0)}
}

func (_c *BusControl_IsConsumed_Call) Run(run func(_a0 func(types.IMessage))) *BusControl_IsConsumed_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func(types.IMessage)))
	})
	return _c
}

func (_c *BusControl_IsConsumed_Call) Return() *BusControl_IsConsumed_Call {
	_c.Call.Return()
	return _c
}

func (_c *BusControl_IsConsumed_Call) RunAndReturn(run func(func(types.IMessage))) *BusControl_IsConsumed_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields: ctx
func (_m *BusControl) Start(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BusControl_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type BusControl_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - ctx context.Context
func (_e *BusControl_Expecter) Start(ctx interface{}) *BusControl_Start_Call {
	return &BusControl_Start_Call{Call: _e.mock.On("Start", ctx)}
}

func (_c *BusControl_Start_Call) Run(run func(ctx context.Context)) *BusControl_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *BusControl_Start_Call) Return(_a0 error) *BusControl_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BusControl_Start_Call) RunAndReturn(run func(context.Context) error) *BusControl_Start_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with given fields:
func (_m *BusControl) Stop() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Stop")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BusControl_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type BusControl_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *BusControl_Expecter) Stop() *BusControl_Stop_Call {
	return &BusControl_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *BusControl_Stop_Call) Run(run func()) *BusControl_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *BusControl_Stop_Call) Return(_a0 error) *BusControl_Stop_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BusControl_Stop_Call) RunAndReturn(run func() error) *BusControl_Stop_Call {
	_c.Call.Return(run)
	return _c
}

// NewBusControl creates a new instance of BusControl. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBusControl(t interface {
	mock.TestingT
	Cleanup(func())
}) *BusControl {
	mock := &BusControl{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
