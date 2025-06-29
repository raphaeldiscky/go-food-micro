// Package es provides a config for the event sourcing.
package es

// Config is a config for the KurrentDB.
type Config struct {
	SnapshotFrequency int64 `json:"snapshotFrequency" validate:"required,gte=0"`
}
