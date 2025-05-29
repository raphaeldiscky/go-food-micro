// Package contracts provides a migration runner.
package contracts

import "context"

// PostgresMigrationRunner is a postgres migration runner.
type PostgresMigrationRunner interface {
	Up(ctx context.Context, version uint) error
	Down(ctx context.Context, version uint) error
}
