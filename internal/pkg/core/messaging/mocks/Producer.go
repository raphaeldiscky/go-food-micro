// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	metadata "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	mock "github.com/stretchr/testify/mock"

	types "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

// Producer is an autogenerated mock type for the Producer type
type Producer struct {
	mock.Mock
}

type Producer_Expecter struct {
	mock *mock.Mock
}

func (_m *Producer) EXPECT() *Producer_Expecter {
	return &Producer_Expecter{mock: &_m.Mock}
}

// IsProduced provides a mock function with given fields: _a0
func (_m *Producer) IsProduced(_a0 func(types.IMessage)) {
	_m.Called(_a0)
}

// Producer_IsProduced_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsProduced'
type Producer_IsProduced_Call struct {
	*mock.Call
}

// IsProduced is a helper method to define mock.On call
//   - _a0 func(types.IMessage)
func (_e *Producer_Expecter) IsProduced(_a0 interface{}) *Producer_IsProduced_Call {
	return &Producer_IsProduced_Call{Call: _e.mock.On("IsProduced", _a0)}
}

func (_c *Producer_IsProduced_Call) Run(run func(_a0 func(types.IMessage))) *Producer_IsProduced_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func(types.IMessage)))
	})
	return _c
}

func (_c *Producer_IsProduced_Call) Return() *Producer_IsProduced_Call {
	_c.Call.Return()
	return _c
}

func (_c *Producer_IsProduced_Call) RunAndReturn(run func(func(types.IMessage))) *Producer_IsProduced_Call {
	_c.Call.Return(run)
	return _c
}

// PublishMessage provides a mock function with given fields: ctx, message, meta
func (_m *Producer) PublishMessage(ctx context.Context, message types.IMessage, meta metadata.Metadata) error {
	ret := _m.Called(ctx, message, meta)

	if len(ret) == 0 {
		panic("no return value specified for PublishMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.IMessage, metadata.Metadata) error); ok {
		r0 = rf(ctx, message, meta)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Producer_PublishMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishMessage'
type Producer_PublishMessage_Call struct {
	*mock.Call
}

// PublishMessage is a helper method to define mock.On call
//   - ctx context.Context
//   - message types.IMessage
//   - meta metadata.Metadata
func (_e *Producer_Expecter) PublishMessage(ctx interface{}, message interface{}, meta interface{}) *Producer_PublishMessage_Call {
	return &Producer_PublishMessage_Call{Call: _e.mock.On("PublishMessage", ctx, message, meta)}
}

func (_c *Producer_PublishMessage_Call) Run(run func(ctx context.Context, message types.IMessage, meta metadata.Metadata)) *Producer_PublishMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.IMessage), args[2].(metadata.Metadata))
	})
	return _c
}

func (_c *Producer_PublishMessage_Call) Return(_a0 error) *Producer_PublishMessage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Producer_PublishMessage_Call) RunAndReturn(run func(context.Context, types.IMessage, metadata.Metadata) error) *Producer_PublishMessage_Call {
	_c.Call.Return(run)
	return _c
}

// PublishMessageWithTopicName provides a mock function with given fields: ctx, message, meta, topicOrExchangeName
func (_m *Producer) PublishMessageWithTopicName(ctx context.Context, message types.IMessage, meta metadata.Metadata, topicOrExchangeName string) error {
	ret := _m.Called(ctx, message, meta, topicOrExchangeName)

	if len(ret) == 0 {
		panic("no return value specified for PublishMessageWithTopicName")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.IMessage, metadata.Metadata, string) error); ok {
		r0 = rf(ctx, message, meta, topicOrExchangeName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Producer_PublishMessageWithTopicName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishMessageWithTopicName'
type Producer_PublishMessageWithTopicName_Call struct {
	*mock.Call
}

// PublishMessageWithTopicName is a helper method to define mock.On call
//   - ctx context.Context
//   - message types.IMessage
//   - meta metadata.Metadata
//   - topicOrExchangeName string
func (_e *Producer_Expecter) PublishMessageWithTopicName(ctx interface{}, message interface{}, meta interface{}, topicOrExchangeName interface{}) *Producer_PublishMessageWithTopicName_Call {
	return &Producer_PublishMessageWithTopicName_Call{Call: _e.mock.On("PublishMessageWithTopicName", ctx, message, meta, topicOrExchangeName)}
}

func (_c *Producer_PublishMessageWithTopicName_Call) Run(run func(ctx context.Context, message types.IMessage, meta metadata.Metadata, topicOrExchangeName string)) *Producer_PublishMessageWithTopicName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.IMessage), args[2].(metadata.Metadata), args[3].(string))
	})
	return _c
}

func (_c *Producer_PublishMessageWithTopicName_Call) Return(_a0 error) *Producer_PublishMessageWithTopicName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Producer_PublishMessageWithTopicName_Call) RunAndReturn(run func(context.Context, types.IMessage, metadata.Metadata, string) error) *Producer_PublishMessageWithTopicName_Call {
	_c.Call.Return(run)
	return _c
}

// NewProducer creates a new instance of Producer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProducer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Producer {
	mock := &Producer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
