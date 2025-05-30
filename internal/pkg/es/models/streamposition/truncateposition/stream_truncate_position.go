// Package truncateposition provides stream truncate position.
package truncateposition

// StreamTruncatePosition an int64 for accepts negative and positive value.
type StreamTruncatePosition int64

// Value returns the value of the stream truncate position.
func (e StreamTruncatePosition) Value() int64 {
	return int64(e)
}

// FromInt64 returns a stream truncate position from an int64.
func FromInt64(position int64) StreamTruncatePosition {
	return StreamTruncatePosition(position)
}
