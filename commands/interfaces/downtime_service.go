package interfaces

import "github.com/spf13/cobra"

// DowntimeService defines the interface for downtime management operations.
type DowntimeService interface {

	// DeclareDowntime initiates a downtime period for specified streams based on the provided command and arguments.
	DeclareDowntime(cmd *cobra.Command, args []string) error

	// StopDowntime ends an active downtime period for specified streams based on the provided command and arguments.
	StopDowntime(cmd *cobra.Command, args []string) error
}
