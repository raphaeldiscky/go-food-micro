// Package typeregistry provides a type registry.
package typeregistry

// Ref:https://stackoverflow.com/questions/23030884/is-there-a-way-to-create-an-instance-of-a-struct-from-a-string
import "reflect"

var typeRegistry = make(map[string]reflect.Type)

// registerType registers a type.
func registerType(typedNil interface{}) {
	t := reflect.TypeOf(typedNil).Elem()
	typeRegistry[t.PkgPath()+"."+t.Name()] = t
}

// MyString is a string.
type MyString string

// init initializes the type registry.
func init() {
	registerType((*MyString)(nil))
}
