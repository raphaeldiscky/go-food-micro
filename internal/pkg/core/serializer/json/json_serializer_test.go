//go:build unit
// +build unit

package json

import (
	"testing"

	"github.com/stretchr/testify/assert"

	reflect "github.com/goccy/go-reflect"
)

const (
	TestName = "John"
	TestAge  = 30
)

// person is a struct that represents a person.
type person struct {
	Name string
	Age  int
}

// currentSerializer is the current serializer.
var currentSerializer = NewDefaultJsonSerializer()

// TestDeserializeUnstructuredDataIntoEmptyInterface tests the deserialize unstructured data into empty interface.
func TestDeserializeUnstructuredDataIntoEmptyInterface(t *testing.T) {
	t.Helper()
	// https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data
	// https://developpaper.com/mapstructure-of-go/
	// https://pkg.go.dev/encoding/json#Unmarshal
	// when we assign an object type to interface is not pointer object or we don't assign interface, defaultLogger unmarshaler can't deserialize it to the object type and serialize it to map[string]interface{}

	// To unmarshal JSON into an interface value, Unmarshal stores map[string]interface{}
	var jsonMap interface{}

	marshal, err := currentSerializer.Marshal(person{TestName, TestAge})
	if err != nil {
		return
	}

	err = currentSerializer.Unmarshal(marshal, &jsonMap)
	if err != nil {
		panic(err)
	}

	t.Log(jsonMap)

	var jsonMapTyped map[string]interface{}
	var ok bool
	jsonMapTyped, ok = jsonMap.(map[string]interface{})
	if !ok {
		t.Fatal("Failed to convert to map[string]interface{}")
	}

	for key, value := range jsonMapTyped {
		t.Log(key, value)
	}

	assert.True(t, reflect.TypeOf(jsonMapTyped).Kind() == reflect.Map)
	assert.True(t, reflect.TypeOf(jsonMapTyped) == reflect.TypeOf(map[string]interface{}(nil)))
	assert.True(t, jsonMapTyped["Name"] == TestName)
	assert.True(t, jsonMapTyped["Age"] == float64(TestAge))
}

// TestDeserializeUnstructuredDataIntoMap tests the deserialize unstructured data into map.
func TestDeserializeUnstructuredDataIntoMap(t *testing.T) {
	// https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data
	t.Helper()

	// https://developpaper.com/mapstructure-of-go/
	// https://pkg.go.dev/encoding/json#Unmarshal
	// when we assign an object type to interface is not pointer object or we don't assign interface, defaultLogger unmarshaler can't deserialize it to the object type and serialize it to map[string]interface{}

	// To unmarshal a JSON object into a map, Unmarshal first establishes a map to use. If the map is nil, Unmarshal allocates a new map. Otherwise Unmarshal reuses the existing map, keeping existing entries. Unmarshal then stores key-value pairs from the JSON object into the map.
	var jsonMap map[string]interface{}

	marshal, err := currentSerializer.Marshal(person{TestName, TestAge})
	if err != nil {
		t.Fatal(err)
	}

	err = currentSerializer.Unmarshal(marshal, &jsonMap)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(jsonMap)

	for key, value := range jsonMap {
		t.Log(key, value)
	}

	assert.True(t, reflect.TypeOf(jsonMap).Kind() == reflect.Map)
	assert.True(t, reflect.TypeOf(jsonMap) == reflect.TypeOf(map[string]interface{}(nil)))
	assert.True(t, jsonMap["Name"] == TestName)
	assert.True(t, jsonMap["Age"] == float64(TestAge))
}

