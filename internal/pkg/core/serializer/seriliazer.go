// Package serializer provides a serializer.
package serializer

// Serializer is an interface that represents a serializer.
type Serializer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	UnmarshalFromJson(data string, v interface{}) error
	DecodeWithMapStructure(
		input interface{},
		output interface{},
	) error
	UnmarshalToMap(data []byte, v *map[string]interface{}) error
	UnmarshalToMapFromJSON(data string, v *map[string]interface{}) error
	PrettyPrint(data interface{}) string
	ColoredPrettyPrint(data interface{}) string
}
