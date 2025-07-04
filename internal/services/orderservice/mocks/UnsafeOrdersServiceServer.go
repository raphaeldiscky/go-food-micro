// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UnsafeOrdersServiceServer is an autogenerated mock type for the UnsafeOrdersServiceServer type
type UnsafeOrdersServiceServer struct {
	mock.Mock
}

type UnsafeOrdersServiceServer_Expecter struct {
	mock *mock.Mock
}

func (_m *UnsafeOrdersServiceServer) EXPECT() *UnsafeOrdersServiceServer_Expecter {
	return &UnsafeOrdersServiceServer_Expecter{mock: &_m.Mock}
}

// mustEmbedUnimplementedOrdersServiceServer provides a mock function with given fields:
func (_m *UnsafeOrdersServiceServer) mustEmbedUnimplementedOrdersServiceServer() {
	_m.Called()
}

// UnsafeOrdersServiceServer_mustEmbedUnimplementedOrdersServiceServer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'mustEmbedUnimplementedOrdersServiceServer'
type UnsafeOrdersServiceServer_mustEmbedUnimplementedOrdersServiceServer_Call struct {
	*mock.Call
}

// mustEmbedUnimplementedOrdersServiceServer is a helper method to define mock.On call
func (_e *UnsafeOrdersServiceServer_Expecter) mustEmbedUnimplementedOrdersServiceServer() *UnsafeOrdersServiceServer_mustEmbedUnimplementedOrdersServiceServer_Call {
	return &UnsafeOrdersServiceServer_mustEmbedUnimplementedOrdersServiceServer_Call{Call: _e.mock.On("mustEmbedUnimplementedOrdersServiceServer")}
}

func (_c *UnsafeOrdersServiceServer_mustEmbedUnimplementedOrdersServiceServer_Call) Run(run func()) *UnsafeOrdersServiceServer_mustEmbedUnimplementedOrdersServiceServer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UnsafeOrdersServiceServer_mustEmbedUnimplementedOrdersServiceServer_Call) Return() *UnsafeOrdersServiceServer_mustEmbedUnimplementedOrdersServiceServer_Call {
	_c.Call.Return()
	return _c
}

func (_c *UnsafeOrdersServiceServer_mustEmbedUnimplementedOrdersServiceServer_Call) RunAndReturn(run func()) *UnsafeOrdersServiceServer_mustEmbedUnimplementedOrdersServiceServer_Call {
	_c.Call.Return(run)
	return _c
}

// NewUnsafeOrdersServiceServer creates a new instance of UnsafeOrdersServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUnsafeOrdersServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *UnsafeOrdersServiceServer {
	mock := &UnsafeOrdersServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
