// Package reflectionHelper provides a reflection helper.
package reflectionHelper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ref: https://gist.github.com/drewolson/4771479
// https://stackoverflow.com/a/60598827/581476

// PersonPublic represents a person public.
type PersonPublic struct {
	Name string
	Age  int
}

// PersonPrivate represents a person private.
type PersonPrivate struct {
	name string
	age  int
}

// Name returns the name of the person.
func (p *PersonPrivate) Name() string {
	return p.name
}

// Age returns the age of the person.
func (p *PersonPrivate) Age() int {
	return p.age
}

// TestFieldValuesForExportedFieldsAndAddressableStruct tests the field values for exported fields and addressable struct.
func TestFieldValuesForExportedFieldsAndAddressableStruct(t *testing.T) {
	p := &PersonPublic{Name: "John", Age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

// TestFieldValuesForExportedFieldsAndUnAddressableStruct tests the field values for exported fields and unaddressable struct.
func TestFieldValuesForExportedFieldsAndUnAddressableStruct(t *testing.T) {
	p := PersonPublic{Name: "John", Age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

// TestFieldValuesForUnExportedFieldsAndAddressableStruct tests the field values for unexported fields and addressable struct.
func TestFieldValuesForUnExportedFieldsAndAddressableStruct(t *testing.T) {
	p := &PersonPrivate{name: "John", age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

// TestFieldValuesForUnExportedFieldsAndUnAddressableStruct tests the field values for unexported fields and unaddressable struct.
func TestFieldValuesForUnExportedFieldsAndUnAddressableStruct(t *testing.T) {
	p := PersonPrivate{name: "John", age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

// TestSetFieldValueForExportedFieldsAndAddressableStruct tests the set field value for exported fields and addressable struct.
func TestSetFieldValueForExportedFieldsAndAddressableStruct(t *testing.T) {
	p := &PersonPublic{}

	SetFieldValueByIndex(p, 0, "John")
	SetFieldValueByIndex(p, 1, 20)

	assert.Equal(t, "John", p.Name)
	assert.Equal(t, 20, p.Age)
}

// TestSetFieldValueForExportedFieldsAndUnAddressableStruct tests the set field value for exported fields and unaddressable struct.
func TestSetFieldValueForExportedFieldsAndUnAddressableStruct(t *testing.T) {
	p := PersonPublic{}

	SetFieldValueByIndex(&p, 0, "John")
	SetFieldValueByIndex(&p, 1, 20)

	assert.Equal(t, "John", p.Name)
	assert.Equal(t, 20, p.Age)
}

// TestSetFieldValueForUnExportedFieldsAndUnAddressableStruct tests the set field value for unexported fields and unaddressable struct.
func TestSetFieldValueForUnExportedFieldsAndUnAddressableStruct(t *testing.T) {
	p := PersonPrivate{}

	SetFieldValueByIndex(&p, 0, "John")
	SetFieldValueByIndex(&p, 1, 20)

	assert.Equal(t, "John", p.name)
	assert.Equal(t, 20, p.age)
}

// TestGetFieldValueForExportedFieldsAndAddressableStruct tests the get field value for exported fields and addressable struct.
func TestGetFieldValueForExportedFieldsAndAddressableStruct(t *testing.T) {
	p := &PersonPublic{Name: "John", Age: 20}

	// field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(p).Elem()
	name := GetFieldValue(v.FieldByName("ShortTypeName")).Interface()
	age := GetFieldValue(v.FieldByName("Age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}

// TestGetFieldValueForUnExportedFieldsAndAddressableStruct tests the get field value for unexported fields and addressable struct.
func TestGetFieldValueForUnExportedFieldsAndAddressableStruct(t *testing.T) {
	p := &PersonPrivate{name: "John", age: 30}

	// field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(p).Elem()
	name := GetFieldValue(v.FieldByName("name")).Interface()
	age := GetFieldValue(v.FieldByName("age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 30, age)
}

// TestGetFieldValueForExportedFieldsAndUnAddressableStruct tests the get field value for exported fields and unaddressable struct.
func TestGetFieldValueForExportedFieldsAndUnAddressableStruct(t *testing.T) {
	p := PersonPublic{Name: "John", Age: 20}

	// field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(&p).Elem()
	name := GetFieldValue(v.FieldByName("ShortTypeName")).Interface()
	age := GetFieldValue(v.FieldByName("Age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}

// TestGetFieldValueForUnExportedFieldsAndUnAddressableStruct tests the get field value for unexported fields and unaddressable struct.
func TestGetFieldValueForUnExportedFieldsAndUnAddressableStruct(t *testing.T) {
	p := PersonPrivate{name: "John", age: 20}

	// field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(&p).Elem()
	name := GetFieldValue(v.FieldByName("name")).Interface()
	age := GetFieldValue(v.FieldByName("age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}

// TestGetFieldValueFromMethodAndAddressableStruct tests the get field value from method and addressable struct.
func TestGetFieldValueFromMethodAndAddressableStruct(t *testing.T) {
	p := &PersonPrivate{name: "John", age: 20}
	name := GetFieldValueFromMethodAndObject(p, "ShortTypeName")

	assert.Equal(t, "John", name.Interface())
}

// TestGetFieldValueFromMethodAndUnAddressableStruct tests the get field value from method and unaddressable struct.
func TestGetFieldValueFromMethodAndUnAddressableStruct(t *testing.T) {
	p := PersonPrivate{name: "John", age: 20}
	name := GetFieldValueFromMethodAndObject(p, "ShortTypeName")

	assert.Equal(t, "John", name.Interface())
}

// TestConvertNoPointerTypeToPointerTypeWithAddr tests the convert no pointer type to pointer type with addr.
func TestConvertNoPointerTypeToPointerTypeWithAddr(t *testing.T) {
	// https://www.geeksforgeeks.org/reflect-addr-function-in-golang-with-examples/

	p := PersonPrivate{name: "John", age: 20}
	v := reflect.ValueOf(&p).Elem()
	pointerType := v.Addr()
	name := pointerType.MethodByName("ShortTypeName").Call(nil)[0].Interface()
	age := pointerType.MethodByName("Age").Call(nil)[0].Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}
