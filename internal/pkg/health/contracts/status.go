// Package contracts provides a health check contracts.
package contracts

const (
	StatusUp   = "up"
	StatusDown = "down"
)

// Status is a struct that represents a status.
type Status struct {
	Status string `json:"status"`
}

// NewStatus is a function that creates a new status.
func NewStatus(err error) Status {
	if err != nil {
		return Status{Status: StatusDown}
	}

	return Status{Status: StatusUp}
}

// IsUp is a function that checks if the status is up.
func (status Status) IsUp() bool {
	return status.Status == StatusUp
}
