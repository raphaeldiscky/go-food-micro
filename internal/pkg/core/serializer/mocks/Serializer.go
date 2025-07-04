// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Serializer is an autogenerated mock type for the Serializer type
type Serializer struct {
	mock.Mock
}

type Serializer_Expecter struct {
	mock *mock.Mock
}

func (_m *Serializer) EXPECT() *Serializer_Expecter {
	return &Serializer_Expecter{mock: &_m.Mock}
}

// ColoredPrettyPrint provides a mock function with given fields: data
func (_m *Serializer) ColoredPrettyPrint(data interface{}) string {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for ColoredPrettyPrint")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(interface{}) string); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Serializer_ColoredPrettyPrint_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ColoredPrettyPrint'
type Serializer_ColoredPrettyPrint_Call struct {
	*mock.Call
}

// ColoredPrettyPrint is a helper method to define mock.On call
//   - data interface{}
func (_e *Serializer_Expecter) ColoredPrettyPrint(data interface{}) *Serializer_ColoredPrettyPrint_Call {
	return &Serializer_ColoredPrettyPrint_Call{Call: _e.mock.On("ColoredPrettyPrint", data)}
}

func (_c *Serializer_ColoredPrettyPrint_Call) Run(run func(data interface{})) *Serializer_ColoredPrettyPrint_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Serializer_ColoredPrettyPrint_Call) Return(_a0 string) *Serializer_ColoredPrettyPrint_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Serializer_ColoredPrettyPrint_Call) RunAndReturn(run func(interface{}) string) *Serializer_ColoredPrettyPrint_Call {
	_c.Call.Return(run)
	return _c
}

// DecodeWithMapStructure provides a mock function with given fields: input, output
func (_m *Serializer) DecodeWithMapStructure(input interface{}, output interface{}) error {
	ret := _m.Called(input, output)

	if len(ret) == 0 {
		panic("no return value specified for DecodeWithMapStructure")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) error); ok {
		r0 = rf(input, output)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Serializer_DecodeWithMapStructure_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DecodeWithMapStructure'
type Serializer_DecodeWithMapStructure_Call struct {
	*mock.Call
}

// DecodeWithMapStructure is a helper method to define mock.On call
//   - input interface{}
//   - output interface{}
func (_e *Serializer_Expecter) DecodeWithMapStructure(input interface{}, output interface{}) *Serializer_DecodeWithMapStructure_Call {
	return &Serializer_DecodeWithMapStructure_Call{Call: _e.mock.On("DecodeWithMapStructure", input, output)}
}

func (_c *Serializer_DecodeWithMapStructure_Call) Run(run func(input interface{}, output interface{})) *Serializer_DecodeWithMapStructure_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(interface{}))
	})
	return _c
}

func (_c *Serializer_DecodeWithMapStructure_Call) Return(_a0 error) *Serializer_DecodeWithMapStructure_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Serializer_DecodeWithMapStructure_Call) RunAndReturn(run func(interface{}, interface{}) error) *Serializer_DecodeWithMapStructure_Call {
	_c.Call.Return(run)
	return _c
}

// Marshal provides a mock function with given fields: v
func (_m *Serializer) Marshal(v interface{}) ([]byte, error) {
	ret := _m.Called(v)

	if len(ret) == 0 {
		panic("no return value specified for Marshal")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}) ([]byte, error)); ok {
		return rf(v)
	}
	if rf, ok := ret.Get(0).(func(interface{}) []byte); ok {
		r0 = rf(v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(v)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Serializer_Marshal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Marshal'
type Serializer_Marshal_Call struct {
	*mock.Call
}

// Marshal is a helper method to define mock.On call
//   - v interface{}
func (_e *Serializer_Expecter) Marshal(v interface{}) *Serializer_Marshal_Call {
	return &Serializer_Marshal_Call{Call: _e.mock.On("Marshal", v)}
}

func (_c *Serializer_Marshal_Call) Run(run func(v interface{})) *Serializer_Marshal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Serializer_Marshal_Call) Return(_a0 []byte, _a1 error) *Serializer_Marshal_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Serializer_Marshal_Call) RunAndReturn(run func(interface{}) ([]byte, error)) *Serializer_Marshal_Call {
	_c.Call.Return(run)
	return _c
}

// PrettyPrint provides a mock function with given fields: data
func (_m *Serializer) PrettyPrint(data interface{}) string {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for PrettyPrint")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(interface{}) string); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Serializer_PrettyPrint_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PrettyPrint'
type Serializer_PrettyPrint_Call struct {
	*mock.Call
}

// PrettyPrint is a helper method to define mock.On call
//   - data interface{}
func (_e *Serializer_Expecter) PrettyPrint(data interface{}) *Serializer_PrettyPrint_Call {
	return &Serializer_PrettyPrint_Call{Call: _e.mock.On("PrettyPrint", data)}
}

func (_c *Serializer_PrettyPrint_Call) Run(run func(data interface{})) *Serializer_PrettyPrint_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Serializer_PrettyPrint_Call) Return(_a0 string) *Serializer_PrettyPrint_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Serializer_PrettyPrint_Call) RunAndReturn(run func(interface{}) string) *Serializer_PrettyPrint_Call {
	_c.Call.Return(run)
	return _c
}

