// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// IMessage is an autogenerated mock type for the IMessage type
type IMessage struct {
	mock.Mock
}

type IMessage_Expecter struct {
	mock *mock.Mock
}

func (_m *IMessage) EXPECT() *IMessage_Expecter {
	return &IMessage_Expecter{mock: &_m.Mock}
}

// GeMessageId provides a mock function with given fields:
func (_m *IMessage) GeMessageId() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GeMessageId")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IMessage_GeMessageId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GeMessageId'
type IMessage_GeMessageId_Call struct {
	*mock.Call
}

// GeMessageId is a helper method to define mock.On call
func (_e *IMessage_Expecter) GeMessageId() *IMessage_GeMessageId_Call {
	return &IMessage_GeMessageId_Call{Call: _e.mock.On("GeMessageId")}
}

func (_c *IMessage_GeMessageId_Call) Run(run func()) *IMessage_GeMessageId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IMessage_GeMessageId_Call) Return(_a0 string) *IMessage_GeMessageId_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IMessage_GeMessageId_Call) RunAndReturn(run func() string) *IMessage_GeMessageId_Call {
	_c.Call.Return(run)
	return _c
}

// GetCreated provides a mock function with given fields:
func (_m *IMessage) GetCreated() time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetCreated")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// IMessage_GetCreated_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCreated'
type IMessage_GetCreated_Call struct {
	*mock.Call
}

// GetCreated is a helper method to define mock.On call
func (_e *IMessage_Expecter) GetCreated() *IMessage_GetCreated_Call {
	return &IMessage_GetCreated_Call{Call: _e.mock.On("GetCreated")}
}

func (_c *IMessage_GetCreated_Call) Run(run func()) *IMessage_GetCreated_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IMessage_GetCreated_Call) Return(_a0 time.Time) *IMessage_GetCreated_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IMessage_GetCreated_Call) RunAndReturn(run func() time.Time) *IMessage_GetCreated_Call {
	_c.Call.Return(run)
	return _c
}

// GetMessageFullTypeName provides a mock function with given fields:
func (_m *IMessage) GetMessageFullTypeName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMessageFullTypeName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IMessage_GetMessageFullTypeName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMessageFullTypeName'
type IMessage_GetMessageFullTypeName_Call struct {
	*mock.Call
}

// GetMessageFullTypeName is a helper method to define mock.On call
func (_e *IMessage_Expecter) GetMessageFullTypeName() *IMessage_GetMessageFullTypeName_Call {
	return &IMessage_GetMessageFullTypeName_Call{Call: _e.mock.On("GetMessageFullTypeName")}
}

func (_c *IMessage_GetMessageFullTypeName_Call) Run(run func()) *IMessage_GetMessageFullTypeName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IMessage_GetMessageFullTypeName_Call) Return(_a0 string) *IMessage_GetMessageFullTypeName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IMessage_GetMessageFullTypeName_Call) RunAndReturn(run func() string) *IMessage_GetMessageFullTypeName_Call {
	_c.Call.Return(run)
	return _c
}

// GetMessageTypeName provides a mock function with given fields:
func (_m *IMessage) GetMessageTypeName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMessageTypeName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IMessage_GetMessageTypeName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMessageTypeName'
type IMessage_GetMessageTypeName_Call struct {
	*mock.Call
}

// GetMessageTypeName is a helper method to define mock.On call
func (_e *IMessage_Expecter) GetMessageTypeName() *IMessage_GetMessageTypeName_Call {
	return &IMessage_GetMessageTypeName_Call{Call: _e.mock.On("GetMessageTypeName")}
}

func (_c *IMessage_GetMessageTypeName_Call) Run(run func()) *IMessage_GetMessageTypeName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IMessage_GetMessageTypeName_Call) Return(_a0 string) *IMessage_GetMessageTypeName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IMessage_GetMessageTypeName_Call) RunAndReturn(run func() string) *IMessage_GetMessageTypeName_Call {
	_c.Call.Return(run)
	return _c
}

// NewIMessage creates a new instance of IMessage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIMessage(t interface {
	mock.TestingT
	Cleanup(func())
}) *IMessage {
	mock := &IMessage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
