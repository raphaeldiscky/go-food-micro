// Package attribute provides a module for the attribute.
package attribute

import (
	"go.opentelemetry.io/otel/attribute"

	json "github.com/goccy/go-json"
)

// Object creates a KeyValue with a interface{} value type.
func Object(k string, v interface{}) attribute.KeyValue {
	marshal, err := json.Marshal(&v)
	if err != nil {
		return attribute.KeyValue{}
	}

	return attribute.KeyValue{Key: attribute.Key(k), Value: attribute.StringValue(string(marshal))}
}
