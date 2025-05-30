// Package readposition provides stream read position.
package readposition

import expectedStreamVersion "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamversion"

// https://github.com/EventStore/EventStore-Client-Dotnet/blob/b8beee7b97ef359316822cb2d00f120bf67bd14d/src/EventStore.Client/StreamPosition.cs
// https://github.com/EventStore/EventStore-Client-Go/blob/1591d047c0c448cacc0468f9af3605572aba7970/esdb/position.go

// StreamReadPosition an int64 for accepts negative and positive value.
type StreamReadPosition int64

// Value returns the value of the stream read position.
func (e StreamReadPosition) Value() int64 {
	return int64(e)
}

// IsEnd returns true if the stream read position is end.
func (e StreamReadPosition) IsEnd() bool {
	return e == End
}

// IsStart returns true if the stream read position is start.
func (e StreamReadPosition) IsStart() bool {
	return e == Start
}

// Next returns the next stream read position.
func (e StreamReadPosition) Next() StreamReadPosition {
	return e + 1
}

// Start is the start stream read position.
const Start StreamReadPosition = 0

// End is the end stream read position.
const End StreamReadPosition = -1

// FromInt64 returns a stream read position from an int64.
func FromInt64(position int64) StreamReadPosition {
	return StreamReadPosition(position)
}

// FromStreamRevision returns a stream read position from a stream version.
func FromStreamRevision(
	streamVersion expectedStreamVersion.ExpectedStreamVersion,
) StreamReadPosition {
	return StreamReadPosition(streamVersion.Value())
}
