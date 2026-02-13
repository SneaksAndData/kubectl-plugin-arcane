package interfaces

import (
	"context"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
)

// DowntimeService defines the interface for downtime management operations.
type DowntimeService interface {

	// DeclareDowntime initiates a downtime period for specified streams based on the provided command and arguments.
	DeclareDowntime(ctx context.Context, parameters models.DowntimeDeclareParameters) error

	// StopDowntime ends an active downtime period for specified streams based on the provided command and arguments.
	StopDowntime(ctx context.Context, parameters models.DowntimeStopParameters) error
}
