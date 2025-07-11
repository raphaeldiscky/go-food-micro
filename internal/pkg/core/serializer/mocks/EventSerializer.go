// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"
	mock "github.com/stretchr/testify/mock"

	reflect "reflect"

	serializer "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer"
)

// EventSerializer is an autogenerated mock type for the EventSerializer type
type EventSerializer struct {
	mock.Mock
}

type EventSerializer_Expecter struct {
	mock *mock.Mock
}

func (_m *EventSerializer) EXPECT() *EventSerializer_Expecter {
	return &EventSerializer_Expecter{mock: &_m.Mock}
}

// ContentType provides a mock function with given fields:
func (_m *EventSerializer) ContentType() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ContentType")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// EventSerializer_ContentType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ContentType'
type EventSerializer_ContentType_Call struct {
	*mock.Call
}

// ContentType is a helper method to define mock.On call
func (_e *EventSerializer_Expecter) ContentType() *EventSerializer_ContentType_Call {
	return &EventSerializer_ContentType_Call{Call: _e.mock.On("ContentType")}
}

func (_c *EventSerializer_ContentType_Call) Run(run func()) *EventSerializer_ContentType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EventSerializer_ContentType_Call) Return(_a0 string) *EventSerializer_ContentType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventSerializer_ContentType_Call) RunAndReturn(run func() string) *EventSerializer_ContentType_Call {
	_c.Call.Return(run)
	return _c
}

// Deserialize provides a mock function with given fields: data, eventType, contentType
func (_m *EventSerializer) Deserialize(data []byte, eventType string, contentType string) (domain.IDomainEvent, error) {
	ret := _m.Called(data, eventType, contentType)

	if len(ret) == 0 {
		panic("no return value specified for Deserialize")
	}

	var r0 domain.IDomainEvent
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, string, string) (domain.IDomainEvent, error)); ok {
		return rf(data, eventType, contentType)
	}
	if rf, ok := ret.Get(0).(func([]byte, string, string) domain.IDomainEvent); ok {
		r0 = rf(data, eventType, contentType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.IDomainEvent)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte, string, string) error); ok {
		r1 = rf(data, eventType, contentType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventSerializer_Deserialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Deserialize'
type EventSerializer_Deserialize_Call struct {
	*mock.Call
}

// Deserialize is a helper method to define mock.On call
//   - data []byte
//   - eventType string
//   - contentType string
func (_e *EventSerializer_Expecter) Deserialize(data interface{}, eventType interface{}, contentType interface{}) *EventSerializer_Deserialize_Call {
	return &EventSerializer_Deserialize_Call{Call: _e.mock.On("Deserialize", data, eventType, contentType)}
}

func (_c *EventSerializer_Deserialize_Call) Run(run func(data []byte, eventType string, contentType string)) *EventSerializer_Deserialize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *EventSerializer_Deserialize_Call) Return(_a0 domain.IDomainEvent, _a1 error) *EventSerializer_Deserialize_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *EventSerializer_Deserialize_Call) RunAndReturn(run func([]byte, string, string) (domain.IDomainEvent, error)) *EventSerializer_Deserialize_Call {
	_c.Call.Return(run)
	return _c
}

// DeserializeObject provides a mock function with given fields: data, eventType, contentType
func (_m *EventSerializer) DeserializeObject(data []byte, eventType string, contentType string) (interface{}, error) {
	ret := _m.Called(data, eventType, contentType)

	if len(ret) == 0 {
		panic("no return value specified for DeserializeObject")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, string, string) (interface{}, error)); ok {
		return rf(data, eventType, contentType)
	}
	if rf, ok := ret.Get(0).(func([]byte, string, string) interface{}); ok {
		r0 = rf(data, eventType, contentType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func([]byte, string, string) error); ok {
		r1 = rf(data, eventType, contentType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventSerializer_DeserializeObject_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeserializeObject'
type EventSerializer_DeserializeObject_Call struct {
	*mock.Call
}

// DeserializeObject is a helper method to define mock.On call
//   - data []byte
//   - eventType string
//   - contentType string
func (_e *EventSerializer_Expecter) DeserializeObject(data interface{}, eventType interface{}, contentType interface{}) *EventSerializer_DeserializeObject_Call {
	return &EventSerializer_DeserializeObject_Call{Call: _e.mock.On("DeserializeObject", data, eventType, contentType)}
}

func (_c *EventSerializer_DeserializeObject_Call) Run(run func(data []byte, eventType string, contentType string)) *EventSerializer_DeserializeObject_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *EventSerializer_DeserializeObject_Call) Return(_a0 interface{}, _a1 error) *EventSerializer_DeserializeObject_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *EventSerializer_DeserializeObject_Call) RunAndReturn(run func([]byte, string, string) (interface{}, error)) *EventSerializer_DeserializeObject_Call {
	_c.Call.Return(run)
	return _c
}

// DeserializeType provides a mock function with given fields: data, eventType, contentType
func (_m *EventSerializer) DeserializeType(data []byte, eventType reflect.Type, contentType string) (domain.IDomainEvent, error) {
	ret := _m.Called(data, eventType, contentType)

	if len(ret) == 0 {
		panic("no return value specified for DeserializeType")
	}

	var r0 domain.IDomainEvent
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, reflect.Type, string) (domain.IDomainEvent, error)); ok {
		return rf(data, eventType, contentType)
	}
	if rf, ok := ret.Get(0).(func([]byte, reflect.Type, string) domain.IDomainEvent); ok {
		r0 = rf(data, eventType, contentType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.IDomainEvent)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte, reflect.Type, string) error); ok {
		r1 = rf(data, eventType, contentType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventSerializer_DeserializeType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeserializeType'
type EventSerializer_DeserializeType_Call struct {
	*mock.Call
}

// DeserializeType is a helper method to define mock.On call
//   - data []byte
//   - eventType reflect.Type
//   - contentType string
func (_e *EventSerializer_Expecter) DeserializeType(data interface{}, eventType interface{}, contentType interface{}) *EventSerializer_DeserializeType_Call {
	return &EventSerializer_DeserializeType_Call{Call: _e.mock.On("DeserializeType", data, eventType, contentType)}
}

func (_c *EventSerializer_DeserializeType_Call) Run(run func(data []byte, eventType reflect.Type, contentType string)) *EventSerializer_DeserializeType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte), args[1].(reflect.Type), args[2].(string))
	})
	return _c
}

