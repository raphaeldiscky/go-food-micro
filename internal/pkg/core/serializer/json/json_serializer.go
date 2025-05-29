// Package json provides a json serializer.
package json

import (
	"log"

	"github.com/TylerBrock/colorjson"
	"github.com/mitchellh/mapstructure"

	json "github.com/goccy/go-json"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer"
)

// jsonSerializer is a struct that represents a json serializer.
type jsonSerializer struct{}

// NewDefaultJsonSerializer is a function that creates a new default json serializer.
func NewDefaultJsonSerializer() serializer.Serializer {
	return &jsonSerializer{}
}

// Marshal is a function that marshals an object.
// https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data
// https://developpaper.com/mapstructure-of-go/
// https://github.com/goccy/go-json

func (s *jsonSerializer) Marshal(v interface{}) ([]byte, error) {
	return Marshal(v)
}

// Unmarshal is a function that unmarshals an object.
// Unmarshal is a wrapper around json.Unmarshal.
// To unmarshal JSON into an interface value, Unmarshal stores in a map[string]interface{}.
func (s *jsonSerializer) Unmarshal(data []byte, v interface{}) error {
	return Unmarshal(data, v)
}

// UnmarshalFromJson is a function that unmarshals an object from json.
// UnmarshalFromJson is a wrapper around json.Unmarshal.
func (s *jsonSerializer) UnmarshalFromJson(data string, v interface{}) error {
	return UnmarshalFromJSON(data, v)
}

// DecodeWithMapStructure is a function that decodes an object with a map structure.
// DecodeWithMapStructure is a wrapper around mapstructure.Decode.
// Decode takes an input structure or map[string]interface{} and uses reflection to translate it to the output structure. output must be a pointer to a map or struct.
// https://pkg.go.dev/github.com/mitchellh/mapstructure#section-readme
func (s *jsonSerializer) DecodeWithMapStructure(
	input interface{},
	output interface{},
) error {
	return DecodeWithMapStructure(input, output)
}

// UnmarshalToMap is a function that unmarshals an object to a map.
func (s *jsonSerializer) UnmarshalToMap(
	data []byte,
	v *map[string]interface{},
) error {
	return UnmarshalToMap(data, v)
}

// UnmarshalToMapFromJSON is a function that unmarshals an object to a map from json.
func (s *jsonSerializer) UnmarshalToMapFromJSON(
	data string,
	v *map[string]interface{},
) error {
	return UnmarshalToMapFromJSON(data, v)
}

// PrettyPrint is a function that pretty prints an object.
func (s *jsonSerializer) PrettyPrint(data interface{}) string {
	return PrettyPrint(data)
}

// ColoredPrettyPrint is a function that pretty prints an object with color.
func (s *jsonSerializer) ColoredPrettyPrint(data interface{}) string {
	return ColoredPrettyPrint(data)
}

// PrettyPrint is a function that pretty prints an object.
func PrettyPrint(data interface{}) string {
	// https://gosamples.dev/pretty-print-json/
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return ""
	}

	return string(val)
}

// ColoredPrettyPrint is a function that pretty prints an object with color.
func ColoredPrettyPrint(data interface{}) string {
	// https://github.com/TylerBrock/colorjson
	var obj map[string]interface{}
	err := json.Unmarshal([]byte(PrettyPrint(data)), &obj)
	if err != nil {
		return ""
	}
	// Make a custom formatter with indent set
	f := colorjson.NewFormatter()
	f.Indent = 4
	val, err := f.Marshal(obj)
	if err != nil {
		return ""
	}

	return string(val)
}

// Marshal is a function that marshals an object.
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal is a function that unmarshals an object.
// Unmarshal is a wrapper around json.Unmarshal.
// To unmarshal JSON into an interface value, Unmarshal stores in a map[string]interface{}.
func Unmarshal(data []byte, v interface{}) error {
	// https://pkg.go.dev/encoding/json#Unmarshal
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	log.Printf("deserialize structure object")

	return nil
}

// UnmarshalFromJSON is a function that unmarshals an object from json.
// UnmarshalFromJSON is a wrapper around json.Unmarshal.
func UnmarshalFromJSON(data string, v interface{}) error {
	err := Unmarshal([]byte(data), v)
	if err != nil {
		return err
	}

	return nil
}

// DecodeWithMapStructure is a function that decodes an object with a map structure.
// DecodeWithMapStructure is a wrapper around mapstructure.Decode.
// Decode takes an input structure or map[string]interface{} and uses reflection to translate it to the output structure. output must be a pointer to a map or struct.
// https://pkg.go.dev/github.com/mitchellh/mapstructure#section-readme
func DecodeWithMapStructure(
	input interface{},
	output interface{},
) error {
	// https://developpaper.com/mapstructure-of-go/
	return mapstructure.Decode(input, output)
}

// UnmarshalToMap is a function that unmarshals an object to a map.
func UnmarshalToMap(
	data []byte,
	v *map[string]interface{},
) error {
	// https://developpaper.com/mapstructure-of-go/
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

// UnmarshalToMapFromJSON is a function that unmarshals an object to a map from json.
func UnmarshalToMapFromJSON(
	data string,
	v *map[string]interface{},
) error {
	return UnmarshalToMap([]byte(data), v)
}
