package interfaces

import (
	"context"

	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
)

// BackfillService defines the backfill operations.
type BackfillService interface {

	// Backfill performs a backfill operation for a stream based on the provided command and arguments.
	Backfill(ctx context.Context, parameters *models.BackfillParameters) error
}
