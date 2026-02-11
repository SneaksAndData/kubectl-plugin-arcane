package interfaces

import "github.com/spf13/cobra"

// StreamService defines the interface for managing streaming operations, including starting, stopping, and backfilling streams.
type StreamService interface {

	// Start initiates a stream based on the provided command and arguments.
	Start(cmd *cobra.Command, args []string) error

	// Stop terminates an active stream based on the provided command and arguments.
	Stop(cmd *cobra.Command, args []string) error

	// Backfill performs a backfill operation for a stream based on the provided command and arguments.
	Backfill(cmd *cobra.Command, args []string) error
}