func (_c *EventSerializer_DeserializeType_Call) Return(_a0 domain.IDomainEvent, _a1 error) *EventSerializer_DeserializeType_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *EventSerializer_DeserializeType_Call) RunAndReturn(run func([]byte, reflect.Type, string) (domain.IDomainEvent, error)) *EventSerializer_DeserializeType_Call {
	_c.Call.Return(run)
	return _c
}

// Serialize provides a mock function with given fields: event
func (_m *EventSerializer) Serialize(event domain.IDomainEvent) (*serializer.EventSerializationResult, error) {
	ret := _m.Called(event)

	if len(ret) == 0 {
		panic("no return value specified for Serialize")
	}

	var r0 *serializer.EventSerializationResult
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.IDomainEvent) (*serializer.EventSerializationResult, error)); ok {
		return rf(event)
	}
	if rf, ok := ret.Get(0).(func(domain.IDomainEvent) *serializer.EventSerializationResult); ok {
		r0 = rf(event)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*serializer.EventSerializationResult)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.IDomainEvent) error); ok {
		r1 = rf(event)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventSerializer_Serialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Serialize'
type EventSerializer_Serialize_Call struct {
	*mock.Call
}

// Serialize is a helper method to define mock.On call
//   - event domain.IDomainEvent
func (_e *EventSerializer_Expecter) Serialize(event interface{}) *EventSerializer_Serialize_Call {
	return &EventSerializer_Serialize_Call{Call: _e.mock.On("Serialize", event)}
}

func (_c *EventSerializer_Serialize_Call) Run(run func(event domain.IDomainEvent)) *EventSerializer_Serialize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.IDomainEvent))
	})
	return _c
}

func (_c *EventSerializer_Serialize_Call) Return(_a0 *serializer.EventSerializationResult, _a1 error) *EventSerializer_Serialize_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *EventSerializer_Serialize_Call) RunAndReturn(run func(domain.IDomainEvent) (*serializer.EventSerializationResult, error)) *EventSerializer_Serialize_Call {
	_c.Call.Return(run)
	return _c
}

// SerializeObject provides a mock function with given fields: event
func (_m *EventSerializer) SerializeObject(event interface{}) (*serializer.EventSerializationResult, error) {
	ret := _m.Called(event)

	if len(ret) == 0 {
		panic("no return value specified for SerializeObject")
	}

	var r0 *serializer.EventSerializationResult
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}) (*serializer.EventSerializationResult, error)); ok {
		return rf(event)
	}
	if rf, ok := ret.Get(0).(func(interface{}) *serializer.EventSerializationResult); ok {
		r0 = rf(event)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*serializer.EventSerializationResult)
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(event)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventSerializer_SerializeObject_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SerializeObject'
type EventSerializer_SerializeObject_Call struct {
	*mock.Call
}

// SerializeObject is a helper method to define mock.On call
//   - event interface{}
func (_e *EventSerializer_Expecter) SerializeObject(event interface{}) *EventSerializer_SerializeObject_Call {
	return &EventSerializer_SerializeObject_Call{Call: _e.mock.On("SerializeObject", event)}
}

func (_c *EventSerializer_SerializeObject_Call) Run(run func(event interface{})) *EventSerializer_SerializeObject_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *EventSerializer_SerializeObject_Call) Return(_a0 *serializer.EventSerializationResult, _a1 error) *EventSerializer_SerializeObject_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *EventSerializer_SerializeObject_Call) RunAndReturn(run func(interface{}) (*serializer.EventSerializationResult, error)) *EventSerializer_SerializeObject_Call {
	_c.Call.Return(run)
	return _c
}

// Serializer provides a mock function with given fields:
func (_m *EventSerializer) Serializer() serializer.Serializer {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Serializer")
	}

	var r0 serializer.Serializer
	if rf, ok := ret.Get(0).(func() serializer.Serializer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(serializer.Serializer)
		}
	}

	return r0
}

// EventSerializer_Serializer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Serializer'
type EventSerializer_Serializer_Call struct {
	*mock.Call
}

// Serializer is a helper method to define mock.On call
func (_e *EventSerializer_Expecter) Serializer() *EventSerializer_Serializer_Call {
	return &EventSerializer_Serializer_Call{Call: _e.mock.On("Serializer")}
}

func (_c *EventSerializer_Serializer_Call) Run(run func()) *EventSerializer_Serializer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EventSerializer_Serializer_Call) Return(_a0 serializer.Serializer) *EventSerializer_Serializer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventSerializer_Serializer_Call) RunAndReturn(run func() serializer.Serializer) *EventSerializer_Serializer_Call {
	_c.Call.Return(run)
	return _c
}

// NewEventSerializer creates a new instance of EventSerializer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventSerializer(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventSerializer {
	mock := &EventSerializer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
