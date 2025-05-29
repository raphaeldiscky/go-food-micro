// Package contracts provides a health check contracts.
package contracts

// Check is a map of statuses.
type Check map[string]Status

// AllUp is a function that checks if all the statuses are up.
func (check Check) AllUp() bool {
	for _, status := range check {
		if !status.IsUp() {
			return false
		}
	}

	return true
}
