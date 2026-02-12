package interfaces

import (
	"context"
	"github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
)

// StreamService defines the interface for managing streaming operations, including starting, stopping, and backfilling streams.
type StreamService interface {

	// Start initiates a stream based on the provided command and arguments.
	Start(ctx context.Context, parameters *models.StartParameters) error

	// Stop terminates an active stream based on the provided command and arguments.
	Stop(ctx context.Context, parameters *models.StopParameters) error

	// Backfill performs a backfill operation for a stream based on the provided command and arguments.
	Backfill(ctx context.Context, parameters *models.BackfillParameters) (*v1.BackfillRequest, error)
}
