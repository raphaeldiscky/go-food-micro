// Package appendresult provides append event result.
package appendresult

// AppendEventsResult is a struct that represents a append event result.
type AppendEventsResult struct {
	GlobalPosition      uint64
	NextExpectedVersion uint64
}

// From creates a new append event result.
func From(globalPosition uint64, nextExpectedVersion uint64) *AppendEventsResult {
	return &AppendEventsResult{
		GlobalPosition:      globalPosition,
		NextExpectedVersion: nextExpectedVersion,
	}
}

// NoOp is a append event result that represents a no operation.
var NoOp = From(0, 0)