// TestDeserializeStructuredDataStruct tests the deserialize structured data into struct.
func TestDeserializeStructuredDataStruct(t *testing.T) {
	// https://pkg.go.dev/encoding/json#Unmarshal
	t.Helper()

	// when we assign object to explicit struct type, defaultLogger unmarshaler can deserialize it to the struct

	// To unmarshal JSON into a struct, Unmarshal matches incoming object keys to the keys used by Marshal (either the struct field name or its tag), preferring an exact match but also accepting a case-insensitive match.
	jsonMap := person{}
	v := reflect.ValueOf(&jsonMap)
	if v.Elem().Kind() == reflect.Interface && v.NumMethod() == 0 {
		t.Log("deserialize to map[string]interface{}")
	} else {
		t.Log("deserialize to struct")
	}

	serializedObj := person{Name: TestName, Age: TestAge}
	marshal, err := currentSerializer.Marshal(serializedObj)
	if err != nil {
		return
	}

	err = currentSerializer.Unmarshal(marshal, &jsonMap)
	if err != nil {
		panic(err)
	}

	assert.True(t, jsonMap.Name == TestName)
	assert.True(t, jsonMap.Age == TestAge)
	assert.True(t, reflect.TypeOf(jsonMap) == reflect.TypeOf(person{}))
	assert.Equal(t, serializedObj, jsonMap)
}

// TestDeserializeStructuredDataStruct2 tests the deserialize structured data into struct.
func TestDeserializeStructuredDataStruct2(t *testing.T) {
	// https://pkg.go.dev/encoding/json#Unmarshal
	t.Helper()

	// when we assign object to explicit struct type, defaultLogger unmarshaler can deserialize it to the struct

	// To unmarshal JSON into a struct, Unmarshal matches incoming object keys to the keys used by Marshal (either the struct field name or its tag), preferring an exact match but also accepting a case-insensitive match.
	var jsonMap interface{} = &person{}

	serializedObj := person{Name: TestName, Age: TestAge}
	marshal, err := currentSerializer.Marshal(serializedObj)
	if err != nil {
		return
	}

	err = currentSerializer.Unmarshal(marshal, jsonMap)
	if err != nil {
		panic(err)
	}

	personMap, ok := jsonMap.(*person)
	if !ok {
		t.Fatal("Failed to convert to *person")
	}

	assert.True(t, personMap.Name == TestName)
	assert.True(t, personMap.Age == TestAge)
	assert.True(t, reflect.TypeOf(jsonMap).Elem() == reflect.TypeOf(person{}))
}

// TestDeserializeStructuredDataPointer tests the deserialize structured data into pointer.
func TestDeserializeStructuredDataPointer(t *testing.T) {
	// https://pkg.go.dev/encoding/json#Unmarshal
	t.Helper()

	// when we assign object to explicit struct type, defaultLogger unmarshaler can deserialize it to the struct

	// To unmarshal JSON into a pointer, Unmarshal first handles the case of the JSON being the JSON literal null. In that case, Unmarshal sets the pointer to nil. Otherwise, Unmarshal unmarshals the JSON into the value pointed at by the pointer. If the pointer is nil, Unmarshal allocates a new value for it to point to.To unmarshal JSON into a struct, Unmarshal matches incoming object keys to the keys used by Marshal (either the struct field name or its tag), preferring an exact match but also accepting a case-insensitive match.
	jsonMap := &person{}
	// var jsonMap *person = nil

	serializedObj := person{Name: TestName, Age: TestAge}
	marshal, err := currentSerializer.Marshal(serializedObj)
	if err != nil {
		return
	}

	err = currentSerializer.Unmarshal(marshal, jsonMap)
	if err != nil {
		panic(err)
	}

	assert.True(t, jsonMap.Name == TestName)
	assert.True(t, jsonMap.Age == TestAge)
	assert.True(t, reflect.TypeOf(jsonMap).Elem() == reflect.TypeOf(person{}))
}

// TestDecodeToMap tests the decode to map.
func TestDecodeToMap(t *testing.T) {
	t.Helper()

	var jsonMap map[string]interface{}

	serializedObj := person{Name: TestName, Age: TestAge}
	marshal, err := currentSerializer.Marshal(serializedObj)
	if err != nil {
		t.Fatal(err)
	}

	// https://pkg.go.dev/encoding/json#Unmarshal
	// To unmarshal a JSON object into a map, Unmarshal first establishes a map to use. If the map is nil, Unmarshal allocates a new map. Otherwise Unmarshal reuses the existing map, keeping existing entries. Unmarshal then stores key-value pairs from the JSON object into the map.
	err = currentSerializer.UnmarshalToMap(marshal, &jsonMap)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, jsonMap["Name"] == TestName)
	assert.True(t, jsonMap["Age"] == float64(TestAge))
}
