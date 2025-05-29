// Package metadata provides metadata.
package metadata

// Metadata is a type that represents a metadata.
type Metadata map[string]interface{}

// ExistsKey is a function that checks if the metadata exists.
func (m Metadata) ExistsKey(key string) bool {
	_, exists := m[key]

	return exists
}

// Get is a function that returns the value of the metadata.
func (m Metadata) Get(key string) interface{} {
	val, exists := m[key]
	if !exists {
		return nil
	}

	return val
}

// Set is a function that sets the value of the metadata.
func (m Metadata) Set(key string, value interface{}) {
	m[key] = value
}

// Keys is a function that returns the keys of the metadata.
func (m Metadata) Keys() []string {
	i := 0
	r := make([]string, len(m))

	for k := range m {
		r[i] = k
		i++
	}

	return r
}

// MapToMetadata is a function that converts a map to a metadata.
func MapToMetadata(data map[string]interface{}) Metadata {
	m := Metadata(data)

	return m
}

// MetadataToMap is a function that converts a metadata to a map.
func MetadataToMap(meta Metadata) map[string]interface{} {
	return meta
}

// FromMetadata is a function that returns the metadata.
func FromMetadata(m Metadata) Metadata {
	if m == nil {
		return Metadata{}
	}

	return m
}
