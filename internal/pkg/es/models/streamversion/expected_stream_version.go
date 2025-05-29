// Package streamversion provides expected stream version.
package streamversion

// https://github.com/EventStore/EventStore-Client-Go/blob/1591d047c0c448cacc0468f9af3605572aba7970/esdb/revision.go
// https://github.com/EventStore/EventStore-Client-Dotnet/blob/b8beee7b97ef359316822cb2d00f120bf67bd14d/src/EventStore.Client/StreamRevision.cs

// ExpectedStreamVersion an int64 for accepts negative and positive value.
type ExpectedStreamVersion int64

const (
	// NoStream is a expected stream version that represents a no stream.
	NoStream ExpectedStreamVersion = -1
	// Any is a expected stream version that represents a any.
	Any ExpectedStreamVersion = -2
	// StreamExists is a expected stream version that represents a stream exists.
	StreamExists ExpectedStreamVersion = -3
)

// FromInt64 returns a expected stream version from an int64.
func FromInt64(expectedVersion int64) ExpectedStreamVersion {
	return ExpectedStreamVersion(expectedVersion)
}

// Next returns the next expected stream version.
func (e ExpectedStreamVersion) Next() ExpectedStreamVersion {
	return e + 1
}

// Value returns the value of the expected stream version.
func (e ExpectedStreamVersion) Value() int64 {
	return int64(e)
}

// IsNoStream returns true if the expected stream version is no stream.
func (e ExpectedStreamVersion) IsNoStream() bool {
	return e == NoStream
}

// IsAny returns true if the expected stream version is any.
func (e ExpectedStreamVersion) IsAny() bool {
	return e == Any
}

// IsStreamExists returns true if the expected stream version is stream exists.
func (e ExpectedStreamVersion) IsStreamExists() bool {
	return e == StreamExists
}
