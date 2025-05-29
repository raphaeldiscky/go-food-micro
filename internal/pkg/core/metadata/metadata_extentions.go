// Package metadata provides metadata.
package metadata

import (
	"time"

	json "github.com/goccy/go-json"
)

// GetString is a function that returns the string value of the metadata.
func (m Metadata) GetString(key string) string {
	val, ok := m.Get(key).(string)
	if ok {
		return val
	}

	return ""
}

// GetTime is a function that returns the time value of the metadata.
func (m Metadata) GetTime(key string) time.Time {
	val, ok := m.Get(key).(time.Time)
	if ok {
		return val
	}

	return *new(time.Time)
}

// ToJSON is a function that returns the json value of the metadata.
func (m Metadata) ToJSON() string {
	marshal, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(marshal)
}