// Unmarshal provides a mock function with given fields: data, v
func (_m *Serializer) Unmarshal(data []byte, v interface{}) error {
	ret := _m.Called(data, v)

	if len(ret) == 0 {
		panic("no return value specified for Unmarshal")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, interface{}) error); ok {
		r0 = rf(data, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Serializer_Unmarshal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unmarshal'
type Serializer_Unmarshal_Call struct {
	*mock.Call
}

// Unmarshal is a helper method to define mock.On call
//   - data []byte
//   - v interface{}
func (_e *Serializer_Expecter) Unmarshal(data interface{}, v interface{}) *Serializer_Unmarshal_Call {
	return &Serializer_Unmarshal_Call{Call: _e.mock.On("Unmarshal", data, v)}
}

func (_c *Serializer_Unmarshal_Call) Run(run func(data []byte, v interface{})) *Serializer_Unmarshal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte), args[1].(interface{}))
	})
	return _c
}

func (_c *Serializer_Unmarshal_Call) Return(_a0 error) *Serializer_Unmarshal_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Serializer_Unmarshal_Call) RunAndReturn(run func([]byte, interface{}) error) *Serializer_Unmarshal_Call {
	_c.Call.Return(run)
	return _c
}

// UnmarshalFromJson provides a mock function with given fields: data, v
func (_m *Serializer) UnmarshalFromJson(data string, v interface{}) error {
	ret := _m.Called(data, v)

	if len(ret) == 0 {
		panic("no return value specified for UnmarshalFromJson")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}) error); ok {
		r0 = rf(data, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Serializer_UnmarshalFromJson_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnmarshalFromJson'
type Serializer_UnmarshalFromJson_Call struct {
	*mock.Call
}

// UnmarshalFromJson is a helper method to define mock.On call
//   - data string
//   - v interface{}
func (_e *Serializer_Expecter) UnmarshalFromJson(data interface{}, v interface{}) *Serializer_UnmarshalFromJson_Call {
	return &Serializer_UnmarshalFromJson_Call{Call: _e.mock.On("UnmarshalFromJson", data, v)}
}

func (_c *Serializer_UnmarshalFromJson_Call) Run(run func(data string, v interface{})) *Serializer_UnmarshalFromJson_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(interface{}))
	})
	return _c
}

func (_c *Serializer_UnmarshalFromJson_Call) Return(_a0 error) *Serializer_UnmarshalFromJson_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Serializer_UnmarshalFromJson_Call) RunAndReturn(run func(string, interface{}) error) *Serializer_UnmarshalFromJson_Call {
	_c.Call.Return(run)
	return _c
}

// UnmarshalToMap provides a mock function with given fields: data, v
func (_m *Serializer) UnmarshalToMap(data []byte, v *map[string]interface{}) error {
	ret := _m.Called(data, v)

	if len(ret) == 0 {
		panic("no return value specified for UnmarshalToMap")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, *map[string]interface{}) error); ok {
		r0 = rf(data, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Serializer_UnmarshalToMap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnmarshalToMap'
type Serializer_UnmarshalToMap_Call struct {
	*mock.Call
}

// UnmarshalToMap is a helper method to define mock.On call
//   - data []byte
//   - v *map[string]interface{}
func (_e *Serializer_Expecter) UnmarshalToMap(data interface{}, v interface{}) *Serializer_UnmarshalToMap_Call {
	return &Serializer_UnmarshalToMap_Call{Call: _e.mock.On("UnmarshalToMap", data, v)}
}

func (_c *Serializer_UnmarshalToMap_Call) Run(run func(data []byte, v *map[string]interface{})) *Serializer_UnmarshalToMap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte), args[1].(*map[string]interface{}))
	})
	return _c
}

func (_c *Serializer_UnmarshalToMap_Call) Return(_a0 error) *Serializer_UnmarshalToMap_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Serializer_UnmarshalToMap_Call) RunAndReturn(run func([]byte, *map[string]interface{}) error) *Serializer_UnmarshalToMap_Call {
	_c.Call.Return(run)
	return _c
}

// UnmarshalToMapFromJSON provides a mock function with given fields: data, v
func (_m *Serializer) UnmarshalToMapFromJSON(data string, v *map[string]interface{}) error {
	ret := _m.Called(data, v)

	if len(ret) == 0 {
		panic("no return value specified for UnmarshalToMapFromJSON")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *map[string]interface{}) error); ok {
		r0 = rf(data, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Serializer_UnmarshalToMapFromJSON_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnmarshalToMapFromJSON'
type Serializer_UnmarshalToMapFromJSON_Call struct {
	*mock.Call
}

// UnmarshalToMapFromJSON is a helper method to define mock.On call
//   - data string
//   - v *map[string]interface{}
func (_e *Serializer_Expecter) UnmarshalToMapFromJSON(data interface{}, v interface{}) *Serializer_UnmarshalToMapFromJSON_Call {
	return &Serializer_UnmarshalToMapFromJSON_Call{Call: _e.mock.On("UnmarshalToMapFromJSON", data, v)}
}

func (_c *Serializer_UnmarshalToMapFromJSON_Call) Run(run func(data string, v *map[string]interface{})) *Serializer_UnmarshalToMapFromJSON_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*map[string]interface{}))
	})
	return _c
}

func (_c *Serializer_UnmarshalToMapFromJSON_Call) Return(_a0 error) *Serializer_UnmarshalToMapFromJSON_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Serializer_UnmarshalToMapFromJSON_Call) RunAndReturn(run func(string, *map[string]interface{}) error) *Serializer_UnmarshalToMapFromJSON_Call {
	_c.Call.Return(run)
	return _c
}

// NewSerializer creates a new instance of Serializer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSerializer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Serializer {
	mock := &Serializer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
