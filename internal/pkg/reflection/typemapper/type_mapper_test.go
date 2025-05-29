// Package typemapper provides a type mapper.
package typemapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetTypeNameByT tests the get type name by t.
func TestGetTypeNameByT(t *testing.T) {
	t.Parallel()
	pointerTypeName := GetGenericTypeNameByT[*Test]()
	nonePointerTypeName := GetGenericTypeNameByT[Test]()

	require.Equal(t, pointerTypeName, "*Test")
	require.Equal(t, nonePointerTypeName, "Test")
}

// TestGetNonePointerTypeNameByT tests the get none pointer type name by t.
func TestGetNonePointerTypeNameByT(t *testing.T) {
	pointerTypeName := GetGenericNonePointerTypeNameByT[*Test]()
	nonePointerTypeName := GetGenericNonePointerTypeNameByT[Test]()

	require.Equal(t, pointerTypeName, "Test")
	require.Equal(t, nonePointerTypeName, "Test")
}

// TestTypeByName tests the type by name.
func TestTypeByName(t *testing.T) {
	s1 := TypeByName("*typemapper.Test")
	s2 := TypeByName("typemapper.Test")
	s3 := TypeByName("*Test")
	s4 := TypeByName("Test")

	assert.NotNil(t, s1)
	assert.NotNil(t, s2)
	assert.NotNil(t, s3)
	assert.NotNil(t, s4)
}

// TestGetTypeName tests the get type name.
func TestGetTypeName(t *testing.T) {
	t1 := Test{A: 10}
	t2 := &Test{A: 10}

	typeName1 := GetTypeName(t1)
	typeName2 := GetTypeName(t2)

	assert.Equal(t, "Test", typeName1)
	assert.Equal(t, "*Test", typeName2)
}

// TestGetFullTypeName tests the get full type name.
func TestGetFullTypeName(t *testing.T) {
	t1 := Test{A: 10}
	t2 := &Test{A: 10}

	typeName1 := GetFullTypeName(t1)
	typeName2 := GetFullTypeName(t2)
	typeName3 := GetGenericFullTypeNameByT[*Test]()
	typeName4 := GetGenericFullTypeNameByT[Test]()

	assert.Equal(t, "typemapper.Test", typeName1)
	assert.Equal(t, "*typemapper.Test", typeName2)
	assert.Equal(t, "*typemapper.Test", typeName3)
	assert.Equal(t, "typemapper.Test", typeName4)
}

// TestInstanceByTypeName tests the instance by type name.
func TestInstanceByTypeName(t *testing.T) {
	s1 := InstanceByTypeName("typemapper.Test").(Test)
	s1.A = 100
	assert.NotNil(t, s1)
	assert.NotZero(t, s1.A)

	s2 := InstanceByTypeName("*typemapper.Test").(*Test)
	s2.A = 100
	assert.NotNil(t, s2)
	assert.NotZero(t, s2.A)

	s3 := InstanceByTypeName("*Test").(*Test)
	assert.NotNil(t, s3)

	s4 := InstanceByTypeName("Test").(Test)
	assert.NotNil(t, s4)
}

// TestInstancePointerByTypeName tests the instance pointer by type name.
func TestInstancePointerByTypeName(t *testing.T) {
	s1 := InstancePointerByTypeName("*typemapper.Test").(*Test)
	s2 := InstancePointerByTypeName("typemapper.Test").(*Test)
	s3 := InstancePointerByTypeName("*Test").(*Test)
	s4 := InstancePointerByTypeName("Test").(*Test)

	assert.NotNil(t, s1)
	assert.NotNil(t, s2)
	assert.NotNil(t, s3)
	assert.NotNil(t, s4)
}

// TestGetTypeFromGeneric tests the get type from generic.
func TestGetTypeFromGeneric(t *testing.T) {
	s1 := GetGenericTypeByT[Test]()
	s2 := GetGenericTypeByT[*Test]()
	s3 := GetGenericTypeByT[ITest]()

	assert.NotNil(t, s1)
	assert.NotNil(t, s2)
	assert.NotNil(t, s3)
}

// TestGenericInstanceByT tests the generic instance by t.
func TestGenericInstanceByT(t *testing.T) {
	t.Parallel()
	s1 := GenericInstanceByT[*Test]()
	s2 := GenericInstanceByT[Test]()

	assert.NotNil(t, s1)
	assert.NotNil(t, s2)
}

// TestTypeByNameAndImplementedInterface tests the type by name and implemented interface.
func TestTypeByNameAndImplementedInterface(t *testing.T) {
	t.Parallel()
	s1 := TypeByNameAndImplementedInterface[ITest]("*typemapper.Test")

	assert.NotNil(t, s1)
}

// TestEmptyInstanceByTypeNameAndImplementedInterface tests the empty instance by type name and implemented interface.
func TestEmptyInstanceByTypeNameAndImplementedInterface(t *testing.T) {
	t.Parallel()
	s1 := EmptyInstanceByTypeNameAndImplementedInterface[ITest]("*typemapper.Test")

	assert.NotNil(t, s1)
}

// TestGetReflectType tests the get reflect type.
func TestGetReflectType(t *testing.T) {
	t.Parallel()
	s1 := GetReflectType(Test{})
	s2 := GetReflectType(&Test{})
	s3 := GetReflectType((*ITest)(nil))

	assert.NotNil(t, s1)
	assert.NotNil(t, s2)
	assert.NotNil(t, s3)
}

// TestGetPackageName tests the get package name.
func TestGetPackageName(t *testing.T) {
	t.Parallel()
	pkName := GetPackageName(&Test{})
	pkName2 := GetPackageName(Test{})

	assert.Equal(t, "typemapper", pkName)
	assert.Equal(t, "typemapper", pkName2)
}

type Test struct {
	A int
}

type ITest interface {
	Method1()
}

func (t *Test) Method1() {
}
